package cuhttp

import (
	groq "github.com/learnLi/groq_client"
	"io"
	"net/http"
	"net/url"
)

type BasicClient struct {
	client *http.Client
}

func NewBasicClient() *BasicClient {
	return &BasicClient{
		client: &http.Client{},
	}
}

func handlerHeaders(req *http.Request, headers groq.Headers) {
	if headers == nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func handlerCookies(req *http.Request, cookies []*http.Cookie) {
	if cookies == nil {
		return
	}
	for _, v := range cookies {
		req.AddCookie(v)
	}
}

func (b BasicClient) Request(method string, url string, headers groq.Headers, cookies []*http.Cookie, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	handlerHeaders(req, headers)
	handlerCookies(req, cookies)
	return b.client.Do(req)
}

func (b BasicClient) SetProxy(proxy string) {
	if proxy == "" {
		return
	}
	parse, err := url.Parse(proxy)
	if err != nil {
		return
	}

	b.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(parse),
	}
}
