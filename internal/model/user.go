package model

type User struct {
	ID       int          `json:"id"`
	Username string       `json:"username"`
	Password string       `json:"password"`
	Admin    bool         `json:"admin"`
	Settings UserSettings `json:"settings"`
	ApiKeys  []ApiKey     `json:"api_keys"`
}

type UserSettings struct {
}

type ApiKey struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Key    string `json:"key"`
}
