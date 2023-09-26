package test_example

import (
	"github.com/RandomEstimate/BaiduPanApi"
	"testing"
)

func TestNewFileUpload(t *testing.T) {
	f := BaiduPanApi.NewFileUploader().WithAccessKey("123.28db5f7251a1a73ff4fa62e2804cf3fd.Y3zkepxqVVlsfuUWfUpX-AT3t6ftkfvC2AGD2Nn.naZL7g").
		WithFilePath("../test/NKNUSDT-trades-2021-05-10.zip").
		WithUploadFilePath("/app_ykc_test/NKNUSDT-trades-2021-05-10.zip").
		WithRtype(BaiduPanApi.Rtype3)

	err := f.Upload()
	if err != nil {
		t.Fatal(err)
	}

}
