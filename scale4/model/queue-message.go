package model

import "encoding/json"

type MsgType string

var MsgTypeSave MsgType = "message_save"

type QueueMessage struct {
	MsgType MsgType            `json:"msg_type"`
	Payload SaveMessagePayload `json:"payload"`
}

func (qm QueueMessage) ToString() string {
	b, _ := json.Marshal(qm)
	return string(b)
}

type SaveMessagePayload struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}
