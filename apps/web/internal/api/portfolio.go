package api

//go:generate go tool mockgen -source=portfolio.go -destination=mocks/portfolio.go -package=mocks

import (
	"context"
	"net/http"

	"github.com/jayleaton/myPortfolio/apps/web/internal/access"
	"github.com/jayleaton/myPortfolio/apps/web/internal/config"
	"github.com/jayleaton/myPortfolio/apps/web/internal/services"
)

// PortfolioService is the port these endpoints call (*services.PortfolioService).
type PortfolioService interface {
	Page(ctx context.Context) services.PortfolioPage
	Followers(ctx context.Context) services.Followers
}

var _ PortfolioService = (*services.PortfolioService)(nil)

// Portfolio serves the public portfolio under /api/core/v1/portfolio: the
// home page's content with its live values, and the X follower count alone
// (what the live stream polls). Both are public and read memory only.
type Portfolio struct {
	svc  PortfolioService
	gate *access.Gate

	page      Endpoint[NoInput, services.PortfolioPage]
	followers Endpoint[NoInput, services.Followers]
}

// NewPortfolio returns the portfolio endpoints.
func NewPortfolio(svc PortfolioService, gate *access.Gate) *Portfolio {
	portfolio := &Portfolio{svc: svc, gate: gate}
	routes := config.APIRoutes
	portfolio.page = Endpoint[NoInput, services.PortfolioPage]{"GET", routes.Portfolio, access.Public, portfolio.read}
	portfolio.followers = Endpoint[NoInput, services.Followers]{"GET", routes.PortfolioFollowers, access.Public, portfolio.readFollowers}
	return portfolio
}

// Routes registers the endpoints on mux.
func (portfolio *Portfolio) Routes(mux *http.ServeMux) {
	mux.Handle(portfolio.page.Pattern(), portfolio.page.Serve(portfolio.gate))
	mux.Handle(portfolio.followers.Pattern(), portfolio.followers.Serve(portfolio.gate))
}

func (portfolio *Portfolio) read(request *http.Request, _ access.Caller, _ NoInput) (services.PortfolioPage, error) {
	return portfolio.svc.Page(request.Context()), nil
}

func (portfolio *Portfolio) readFollowers(request *http.Request, _ access.Caller, _ NoInput) (services.Followers, error) {
	return portfolio.svc.Followers(request.Context()), nil
}
