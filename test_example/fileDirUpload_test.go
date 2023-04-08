package test_example

import (
	"BaiduPanApi"
	"testing"
)

func TestFileDirUploader_Upload(t *testing.T) {

	f := BaiduPanApi.NewFileDirUpload().WithAccessKey("123.e7ea833a49011676bd8b2e28364e24b7.YHBhB2uZtaH5iUmzOxMf7TrpbojbPl_4RiLOAXw.PiApnA").WithFileDirPath(
		"../").WithUploadFileDirPath("/app_ykc_test")
	f.Upload()

}
