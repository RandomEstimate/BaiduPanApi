package test_example

import (
	"BaiduPanApi"
	"testing"
)

func TestFileDirDownload(t *testing.T) {

	BaiduPanApi.NewFileDirDownloader().WithAccessKey("123.e7ea833a49011676bd8b2e28364e24b7.YHBhB2uZtaH5iUmzOxMf7TrpbojbPl_4RiLOAXw.PiApnA").
		WithFileDirPath("/app_ykc_test").WithDownloadFileDirPath("../test").Download()
}
