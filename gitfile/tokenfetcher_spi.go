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
	"encoding/base64"
	"fmt"
	"math/rand"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.uber.org/zap"

	"github.com/mshaposhnik/service-provider-integration-scm-file-retriever/gitfile/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SpiTokenFetcher token fetcher implementation that looks for token in the specific ENV variable.
type SpiTokenFetcher struct {
	k8sClient client.Client
}

func (s *SpiTokenFetcher) BuildHeader(ctx context.Context, repoUrl string, loginCallback func(url string)) (HeaderStruct, error) {

	var namespaceName = "default"
	var tBindingName = "file-retriever-binging-" + RandStringBytes(6)
	var secretName = "file-retriever-secret-" + RandStringBytes(5)

	// create binding
	newBinding := &v1beta1.SPIAccessTokenBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tBindingName,
			Namespace: namespaceName,
		},
		Spec: v1beta1.SPIAccessTokenBindingSpec{
			RepoUrl: repoUrl,
			Permissions: v1beta1.Permissions{
				Required: []v1beta1.Permission{
					{
						Type: v1beta1.PermissionTypeReadWrite,
						Area: v1beta1.PermissionAreaRepository,
					},
				},
				AdditionalScopes: []string{"api"},
			},
			Secret: v1beta1.SecretSpec{
				Name: secretName,
				Type: corev1.SecretTypeBasicAuth,
			},
		},
	}
	err := s.k8sClient.Create(ctx, newBinding)
	if err != nil {
		zap.L().Error("Error creating Token Binding item:", zap.Error(err))
		return HeaderStruct{}, err
	}

	// now re-reading SPI TB to get updated fields
	var tokenName string
	for {
		readBinding := &v1beta1.SPIAccessTokenBinding{}
		err = s.k8sClient.Get(ctx, client.ObjectKey{Namespace: namespaceName, Name: tBindingName}, readBinding)
		if err != nil {
			zap.L().Error("Error reading TB item:", zap.Error(err))
			return HeaderStruct{}, err
		}
		tokenName = readBinding.Status.LinkedAccessTokenName
		if tokenName != "" {
			break
		}
	}
	zap.L().Info(fmt.Sprintf("Access Token to watch: %s", tokenName))

	// now try read SPI Token to get link
	var url string
	var loginCalled = false
	for {
		readToken := &v1beta1.SPIAccessToken{}
		_ = s.k8sClient.Get(ctx, client.ObjectKey{Namespace: namespaceName, Name: tokenName}, readToken)
		if readToken.Status.Phase == v1beta1.SPIAccessTokenPhaseAwaitingTokenData && !loginCalled {
			url = readToken.Status.OAuthUrl
			zap.L().Info(fmt.Sprintf("URL to OAUth: %s", url))
			go loginCallback(url)
			loginCalled = true
		} else if readToken.Status.Phase == v1beta1.SPIAccessTokenPhaseReady {
			// now we can exit the loop and read the secret
			break
		}
	}

	tokenSecret := &corev1.Secret{}
	err = s.k8sClient.Get(ctx, client.ObjectKey{Namespace: namespaceName, Name: secretName}, tokenSecret)
	if err != nil {
		zap.L().Error("Error reading Token Secret item:", zap.Error(err))
		return HeaderStruct{}, err
	}

	var data []byte
	_, err = base64.StdEncoding.Decode(data, tokenSecret.Data["password"])
	if err != nil {
		zap.L().Error("Base64 Error:", zap.Error(err))
	}
	return HeaderStruct{Authorization: "Bearer " + string(data)}, nil
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

const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
