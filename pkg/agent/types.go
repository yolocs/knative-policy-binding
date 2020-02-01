/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package agent

import "net/http"

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
