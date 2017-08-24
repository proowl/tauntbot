package telegram


type User struct {
	Id int64
	First_name string
	Username string
	Language_code string
}

type Room struct {
	Id int64
	First_name string
	Username string 
	Type string
}

type Message struct {
	Message_id int64
	From User
	Chat Room
	Fate int64
	Text string
}

type Update struct {
	Update_id int64
	Message Message
}

type UpdateResponse struct {
	Ok bool
	Result []Update
}

type OutgoingMessage struct {
	ChatId int64 `json:"chat_id"`
	Text string `json:"text"`
	DisableNotification bool `json:"disable_notification"`
}
