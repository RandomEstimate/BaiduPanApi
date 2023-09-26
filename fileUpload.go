package BaiduPanApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const fileUploadUrlPreCreate = "https://pan.baidu.com/rest/2.0/xpan/file?method=precreate"
const fileUploadUrlSuperFile2 = "https://d.pcs.baidu.com/rest/2.0/pcs/superfile2?method=upload"
const fileUploadUrlCreate = "https://pan.baidu.com/rest/2.0/xpan/file?method=create"

type FileUploader struct {
	accessKey      string
	filePath       string
	uploadFilePath string
	rtype          Rtype
}

type Rtype int

// 文件命名策略。
//1 表示当path冲突时，进行重命名
//2 表示当path冲突且block_list不同时，进行重命名
//3 当云端存在同名文件时，对该文件进行覆盖

const Rtype1 Rtype = 1
const Rtype2 Rtype = 2
const Rtype3 Rtype = 3

func NewFileUploader() *FileUploader {
	return &FileUploader{}
}

func (f *FileUploader) WithAccessKey(accessKey string) *FileUploader {
	f.accessKey = accessKey
	return f
}

func (f *FileUploader) WithFilePath(filePath string) *FileUploader {
	f.filePath = filePath
	return f
}

func (f *FileUploader) WithUploadFilePath(uploadFilePath string) *FileUploader {
	f.uploadFilePath = uploadFilePath
	return f
}

func (f *FileUploader) WithRtype(rtype Rtype) *FileUploader {
	f.rtype = rtype
	return f
}

func (f *FileUploader) String() string {
	return fmt.Sprintf("FileUploader{accessKey: %v, filePath: %v, uploadFilePath: %v}",
		f.accessKey, f.filePath, f.uploadFilePath)
}

// baiduPCSPreCreate 预上传
func (f *FileUploader) baiduPCSPreCreate(checksums FileChecksums) (string, string, string, error) {
	fileInfo, err := os.Stat(f.filePath)
	if err != nil {
		return "", "", "", fmt.Errorf("error checking file stat: %w", err)
	}

	if err != nil {
		return "", "", "", fmt.Errorf("failed to calculate checksum of the file: %w", err)
	}

	// 构造请求Body
	formData := url.Values{}
	formData.Add("path", f.uploadFilePath)
	formData.Add("size", strconv.FormatInt(checksums.totalSize, 10))
	formData.Add("block_list", fmt.Sprintf(`["%s"]`, strings.Join(checksums.checksums, `","`)))
	formData.Add("isdir", "0")
	formData.Add("rtype", fmt.Sprint(f.rtype))
	formData.Add("autoinit", "1")
	formData.Add("local_ctime", strconv.FormatInt(fileInfo.ModTime().Unix(), 10))
	formData.Add("local_mtime", strconv.FormatInt(fileInfo.ModTime().Unix(), 10))
	if len(checksums.checksums) > 1 {
		formData.Add("content-md5", checksums.contentMd5)
		//formData.Add("slice-md5", checksums.sliceMd5)
	}

	data := formData.Encode()

	createReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf(fileUploadUrlPreCreate+"&access_token=%s", f.accessKey), strings.NewReader(data))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create precreate request: %w", err)
	}

	client := &http.Client{}

	resp, err := client.Do(createReq)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to do precreate request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read precreate response body: %w", err)
	}

	type createResp struct {
		// N1-MTE3LjYxLjguMjUxOjE2ODA4NDg0OTU6OTE1NzkxNjkwNTA2MjYyOTg0Mg==
		UploadID string `json:"uploadid"`
	}

	var createRes struct {
		Errno      int           `json:"errno"`
		Path       string        `json:"path"`
		UploadID   string        `json:"uploadid"`
		ReturnType int           `json:"return_type"`
		BlockList  []interface{} `json:"block_list"`
	}
	err = json.Unmarshal(respBody, &createRes)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to unmarshal create response body: %w", err)
	}
	if createRes.Errno != 0 {
		_, file, line, ok := runtime.Caller(0)
		if !ok {
			file = "unknown"
			line = 0
		}
		return "", "", "", fmt.Errorf("%s %d baiduPCSPreCreate errno not equal to 0 : %v", file, line, err)
	}

	return createRes.Path, createRes.UploadID, fmt.Sprint(createRes.BlockList), nil

}

// baiduPCSSuperFile2 上传单一分片到百度云
func (f *FileUploader) baiduPCSSuperFile2(uploadID string, checksums FileChecksums) error {

	file, err := os.Open(f.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var offset int64
	var partSeq int
	for offset < checksums.totalSize { // 一次发送一个分块

		chunkSizeBucket := chunkSize
		if checksums.totalSize-offset < int64(chunkSize) {
			chunkSizeBucket = int(checksums.totalSize - offset)
		}

		buf := make([]byte, chunkSizeBucket)
		_, err := file.ReadAt(buf, offset)
		if err != nil {
			return err
		}
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		// 添加文件内容到表单
		part, err := writer.CreateFormFile("file", f.filePath)
		if err != nil {
			return err
		}
		part.Write(buf)

		err = writer.Close()
		if err != nil {
			return err
		}

		formData := url.Values{}
		formData.Add("access_token", f.accessKey)
		formData.Add("type", "tmpfile")
		formData.Add("path", f.uploadFilePath)
		formData.Add("uploadid", uploadID)
		formData.Add("partseq", strconv.Itoa(partSeq))

		req, err := http.NewRequest("POST", fileUploadUrlSuperFile2+"&"+formData.Encode(), body)
		if err != nil {
			return err
		}

		//queryParams := req.URL.Query()
		//queryParams.Add("access_token", f.accessKey)
		//queryParams.Add("type", "tmpfile")
		//queryParams.Add("path", f.uploadFilePath)
		//queryParams.Add("uploadid", uploadID)
		//queryParams.Add("partseq", strconv.Itoa(partSeq))
		//req.URL.RawQuery = queryParams.Encode()

		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Content-Length", strconv.Itoa(body.Len()))

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			_, file, line, ok := runtime.Caller(0)
			if !ok {
				file = "unknown"
				line = 0
			}
			return fmt.Errorf("%s %v upload failed, code: %d, resp: %s", file, line, resp.StatusCode, string(bodyBytes))
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var reponseData struct {
			Errno int    `json:"errno"`
			MD5   string `json:"md5"`
		}

		err = json.Unmarshal(respBody, &reponseData)
		if err != nil {
			return err
		}

		offset += int64(chunkSize)
		partSeq++
	}

	return nil
}

// baiduPCSCreate 将文件完整重组
func (f *FileUploader) baiduPCSCreate(uploadId string, checksums FileChecksums) error {

	// 构造请求Body
	formData := url.Values{}
	formData.Add("path", f.uploadFilePath)
	formData.Add("size", strconv.FormatInt(int64(checksums.totalSize), 10))
	formData.Add("isdir", "0")
	formData.Add("uploadid", uploadId)
	formData.Add("mode", "1")
	formData.Add("rtype", fmt.Sprint(f.rtype))
	formData.Add("block_list", fmt.Sprintf(`["%s"]`, strings.Join(checksums.checksums, `","`)))
	data := formData.Encode()

	createReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf(fileUploadUrlCreate+"&access_token=%s", f.accessKey), strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create create request: %w", err)
	}

	client := &http.Client{}

	resp, err := client.Do(createReq)
	if err != nil {
		return fmt.Errorf("failed to do create request: %w", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var reponseData struct {
		Errno int `json:"errno"`
	}

	err = json.Unmarshal(respBody, &reponseData)
	if err != nil {
		return err
	}

	if reponseData.Errno != 0 {
		return fmt.Errorf("baiduPCSCreate errno not equal to 0 .%v", string(respBody))
	}

	return nil
}

func (f *FileUploader) Upload() error {
	const maxRetries = 3
	var err error

	for retries := 0; retries <= maxRetries; retries++ {
		if retries > 0 {
			time.Sleep(5 * time.Second) // 等待 5 秒钟再尝试
		}

		if err = VerifyFileExists(f.filePath); err != nil {
			log.Printf("Upload verifyFile error : %v", err)
			continue
		}

		checksumFile, err := ChecksumFile(f.filePath)
		if err != nil {
			log.Printf("Upload ChecksumFile error : %v", err)
			continue
		}

		// 文件预上传
		_, uploadId, _, err := f.baiduPCSPreCreate(checksumFile)
		if err != nil {
			log.Printf("Upload baiduPCSPreCreate error : %v", err)
			continue
		}

		// 文件单一分片上传
		err = f.baiduPCSSuperFile2(uploadId, checksumFile)
		if err != nil {
			log.Printf("Upload baiduPCSSuperFile2 error : %v", err)
			continue
		}

		for retries2 := 0; retries2 <= maxRetries; retries2++ {
			time.Sleep(time.Second)
			// 文件完成生成
			err = f.baiduPCSCreate(uploadId, checksumFile)
			if err != nil {
				log.Printf("Upload baiduPCSCreate error : %v", err)
				continue
			}
		}

		if err != nil {
			return err
		}

		return nil // 成功上传，直接返回
	}

	log.Printf("Uploaded failed after %d retries: %v", maxRetries, err)
	return err
}
