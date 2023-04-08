package BaiduPanApi

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type FileList struct {
	Method      string `json:"method"`
	Dir         string `json:"dir"`
	Order       string `json:"order"`
	Start       int    `json:"start"`
	Limit       int    `json:"limit"`
	Folder      int    `json:"folder"`
	AccessToken string `json:"access_token"`
	Desc        int    `json:"desc"`

	URL string
}

func NewFileList() *FileList {
	return &FileList{
		Method: "list",
		Order:  "name",
		Start:  0,
		Limit:  1000,
		Folder: 0,
		Desc:   0,
		URL:    "https://pan.baidu.com/rest/2.0/xpan/file",
	}
}

type BaiduPanFile struct {
	ServerFilename string `json:"server_filename"`
	Privacy        int    `json:"privacy"`
	Category       int    `json:"category"`
	Unlist         int    `json:"unlist"`
	FSID           int64  `json:"fs_id"`
	DirEmpty       int    `json:"dir_empty"`
	ServerATime    int    `json:"server_atime"`
	ServerCTime    int    `json:"server_ctime"`
	LocalMTime     int    `json:"local_mtime"`
	Size           int    `json:"size"`
	IsDir          int    `json:"isdir"`
	Share          int    `json:"share"`
	Path           string `json:"path"`
	LocalCTime     int    `json:"local_ctime"`
	ServerMTime    int    `json:"server_mtime"`
	Empty          int    `json:"empty"`
	OperID         int    `json:"oper_id"`
}

type BaiduPanListResponse struct {
	ErrorNo   int            `json:"errno"`
	GUIDInfo  string         `json:"guid_info"`
	FileList  []BaiduPanFile `json:"list"`
	ExtraInfo map[string]interface{}
}

func (api *FileList) WithMethod(method string) *FileList {
	api.Method = method
	return api
}

func (api *FileList) WithDir(dir string) *FileList {
	api.Dir = dir
	return api
}

func (api *FileList) WithOrder(order string) *FileList {
	api.Order = order
	return api
}

func (api *FileList) WithStart(start int) *FileList {
	api.Start = start
	return api
}

func (api *FileList) WithLimit(limit int) *FileList {
	api.Limit = limit
	return api
}

func (api *FileList) WithFolder(folder int) *FileList {
	api.Folder = folder
	return api
}

func (api *FileList) WithAccessToken(token string) *FileList {
	api.AccessToken = token
	return api
}

func (api *FileList) WithDesc(desc int) *FileList {
	api.Desc = desc
	return api
}

func (api *FileList) SendRequest() (*BaiduPanListResponse, error) {
	// 构造GET请求参数
	values := url.Values{}
	values.Add("method", api.Method)
	values.Add("dir", api.Dir)
	values.Add("order", api.Order)
	values.Add("start", fmt.Sprintf("%d", api.Start))
	values.Add("limit", fmt.Sprintf("%d", api.Limit))
	values.Add("folder", fmt.Sprintf("%d", api.Folder))
	values.Add("access_token", api.AccessToken)
	values.Add("desc", fmt.Sprintf("%d", api.Desc))
	url := api.URL + "?" + values.Encode()

	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取HTTP响应的Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析JSON格式的响应
	var response BaiduPanListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *FileList) GetAllList() ([]*BaiduPanListResponse, error) {
	responseList := make([]*BaiduPanListResponse, 0)

	response, err := api.SendRequest()
	if err != nil {
		return nil, err
	}

	responseList = append(responseList, response)
	for _, file := range response.FileList {
		if file.IsDir == 1 {
			apiCopy, err := deepCopy(api)
			if err != nil {
				return nil, err
			}
			apiCopy.WithDir(file.Path)
			list, err := apiCopy.GetAllList()
			if err != nil {
				return nil, err
			}
			responseList = append(responseList, list...)
		}
	}
	return responseList, nil
}

func deepCopy(src *FileList) (*FileList, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// 编码源对象
	err := enc.Encode(src)
	if err != nil {
		return nil, err
	}

	// 解码副本
	dec := gob.NewDecoder(&buf)
	var dest FileList
	err = dec.Decode(&dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}
