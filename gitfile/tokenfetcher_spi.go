// Copyright (c) 2022 Red Hat, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitfile

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.uber.org/zap"

	"github.com/mshaposhnik/service-provider-integration-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SpiTokenFetcher token fetcher implementation that looks for token in the specific ENV variable.
type SpiTokenFetcher struct {
	k8sClient client.Client
}

func (s *SpiTokenFetcher) BuildHeader(repoUrl string) HeaderStruct {
	newtb := &v1beta1.SPIAccessTokenBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "new-test-binding",
			Namespace: "default",
		},
		Spec: v1beta1.SPIAccessTokenBindingSpec{
			RepoUrl:     repoUrl,
			Permissions: v1beta1.Permissions{},
			Secret: v1beta1.SecretSpec{
				Name: "token-secret",
				Type: corev1.SecretTypeBasicAuth,
			},
		},
	}
	err := s.k8sClient.Create(context.Background(), newtb)
	if err != nil {
		zap.L().Error("Error creating item:", zap.Error(err))
		return HeaderStruct{}
	}
	return HeaderStruct{}
}

func newSpiTokenFetcher() *SpiTokenFetcher {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	scheme := runtime.NewScheme()
	if err = corev1.AddToScheme(scheme); err != nil {
		panic(err.Error())
	}

	if err = v1beta1.AddToScheme(scheme); err != nil {
		panic(err.Error())
	}

	// creates the client
	k8sClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		panic(err.Error())
	}
	return &SpiTokenFetcher{k8sClient: k8sClient}
}
