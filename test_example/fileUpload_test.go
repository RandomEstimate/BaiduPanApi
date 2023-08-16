package test_example

import (
	"github.com/RandomEstimate/BaiduPanApi"
	"testing"
)

func TestNewFileUpload(t *testing.T) {
	f := BaiduPanApi.NewFileUploader().WithAccessKey("123.12461983a243d0f7fb2815c4f34ffbb4.Y7lcNmaszewLAm5ztiMp3viArcu7IGLnlbXm2Sp.9ahPbg").
		WithFilePath("../test/a.go").
		WithUploadFilePath("/app_ykc_test/a.go").
		WithRtype(BaiduPanApi.Rtype3)

	err := f.Upload()
	if err != nil {
		t.Fatal(err)
	}

}
