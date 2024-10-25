package webmock

import (
	"github.com/movie-guru/pkg/types"
)

type SessionInfo struct {
	ID            string
	User          string
	Authenticated bool
}

type LoginBody struct {
	InviteCode string `json:"inviteCode" omitempty`
}

type PrefBody struct {
	Content *types.UserProfile `json:"content"`
}

type ChatRequest struct {
	Content string `json:"content"`
}

type Dependencies struct {
	CurrentProbMetrics *MetricsProb
}

type MetricsProb struct {
	ChatSuccess float32
	ChatSafety  float32

	ChatEngaged      float32
	ChatRejected     float32
	ChatUnclassified float32
	ChatAcknowledged float32

	ChatSPositive     float32
	ChatSNegative     float32
	ChatSNeutral      float32
	ChatSUnclassified float32

	ChatLatencyMinMS int
	ChatLatencyMaxMS int

	LoginSuccess      float32
	LoginLatencyMinMS int
	LoginLatencyMaxMS int

	StartupSuccess      float32
	StartupLatencyMinMS int
	StartupLatencyMaxMS int

	PrefUpdateSuccess      float32
	PrefUpdateLatencyMinMS int
	PrefUpdateLatencyMaxMS int

	PrefGetSuccess      float32
	PrefGetLatencyMinMS int
	PrefGetLatencyMaxMS int
}
