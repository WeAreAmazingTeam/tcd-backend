package logs

import (
	"time"

	"github.com/WeAreAmazingTeam/tcd-backend/constant"
)

type (
	ActivityLog struct {
		ID        int    `json:"id"`
		Content   string `json:"content"`
		UserAgent string `json:"user_agent"`
		IpAddress string `json:"ip_address"`
		constant.CreatedDeleted
	}

	ActivityWebhook struct {
		ID            int       `json:"id"`
		Endpoint      string    `json:"endpoint"`
		TriggeredFrom string    `json:"triggered_from"`
		Properties    string    `json:"properties"`
		CreatedAt     time.Time `json:"created_at"`
	}
)
