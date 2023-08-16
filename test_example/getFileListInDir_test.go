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
		AccessToken: "123.12461983a243d0f7fb2815c4f34ffbb4.Y7lcNmaszewLAm5ztiMp3viArcu7IGLnlbXm2Sp.9ahPbg",
		Dir:         "/BinanceTradeData/IMX",
		Limit:       1000,
	})

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for _, v := range resp.List {
		fmt.Println(v.ServerFilename)
	}

}
