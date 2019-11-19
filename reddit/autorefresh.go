package reddit

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type autoRefresh struct {
	refreshTokenTimer time.Duration
	cancelAutorefresh context.CancelFunc
	context           context.Context
}

func newAutoRefresh() *autoRefresh {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &autoRefresh{
		cancelAutorefresh: cancel,
		context:           ctx,
	}
}

func (ar *autoRefresh) setRefreshTimerFromToken(token *oauth2.Token) {
	ar.refreshTokenTimer = token.Expiry.Sub(time.Now())
}

func (ar *autoRefresh) autoRefresh(a *appClient) {
	timer := time.NewTimer(ar.refreshTokenTimer)
	select {
	case <-ar.context.Done():
		log.Println("cancelling token auto refresh from context")
		return
	case <-timer.C:
		break
	}

	err := a.authorize()

	if err != nil {
		log.Println(err)
	}

	ar.autoRefresh(a)
}
