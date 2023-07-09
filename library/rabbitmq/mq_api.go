package rabbitmq

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

var Api *RbApi

// RabbitMQ HTTP API docs:
// https://rawcdn.githack.com/rabbitmq/rabbitmq-management/v3.8.3/priv/www/api/index.html
type RbApi struct {
	client *fasthttp.Client
	uri    *fasthttp.URI
}

type RbApiConfig struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Vhost    string `yaml:"Vhost"`
	ApiPort  string `yaml:"ApiPort"`
}

func NewRabbitMQApi(config RbApiConfig) *RbApi {
	client := fasthttp.Client{}
	uri := &fasthttp.URI{}
	apiPort := "15672"
	if config.ApiPort != "" {
		apiPort = config.ApiPort
	}
	uri.SetHost(config.Host + ":" + apiPort)
	uri.SetUsername(config.User)
	uri.SetPassword(config.Password)
	Api = &RbApi{
		client: &client,
		uri:    uri,
	}

	return Api
}

// The procedure that actually calls the HTTP request
func (r *RbApi) do(method string, params interface{}) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(r.uri.String())
	req.Header.SetMethod(method)
	basicAuth := fmt.Sprintf("%s:%s", r.uri.Username(), r.uri.Password())
	encode := base64.StdEncoding.EncodeToString([]byte(basicAuth))
	req.Header.Add("Authorization", "Basic "+encode)
	req.Header.SetContentType("application/json")
	// 添加参数
	if params != nil {
		params, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		req.AppendBody(params)
	}

	resp := fasthttp.AcquireResponse()
	err := r.client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	// 状态码非2XX 返回错误信息
	code := strconv.Itoa(resp.StatusCode())
	if !strings.HasPrefix(string(code), "2") {
		requestInfo, _ := json.Marshal(map[string]interface{}{
			"status_code": code,
			"body":        string(resp.Body()),
		})
		return resp, errors.New(string(requestInfo))
	}
	return resp, nil
}

// adds a new Vhost to RabbitMQ.
func (r *RbApi) AddVhost(vhost string) (*fasthttp.Response, error) {
	r.uri.SetPath("/api/vhosts/" + vhost)
	return r.do(fasthttp.MethodPut, nil)
}

// delete a Vhost from RabbitMQ.
func (r *RbApi) DeleteVhost(vhost string) (*fasthttp.Response, error) {
	r.uri.SetPath("/api/vhosts/" + vhost)
	return r.do(fasthttp.MethodDelete, nil)
}

// list Vhosts from RabbitMQ.
func (r *RbApi) ListVhosts() (*fasthttp.Response, error) {
	r.uri.SetPath("/api/vhosts/")
	return r.do(fasthttp.MethodGet, nil)
}

// list Vhosts from RabbitMQ.
func (r *RbApi) SetPermissions(vhost, user string) (*fasthttp.Response, error) {
	params := map[string]interface{}{
		"configure": ".*",
		"write":     ".*",
		"read":      ".*",
	}
	r.uri.SetPath(fmt.Sprintf("/api/permissions/%s/%s", vhost, user))
	return r.do(fasthttp.MethodPut, params)
}

// Publish a message to a specified exchange with routeKey and Vhost
func (r *RbApi) Publish(vhost, routeKey, payload string, exchangeName string) (*fasthttp.Response, error) {
	r.uri.SetPath(fmt.Sprintf("/api/exchanges/%s/%s/publish", vhost, exchangeName))
	params := map[string]interface{}{
		"properties":       map[int]int{},
		"payload":          payload,
		"routing_key":      routeKey,
		"payload_encoding": "string",
	}
	return r.do(fasthttp.MethodPost, params)
}
