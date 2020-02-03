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

package eventpolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	cehttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	security "github.com/yolocs/knative-policy-binding/pkg/apis/security/v1alpha1"
)

var (
	codecV03 *cehttp.CodecV03 = &cehttp.CodecV03{}
	codecV1  *cehttp.CodecV1  = &cehttp.CodecV1{}
)

func coreSetters(e *cloudevents.Event) map[string]func(string) {
	return map[string]func(string){
		"id":          e.SetID,
		"type":        e.SetType,
		"source":      e.SetSource,
		"subject":     e.SetSubject,
		"dataschema":  e.SetDataSchema,
		"contenttype": e.SetDataContentType,
		"encoding":    e.SetDataContentEncoding,
	}
}

func MakeOpenPolicyRule(eventRules [][]security.EventPolicyRule) string {
	rgs := []string{}
	for _, eg := range eventRules {
		rgs = append(rgs, makeRuleGroup(eg))
	}
	return strings.Join(rgs, "\n")
}

func makeRuleGroup(rules []security.EventPolicyRule) string {
	v1Ev := &cloudevents.Event{}
	v1Ev.SetSpecVersion(cloudevents.CloudEventsVersionV1)
	v1Setters := coreSetters(v1Ev)

	v03Ev := &cloudevents.Event{}
	v03Ev.SetSpecVersion(cloudevents.CloudEventsVersionV03)
	v03Setters := coreSetters(v03Ev)

	m := make(map[string]security.EventPolicyRule)

	for _, r := range rules {
		m[r.Name] = r
		if strings.HasPrefix(r.Name, "ext:") {
			v1Ev.SetExtension(strings.TrimPrefix(r.Name, "ext:"), r.Name)
			v03Ev.SetExtension(strings.TrimPrefix(r.Name, "ext:"), r.Name)
		} else {
			if setFunc, ok := v1Setters[strings.ToLower(r.Name)]; ok {
				setFunc(r.Name)
			}
			if setFunc, ok := v03Setters[strings.ToLower(r.Name)]; ok {
				setFunc(r.Name)
			}
		}
	}

	v1B := &strings.Builder{}
	if msg, err := codecV1.Encode(context.Background(), *v1Ev); err == nil {
		v1B.WriteString(fmt.Sprintf("  input.httpRequest.header[%q][_] == %q\n", "Ce-Specversion", cloudevents.CloudEventsVersionV1))
		httpMsg := msg.(*cehttp.Message)
		for k, h := range httpMsg.Header {
			for _, v := range h {
				if r, ok := m[v]; ok {
					if r.ExactMatch != "" {
						v1B.WriteString(fmt.Sprintf("  input.httpRequest.header[%q][_] == %q\n", k, r.ExactMatch))
					}
					if r.PrefixMatch != "" {
						v1B.WriteString(fmt.Sprintf("  startswith(input.httpRequest.header[%q][_], %q)\n", k, r.PrefixMatch))
					}
					if r.SuffixMatch != "" {
						v1B.WriteString(fmt.Sprintf("  endswith(input.httpRequest.header[%q][_], %q)\n", k, r.SuffixMatch))
					}
					if r.ContainsMatch != "" {
						v1B.WriteString(fmt.Sprintf("  contains(input.httpRequest.header[%q][_], %q)\n", k, r.ContainsMatch))
					}
				}
			}
		}
	}

	v03B := &strings.Builder{}
	if msg, err := codecV03.Encode(context.Background(), *v03Ev); err == nil {
		v03B.WriteString(fmt.Sprintf("  input.httpRequest.header[%q][_] == %q\n", "Ce-Specversion", cloudevents.CloudEventsVersionV03))
		httpMsg := msg.(*cehttp.Message)
		for k, h := range httpMsg.Header {
			for _, v := range h {
				if r, ok := m[v]; ok {
					if r.ExactMatch != "" {
						v03B.WriteString(fmt.Sprintf("  input.httpRequest.header[%q][_] == %q\n", k, r.ExactMatch))
					}
					if r.PrefixMatch != "" {
						v03B.WriteString(fmt.Sprintf("  startswith(input.httpRequest.header[%q][_], %q)\n", k, r.PrefixMatch))
					}
					if r.SuffixMatch != "" {
						v03B.WriteString(fmt.Sprintf("  endswith(input.httpRequest.header[%q][_], %q)\n", k, r.SuffixMatch))
					}
					if r.ContainsMatch != "" {
						v03B.WriteString(fmt.Sprintf("  contains(input.httpRequest.header[%q][_], %q)\n", k, r.ContainsMatch))
					}
				}
			}
		}
	}

	return strings.Join([]string{wrapRule(v1B.String()), wrapRule(v03B.String())}, "\n")
}

func wrapRule(r string) string {
	return fmt.Sprintf("allow {\n%s}\n", r)
}
