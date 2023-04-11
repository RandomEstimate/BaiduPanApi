package BaiduPanApi

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// OAuthClient 类型定义
type OAuthClient struct {
	responseType string
	clientId     string
	redirectUri  string
	scope        string
	cookie       string
}

func NewOAuthClient() *OAuthClient {
	return &OAuthClient{}
}

// WithResponseType 设置 responseType 属性值
func (c *OAuthClient) WithResponseType(responseType string) *OAuthClient {
	c.responseType = responseType
	return c
}

// WithClientId 设置 clientId 属性值
func (c *OAuthClient) WithClientId(clientId string) *OAuthClient {
	c.clientId = clientId
	return c
}

// WithRedirectUri 设置 redirectUri 属性值
func (c *OAuthClient) WithRedirectUri(redirectUri string) *OAuthClient {
	c.redirectUri = redirectUri
	return c
}

// WithScope 设置 scope 属性值
func (c *OAuthClient) WithScope(scope string) *OAuthClient {
	c.scope = scope
	return c
}

// WithCookie 设置 cookie 属性值
func (c *OAuthClient) WithCookie(cookie string) *OAuthClient {
	c.cookie = cookie
	return c
}

// GetAccessToken 获取 access_token 方法
func (c *OAuthClient) GetAccessToken() (string, error) {
	reqUrl := "https://openapi.baidu.com/oauth/2.0/authorize?"
	params := url.Values{}
	params.Add("response_type", c.responseType)
	params.Add("client_id", c.clientId)
	params.Add("redirect_uri", c.redirectUri)
	params.Add("scope", c.scope)

	// 发送 GET 请求
	req, err := http.NewRequest(http.MethodGet, reqUrl+params.Encode(), nil)
	if err != nil {
		return "", err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	for {
		resp, err := client.Do(req) // 发送请求
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()
		// 获取响应头部 Localation 字段

		buf, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(buf))

		location := resp.Request.Response.Header.Get("Location")
		//if location == "" {
		//	//location = resp.
		//}

		if location != "" {
			// 解析 Location 地址中是否有 access_token 字段
			parsedUrl, err := url.Parse(location)
			if err == nil {
				queryValues := parsedUrl.Query()
				access_token := queryValues.Get("access_token")
				if access_token != "" {
					return access_token, nil
				} else {
					// 更新 url 继续请求
					req, _ = http.NewRequest(http.MethodGet, location, nil)
				}
			}
		} else {
			// 没有 Location 字段表示请求结束，返回错误
			return "err", nil
		}
	}
}
