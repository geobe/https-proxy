# https-proxy
A reverse proxy in golang that routes https requests to http servers

##Native or Cross Compile

go env -w GOOS=linux GOARCH=arm GOARM=7

go env -w GOOS=windows GOARCH=amd64 

go env -u GOARM

go build -o ..\pkg\secure-proxy github.com\geobe\https-proxy\go\server\secure-proxy.go

## Build Tests
go build -o ..\pkg\dummy-server github.com\geobe\https-proxy\go\main\dummyServer.go

go build -o ..\pkg\redirect-demo github.com\geobe\https-proxy\go\main\portRedirectDemoServer.go
## Needed files
config.json

View template directory
cd ../