module github.com/go-4devs/httpclient/json

go 1.12

replace (
	github.com/go-4devs/httpclient/dc => ../dc
	github.com/go-4devs/httpclient/decoder => ../decoder
	github.com/go-4devs/httpclient/transport => ../transport
)

require (
	github.com/go-4devs/httpclient/apierrors v0.0.0-20190801080035-86bc08daf4fd // indirect
	github.com/go-4devs/httpclient/dc v0.0.0-00010101000000-000000000000
	github.com/go-4devs/httpclient/decoder v0.0.0-00010101000000-000000000000
)
