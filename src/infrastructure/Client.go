package infrastructure

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/Rafael24595/go-api-core/src/commons"
	"github.com/Rafael24595/go-api-core/src/commons/collection"
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/domain/body"
	"github.com/Rafael24595/go-api-core/src/domain/cookie"
	"github.com/Rafael24595/go-api-core/src/domain/header"
)

type HttpClient struct {
}

func Client() *HttpClient {
	return &HttpClient{}
}

func (c *HttpClient) Fetch(request domain.Request) (*domain.Response, commons.ApiError) {
	req, err := c.makeRequest(request)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	start := time.Now().UnixMilli()
	resp, err := client.Do(req)
	end := time.Now().UnixMilli()
	if err != nil {
		return nil, commons.ApiErrorFromCause(500, "Cannot execute HTTP request", err)
	}

	response, err := c.makeResponse(start, end, request, *resp)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *HttpClient) makeRequest(operation domain.Request) (*http.Request, commons.ApiError) {
	method := operation.Method.String()
	url := operation.Uri

	var body io.Reader
	if !operation.Body.Empty() && method != "GET" && method != "HEAD" {
		body = bytes.NewBuffer(operation.Body.Bytes)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, commons.ApiErrorFromCause(500, "Cannot build HTTP request", err)
	}

	req.Header = collection.MapMap(collection.FromMap(operation.Headers.Headers), func(key string, value header.Header) []string {
		return value.Header
	}).Collect()

	return req, nil
}

func (c *HttpClient) makeResponse(start int64, end int64, req domain.Request, resp http.Response) (*domain.Response, commons.ApiError) {
	defer resp.Body.Close()

	bodyResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, commons.ApiErrorFromCause(500, "Failed to read response", err)
	}

	headers := header.Headers{
		Headers: collection.MapMap(collection.FromMap(resp.Header), func(key string, value []string) header.Header {
			return header.Header{
				Active: true,
				Header: value,
			}
		}).Collect(),
	}

	cookies := cookie.Cookies{
		Cookies: make(map[string]cookie.Cookie),
	}

	
	if setCookie, ok := headers.Headers["Set-Cookie"]; ok && len(setCookie.Header) > 0 {
		for _, c := range setCookie.Header {
			parsed, err := cookie.CookieFromString(c)
			if err != nil {
				return nil, err
			}
			cookies.Cookies[parsed.Code] = *parsed
		}
	}

	bodyData := body.Body{
		ContentType: body.None,
		Bytes:       bodyResponse,
	}

	return &domain.Response{
		Request: req.Id,
		Date:    start,
		Time:    end - start,
		Status:  int16(resp.StatusCode),
		Headers: headers,
		Cookies: cookies,
		Body:    bodyData,
		Size:    len(bodyResponse),
	}, nil
}
