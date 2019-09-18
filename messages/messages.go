package messages

type Message struct {
	ClientIP string `json:"clientIP"`
	Data     string `json:"data"`
}
