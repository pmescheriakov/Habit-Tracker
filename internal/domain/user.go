package domain

type User struct {
	Id     int    `json:"Id"`
	Login  string `json:"Login"`
	Name   string `json:"Name"`
	Status bool   `json:"Status"`
}
