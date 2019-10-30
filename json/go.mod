module github.com/go-4devs/httpclient/json

go 1.12

replace (
	github.com/go-4devs/httpclient/dc => ../dc
	github.com/go-4devs/httpclient/decoder => ../decoder
)

require (
	github.com/go-4devs/httpclient/dc v0.0.1
	github.com/go-4devs/httpclient/decoder v0.0.0-20191030085833-a0493e492141
)
