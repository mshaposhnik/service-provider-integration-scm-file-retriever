module service-provider-integration-scm-file-retriever-server

go 1.17

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/redhat-appstudio/service-provider-integration-scm-file-retriever v0.4.6
)

replace github.com/redhat-appstudio/service-provider-integration-scm-file-retriever => github.com/mshaposhnik/service-provider-integration-scm-file-retriever v0.0.0-20220429113239-676312f656db
