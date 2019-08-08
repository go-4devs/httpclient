module github.com/go-4devs/httpclient/dc

go 1.12

require (
	github.com/go-4devs/httpclient v0.0.0-20190729052847-527e15269a9c
	github.com/go-4devs/httpclient/apierrors v0.0.0-20190806073327-3cd2c74fdee1
	github.com/go-4devs/httpclient/decoder v0.0.0-20190806073327-3cd2c74fdee1
	github.com/go-4devs/httpclient/transport v0.0.0-00010101000000-000000000000
)

replace (
	github.com/go-4devs/httpclient => ../
	github.com/go-4devs/httpclient/apierror => ../apierror
	github.com/go-4devs/httpclient/transport => ../transport
)
