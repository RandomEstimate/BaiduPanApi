package test_example

import (
	"BaiduPanApi"
	"fmt"
	"testing"
)

func TestFileInfo(t *testing.T) {
	fi := BaiduPanApi.NewFileInfo().
		WithAccessToken("123.e7ea833a49011676bd8b2e28364e24b7.YHBhB2uZtaH5iUmzOxMf7TrpbojbPl_4RiLOAXw.PiApnA").
		WithFSIds([]int64{885019128748895, 166071588707972})

	response, err := fi.SendRequest()
	if err != nil {
		fmt.Println("Get Error:", err.Error())
		return
	}

	for _, item := range response.List {
		fmt.Printf("FileName:%s, FsId:%d, DLink:%s\n", item.FileName, item.FsId, item.DLink)
	}
}
