package test_example

import (
	"fmt"
	"github.com/RandomEstimate/BaiduPanApi"
	"path/filepath"
	"testing"
)

func TestFileDownload(t *testing.T) {

	// 获取文件Dlink
	fi := BaiduPanApi.NewFileInfo().
		WithAccessToken("123.12461983a243d0f7fb2815c4f34ffbb4.Y7lcNmaszewLAm5ztiMp3viArcu7IGLnlbXm2Sp.9ahPbg").
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
		// 计算相对路径
		rePath, err := filepath.Rel(root, response.Path)
		if err != nil {
			panic(err)
		}
		targetPath := filepath.Join(downloadPath, rePath)
		err = BaiduPanApi.NewFileDownloader().WithAccessKey("123.12461983a243d0f7fb2815c4f34ffbb4.Y7lcNmaszewLAm5ztiMp3viArcu7IGLnlbXm2Sp.9ahPbg").
			WithDlink(response.DLink).
			WithDownloadPath(targetPath).Download()
		if err != nil {
			panic(err)
		}

	}

}
