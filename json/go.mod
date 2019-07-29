module github.com/go-4devs/httpclient/json

go 1.12

require (
	github.com/go-4devs/httpclient/dc v0.0.0-00010101000000-000000000000
	github.com/go-4devs/httpclient/transport v0.0.0-00010101000000-000000000000
)

replace (
	github.com/go-4devs/httpclient/dc => ../dc
	github.com/go-4devs/httpclient/transport => ../transport
)
