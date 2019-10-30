module github.com/go-4devs/httpclient/dc

go 1.12

replace (
	github.com/go-4devs/httpclient => ../
	github.com/go-4devs/httpclient/apierrors => ../apierrors
	github.com/go-4devs/httpclient/decoder => ../decoder
	github.com/go-4devs/httpclient/transport => ../transport
)

require (
	github.com/go-4devs/httpclient v0.0.2
	github.com/go-4devs/httpclient/apierrors v0.0.0-20191030085833-a0493e492141
	github.com/go-4devs/httpclient/decoder v0.0.0-20191030085833-a0493e492141
	github.com/go-4devs/httpclient/transport v0.0.1
	github.com/stretchr/testify v1.4.0
)
