package BaiduPanApi

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sync"
)

type FileDirUploader struct {
	accessKey         string
	fileDirPath       string
	uploadFileDirPath string
}

func NewFileDirUpload() *FileDirUploader {
	return &FileDirUploader{}
}

func (f *FileDirUploader) WithAccessKey(accessKey string) *FileDirUploader {
	f.accessKey = accessKey
	return f
}

func (f *FileDirUploader) WithFileDirPath(fileDirPath string) *FileDirUploader {
	f.fileDirPath = fileDirPath
	return f
}

func (f *FileDirUploader) WithUploadFileDirPath(uploadFileDirPath string) *FileDirUploader {
	f.uploadFileDirPath = uploadFileDirPath
	return f
}

func (f *FileDirUploader) Upload() {

	// 遍历目录下的所有文件
	relPathList := make([]string, 0)
	absPathList := make([]string, 0)

	err := filepath.Walk(f.fileDirPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			relPath, err := filepath.Rel(f.fileDirPath, path)
			if err != nil {
				return fmt.Errorf("filepath.Rel cannot parse path. %v", err)
			}

			relPathList = append(relPathList, relPath)
			absPathList = append(absPathList, path)
			return nil
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// 将文件多线程上传
	worker := make(chan struct{}, 3)
	wait := sync.WaitGroup{}

	for index, relPath := range relPathList {
		fileUploader := NewFileUploader().WithAccessKey(f.accessKey).WithFilePath(absPathList[index]).WithUploadFilePath(filepath.Join(f.uploadFileDirPath, relPath))
		worker <- struct{}{}
		wait.Add(1)
		go func() {
			err := fileUploader.Upload()
			if err != nil {
				log.Fatal(err)
			}
			<-worker
			wait.Done()
			if err == nil {
				log.Println(fileUploader.String(), "上传成功")
			}
		}()
	}

	wait.Wait()

}
