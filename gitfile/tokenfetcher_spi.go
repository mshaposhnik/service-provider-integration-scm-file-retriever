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
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mshaposhnik/service-provider-integration-scm-file-retriever/gitfile/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SpiTokenFetcher token fetcher implementation that looks for token in the specific ENV variable.
type SpiTokenFetcher struct {
	k8sClient client.Client
	namespace string
}

const (
	namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	letterBytes   = "abcdefghijklmnopqrstuvwxyz1234567890"
)

func (s *SpiTokenFetcher) BuildHeader(ctx context.Context, repoUrl string, loginCallback func(url string)) (*HeaderStruct, error) {

	var tBindingName = "file-retriever-binging-" + randStringBytes(6)
	var secretName = "file-retriever-secret-" + randStringBytes(5)

	// create binding
	newBinding := newSPIATB(tBindingName, s, repoUrl, secretName)
	err := s.k8sClient.Create(ctx, newBinding)
	if err != nil {
		zap.L().Error("Error creating Token Binding item:", zap.Error(err))
		return nil, err
	}

	// now re-reading SPITokenBinding to get updated fields
	var tokenName string
	for {
		readBinding := &v1beta1.SPIAccessTokenBinding{}
		err = s.k8sClient.Get(ctx, client.ObjectKey{Namespace: s.namespace, Name: tBindingName}, readBinding)
		if err != nil {
			zap.L().Error("Error reading TB item:", zap.Error(err))
			return nil, err
		}
		tokenName = readBinding.Status.LinkedAccessTokenName
		if tokenName != "" {
			break
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("task is cancelled")
		default:
			time.Sleep(200 * time.Millisecond)
		}
	}
	zap.L().Info(fmt.Sprintf("Access Token to watch: %s", tokenName))

	// now try read SPIAccessToken to get link
	var url string
	var loginCalled = false
	for {
		readToken := &v1beta1.SPIAccessToken{}
		_ = s.k8sClient.Get(ctx, client.ObjectKey{Namespace: s.namespace, Name: tokenName}, readToken)
		if readToken.Status.Phase == v1beta1.SPIAccessTokenPhaseAwaitingTokenData && !loginCalled {
			url = readToken.Status.OAuthUrl
			zap.L().Info(fmt.Sprintf("URL to OAUth: %s", url))
			go loginCallback(url)
			loginCalled = true
		} else if readToken.Status.Phase == v1beta1.SPIAccessTokenPhaseReady {
			// now we can exit the loop and read the secret
			break
		}
		time.Sleep(200 * time.Millisecond)
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("task is cancelled")
		default:
			time.Sleep(200 * time.Millisecond)
		}
	}

	// reading token secret
	tokenSecret := &corev1.Secret{}
	err = s.k8sClient.Get(ctx, client.ObjectKey{Namespace: s.namespace, Name: secretName}, tokenSecret)
	if err != nil {
		zap.L().Error("Error reading Token Secret item:", zap.Error(err))
		return nil, err
	}
	return &HeaderStruct{Authorization: "Bearer " + string(tokenSecret.Data["password"])}, nil
}

func newSPIATB(tBindingName string, s *SpiTokenFetcher, repoUrl string, secretName string) *v1beta1.SPIAccessTokenBinding {
	newBinding := &v1beta1.SPIAccessTokenBinding{
		ObjectMeta: metav1.ObjectMeta{Name: tBindingName, Namespace: s.namespace},
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
	return newBinding
}

func newSpiTokenFetcher() *SpiTokenFetcher {
	namespace, err := ioutil.ReadFile(namespaceFile)
	if err != nil {
		panic(err.Error())
	}

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
	return &SpiTokenFetcher{k8sClient: k8sClient, namespace: string(namespace)}
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
