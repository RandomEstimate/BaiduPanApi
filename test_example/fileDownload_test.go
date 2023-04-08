package test_example

import (
	"BaiduPanApi"
	"fmt"
	"path/filepath"
	"testing"
)

func TestFileDownload(t *testing.T) {

	// 获取文件Dlink
	fi := BaiduPanApi.NewFileInfo().
		WithAccessToken("123.e7ea833a49011676bd8b2e28364e24b7.YHBhB2uZtaH5iUmzOxMf7TrpbojbPl_4RiLOAXw.PiApnA").
		WithFSIds([]int64{166071588707972})

	myResponse, err := fi.SendRequest()
	if err != nil {
		fmt.Println("Get Error:", err.Error())
		return
	}

	// 将文件下载到test文件夹下
	root := "/app_ykc_test"
	downloadPath := "../test"
	for _, response := range myResponse.List {
		rePath, err := filepath.Rel(root, response.Path)
		if err != nil {
			panic(err)
		}
		targetPath := filepath.Join(downloadPath, rePath)
		err = BaiduPanApi.NewFileDownloader().WithAccessKey("123.e7ea833a49011676bd8b2e28364e24b7.YHBhB2uZtaH5iUmzOxMf7TrpbojbPl_4RiLOAXw.PiApnA").
			WithDlink(response.DLink).
			WithDownloadPath(targetPath).Download()
		if err != nil {
			panic(err)
		}

	}

}
