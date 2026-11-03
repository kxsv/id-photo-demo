// Package zjzapi 鹧应证件照 API 公共客户端。
//
// 所有接口示例共用的请求封装：
// - Post()         以 POST 表单方式调用接口
// - IsSuccess()    判断业务是否成功（code 为 0）
// - ImageToBase64() 读取本地图片并转 Base64
// - DownloadImage()  下载接口返回的结果图到本地
package zjzapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	requestTimeout = 60 * time.Second
	maxImageBytes  = 8 << 20
)

// BaseURL API 网关地址
const BaseURL = "https://api.zjzapi.com"

// Response 通用返回结构。
type Response struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// IsSuccess 业务是否成功（code 为 0）。
func (r *Response) IsSuccess() bool {
	return r.Code == 0
}

// FlexString 既能接收 JSON 字符串也能接收 JSON 数字的字符串字段。
//
// 接口部分字段（如 pic_id、order_no、手机号、规格的像素/dpi 等）有时返回数字、
// 有时返回字符串，统一按字符串读取，避免 json: cannot unmarshal number into ... of type string。
type FlexString string

// UnmarshalJSON 兼容字符串与数字两种取值。
func (s *FlexString) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		*s = ""
		return nil
	}
	if b[0] == '"' {
		var str string
		if err := json.Unmarshal(b, &str); err != nil {
			return err
		}
		*s = FlexString(str)
		return nil
	}
	// 数字：原样保留为字符串（如 12345 -> "12345"，10.00 -> "10.00"）。
	*s = FlexString(string(b))
	return nil
}

// Post 调用接口。
func Post(path string, params url.Values) (*Response, error) {
	client := &http.Client{Timeout: requestTimeout}
	resp, err := client.PostForm(BaseURL+"/"+strings.TrimPrefix(path, "/"), params)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取返回失败: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("返回内容不是合法 JSON: %s", string(body))
	}
	return &result, nil
}

// ImageToBase64 读取本地图片并转 Base64（不带 data:image/... 前缀）。
func ImageToBase64(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("图片不存在: %s", path)
		}
		return "", fmt.Errorf("图片读取失败: %w", err)
	}
	if len(content) > maxImageBytes {
		return "", fmt.Errorf("图片不能超过 8M: %s", path)
	}
	return base64.StdEncoding.EncodeToString(content), nil
}

// DownloadImage 下载接口返回的结果图到本地。
func DownloadImage(imageURL, saveDir, name string) (string, error) {
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", fmt.Errorf("目录创建失败: %w", err)
	}

	client := &http.Client{Timeout: requestTimeout}
	resp, err := client.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}

	file := filepath.Join(saveDir, name)
	if err := os.WriteFile(file, content, 0644); err != nil {
		return "", fmt.Errorf("文件写入失败: %w", err)
	}
	return file, nil
}

// SaveImages 下载一组"名称 -> 图片URL"的结果图到本地目录，并打印每个结果；
// 单个图片下载失败不影响其余图片，仅打印错误后继续。
func SaveImages(saveDir string, images map[string]string) {
	if len(images) == 0 {
		fmt.Println("返回数据中没有图片列表。")
		return
	}
	for name, imageURL := range images {
		file, err := DownloadImage(imageURL, saveDir, name+".jpg")
		if err != nil {
			fmt.Printf("%s 下载失败：%v\n", name, err)
			continue
		}
		fmt.Printf("%s：%s\n", name, file)
	}
}
