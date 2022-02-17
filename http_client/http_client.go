package http_client

import (
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 5, //超时时间
}

func init() {
}

func Exec(method string, url string, header map[string]string, cookies []*http.Cookie) (*http.Response, error) {
	var body io.Reader
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k, v := range header {
			request.Header.Add(k, v)
		}
	}
	if cookies != nil {
		for _, v := range cookies {
			request.AddCookie(v)

		}
	}
	return client.Do(request)
}


