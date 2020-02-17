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
	"fmt"
	"strings"
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

type PolicyBuilder struct {
	rules []*RuleBuilder
}

func NewPolicyBuilder() *PolicyBuilder {
	return &PolicyBuilder{rules: []*RuleBuilder{}}
}

func (pb *PolicyBuilder) NewRule() *RuleBuilder {
	ret := &RuleBuilder{strs: &strings.Builder{}}
	pb.rules = append(pb.rules, ret)
	return ret
}

func (pb *PolicyBuilder) String() string {
	var rules []string
	for _, r := range pb.rules {
		rules = append(rules, r.String())
	}
	combined := strings.Join(rules, "\n")
	return GenerateFromTemplate(combined)
}

type RuleBuilder struct {
	index int
	strs  *strings.Builder
}

func (rb *RuleBuilder) AppendOneOf(path string, allowed []string) {
	prefix := []string{}
	suffix := []string{}
	reg := []string{}
	for _, v := range allowed {
		if strings.HasSuffix(v, "*") {
			prefix = append(prefix, fmt.Sprintf("%q", strings.TrimSuffix(v, "*")))
		} else if strings.HasPrefix(v, "*") {
			suffix = append(suffix, fmt.Sprintf("%q", strings.TrimPrefix(v, "*")))
		} else {
			reg = append(reg, fmt.Sprintf("%q", v))
		}
	}

	if len(prefix) > 0 {
		rb.strs.WriteString(fmt.Sprintf("pre%d := [%s]\n", rb.index, strings.Join(prefix, ",")))
		rb.strs.WriteString(fmt.Sprintf("startswith(%s, pre%d[_])\n", path, rb.index))
	}
	if len(suffix) > 0 {
		rb.strs.WriteString(fmt.Sprintf("suff%d := [%s]\n", rb.index, strings.Join(suffix, ",")))
		rb.strs.WriteString(fmt.Sprintf("endswith(%s, suff%d[_])\n", path, rb.index))
	}
	if len(reg) > 0 {
		rb.strs.WriteString(fmt.Sprintf("reg%d := [%s]\n", rb.index, strings.Join(reg, ",")))
		rb.strs.WriteString(fmt.Sprintf("re_match(reg%d[_], %s)\n", rb.index, path))
	}
	rb.index++
}

func (rb *RuleBuilder) String() string {
	return fmt.Sprintf("allow {\n%s}", rb.strs.String())
}
