package test_example

import (
	"fmt"
	"testing"
)
import (
	"github.com/RandomEstimate/BaiduPanApi"
)

func TestGetFileListInDir(t *testing.T) {
	obj := BaiduPanApi.NewGetFileListInDir()

	resp, err := obj.Do(&BaiduPanApi.ListRequest{
		Method:      "list",
		AccessToken: "123.e7ea833a49011676bd8b2e28364e24b7.YHBhB2uZtaH5iUmzOxMf7TrpbojbPl_4RiLOAXw.PiApnA",
		Dir:         "/BinanceTradeData/ETC",
	})

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for _, v := range resp.List {
		fmt.Println(v.ServerFilename)
	}

}
