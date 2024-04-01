package main
type Location struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Person struct {
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	PasswordHash []byte `json:"-"`
}
