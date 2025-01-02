package model

type Message struct {
	ID      int    `json:"id"`
	From    int    `json:"from"`
	To      int    `json:"to"`
	Content string `json:"content"`
}
