package BaiduPanApi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type FileDownloader struct {
	url          string
	downloadPath string
	accessKey    string
}

func NewFileDownloader() *FileDownloader {
	return &FileDownloader{}
}
func (f *FileDownloader) WithDlink(dlink string) *FileDownloader {
	f.url = dlink
	return f
}

func (f *FileDownloader) WithDownloadPath(downloadPath string) *FileDownloader {
	f.downloadPath = downloadPath
	return f
}

func (f *FileDownloader) WithAccessKey(accessKey string) *FileDownloader {
	f.accessKey = accessKey
	return f
}

func (f *FileDownloader) Download() error {
	client := &http.Client{}

	values := url.Values{}
	values.Add("access_token", f.accessKey)

	req, err := http.NewRequest("GET", f.url+"&"+values.Encode(), nil)
	if err != nil {
		return err
	}

	values.Encode()

	req.Header.Set("User-Agent", "pan.baidu.com")
	req.Header.Set("Host", "d.pcs.baidu.com")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = createDirectoryIfNotExists(filepath.Dir(f.downloadPath))
	if err != nil {
		return err
	}

	file, err := os.Create(f.downloadPath)
	if err != nil {
		return err
	}

	defer file.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("%s downloaded %d bytes.\n", f.downloadPath, size)

	return nil
}

func createDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
