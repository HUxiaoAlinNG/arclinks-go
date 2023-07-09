package httpctx

import (
	"context"
	"encoding/json"
	"go.elastic.co/apm/module/apmhttp"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
	"net/url"
	"strings"
)

var tracingClient = apmhttp.WrapClient(http.DefaultClient)

func Get(ctx context.Context, url string, header map[string]string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if header != nil {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}
	return ctxhttp.Do(ctx, tracingClient, request)
}

func PostJson(ctx context.Context, url string, body interface{}, header map[string]string) (*http.Response, error) {
	return DoJson(ctx, "POST", url, body, header)
}

func PutJson(ctx context.Context, url string, body interface{}, header map[string]string) (*http.Response, error) {
	return DoJson(ctx, "PUT", url, body, header)
}

func DoJson(ctx context.Context, method string, url string, body interface{}, header map[string]string) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}
	if header != nil {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}
	request.Header.Set("Content-Type", "application/json")
	return ctxhttp.Do(ctx, tracingClient, request)
}

func PostForm(ctx context.Context, url string, data url.Values, header map[string]string) (*http.Response, error) {
	request, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	if header != nil {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return ctxhttp.Do(ctx, tracingClient, request)
}

func Delete(ctx context.Context, url string, header map[string]string) (*http.Response, error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	if header != nil {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}
	return ctxhttp.Do(ctx, tracingClient, request)
}
