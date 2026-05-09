// Package auth normalises plugin auth settings into a single descriptor
// and builds an OAuth2-aware HTTP client backed by grafana-google-sdk-go.
package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/setting"
)

// Auth type constants matching @grafana/google-sdk's GoogleAuthType values.
const (
	TypeJWT = "jwt"
	TypeGCE = "gce"
	TypeWIF = "workloadIdentityFederation"
)

// Resolved is the canonical auth descriptor consumed by the token-provider builder.
type Resolved struct {
	Type        string
	ClientEmail string
	TokenURI    string
	PrivateKey  []byte
	Project     string
}

// Resolve normalises plugin settings into a Resolved descriptor. Explicit
// fields (ClientEmail/TokenURI/PrivateKey) take precedence; if absent it
// falls back to parsing the legacy `secureJsonData.jwt` JSON blob so existing
// datasources keep working without re-configuration.
func Resolve(s *setting.DatasourceSecretSettings) (*Resolved, error) {
	if s == nil {
		return nil, errors.New("auth: settings are nil")
	}
	authType := s.AuthenticationType
	if authType == "" {
		authType = TypeJWT
	}
	switch authType {
	case TypeJWT:
		return resolveJWT(s)
	case TypeGCE, TypeWIF:
		return nil, fmt.Errorf("auth: %q is not yet supported", authType)
	default:
		return nil, fmt.Errorf("auth: unknown authentication type %q", authType)
	}
}

func resolveJWT(s *setting.DatasourceSecretSettings) (*Resolved, error) {
	if s.PrivateKey != "" {
		if s.ClientEmail == "" || s.TokenURI == "" {
			return nil, errors.New("auth: clientEmail and tokenUri are required when privateKey is set")
		}
		return &Resolved{
			Type:        TypeJWT,
			ClientEmail: s.ClientEmail,
			TokenURI:    s.TokenURI,
			PrivateKey:  []byte(s.PrivateKey),
			Project:     s.DefaultProject,
		}, nil
	}

	if s.JWT == "" {
		return nil, errors.New("auth: no credentials configured (set privateKey or upload a JWT JSON)")
	}
	var blob struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
		TokenURI    string `json:"token_uri"`
		ProjectID   string `json:"project_id"`
	}
	if err := json.Unmarshal([]byte(s.JWT), &blob); err != nil {
		return nil, fmt.Errorf("auth: parsing legacy JWT JSON: %w", err)
	}
	if blob.ClientEmail == "" || blob.PrivateKey == "" || blob.TokenURI == "" {
		return nil, errors.New("auth: legacy JWT JSON missing client_email/private_key/token_uri")
	}
	return &Resolved{
		Type:        TypeJWT,
		ClientEmail: blob.ClientEmail,
		TokenURI:    blob.TokenURI,
		PrivateKey:  []byte(blob.PrivateKey),
		Project:     blob.ProjectID,
	}, nil
}
