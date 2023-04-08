package test_example

import (
	"fmt"
	"github.com/RandomEstimate/BaiduPanApi"
	"testing"
)

func TestAuth(t *testing.T) {
	// 创建 OAuthClient 对象并设置参数值
	client := BaiduPanApi.NewOAuthClient()
	token, err := client.WithResponseType("token").
		WithClientId("hnsWBK8TBrhhzxUqT0p3yEwgLqAQW4GV").
		WithRedirectUri("oob").
		WithScope("basic,netdisk").
		GetAccessToken()

	// 输出结果
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(token)
	}
}
