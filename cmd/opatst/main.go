package main

import (
	"context"
	"os"

	"github.com/cloudflare/cfssl/log"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

var str string = `package security.knative.dev

default allow = false

allow {
  input.httpRequest.header["Ce-Specversion"][_] == "1.0"
  contains(input.httpRequest.header["Ce-Id"][_], "cshou")
  endswith(input.httpRequest.header["Ce-Custom"][_], "haha")
  input.httpRequest.header["Ce-Type"][_] == "cloud-storage"
}
`

// startswith(input.httpRequest.header["Ce-Source"][_], "Google")

type Input struct {
	Value string `json:"value,omitempty"`
}

func main() {
	compiler, err := ast.CompileModules(map[string]string{
		"zzzz": str,
	})
	if err != nil {
		log.Errorf("failed to parse policy: %v", err)
		os.Exit(1)
	}

	r := rego.New(
		rego.Query("data.security.knative.dev.allow"),
		rego.Compiler(compiler),
		rego.Dump(os.Stdout),
	)

	ctx := context.Background()

	pq, err := r.PrepareForEval(ctx)
	if err != nil {
		log.Errorf("failed to prepare for eval: %v", err)
		os.Exit(1)
	}

	i := &Input{Value: "xbc-hahaha"}
	rs, err := pq.Eval(ctx, rego.EvalInput(i))

	if err != nil {
		log.Errorf("failed to evaluate input and will fail-close: %v", err)
		os.Exit(1)
	}

	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		os.Exit(1)
	}

	log.Infof("Eval result: %v", rs[0])
}
