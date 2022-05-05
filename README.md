# service-provider-integration-scm-file-retriever
[![Build](https://github.com/redhat-appstudio/service-provider-integration-scm-file-retriever/actions/workflows/build.yml/badge.svg?branch=main&event=push)](https://github.com/redhat-appstudio/service-provider-integration-scm-file-retriever/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/redhat-appstudio/service-provider-integration-scm-file-retriever/branch/main/graph/badge.svg?token=MiQMw3V0wG)](https://codecov.io/gh/redhat-appstudio/service-provider-integration-scm-file-retriever)
Library for downloading files from a source code management sites

### About

This repository contains a library for retrieving files from various source management systems using a repository and file paths as the primary form of input.

The main idea is to allow users to download files from a different SCM providers without the necessity of knowing their APIs and/or download endpoints,
as well as take care about the authentication.

### Usage

Import 

```
import (
  "github.com/redhat-appstudio/service-provider-integration-scm-file-retriever/gitfile"
)
```


The main function signature looks as follows:  

```
func getFileContents(ctx context.Context, namespace, repoUrl, filepath, ref string, callback func(ctx context.Context, url string)) (io.ReadCloser, error) 
```
It expects the user namespace name to perform AccessToken related operations, three file location parameters, from which repository URL and path to file are mandatory , and optional ref for the branch/tags.
Function type parameter is a callback used when user authentication is needed, that function will be called with the URL to OAuth service, on which user need to be redirected, and can be controlled using the context.

### URL and path formats
Repository URLs may or may not contain `.git` suffixes. Paths are usual `/a/b/filename` format. Optional `ref` may
contain commit id, tag or branch name.

### Supported SCM providers

 - GitHub



## Demo server application

For the preview and testing purposes, there is demo server application developed, which consists of API endpoint,
simple UI page and websocket connection mechanism. It's source code located under `server` module.

### Building demo server application 

Simplest way to build demo server app is to use docker based build. Simply run `docker build server -t <image_tag>` from the root of repository,
and demo application image will be built.

### Deploying demo server application

There is a bunch of helpful scripts located at `server/hack` which can be used for different deployment scenarios.
The general prerequisite for all deployment types is to have `SPI_GITHUB_CLIENT_ID` and `SPI_GITHUB_CLIENT_SECRET` environment variables to be set locally, containing
correct values from registered GitHub OAuth application. 

#### Deploying on Kubernetes
  ...in progress

#### Deploying on Openshift
 Entry point for Openshift deployment is a `/server/hack/12_oc_deploy.sh` script. Please, note that script should
be executed from the root project folder and not directly from `hack` folder. When executed, script performs installation
of spi-controller, spi-oauth-service and spi-file-retriever-server deployments, so it's not necessary
to pre-install something before neath. Script also performs Vault storage initialization and unseal.
As a result of script execution, there must be three successful deployments in the `spi-system` project,
and the `oauth-secret` secret must be created and filled with correct OAuth authentication data.
Do not forget to align the server hostname in the secret and OAuth application callback URL on GitHub after installation.  

   
### Known  

