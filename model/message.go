package model

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
)

const (
	MessageTypeText           = "text"
	MessageTypeImage          = "image"
	MessageTypeFile           = "file"
	MessageTypeIndicatorStart = "indicatorStart"
	MessageTypeIndicatorEnd   = "indicatorEnd"
	MessageTypeUpdateRoomUser = "updateRoomUser"
)

type Messages struct {
	Messages []*Message `json:"messages" db:"-"`
	AllCount int64      `json:"allCount" db:"all_count"`
}

type Message struct {
	ID        uint64         `json:"-" db:"id"`
	MessageID string         `json:"messageId" db:"message_id,notnull"`
	RoomID    string         `json:"roomId" db:"room_id,notnull"`
	UserID    string         `json:"userId" db:"user_id,notnull"`
	Type      string         `json:"type,omitempty" db:"type"`
	EventName string         `json:"eventName,omitempty" db:"-"`
	Payload   utils.JSONText `json:"payload" db:"payload"`
	Role      int32          `json:"role,omitempty" db:"role,notnull"`
	Created   int64          `json:"created" db:"created,notnull"`
	Modified  int64          `json:"modified" db:"modified,notnull"`
	Deleted   int64          `json:"-" db:"deleted,notnull"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		MessageID string         `json:"messageId"`
		RoomID    string         `json:"roomId"`
		UserID    string         `json:"userId"`
		Type      string         `json:"type"`
		EventName string         `json:"eventName,omitempty"`
		Payload   utils.JSONText `json:"payload"`
		Role      int32          `json:"role"`
		Created   string         `json:"created"`
		Modified  string         `json:"modified"`
	}{
		MessageID: m.MessageID,
		RoomID:    m.RoomID,
		UserID:    m.UserID,
		Type:      m.Type,
		EventName: m.EventName,
		Payload:   m.Payload,
		Role:      m.Role,
		Created:   time.Unix(m.Created, 0).In(l).Format(time.RFC3339),
		Modified:  time.Unix(m.Modified, 0).In(l).Format(time.RFC3339),
	})
}

type ResponseMessages struct {
	MessageIds []string         `json:"messageIds,omitempty"`
	Errors     []*ProblemDetail `json:"errors,omitempty"`
}

type PayloadText struct {
	Text string `json:"text"`
}

type PayloadImage struct {
	Mime         string `json:"mime"`
	Filename     string `json:"filename"`
	SourceUrl    string `json:"sourceUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type PayloadUsers struct {
	Users []string `json:"users"`
}

func (m *Message) IsValid() *ProblemDetail {
	if m.MessageID != "" && !utils.IsValidID(m.MessageID) {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "messageId",
					Reason: "messageId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	if m.Payload == nil {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "payload",
					Reason: "payload is empty.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	if m.Type == MessageTypeText {
		var pt PayloadText
		json.Unmarshal(m.Payload, &pt)
		if pt.Text == "" {
			return &ProblemDetail{
				Message: "Invalid params",
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "payload",
						Reason: "Text type needs text.",
					},
				},
				Status: http.StatusBadRequest,
			}
		}
	}

	if m.Type == MessageTypeImage {
		var pi PayloadImage
		json.Unmarshal(m.Payload, &pi)
		if pi.Mime == "" {
			return &ProblemDetail{
				Message: "Invalid params",
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "payload",
						Reason: "Image type needs mime.",
					},
				},
				Status: http.StatusBadRequest,
			}
		}

		if pi.SourceUrl == "" {
			return &ProblemDetail{
				Status: http.StatusBadRequest,
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "payload",
						Reason: "Image type needs sourceUrl.",
					},
				},
				Message: "Invalid params",
			}
		}
	}

	return nil
}

func (m *Message) BeforeSave() {
	if m.MessageID == "" {
		m.MessageID = utils.GenerateUUID()
	}

	if m.Role == 0 {
		m.Role = utils.RoleGeneral
	}

	nowTimestamp := time.Now().Unix()
	if m.Created == 0 {
		m.Created = nowTimestamp
	}
	if m.Modified == 0 {
		m.Modified = nowTimestamp
	}
}