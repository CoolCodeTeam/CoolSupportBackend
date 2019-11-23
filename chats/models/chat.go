package models

type Chat struct {
	ID            uint64
	Name          string
	UserID uint64
	SupportID uint64
}

type ResponseChatsArray struct {
	Chats      []Chat
}

func NewChatModel(Name string, ID1 uint64, ID2 uint64) *Chat {
	return &Chat{
		ID:            0,
		Name:          Name,
		UserID:ID1,
		SupportID:ID2,
	}
}
