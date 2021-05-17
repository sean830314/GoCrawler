package httputil

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
	UNAUTHORIZED   = 401
)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "請求參數錯誤",
	UNAUTHORIZED:   "unauthorized",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
