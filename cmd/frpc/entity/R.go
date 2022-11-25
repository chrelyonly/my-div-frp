package entity

// R 转换结构
type R struct {
	Code    int32          `json:"code"`
	Data    map[string]any `json:"data"`
	Success bool           `json:"success"`
	Msg     string         `json:"msg"`
}
