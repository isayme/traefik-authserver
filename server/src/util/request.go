package util

import (
	"context"
	"net/http"

	"github.com/isayme/go-request"
)

var requestOption *request.Option

func init() {
	requestOption = request.New()
	requestOption.UserAgentPrefix = UserAgent
}

// Request 对外发起请求
func Request(ctx context.Context, method, url string, header http.Header, body interface{}, out interface{}) (err error) {
	_, err = requestOption.Request(ctx, method, url, header, body, out)
	return err
}
