package requestor

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/hugebear-io/gofiber/fabric"
)

func PrepareHttpRequest(method, url string, headers map[string]string, data interface{}) (*http.Request, error) {
	var req *http.Request
	var err error

	if data != nil {
		encoded_data, err := sonic.Marshal(data)
		if err != nil {
			return nil, err
		}

		buffered_data := bytes.NewBuffer(encoded_data)
		req, err = http.NewRequest(method, url, buffered_data)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	return req, nil
}

func PrepareHttpResponse[R interface{}](req *http.Request) (*R, int, error) {
	// request to endpoint
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	// read a bytes of data
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// check empty response body
	if len(resBody) == 0 {
		return nil, res.StatusCode, nil
	}

	// decode response body
	var result R
	if err := fabric.Recast(resBody, &result); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &result, res.StatusCode, nil
}

func HttpGet[R any](url string, headers map[string]string) (*R, int, error) {
	req, err := PrepareHttpRequest(http.MethodGet, url, headers, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return PrepareHttpResponse[R](req)
}

func HttpPost[R any](url string, headers map[string]string, data interface{}) (*R, int, error) {
	fmt.Println(url)
	req, err := PrepareHttpRequest(http.MethodPost, url, headers, data)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return PrepareHttpResponse[R](req)
}

func HttpPut[R any](url string, headers map[string]string, data interface{}) (*R, int, error) {
	req, err := PrepareHttpRequest(http.MethodPut, url, headers, data)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return PrepareHttpResponse[R](req)
}

func HttpPatch[R any](url string, headers map[string]string, data interface{}) (*R, int, error) {
	req, err := PrepareHttpRequest(http.MethodPatch, url, headers, data)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return PrepareHttpResponse[R](req)
}

func HttpDelete[R any](url string, headers map[string]string, data interface{}) (*R, int, error) {
	req, err := PrepareHttpRequest(http.MethodDelete, url, headers, data)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return PrepareHttpResponse[R](req)
}
