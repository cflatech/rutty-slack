package rutty

// RequestData Ruttyに投げるRequestData
type RequestData struct {
	Code string `json:"code"`
}

// ResponseData Ruttyから返ってきたレスポンス
type ResponseData struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Rc     int    `json:"rc"`
}
