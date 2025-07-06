package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type HTTPClient struct {
	httpClient *http.Client
}

func NewHttpClient() (*HTTPClient, error) {

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{Jar: jar}

	return &HTTPClient{httpClient: &httpClient}, nil
}

func (httpC *HTTPClient) PostRequest(uri string, requestPayload map[string]interface{}, headers map[string]string) (int, []byte, error) {

	b, err := json.Marshal(requestPayload)
	if err != nil {
		return 0, nil, err
	}

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		return 0, nil, err
	}

	for headerKey, headerValue := range headers {
		request.Header.Set(headerKey, headerValue)
	}

	response, err := httpC.httpClient.Do(request)
	if err != nil {
		return 0, nil, err
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, bodyBytes, nil
}

func (httpC *HTTPClient) GetRequest(uri string, headers map[string]string) (int, []byte, error) {

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return 0, nil, err
	}

	for headerKey, headerValue := range headers {
		request.Header.Set(headerKey, headerValue)
	}

	response, err := httpC.httpClient.Do(request)
	if err != nil {
		return 0, nil, err
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, bodyBytes, nil
}

func (httpC *HTTPClient) GetCookieValue(baseURL, cookieName string) (string, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	cookies := httpC.httpClient.Jar.Cookies(url)

	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			return cookie.Value, nil
		}
	}

	return "", nil
}
