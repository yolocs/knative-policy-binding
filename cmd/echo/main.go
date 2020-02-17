package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type DecisionRequest struct {
	Source      Source              `json:"source,omitempty"`
	Protocol    string              `json:"protocol,omitempty"`
	HTTPRequest *PartialHTTPRequest `json:"httpRequest,omitempty"`
}

type Source struct {
	Issuer   string `json:"issuer,omitempty"`
	Audience string `json:"audience,omitempty"`
	Identity string `json:"identity,omitempty"`
}

type PartialHTTPRequest struct {
	Method        string                 `json:"method,omitempty"`
	Host          string                 `json:"host,omitempty"`
	Path          string                 `json:"path,omitempty"`
	Header        http.Header            `json:"header,omitempty"`
	ContentLength int64                  `json:"contentLength,omitempty"`
	RemoteAddr    string                 `json:"remoteAddr,omitempty"`
	Body          map[string]interface{} `json:"body,omitempty"`
}

type DecisionResponse struct {
	Allow bool `json:"allow"`
}

type config struct {
	DeciderURL string `envconfig:"K_POLICY_DECIDER"`
}

func main() {
	var env config
	if err := envconfig.Process("", &env); err != nil {
		log.Fatal("Failed to process env var", zap.Error(err))
	}

	deciderClient := http.DefaultClient

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Received request %s %q\n", req.Method, req.URL.Path)

		if env.DeciderURL != "" {
			partial := &PartialHTTPRequest{
				Method:        req.Method,
				Host:          req.Host,
				Path:          req.URL.Path,
				ContentLength: req.ContentLength,
				Header:        req.Header,
				RemoteAddr:    req.RemoteAddr,
			}
			dr := &DecisionRequest{
				Protocol:    "http",
				HTTPRequest: partial,
			}

			// ignore body for now.

			b, err := json.Marshal(dr)
			if err != nil {
				log.Printf("failed to marshal decision request: %v", err)
				// fail close
				w.WriteHeader(http.StatusForbidden)
				return
			}

			dresp, err := deciderClient.Post(env.DeciderURL, "application/json", bytes.NewReader(b))
			if err != nil {
				// fail close
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if dresp.StatusCode != http.StatusOK {
				log.Printf("Decision server error: %v\n", dresp.StatusCode)
				// fail close
				w.WriteHeader(http.StatusForbidden)
				return
			}

			rb, err := ioutil.ReadAll(dresp.Body)
			if err != nil {
				log.Printf("Failed to read decision response: %v\n", err)
				// fail close
				w.WriteHeader(http.StatusForbidden)
				return
			}
			var decision DecisionResponse
			if err := json.Unmarshal(rb, &decision); err != nil {
				log.Printf("Failed to unmarshal decision response: %v\n", err)
				// fail close
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if !decision.Allow {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer req.Body.Close()

		respStr := fmt.Sprintf("You made request %s %s%s with:\n===\nHeaders: %v\n===\nBody: %s\n", req.Method, req.Host, req.URL.Path, req.Header, string(b))
		w.Write([]byte(respStr))
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":5678", nil)
}
