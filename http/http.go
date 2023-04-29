package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HttpRequest struct {
	Method   string
	Url      string
	Form     *map[string]string
	Params   *map[string]string
	Headers  *map[string]string
	Json     *map[string]interface{}
	Body     []byte
	request  *http.Request
	Response *http.Response
	Result   []byte
	//LoggerPrint bool
	cookies []*http.Cookie
}

func (h *HttpRequest) SetHeader(key string, value string) *HttpRequest {
	if h.Headers == nil {
		h.Headers = &map[string]string{}
	}
	(*h.Headers)[key] = value
	return h
}

func (h *HttpRequest) SetParams(key string, value string) *HttpRequest {
	if h.Params == nil {
		h.Params = &map[string]string{}
	}
	(*h.Params)[key] = value
	return h
}

func (h *HttpRequest) SetForm(key string, value string) *HttpRequest {
	if h.Form == nil {
		h.Form = &map[string]string{}
	}
	(*h.Form)[key] = value
	return h
}

func (h *HttpRequest) SetJson(key string, value interface{}) *HttpRequest {
	if h.Json == nil {
		h.Json = &map[string]interface{}{}
	}
	(*h.Json)[key] = value
	return h
}

func (h *HttpRequest) JsonDumps() ([]byte, error) {
	row, err := json.Marshal(h.Json)
	return row, err
}

func (h *HttpRequest) SetBody(data []byte) {
	h.Body = data
}

func (h *HttpRequest) SetCookie(cookie *http.Cookie) {
	h.cookies = append(h.cookies, cookie)
}

func (h *HttpRequest) UrlParse() {
	path, _ := url.Parse(h.Url)
	h.Url = fmt.Sprintf("%s://%s%s", path.Scheme, path.Host, path.Path)

	if len(path.RawQuery) > 0 {
		for _, item := range strings.Split(path.RawQuery, "&") {
			p := strings.Split(item, "=")
			tmp, _ := url.QueryUnescape(p[1])
			h.SetParams(p[0], tmp)
		}
	}
}

func (h *HttpRequest) Do() error {

	if err := h.structureRequest(); err != nil {
		return err
	}

	var (
		client http.Client
		err    error
	)

	if h.Response, err = client.Do(h.request); err != nil {
		return err
	}

	defer func() {
		_ = h.Response.Body.Close()
	}()
	h.Result, err = ioutil.ReadAll(h.Response.Body)
	if err != nil {
		return err
	}

	//fmt.Println(string(h.Result))
	if h.Response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status %d error!", h.Response.StatusCode))
	}
	return err
}

func (h *HttpRequest) To(v interface{}) error {
	if err := json.Unmarshal(h.Result, v); err != nil {
		return err
	}
	return nil
}

func (h *HttpRequest) structureRequest() error {

	var err error
	baseUrl, err := url.Parse(h.Url)
	if err != nil {
		return err
	}

	if len(h.Method) <= 0 {
		h.Method = "GET"
	}

	p := url.Values{}
	if h.Params != nil {
		for k, v := range *h.Params {
			p.Add(k, v)
		}
	}
	baseUrl.RawQuery = p.Encode()

	var body io.Reader

	if h.Body != nil {
		body = bytes.NewBuffer(h.Body)
	} else if h.Json != nil {
		tmpBody, err := h.JsonDumps()
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(tmpBody)
	} else {
		d := url.Values{}
		if h.Form != nil {
			for k, v := range *h.Form {
				d.Add(k, v)
			}
		}
		body = bytes.NewBufferString(d.Encode())
	}

	u := baseUrl.String()

	h.request, err = http.NewRequest(h.Method, u, body)
	if err != nil {
		return err
	}

	if len(h.cookies) > 0 {
		for _, v := range h.cookies {
			h.request.AddCookie(v)
		}
	}

	if h.Headers != nil {
		for k, v := range *h.Headers {
			h.request.Header.Set(k, v)
		}
	}

	return nil
}
