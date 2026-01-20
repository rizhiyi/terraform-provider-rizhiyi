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

	_ = c.HTTPClient.Jar.Cookies(&u)
}

// BuildRizhiyiURL Http request path
func (c *Client) BuildRizhiyiURL(parametersValues url.Values, urlPathParts ...string) url.URL {
	buildPath := "/api"
	
	// 如果第一个参数不是 "v2" 或 "v3"，则默认添加 "v2"
	hasVersion := false
	if len(urlPathParts) > 0 {
		if urlPathParts[0] == "v2" || urlPathParts[0] == "v3" {
			hasVersion = true
		}
	}
	
	if !hasVersion {
		buildPath = path.Join(buildPath, "v2")
	}

	for _, pathPart := range urlPathParts {
		pathPart = strings.ReplaceAll(pathPart, " ", "+")
		buildPath = path.Join(buildPath, pathPart)
	}
	// 确保路径以 / 结尾
	if !strings.HasSuffix(buildPath, "/") {
		buildPath = buildPath + "/"
	}

	httpScheme := getEnv(envVarHTTPScheme, defaultScheme)
	if parametersValues == nil {
		parametersValues = url.Values{}
	}
	
	host := c.Host
	if strings.HasPrefix(host, "http://") {
		host = strings.TrimPrefix(host, "http://")
	} else if strings.HasPrefix(host, "https://") {
		host = strings.TrimPrefix(host, "https://")
	}
	host = strings.TrimRight(host, "/")

	// To avoid http response truncation
	parametersValues.Set("count", "-1")
	return url.URL{
		Scheme:   httpScheme,
		Host:     host,
		Path:     buildPath,
		RawQuery: parametersValues.Encode(),
	}
}

// get resource id by name

func (c *Client) GetResourceIdByName(name string, resourceNameParts ...string) (id string, err error) {
	endpoint := c.BuildRizhiyiURL(nil, resourceNameParts...)
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

	// 兼容 v2 (objects, resources) 和 v3 (list)
	var list []interface{}
	if v, ok := data["list"].([]interface{}); ok {
		list = v
	} else if v, ok := data["objects"].([]interface{}); ok {
		list = v
	} else if v, ok := data["resources"].([]interface{}); ok {
		list = v
	}

	if list == nil {
		return "", nil
	}

	for _, item := range list {
		object, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if object["name"].(string) == name {
			// 处理 ID 可能为 float64 或 string 的情况
			switch v := object["id"].(type) {
			case float64:
				return strconv.Itoa(int(v)), nil
			case string:
				return v, nil
			case int:
				return strconv.Itoa(v), nil
			default:
				return fmt.Sprintf("%v", v), nil
			}
		}
	}
	return "", nil
}

// GetResourceById get resource detail by id
func (c *Client) GetResourceById(id string, resourceNameParts ...string) (data map[string]interface{}, err error) {
	endpoint := c.BuildRizhiyiURL(nil, append(resourceNameParts, id)...)
	response, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var respData map[string]interface{}
	err = json.Unmarshal(responseBody, &respData)
	if err != nil {
		return nil, err
	}

	// 兼容 v2 (包装在 object 字段中) 和 v3 (直接返回对象)
	if object, ok := respData["object"].(map[string]interface{}); ok {
		return object, nil
	}

	// 如果没有 object 字段，则可能直接返回了对象 (v3)
	// 简单校验一下是否包含 id 或 name 字段
	if _, ok := respData["id"]; ok {
		return respData, nil
	}
	if _, ok := respData["name"]; ok {
		return respData, nil
	}

	return nil, fmt.Errorf("resource not found or invalid response: %s", id)
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
