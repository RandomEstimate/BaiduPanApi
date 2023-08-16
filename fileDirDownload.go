package BaiduPanApi

import (
	"fmt"
	"path/filepath"
)

type FileDirDownloader struct {
	accessKey           string
	fileDirPath         string
	downloadFileDirPath string
}

func NewFileDirDownloader() *FileDirDownloader {
	return &FileDirDownloader{}
}

func (f *FileDirDownloader) WithAccessKey(accessKey string) *FileDirDownloader {
	f.accessKey = accessKey
	return f
}

func (f *FileDirDownloader) WithFileDirPath(fileDirPath string) *FileDirDownloader {
	f.fileDirPath = fileDirPath
	return f
}

func (f *FileDirDownloader) WithDownloadFileDirPath(downloadFileDirPath string) *FileDirDownloader {
	f.downloadFileDirPath = downloadFileDirPath
	return f
}

func (f *FileDirDownloader) Download() {

	// 获取文件夹下文件的所有fsid
	api := NewFileList().
		WithDir(f.fileDirPath).
		WithAccessToken(f.accessKey)

	response, err := api.GetAllList()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// 获取Dlink
	fsIdList := make([]int64, 0)
	for _, obj := range response {
		for _, list := range obj.FileList {
			if list.FSID != 0 {
				fsIdList = append(fsIdList, list.FSID)
			}
		}
	}

	fi := NewFileInfo().
		WithAccessToken(f.accessKey).
		WithFSIds(fsIdList)

	response2, err := fi.SendRequest()
	if err != nil {
		fmt.Println("Get Error:", err.Error())
		return
	}

	// 将文件下载到文件夹下
	for _, response := range response2.List {
		if response.DLink != "" {
			rePath, err := filepath.Rel(f.fileDirPath, response.Path)
			if err != nil {
				panic(err)
			}
			targetPath := filepath.Join(f.downloadFileDirPath, rePath)
			err = NewFileDownloader().WithAccessKey(f.accessKey).
				WithDlink(response.DLink).
				WithDownloadPath(targetPath).Download()
			if err != nil {
				panic(err)
			}
		}

	}

}
