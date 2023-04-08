package test_example

import (
	"BaiduPanApi"
	"testing"
)

func TestNewFileUpload(t *testing.T) {
	f := BaiduPanApi.NewFileUploader().WithAccessKey("123.e7ea833a49011676bd8b2e28364e24b7.YHBhB2uZtaH5iUmzOxMf7TrpbojbPl_4RiLOAXw.PiApnA").WithFilePath("./test/1.txt").WithUploadFilePath("/app_ykc_test/2.txt")

	err := f.Upload()
	if err != nil {
		t.Fatal(err)
	}

}