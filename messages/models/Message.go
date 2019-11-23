package models

//1 - сообщение
//2 - чувак набирает

type Message struct {
	ID          uint64 `json:"id"`
	Text        string `json:"text"`
	MessageTime string `json:"message_time"`
	ChatID      uint64 `json:"chat_id"`
	IsSupp      bool   `json:"is_support"`
}

type Messages struct {
	Messages []*Message
}
