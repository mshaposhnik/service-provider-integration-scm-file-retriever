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
func getFileContents(ctx context.Context, repoUrl, filepath, ref string, callback func(url string)) (io.ReadCloser, error) 
```
It expects three file location parameters, from which repository URL and path to file are mandatory , and optional ref for the branch/tags.
Function type parameter is a callback used when user authentication is needed, that function will be called with the URL to OAuth service, on which user need to be redirected.

### URL and path formats
Repository URLs may or may not contain `.git` suffixes. Paths are usual `/a/b/filename` format. Optional `ref` may
contain commit id, tag or branch name.

### Supported SCM providers

 - GitHub
