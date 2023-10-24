package types

type SubscribeRequest struct {
	Email string `json:"email"`
}
type SubscribeRequestBody struct {
	Email  string   `json:"email_address"`
	Status string   `json:"status"`
	Tags   []string `json:"tags"`
}
