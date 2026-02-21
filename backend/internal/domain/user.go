package domain

type User struct {
	ID         int64  `json:"id"`
	GroupAlias string `json:"group_alias"`
}

type UserLogin struct {
	UserID     string `json:"user_id"`
	Password   string `json:"password"`
	GroupAlias string `json:"group_alias"`
}

type UserRegister struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}
