package model

type UserDetails struct {
	UserId      int64    // User ID
	Username    string   // User 名稱
	Email       string   // Email
	Password    string   // User 密碼
	Authorities []string // User 可以存取的權限
}
