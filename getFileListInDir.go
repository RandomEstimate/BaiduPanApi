package BaiduPanApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

const apiUrl = "https://pan.baidu.com/rest/2.0/xpan/file"

type GetFileListInDirResponse struct {
	Errno    int    `json:"errno"`
	GuidInfo string `json:"guid_info"`
	List     []struct {
		ServerFilename string `json:"server_filename"`
		Privacy        int    `json:"privacy"`
		Category       int    `json:"category"`
		Unlist         int    `json:"unlist"`
		FsID           int    `json:"fs_id"`
		DirEmpty       int    `json:"dir_empty"`
		ServerAtime    int    `json:"server_atime"`
		ServerCtime    int    `json:"server_ctime"`
		LocalMtime     int    `json:"local_mtime"`
		Size           int    `json:"size"`
		Isdir          int    `json:"isdir"`
		Share          int    `json:"share"`
		Path           string `json:"path"`
		LocalCtime     int    `json:"local_ctime"`
		ServerMtime    int    `json:"server_mtime"`
		Empty          int    `json:"empty"`
		OperID         int    `json:"oper_id"`
	} `json:"list"`
}

type ListRequest struct {
	Method      string
	Dir         string
	Order       string
	Start       int
	Limit       int
	Web         string
	Folder      int
	AccessToken string
	Description int
}

type GetFileListInDir struct {
}

func NewGetFileListInDir() *GetFileListInDir {
	return &GetFileListInDir{}
}

func (g *GetFileListInDir) Do(listRequest *ListRequest) (*GetFileListInDirResponse, error) {

	v := url.Values{}
	v.Add("method", listRequest.Method)
	v.Add("access_token", listRequest.AccessToken)
	v.Add("dir", listRequest.Dir)
	if listRequest.Order != "" {
		v.Add("order", listRequest.Order)
	}
	if listRequest.Start != 0 {
		v.Add("start", fmt.Sprint(listRequest.Start))
	}
	if listRequest.Limit != 0 {
		v.Add("limit", fmt.Sprint(listRequest.Limit))
	}
	if listRequest.Folder != 0 {
		v.Add("folder", fmt.Sprint(listRequest.Folder))
	}
	if listRequest.Description != 0 {
		v.Add("description", fmt.Sprint(listRequest.Description))
	}
	if listRequest.Web != "" {
		v.Add("web", listRequest.Web)
	}

	encodedStr := v.Encode()

	url := apiUrl + "?" + encodedStr
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "pan.baidu.com")
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	responseJson := &GetFileListInDirResponse{}
	err = json.Unmarshal(buf, responseJson)
	if err != nil {
		return nil, err
	}

	return responseJson, nil

}

func addValues(v *url.Values, prefix string, val interface{}) error {
	value := reflect.ValueOf(val)
	if !value.IsValid() {
		return nil
	}

	valueType := value.Type()
	if valueType.Kind() == reflect.Ptr {
		value = value.Elem()
		valueType = valueType.Elem()
	}

	switch valueType.Kind() {
	case reflect.Struct:
		for i := 0; i < valueType.NumField(); i++ {
			field := valueType.Field(i)
			fieldName := field.Tag.Get("form")
			if fieldName == "" {
				fieldName = field.Name
			}
			if fieldName == "-" {
				continue
			}
			if prefix != "" {
				fieldName = prefix + "." + fieldName
			}
			if err := addValues(v, fieldName, value.Field(i).Interface()); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			if err := addValues(v, prefix+"[]", value.Index(i).Interface()); err != nil {
				return err
			}
		}
	case reflect.String:
		v.Add(prefix, value.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.Add(prefix, strconv.FormatInt(value.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.Add(prefix, strconv.FormatUint(value.Uint(), 10))
	case reflect.Float32:
		v.Add(prefix, strconv.FormatFloat(value.Float(), 'f', 4, 32))
	case reflect.Float64:
		v.Add(prefix, strconv.FormatFloat(value.Float(), 'f', 4, 64))
	case reflect.Bool:
		v.Add(prefix, strconv.FormatBool(value.Bool()))
	default:
		return fmt.Errorf("unsupported type %s", valueType.Kind())
	}

	return nil
}
