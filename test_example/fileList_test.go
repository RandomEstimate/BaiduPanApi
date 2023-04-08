package test_example

import (
	"fmt"
	"github.com/RandomEstimate/BaiduPanApi"
	"testing"
)

func TestFileList(t *testing.T) {
	api := BaiduPanApi.NewFileList().WithMethod("list").
		WithDir("/app_ykc_test").
		WithFolder(0).
		WithAccessToken("123.e7ea833a49011676bd8b2e28364e24b7.YHBhB2uZtaH5iUmzOxMf7TrpbojbPl_4RiLOAXw.PiApnA")

	response, err := api.SendRequest()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, f := range response.FileList {
			fmt.Printf("File: %s, Size: %d, Path: %v, IsDir: %v,FSIds: %d \n", f.ServerFilename, f.Size, f.Path, f.IsDir, f.FSID)
		}
	}
}

func TestFileListAll(t *testing.T) {
	// 循环递归调用 获取目录下所有 fsid
	api := BaiduPanApi.NewFileList().WithMethod("list").
		WithDir("/app_ykc_test").
		WithFolder(0).
		WithAccessToken("123.e7ea833a49011676bd8b2e28364e24b7.YHBhB2uZtaH5iUmzOxMf7TrpbojbPl_4RiLOAXw.PiApnA")

	list, err := api.GetAllList()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range list {
		for _, fileInfo := range v.FileList {
			fmt.Printf("File: %s, Size: %d, Path: %v, IsDir: %v,FSIds: %d \n", fileInfo.ServerFilename, fileInfo.Size, fileInfo.Path, fileInfo.IsDir, fileInfo.FSID)
		}
	}

}
