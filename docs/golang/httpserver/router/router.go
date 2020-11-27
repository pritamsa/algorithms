package router

import (
	"net/http"

	"github.com/bmizerany/pat"
)

const (
	AuthinitPath         = "/authinit"
	AuthorizePagePath    = "/authorize"
	AuthorizePostPath    = "/authorize"
	AccountPostPath      = "/account"
	RefreshPostPath      = "/refresh"
	SignOutPath          = "/signout"
	SignOutAllPath       = "/signout/all"
	GuestTokenPath       = "/authtoken/guest"
	ForgotPasswordPath   = "/password/forgot"
	ConfirmPasswordPath  = "/password/confirm"
	UpdatePasswordPath   = "/password/update"
	ChallengeSendPath    = "/challenge/send"
	ChallengeRespondPath = "/challenge/respond"
)

type Router interface {
	AddRoute(method string, path string, h http.Handler) *router
	GetRoutes() http.Handler
	SetPrefix(prefix string) *router
	GetPrefix() string
}

type router struct {
	p      *pat.PatternServeMux
	prefix string
}

// NewRouter creates a new router.
func NewRouter(notFound http.Handler) Router {
	p := pat.New()
	p.NotFound = notFound
	return &router{
		p: p,
	}
}

func (r *router) AddRoute(method string, pattern string, h http.Handler) *router {
	pattern = r.GetPrefix() + pattern
	r.p.Add(method, pattern, h)
	if method == "GET" {
		r.p.Add("HEAD", pattern, h)
	}
	return r
}

func (r *router) GetRoutes() http.Handler {
	return r.p
}

func (r *router) SetPrefix(prefix string) *router {
	r.prefix = prefix
	return r
}

func (r router) GetPrefix() string {
	return r.prefix
}
