package utils

import (
	"bytes"
	"fmt"
	"github.com/kosotd/go-service-base/utils/domain"
	"github.com/kosotd/go-service-base/utils/domain/params_type"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func DoRequest(method domain.Method) (int, []byte, error) {
	client := http.Client{}
	if method.Timeout > 0 {
		client.Timeout = method.Timeout
	}

	var body io.Reader
	paramsType, err := params_type.FromString(method.ParamsType)
	if err != nil {
		return 0, nil, errors.Wrap(err, "utils.DoRequest -> params_type.FromString")
	}

	switch paramsType {
	case params_type.FormUrlencoded:
		data := url.Values{}
		for k, v := range method.Params {
			data.Set(k, v)
		}
		body = strings.NewReader(data.Encode())
	case params_type.JsonBody:
		body = bytes.NewReader(method.Body)
	case params_type.PathParams:
		for k, v := range method.Params {
			method.Url = strings.Replace(method.Url, fmt.Sprintf("{%s}", k), v, -1)
		}
	}

	req, err := http.NewRequest(method.Method, method.Url, body)
	if err != nil {
		return 0, nil, errors.Wrap(err, "utils.DoRequest -> http.NewRequest")
	}

	for k, v := range method.Headers {
		req.Header.Set(k, v)
	}
	if paramsType == params_type.FormUrlencoded {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if paramsType == params_type.JsonBody {
		req.Header.Set("Content-Type", "application/json")
	} else if paramsType == params_type.RawQuery {
		q := req.URL.Query()
		for k, v := range method.Params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, errors.Wrap(err, "utils.DoRequest -> client.Do")
	}
	defer CloseSafe(resp.Body)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, errors.Wrap(err, "utils.DoRequest -> ioutil.ReadAll")
	}

	return resp.StatusCode, respBody, nil
}
