package main

type Data struct {
	Type     string   `json:"type"`
	Ip       string   `json:"ip"`
	User     string   `json:"user"`
	Content  string   `json:"content"`
	UserList []string `json:"user_list"`
}
