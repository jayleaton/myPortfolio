package ui

import (
	"context"

	"github.com/starfederation/datastar-go/datastar"

	"github.com/jayleaton/myPortfolio/apps/web/internal/apiclient"
	"github.com/jayleaton/myPortfolio/apps/web/ui/routes"
)

// The live follower count keeps its stream open across restarts and
// deploys.
//
// A page opens routes.LiveFollowers from data-init with Datastar's default
// retry ("auto"). That gives up when a stream ends cleanly (a proxy closing a
// long request, a server shutting down), on any non-200 answer (an edge's 502
// while no instance is up) and after 10 failed attempts, which would leave the
// count frozen until a reload. So that first request is only a handshake: it
// sends the current count, appends liveStreamElement to the body, and ends.
// The appended element holds the real stream with retry "always" and no
// attempt limit: Datastar reopens it after every end, error or non-200,
// backing off from 1 s to 30 s, and starts the backoff over after each
// successful connect.

// liveStreamParam marks a request from liveStreamElement (the managed stream)
// rather than the page's data-init handshake.
const liveStreamParam = "stream"

// liveStreamElement is the hidden element whose data-init owns the stream.
const liveStreamElement = `<div id="pf-live-stream" hidden data-init="@get('` + routes.LiveFollowers + `?` + liveStreamParam +
	`=1', {retry: 'always', retryMaxCount: Infinity})"></div>`

// liveHandshake answers the page's first request: the current count (when
// the API has one) and the element that opens the managed stream.
func (ui *UI) liveHandshake(ctx context.Context, sse *datastar.ServerSentEventGenerator) {
	if followers, err := ui.api.PortfolioFollowers(ctx); err != nil {
		ui.log.WarnContext(ctx, "ui: live followers", "err", err)
	} else if err := sse.MarshalAndPatchSignals(followerSignals(followers)); err != nil {
		return
	}
	_ = sse.PatchElements(liveStreamElement, datastar.WithSelector("body"), datastar.WithModeAppend())
}

// followersChanged reports whether a stream must send followers.
func followersChanged(sent, followers apiclient.Followers) bool {
	return followers.Count != sent.Count || followers.Known != sent.Known
}

// Shutdown ends the live streams cleanly, so http.Server.Shutdown does not wait
// out its timeout on them and pages reconnect (retry "always") to the next
// instance instead of seeing a broken connection. Register it with
// http.Server.RegisterOnShutdown; it is safe to call more than once.
func (ui *UI) Shutdown() {
	ui.closeOnce.Do(func() { close(ui.closing) })
}
