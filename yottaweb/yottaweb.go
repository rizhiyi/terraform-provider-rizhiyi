package yottaweb

// setting yottaweb request client
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	defaultScheme    = "http"
	MethodGet        = "GET"
	MethodPost       = "POST"
	MethodPut        = "PUT"
	MethodPatch      = "PATCH"
	MethodDelete     = "DELETE"
	envVarHTTPScheme = "HTTPScheme"
)

// Client is the yottaweb API client
type Client struct {
	Host          string
	Authorization string
	HTTPClient    *http.Client
}

// NewClient creates a new yottaweb client
func NewClient(host, auth string) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		Host:          host,
		Authorization: auth,
		HTTPClient: &http.Client{
			Jar: jar,
		},
	}
}

// Request create request header
func (c *Client) Request(httpMethod, requestUrl string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(httpMethod, requestUrl, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Basic "+c.Authorization)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")

	// 自动从 CookieJar 中提取 csrftoken 并设置到 Header
	if c.HTTPClient.Jar != nil {
		u, _ := url.Parse(requestUrl)
		for _, cookie := range c.HTTPClient.Jar.Cookies(u) {
			if cookie.Name == "csrftoken" {
				request.Header.Set("X-CSRFToken", cookie.Value)
				fmt.Println("Request set X-CSRFToken for", httpMethod, requestUrl)
				break
			}
		}
	}

	return request, nil
}

// DoRequest execute http request
func (c *Client) DoRequest(method string, requestURL url.URL, body map[string]interface{}) (*http.Response, error) {
	if method == MethodPost || method == MethodPut || method == MethodPatch || method == MethodDelete {
		c.ensureCSRFCookie(requestURL)
	}

	var bodyData io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		bodyData = bytes.NewBuffer(jsonData)
	}
	request, err := c.Request(method, requestURL.String(), bodyData)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	// 非 2xx 直接认为是 API 错误
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (Status: %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// 写操作如果返回 HTML，很可能是登录页或 CSRF 错误页面，而不是正常 JSON
	if (method == MethodPost || method == MethodPut || method == MethodPatch || method == MethodDelete) &&
		strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (Unexpected HTML response): %s", string(bodyBytes))
	}

	return resp, nil
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		if value == "https" {
			return value
		}
	}
	return defaultValue
}

// Do sends out request and returns HTTP response
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}

func (c *Client) ensureCSRFCookie(requestURL url.URL) {
	if c.HTTPClient == nil || c.HTTPClient.Jar == nil {
		return
	}

	u := requestURL
	u.Path = "/"
	u.RawQuery = ""

	existingCookies := c.HTTPClient.Jar.Cookies(&u)
	fmt.Println("ensureCSRFCookie before request, url:", u.String(), "cookies:", existingCookies)
	for _, cookie := range existingCookies {
		if cookie.Name == "csrftoken" && cookie.Value != "" {
			return
		}
	}

	req, err := http.NewRequest(MethodGet, u.String(), nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Basic "+c.Authorization)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := c.Do(req)
	if err != nil {
		return
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	newCookies := c.HTTPClient.Jar.Cookies(&u)
	fmt.Println("ensureCSRFCookie after request, url:", u.String(), "cookies:", newCookies)
}

// BuildRizhiyiURL Http request path
func (c *Client) BuildRizhiyiURL(parametersValues url.Values, urlPathParts ...string) url.URL {
	buildPath := "/api/v2"
	for _, pathPart := range urlPathParts {
		pathPart = strings.ReplaceAll(pathPart, " ", "+")
		buildPath = path.Join(buildPath, pathPart)
		buildPath = buildPath + "/"
	}
	httpScheme := getEnv(envVarHTTPScheme, defaultScheme)
	if parametersValues == nil {
		parametersValues = url.Values{}
	}
	// To avoid http response truncation
	parametersValues.Set("count", "-1")
	return url.URL{
		Scheme:   httpScheme,
		Host:     c.Host,
		Path:     buildPath,
		RawQuery: parametersValues.Encode(),
	}
}

// get resource id by name

func (c *Client) GetResourceIdByName(name string, resourceName string) (id string, err error) {
	var app_id = ""

	endpoint := c.BuildRizhiyiURL(nil, resourceName)
	response, err := c.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return "", err
	}
	objects, ok := data["objects"].([]interface{})
	if !ok {
		return "", nil
	}
	for _, obj := range objects {
		object := obj.(map[string]interface{})
		if object["name"].(string) == name {
			app_id = strconv.Itoa(int(object["id"].(float64)))
		}
	}
	return app_id, nil
}

// Get func
func (c *Client) Get(getURL url.URL) (*http.Response, error) {
	return c.DoRequest(MethodGet, getURL, nil)
}

// Post func
func (c *Client) Post(postURL url.URL, body map[string]interface{}) (*http.Response, error) {
	return c.DoRequest(MethodPost, postURL, body)
}

// Put func
func (c *Client) Put(putURL url.URL, body map[string]interface{}) (*http.Response, error) {
	return c.DoRequest(MethodPut, putURL, body)
}

// Delete func
func (c *Client) Delete(deleteURL url.URL) (*http.Response, error) {
	return c.DoRequest(MethodDelete, deleteURL, nil)
}
