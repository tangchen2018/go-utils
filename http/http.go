package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	Method      string
	Url         string
	Host        string
	Timeout     time.Duration
	Header      http.Header
	Client      *http.Client
	Transport   *http.Transport
	requestType RequestType
	contentType string
	QueryParams url.Values
	Body        map[string]interface{}
	bodySize    int
	Result      []byte
	Response    *http.Response
}

func New(opts ...RequestOption) *Request {
	req := &Request{
		Method: GET,
		Client: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
				DisableKeepAlives: true,
				Proxy:             http.ProxyFromEnvironment,
			},
		},
		Transport:   nil,
		bodySize:    10, // default is 10MB
		Header:      make(http.Header),
		QueryParams: make(url.Values),
		requestType: TypeJSON,
	}
	for _, option := range opts {
		option(req)
	}
	return req
}

func (p *Request) SetBody(key string, value interface{}) *Request {
	if p.Body == nil {
		p.Body = make(map[string]interface{})
	}
	p.Body[key] = value
	return p
}

func (p *Request) structureRequest() (io.Reader, error) {

	if len(p.Method) <= 0 {
		p.Method = GET
	}

	var (
		body io.Reader
		bw   *multipart.Writer
	)
	// multipart-form-data
	if p.requestType == TypeMultipartFormData {
		body = &bytes.Buffer{}
		bw = multipart.NewWriter(body.(io.Writer))
	}

	if p.QueryParams != nil && len(p.QueryParams) > 0 {
		p.Url = fmt.Sprintf("%s?%s", p.Url, p.QueryParams.Encode())
	}

	switch p.Method {
	case GET:
		switch p.requestType {
		case TypeJSON:
			p.contentType = types[TypeJSON]
		case TypeForm, TypeFormData, TypeUrlencoded:
			p.contentType = types[TypeForm]
		case TypeMultipartFormData:
			p.contentType = bw.FormDataContentType()
		case TypeXML:
			p.contentType = types[TypeXML]
		default:
			return body, errors.New("Request type Error ")
		}
	case POST, PUT, DELETE, PATCH:
		switch p.requestType {
		case TypeJSON:
			if p.Body != nil {
				if bodyTmp, err := json.Marshal(p.Body); err != nil {
					return body, errors.New("marshal error")
				} else {
					body = strings.NewReader(string(bodyTmp))
				}
			}
			p.contentType = types[TypeJSON]
		case TypeForm, TypeFormData, TypeUrlencoded:
			body = strings.NewReader(FormatURLParam(p.Body))
			p.contentType = types[TypeForm]

		case TypeMultipartFormData:
			for k, v := range p.Body {
				// file 参数
				if file, ok := v.(*File); ok {
					fw, err := bw.CreateFormFile(k, file.Name)
					if err != nil {
						return body, err
					}
					_, _ = fw.Write(file.Content)
					continue
				}
				// text 参数
				vs, ok2 := v.(string)
				if ok2 {
					_ = bw.WriteField(k, vs)
				} else if ss := convertToString(v); ss != "" {
					_ = bw.WriteField(k, ss)
				}
			}
			_ = bw.Close()
			p.contentType = bw.FormDataContentType()
		case TypeXML:
			body = strings.NewReader(FormatURLParam(p.Body))
			p.contentType = types[TypeXML]
		default:
			return body, errors.New("Request type Error ")
		}
	default:
		return body, errors.New("Only support GET and POST and PUT and DELETE ")
	}
	return body, nil
}

func (p *Request) Do() error {

	var (
		err  error
		body io.Reader
		req  *http.Request
	)

	if body, err = p.structureRequest(); err != nil {
		return err
	}

	if req, err = http.NewRequestWithContext(context.Background(), p.Method, p.Url, body); err != nil {
		return err
	}

	req.Header = p.Header

	req.Header.Set("Content-Type", p.contentType)
	if p.Transport != nil {
		p.Client.Transport = p.Transport
	}
	if p.Host != "" {
		req.Host = p.Host
	}
	if p.Timeout > 0 {
		p.Client.Timeout = p.Timeout
	}

	p.Response, err = p.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = p.Response.Body.Close()
	}()

	p.Result, err = ioutil.ReadAll(io.LimitReader(p.Response.Body, int64(p.bodySize<<20))) // default 10MB change the size you want
	if err != nil {
		return err
	}
	return nil
}
