# JJ Eaton — portfolio

The source of [jayleaton.com](https://www.jayleaton.com): inference engines, routers and GPU recipes. It is one static Go
binary built on [Go Fast Stack](https://github.com/jayleaton/go-fast-stack): server-rendered [templ](https://templ.guide)
pages, [Datastar](https://data-star.dev) for live updates, Postgres for a small cache, deployed on Railway.

Working with agents? Read [AGENTS.md](AGENTS.md) (the template's architecture, test boundaries and verification rules)
and [CONTEXT.md](CONTEXT.md).

## What is where

| What | Where |
|---|---|
| All copy: featured work, network, socials, archive, section labels, share text | [`apps/web/internal/content/data.go`](apps/web/internal/content/data.go) (types in `content.go`) |
| Live values: X follower count, latest YouTube videos | [`apps/web/internal/services/portfolio.go`](apps/web/internal/services/portfolio.go), sources in `followers.go` and `videos.go` |
| JSON API (`/api/core/v1/portfolio`, `/portfolio/followers`) | [`apps/web/internal/api/portfolio.go`](apps/web/internal/api/portfolio.go) |
| Page, one templ file per section | `apps/web/ui/portfolio_*.templ`, handler and diagram layout in `apps/web/ui/portfolio.go` |
| Design tokens and all page styles | [`apps/web/ui/styles/portfolio.css`](apps/web/ui/styles/portfolio.css) |
| Page script (count-up, readouts, reveals) | `apps/web/ui/static/portfolio.js` |
| Images, fonts, share card (`og.png`, 1200×630) | `apps/web/ui/static/{img,fonts}/`, `apps/web/ui/static/og.png` |
| Railway | [`deploy/platform.yaml`](deploy/platform.yaml) (project `jayleaton-portfolio`), [docs/deploy.md](docs/deploy.md) |

To change the site's text, edit `data.go` and rebuild; nothing else hard-codes copy.

## Live data

Pages never wait on X or YouTube. `PortfolioService.Run` (started in `apps/web/cmd/server`) loads the last good values
from the `cache_entry` table, then refreshes in the background:

- **X followers**, every 15 minutes, from the free public JSON at `api.fxtwitter.com/<handle>` (`user.followers`), with
  `api.vxtwitter.com` as a fallback. There is no paid X API and no key. A failed refresh keeps the last good count;
  before the first success the number is hidden. Open pages get the new count over a Datastar SSE stream
  (`/live/followers`) and animate up to it. The stream survives restarts and deploys: the page's first request is a
  handshake that appends an element whose stream retries forever (Datastar `retry: 'always'`), and a shutting-down
  server ends streams cleanly so pages reconnect to the next instance (`apps/web/ui/live.go`).
- **YouTube**, hourly: the channel's `/videos` page (`ytInitialData`), with the RSS feed as a fallback (it has returned
  404 since late 2026). Until the first success the page shows `FallbackVideos` from `data.go`.

Both parsers are tested against recorded responses in `apps/web/internal/services/testdata/`; tests make no network calls.

## Run it locally

You need Go (the version in `mise.toml`, which `make` fetches) and Docker:

```sh
make db-up && make migrate   # local Postgres, then auth + app migrations
make dev                     # http://localhost:3000 (PORT in .env to change)
```

No Docker? Any Postgres works: set `DATABASE_URL` in `.env` (created from `.env.example` on first run), then
`make migrate dev`.

## Regenerate after editing

```sh
make generate   # templ (*.templ -> *_templ.go), sqlc, mocks
make css        # Tailwind standalone CLI: styles/input.css + portfolio.css -> static/app.css
```

Templ comes first, then CSS. Both outputs are committed; `make check-generated` fails when they are stale.

## Checks

`make lint test build` while you work and `make ci` before a pull request. See [AGENTS.md](AGENTS.md) "Verifying" and
[docs/CI.md](docs/CI.md).

## Production domain

Production's `PUBLIC_URL` is `https://www.jayleaton.com` (`deploy/platform.yaml` `per_environment.production`), so the
canonical, Open Graph and JSON-LD URLs name www. The server answers the bare `jayleaton.com` with a 308 to the same path on
www (`wwwRedirect` in `apps/web/internal/server/middleware.go`); the Railway domain and health checks are served as is.

## Template features kept but switched off

The template's sign-in, dashboard and account pages are still in the code but not mounted (`config.Features.AccountPages`).
The admin app, backups and the Expo app are template workspaces that are not deployed (`deploy/platform.yaml`:
`admin: false`, `backups.enabled: false`).

Email is off too. With no accounts the site never sends mail, so `config.Features.Email` follows `AccountPages`: production
boots without `RESEND_API_KEY` or `FROM_EMAIL`, the mailer is a no-op that logs only the message kind, and
`deploy/platform.yaml` has `email.enabled: false` (no sender, no Resend prompt in `railway-setup`). `make release-smoke`
boots the production binary with no email variables to prove it. Turning accounts back on turns email back on, and then
needs a verified sender.
