package todo

import "fmt"

// 定義錯誤結構
type AppError struct {
	Code    string      `json:"code"` // 錯誤代碼
	Message string      `json:"msg"`  // 錯誤訊息
	Ret     interface{} `json:"ret"`  // 返回額外的資料
}

// 定義成功響應結構
type SuccessResponse struct {
	Code    string      `json:"code"` // 狀態碼
	Message string      `json:"msg"`  // 狀態訊息
	Ret     interface{} `json:"ret"`  // 返回額外的資料
}

// 定義錯誤碼常量
const (
	Error6000   = "6000" // 無法綁定 JSON 請求體
	Error6001   = "6001" // 任務內容不能為空
	Error6002   = "6002" // 創建 Todo 項目失敗
	Error6003   = "6003" // 獲取 Todo 項目失敗
	Error6004   = "6004" // 更新 Todo 項目失敗
	Error6005   = "6005" // Todo 項目未找到
	Error6006   = "6006" // 刪除 Todo 項目失敗
	SuccessCode = "0"    // 成功狀態代碼
)

// 錯誤碼對應的錯誤消息
var ErrorMap = map[string]string{
	Error6000: "無法綁定 JSON 請求體",
	Error6001: "任務內容不能為空",
	Error6002: "創建 Todo 項目失敗",
	Error6003: "獲取 Todo 項目失敗",
	Error6004: "更新 Todo 項目失敗",
	Error6005: "Todo 項目未找到",
	Error6006: "刪除 Todo 項目失敗",
}

// 創建 AppError 的函數
func NewAppError(code string, ret interface{}) *AppError {
	message, exists := ErrorMap[code]
	if !exists {
		message = "未知錯誤"
	}
	return &AppError{Code: code, Message: message, Ret: ret}
}

// 創建成功響應的函數
func NewSuccessResponse(ret interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code:    SuccessCode,
		Message: "OK",
		Ret:     ret,
	}
}

// 將錯誤信息格式化
func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s", e.Code, e.Message)
}
