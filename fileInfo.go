package BaiduPanApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FileInfo struct {
	Method      string  `json:"method,omitempty"`
	AccessToken string  `json:"access_token,omitempty"`
	FSIds       []int64 `json:"fsids,omitempty"`
	Thumb       int     `json:"thumb,omitempty"`
	DLink       int     `json:"dlink,omitempty"`
	Extra       int     `json:"extra,omitempty"`

	URL string
}

func (fi *FileInfo) WithAccessToken(accessToken string) *FileInfo {
	fi.AccessToken = accessToken
	return fi
}

func (fi *FileInfo) WithFSIds(fsIds []int64) *FileInfo {
	fi.FSIds = fsIds
	return fi
}

func (fi *FileInfo) WithThumb(thumb int) *FileInfo {
	fi.Thumb = thumb
	return fi
}

func (fi *FileInfo) WithDLink(dlink int) *FileInfo {
	fi.DLink = dlink
	return fi
}

func (fi *FileInfo) WithExtra(extra int) *FileInfo {
	fi.Extra = extra
	return fi
}

func NewFileInfo() *FileInfo {
	return &FileInfo{
		Method: "filemetas",
		DLink:  1,
		Thumb:  0,
		Extra:  0,
		URL:    "http://pan.baidu.com/rest/2.0/xpan/multimedia",
	}
}

func (fi *FileInfo) SendRequest() (*BaiduPanInfoResponse, error) {
	// 拼接URL参数
	params := make(map[string]string)
	if fi.Method != "" {
		params["method"] = fi.Method
	}
	if fi.AccessToken != "" {
		params["access_token"] = fi.AccessToken
	}
	if len(fi.FSIds) > 0 {
		fsIdsStr := "%5B"
		for i, fsid := range fi.FSIds {
			fsIdsStr += fmt.Sprintf("%d", fsid)
			if i < len(fi.FSIds)-1 {
				fsIdsStr += "%2C"
			}
		}
		fsIdsStr += "%5D"
		params["fsids"] = fsIdsStr
	}
	if fi.Thumb != 0 {
		params["thumb"] = fmt.Sprintf("%d", fi.Thumb)
	}
	if fi.DLink != 0 {
		params["dlink"] = fmt.Sprintf("%d", fi.DLink)
	}
	if fi.Extra != 0 {
		params["extra"] = fmt.Sprintf("%d", fi.Extra)
	}

	queryString := ""
	for k, v := range params {
		if v != "" {
			queryString += k + "=" + v + "&"
		}
	}
	if queryString != "" {
		queryString = "?" + queryString[:len(queryString)-1]
	}

	// 发送GET请求
	url := fi.URL + queryString
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP GET Error:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read Body Error:", err.Error())
		return nil, err
	}

	myResponse := new(BaiduPanInfoResponse)
	err = json.Unmarshal(body, myResponse)
	if err != nil {
		fmt.Println("json.Unmarshal Error:", err.Error())
		return nil, err
	}

	return myResponse, nil
}

type BaiduPanInfoResponse struct {
	ErrMsg string `json:"errmsg"`
	ErrNo  int    `json:"errno"`
	List   []struct {
		Category    int    `json:"category"`
		DateTaken   int64  `json:"date_taken"`
		DLink       string `json:"dlink"`
		FileName    string `json:"filename"`
		FsId        int64  `json:"fs_id"`
		Height      int    `json:"height"`
		IsDir       int    `json:"isdir"`
		Md5         string `json:"md5"`
		OperId      int64  `json:"oper_id"`
		Path        string `json:"path"`
		ServerCtime int64  `json:"server_ctime"`
		ServerMtime int64  `json:"server_mtime"`
		Size        int64  `json:"size"`
		Thumbs      struct {
			Icon string `json:"icon"`
		} `json:"thumbs"`
	} `json:"list"`
}
