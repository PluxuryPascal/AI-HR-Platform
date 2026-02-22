package domain

type User struct {
	ID         string `json:"id"`
	GroupAlias string `json:"group_alias"`
}

type UserLogin struct {
	ID         string `json:"user_id"`
	Password   string `json:"password"`
	GroupAlias string `json:"group_alias"`
}

type UserRegister struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}
