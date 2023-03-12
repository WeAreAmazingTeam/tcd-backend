package logs

import "github.com/WeAreAmazingTeam/tcd-backend/user"

type (
	RequestCreateActivityLog struct {
		Content string `json:"content" binding:"required"`
	}

	RequestGetActivityLogByID struct {
		ID int `uri:"id" binding:"required"`
	}

	RequestDeleteActivityLog struct {
		user.User
	}

	RequestCreateActivityWebhook struct {
		Endpoint      string `json:"endpoint"`
		TriggeredFrom string `json:"triggered_from"`
		Properties    string `json:"properties"`
	}
)
