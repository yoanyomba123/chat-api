package model

import (
	"net/http"

	scpb "github.com/swagchat/protobuf"
)

const (
	PlatformIOS = iota + 1
	PlatformAndroid
	PlatformEnd
)

type Devices struct {
	Devices []*Device `json:"devices"`
}

type Device struct {
	scpb.Device
}

func (d *Device) ConvertToPbDevice() *scpb.Device {
	pbDevice := &scpb.Device{
		UserID:               d.UserID,
		Platform:             d.Platform,
		Token:                d.Token,
		NotificationDeviceID: d.NotificationDeviceID,
	}

	return pbDevice
}

// type Device struct {
// 	UserID               string `json:"userId,omitempty" db:"user_id,notnull"`
// 	Platform             int32  `json:"platform,omitempty" db:"platform,notnull"`
// 	Token                string `json:"token,omitempty" db:"token,notnull"`
// 	NotificationDeviceID string `json:"notificationDeviceId,omitempty" db:"notification_device_id"`
// }

func IsValidDevicePlatform(platform int) bool {
	return platform > 0 && platform < int(PlatformEnd)
}

func (d *Device) IsValid() *ProblemDetail {
	if !(d.Platform > 0 && d.Platform < int32(PlatformEnd)) {
		return &ProblemDetail{
			Message: "Request error",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "device.platform",
					Reason: "platform is invalid. Currently only 1(iOS) and 2(Android) are supported.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	if d.Token == "" {
		return &ProblemDetail{
			Message: "Request error",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "token",
					Reason: "token is required, but it's empty.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	return nil
}
