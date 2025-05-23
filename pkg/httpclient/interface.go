package httpclient

import (
	"context"
	"io"
	"net/http"
)

type IHTTPClient[TClient any] interface {
	// Request 通用請求方法
	Request(ctx context.Context, param RequestParam) (response *http.Response, err error)
	// Client 原生客戶端獲取
	Client() TClient
	// UploadFile 文件上傳方法
	UploadFile(ctx context.Context, param FileUploadParam) (*http.Response, error)
}

// RequestParam 通用請求參數
type RequestParam struct {
	BaseParam
	Body interface{} // 可以是任何類型的 body
}

// FileUploadParam 文件上傳參數
type FileUploadParam struct {
	BaseParam

	Files    map[string]FileInfo // key: form field name, value: file info
	FormData map[string]string   // 額外的表單數據
}

type BaseParam struct {
	URL             string
	Method          string
	Headers         map[string]string
	PathParams      map[string]string
	QueryParams     map[string]string
	SuccessResponse interface{}
	ErrorResponse   interface{}
	// DisableResponseData 避免資料太多也在印
	DisableLogResponseData bool
}

// FileInfo 文件信息
type FileInfo struct {
	FileName string
	Reader   io.Reader
	// 或者使用文件路径
	FilePath string
}
