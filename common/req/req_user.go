package req

type LoginRequest struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}
