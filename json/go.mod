module github.com/go-4devs/httpclient/json

go 1.12

replace (
	github.com/go-4devs/httpclient/dc => ../dc
	github.com/go-4devs/httpclient/decoder => ../decoder
)

require (
	github.com/go-4devs/httpclient/dc v0.0.1
	github.com/go-4devs/httpclient/decoder v0.0.0-20190822111646-5ff1b8dab5b4
)
