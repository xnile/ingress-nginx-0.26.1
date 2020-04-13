/*
Copyright 2015 The Kubernetes Authors.

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

package defaultbackend

import (
	"testing"

	api "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/ingress-nginx/internal/ingress/annotations/parser"
	"k8s.io/ingress-nginx/internal/ingress/errors"
	"k8s.io/ingress-nginx/internal/ingress/resolver"

	"k8s.io/apimachinery/pkg/util/intstr"
)

func buildIngress() *networking.Ingress {
	defaultBackend := networking.IngressBackend{
		ServiceName: "default-backend",
		ServicePort: intstr.FromInt(80),
	}

	return &networking.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
		Spec: networking.IngressSpec{
			Backend: &networking.IngressBackend{
				ServiceName: "default-backend",
				ServicePort: intstr.FromInt(80),
			},
			Rules: []networking.IngressRule{
				{
					Host: "foo.bar.com",
					IngressRuleValue: networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Path:    "/foo",
									Backend: defaultBackend,
								},
							},
						},
					},
				},
			},
		},
	}
}

type mockService struct {
	resolver.Mock
}

// GetService mocks the GetService call from the defaultbackend package
func (m mockService) GetService(name string) (*api.Service, error) {
	if name != "default/demo-service" {
		return nil, errors.Errorf("there is no service with name %v", name)
	}

	return &api.Service{
		ObjectMeta: meta_v1.ObjectMeta{
			Namespace: api.NamespaceDefault,
			Name:      "demo-service",
		},
	}, nil
}

func TestAnnotations(t *testing.T) {
	ing := buildIngress()

	data := map[string]string{}
	data[parser.GetAnnotationWithPrefix("default-backend")] = "demo-service"
	ing.SetAnnotations(data)

	fakeService := &mockService{}
	i, err := NewParser(fakeService).Parse(ing)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	svc, ok := i.(*api.Service)
	if !ok {
		t.Errorf("expected *api.Service but got %v", svc)
	}
	if svc.Name != "demo-service" {
		t.Errorf("expected %v but got %v", "demo-service", svc.Name)
	}
}
