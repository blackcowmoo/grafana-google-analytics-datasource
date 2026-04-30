package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-google-sdk-go/pkg/tokenprovider"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
)

// NewHTTPClient returns an HTTP client whose transport injects an OAuth2
// access token derived from the resolved auth descriptor. Pass the result
// to `option.WithHTTPClient(...)` when constructing a google-api service.
func NewHTTPClient(ctx context.Context, r *Resolved, scopes []string) (*http.Client, error) {
	provider, err := newTokenProvider(r, scopes)
	if err != nil {
		return nil, err
	}
	opts := httpclient.Options{
		Middlewares: []httpclient.Middleware{tokenprovider.AuthMiddleware(provider)},
	}
	return httpclient.New(opts)
}

func newTokenProvider(r *Resolved, scopes []string) (tokenprovider.TokenProvider, error) {
	if r == nil {
		return nil, fmt.Errorf("auth: resolved descriptor is nil")
	}
	switch r.Type {
	case TypeJWT:
		return tokenprovider.NewJwtAccessTokenProvider(tokenprovider.Config{
			Scopes: scopes,
			JwtTokenConfig: &tokenprovider.JwtTokenConfig{
				Email:      r.ClientEmail,
				URI:        r.TokenURI,
				PrivateKey: r.PrivateKey,
			},
		}), nil
	}
	return nil, fmt.Errorf("auth: token provider for %q not implemented", r.Type)
}
