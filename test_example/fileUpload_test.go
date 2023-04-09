package test_example

import (
	"github.com/RandomEstimate/BaiduPanApi"
	"testing"
)

func TestNewFileUpload(t *testing.T) {
	f := BaiduPanApi.NewFileUploader().WithAccessKey("123.e7ea833a49011676bd8b2e28364e24b7.YHBhB2uZtaH5iUmzOxMf7TrpbojbPl_4RiLOAXw.PiApnA").
		WithFilePath("../test/ONTUSDT-trades-2023-04-06.zip").
		WithUploadFilePath("/BinanceTradeData/ONT/ONTUSDT-trades-2023-04-06.zip")

	err := f.Upload()
	if err != nil {
		t.Fatal(err)
	}

}
