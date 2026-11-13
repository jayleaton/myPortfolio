package api_test

import (
	"net/http"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jayleaton/myPortfolio/apps/web/internal/api"
	"github.com/jayleaton/myPortfolio/apps/web/internal/api/apitest"
	"github.com/jayleaton/myPortfolio/apps/web/internal/api/mocks"
	"github.com/jayleaton/myPortfolio/apps/web/internal/services"
)

var (
	jjPage = services.PortfolioPage{
		Portfolio: services.Portfolio{Name: "JJ Eaton", Featured: []services.Project{{Slug: "tensorfold", Name: "TensorFold"}}},
		Videos:    services.Videos{Items: []services.Video{{ID: "WCrknrNsmjQ", Title: "Latest"}}, Live: true},
		Followers: services.Followers{Handle: "jayleaton", Count: 2367, Known: true},
	}
)

func newPortfolioRouter(test *testing.T) (*mocks.MockPortfolioServiceMockRecorder, http.Handler) {
	service := mocks.NewMockPortfolioService(gomock.NewController(test))
	mux := http.NewServeMux()
	api.NewPortfolio(service, apitest.Gate()).Routes(mux)
	return service.EXPECT(), mux
}

// The portfolio is public: signed out and signed in read the same data in
// the envelope.
func TestPortfolioRoutes(test *testing.T) {
	test.Parallel()
	cases := map[string]struct {
		as           apitest.Session
		method, path string
		expect       func(service *mocks.MockPortfolioServiceMockRecorder)
		status       int
		check        func(test *testing.T, response apitest.Response)
	}{
		"page signed out": {
			as: apitest.Anonymous, method: "GET", path: apiRoutes.Portfolio, status: http.StatusOK,
			expect: func(service *mocks.MockPortfolioServiceMockRecorder) { service.Page(anyCtx).Return(jjPage) },
			check: func(test *testing.T, response apitest.Response) {
				got := apitest.Data[services.PortfolioPage](test, response)
				if got.Portfolio.Name != "JJ Eaton" || len(got.Portfolio.Featured) != 1 || got.Followers.Count != 2367 || !got.Videos.Live {
					test.Errorf("page = %+v", got)
				}
			},
		},
		"page signed in": {
			as: apitest.AsUser(ada), method: "GET", path: apiRoutes.Portfolio, status: http.StatusOK,
			expect: func(service *mocks.MockPortfolioServiceMockRecorder) { service.Page(anyCtx).Return(jjPage) },
		},
		"followers": {
			as: apitest.Anonymous, method: "GET", path: apiRoutes.PortfolioFollowers, status: http.StatusOK,
			expect: func(service *mocks.MockPortfolioServiceMockRecorder) {
				service.Followers(anyCtx).Return(jjPage.Followers)
			},
			check: func(test *testing.T, response apitest.Response) {
				if got := apitest.Data[services.Followers](test, response); got.Count != 2367 || !got.Known || got.Handle != "jayleaton" {
					test.Errorf("followers = %+v", got)
				}
			},
		},
		"followers not fetched yet": {
			as: apitest.Anonymous, method: "GET", path: apiRoutes.PortfolioFollowers, status: http.StatusOK,
			expect: func(service *mocks.MockPortfolioServiceMockRecorder) {
				service.Followers(anyCtx).Return(services.Followers{Handle: "jayleaton"})
			},
			check: func(test *testing.T, response apitest.Response) {
				if got := apitest.Data[services.Followers](test, response); got.Known {
					test.Errorf("followers = %+v, want unknown", got)
				}
			},
		},
	}
	for name, testCase := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			service, router := newPortfolioRouter(test)
			if testCase.expect != nil {
				testCase.expect(service)
			}
			response := apitest.Do(test, router, testCase.as, testCase.method, testCase.path, "")
			if response.Status != testCase.status {
				test.Fatalf("status = %d, want %d: %v", response.Status, testCase.status, response.Error)
			}
			if testCase.check != nil {
				testCase.check(test, response)
			}
		})
	}
}

// The UI's in-process calls reach the service; without the endpoints the
// client answers empty values.
func TestClientPortfolio(test *testing.T) {
	test.Parallel()
	service := mocks.NewMockPortfolioService(gomock.NewController(test))
	service.EXPECT().Page(anyCtx).Return(jjPage)
	service.EXPECT().Followers(anyCtx).Return(jjPage.Followers)
	users := mocks.NewMockUserService(gomock.NewController(test))
	client := api.NewClient(api.NewUsers(users, apitest.Gate()), nil).WithPortfolio(api.NewPortfolio(service, apitest.Gate()))
	ctx := signedIn(test, apitest.Anonymous)
	page, err := client.Portfolio(ctx)
	if err != nil || page.Portfolio.Name != "JJ Eaton" {
		test.Fatalf("Portfolio = %+v, %v", page, err)
	}
	followers, err := client.PortfolioFollowers(ctx)
	if err != nil || followers.Count != 2367 {
		test.Fatalf("PortfolioFollowers = %+v, %v", followers, err)
	}
	bare := api.NewClient(api.NewUsers(users, apitest.Gate()), nil)
	if page, err := bare.Portfolio(ctx); err != nil || page.Portfolio.Name != "" {
		test.Fatalf("bare Portfolio = %+v, %v", page, err)
	}
	if followers, err := bare.PortfolioFollowers(ctx); err != nil || followers.Known {
		test.Fatalf("bare PortfolioFollowers = %+v, %v", followers, err)
	}
}
