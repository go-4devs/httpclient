package retry

import (
	"net/http"
	"time"

	"github.com/go-4devs/httpclient/transport"
)

type config struct {
	backOff        func(uint) time.Duration
	checkRetriable []func(res *http.Response) bool
}

func (c *config) isRetriable(res *http.Response) bool {
	if res == nil {
		return false
	}
	for _, f := range c.checkRetriable {
		if f(res) {
			return true
		}
	}

	return false
}

// Option configure retry
type Option func(c *config)

// WithStatusCode5XX check when status code more or equals 500
func WithStatusCode5XX() Option {
	return WithRetriable(func(res *http.Response) bool {
		return res.StatusCode >= 500
	})
}

// WithStatusCode set error codes
func WithStatusCode(codes ...int) Option {
	return WithRetriable(func(res *http.Response) bool {
		for _, c := range codes {
			if res.StatusCode == c {
				return true
			}
		}
		return false
	})
}

// WithBackOff configure BackOff
func WithBackOff(fn func(uint) time.Duration) Option {
	return func(c *config) {
		c.backOff = fn
	}
}

// WithBackOffLinear set linear back off
func WithBackOffLinear(timeout time.Duration) Option {
	return WithBackOff(func(uint) time.Duration {
		return timeout
	})
}

// WithRetriable check response to retry
// nolint: bodyclose
func WithRetriable(hr ...func(*http.Response) bool) Option {
	return func(c *config) {
		c.checkRetriable = append(c.checkRetriable, hr...)
	}
}

// New create new retry middleware
func New(retry uint, opts ...Option) transport.Middleware {
	cfg := &config{}
	WithBackOffLinear(time.Millisecond * 20)(cfg)

	for _, o := range opts {
		o(cfg)
	}

	if len(cfg.checkRetriable) == 0 {
		WithStatusCode5XX()(cfg)
	}

	return func(r *http.Request, n func(r *http.Request) (*http.Response, error)) (*http.Response, error) {
		var do uint
		res, err := n(r)
		for retry > do && (err != nil || cfg.isRetriable(res)) {
			select {
			case <-r.Context().Done():
				return res, err
			case <-time.After(cfg.backOff(do)):
				if res != nil {
					_ = res.Body.Close()
				}
			}
			do++
			res, err = n(r)
			if err != nil {
				continue
			}
		}

		return res, err
	}
}
