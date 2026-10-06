package domain

type Message struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Body   string `json:"body"`
}
