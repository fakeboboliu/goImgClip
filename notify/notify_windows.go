package notify

import (
	"fmt"
	"gopkg.in/toast.v1"
)

const (
	appID = "goImgClip"
)

func doToast(title, msg string) {
	notification := toast.Notification{
		AppID: appID,
		Title: title,
	}
	if msg != "" {
		notification.Message = msg
	}
	notification.Push()
}

func Error(msg string) {
	doToast("[goImgClip] Error", msg)
}

func Action(msg string) {
	doToast("[goImgClip] Action", msg)
}

func Success(targetName string) {
	doToast("Upload success", fmt.Sprint("image uploaded to ", targetName, ", URL is in your clipboard now."))
}
