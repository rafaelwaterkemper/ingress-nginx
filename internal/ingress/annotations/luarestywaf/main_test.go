/*
Copyright 2017 The Kubernetes Authors.

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

package luarestywaf

import (
	"testing"

	api "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/ingress-nginx/internal/ingress/annotations/parser"
	"k8s.io/ingress-nginx/internal/ingress/resolver"
)

func TestParse(t *testing.T) {
	luaRestyWAFAnnotation := parser.GetAnnotationWithPrefix("lua-resty-waf")
	luaRestyWAFDebugAnnotation := parser.GetAnnotationWithPrefix("lua-resty-waf-debug")
	luaRestyWAFIgnoredRuleSetsAnnotation := parser.GetAnnotationWithPrefix("lua-resty-waf-ignore-rulesets")

	ap := NewParser(&resolver.Mock{})
	if ap == nil {
		t.Fatalf("expected a parser.IngressAnnotation but returned nil")
	}

	testCases := []struct {
		annotations map[string]string
		expected    *Config
	}{
		{nil, &Config{}},
		{map[string]string{}, &Config{}},

		{map[string]string{luaRestyWAFAnnotation: "true"}, &Config{Enabled: true, Debug: false, IgnoredRuleSets: []string{}}},
		{map[string]string{luaRestyWAFDebugAnnotation: "true"}, &Config{Enabled: false, Debug: false}},

		{map[string]string{luaRestyWAFAnnotation: "true", luaRestyWAFDebugAnnotation: "true"}, &Config{Enabled: true, Debug: true, IgnoredRuleSets: []string{}}},
		{map[string]string{luaRestyWAFAnnotation: "true", luaRestyWAFDebugAnnotation: "false"}, &Config{Enabled: true, Debug: false, IgnoredRuleSets: []string{}}},
		{map[string]string{luaRestyWAFAnnotation: "false", luaRestyWAFDebugAnnotation: "true"}, &Config{Enabled: false, Debug: true, IgnoredRuleSets: []string{}}},

		{map[string]string{
			luaRestyWAFAnnotation:                "true",
			luaRestyWAFDebugAnnotation:           "true",
			luaRestyWAFIgnoredRuleSetsAnnotation: "ruleset1, ruleset2 ruleset3,   another.ruleset"},
			&Config{Enabled: true, Debug: true, IgnoredRuleSets: []string{"ruleset1", "ruleset2", "ruleset3", "another.ruleset"}}},
	}

	ing := &extensions.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
		Spec: extensions.IngressSpec{},
	}

	for _, testCase := range testCases {
		ing.SetAnnotations(testCase.annotations)
		result, _ := ap.Parse(ing)
		config := result.(*Config)
		if !config.Equal(testCase.expected) {
			t.Errorf("expected %v but returned %v, annotations: %s", testCase.expected, result, testCase.annotations)
		}
	}
}