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

package opa

import (
	"bytes"
	"text/template"
)

func init() {
	rt, err := template.New("policy").Parse("package security.knative.dev\n\ndefault allow = false\n\n{{.CustomRules}}")
	if err != nil {
		panic(err)
	}
	regoTemplate = rt
}

var regoTemplate *template.Template

type PolicyTemplate struct {
	CustomRules string
}

func GenerateFromTemplate(rules string) string {
	buf := bytes.NewBuffer([]byte{})
	if err := regoTemplate.Execute(buf, &PolicyTemplate{CustomRules: rules}); err != nil {
		panic(err)
	}
	return buf.String()
}
