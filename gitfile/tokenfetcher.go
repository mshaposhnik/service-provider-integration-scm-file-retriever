// Copyright (c) 2021 - 2022 Red Hat, Inc.
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

// HeaderStruct is the simple struct to carry authentication string from different suppliers
type HeaderStruct struct {
	Authorization string `json:"Authorization"`
}

// TokenFetcher is the interface for the authentication token suppliers which are provides tokens as a HeaderStruct
// instances
type TokenFetcher interface {
	BuildHeader(repoUrl string) HeaderStruct
}

func buildAuthHeader(repoUrl string, fetcher TokenFetcher) HeaderStruct {
	headerStruct := fetcher.BuildHeader(repoUrl)
	if len(headerStruct.Authorization) > 0 {
		return headerStruct
	}
	return HeaderStruct{}
}
