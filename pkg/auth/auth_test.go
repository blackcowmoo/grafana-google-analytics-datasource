package auth

import (
	"strings"
	"testing"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/setting"
)

const sampleJWTJSON = `{
	"type": "service_account",
	"project_id": "demo-project",
	"private_key_id": "kid",
	"private_key": "-----BEGIN PRIVATE KEY-----\nFAKE\n-----END PRIVATE KEY-----\n",
	"client_email": "demo@demo-project.iam.gserviceaccount.com",
	"client_id": "1",
	"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "x",
	"client_x509_cert_url": "y"
}`

func TestResolve_DefaultsAuthTypeToJWT(t *testing.T) {
	got, err := Resolve(&setting.DatasourceSecretSettings{JWT: sampleJWTJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Type != TypeJWT {
		t.Errorf("Type = %q, want %q", got.Type, TypeJWT)
	}
}

func TestResolve_LegacyJWTBlob(t *testing.T) {
	got, err := Resolve(&setting.DatasourceSecretSettings{JWT: sampleJWTJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ClientEmail != "demo@demo-project.iam.gserviceaccount.com" {
		t.Errorf("ClientEmail = %q", got.ClientEmail)
	}
	if got.TokenURI != "https://oauth2.googleapis.com/token" {
		t.Errorf("TokenURI = %q", got.TokenURI)
	}
	if got.Project != "demo-project" {
		t.Errorf("Project = %q", got.Project)
	}
	if !strings.Contains(string(got.PrivateKey), "BEGIN PRIVATE KEY") {
		t.Errorf("PrivateKey did not propagate: %q", string(got.PrivateKey))
	}
}

func TestResolve_NewExplicitFields(t *testing.T) {
	got, err := Resolve(&setting.DatasourceSecretSettings{
		AuthenticationType: TypeJWT,
		ClientEmail:        "new@demo.iam.gserviceaccount.com",
		TokenURI:           "https://oauth2.googleapis.com/token",
		DefaultProject:     "new-proj",
		PrivateKey:         "PEM",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ClientEmail != "new@demo.iam.gserviceaccount.com" {
		t.Errorf("ClientEmail = %q", got.ClientEmail)
	}
	if string(got.PrivateKey) != "PEM" {
		t.Errorf("PrivateKey = %q", string(got.PrivateKey))
	}
	if got.Project != "new-proj" {
		t.Errorf("Project = %q", got.Project)
	}
}

func TestResolve_NewFieldsTakePrecedenceOverLegacy(t *testing.T) {
	got, err := Resolve(&setting.DatasourceSecretSettings{
		ClientEmail: "explicit@demo.iam.gserviceaccount.com",
		TokenURI:    "https://oauth2.googleapis.com/token",
		PrivateKey:  "EXPLICIT",
		JWT:         sampleJWTJSON, // should be ignored
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ClientEmail != "explicit@demo.iam.gserviceaccount.com" {
		t.Errorf("legacy JWT was preferred over explicit fields: %q", got.ClientEmail)
	}
}

func TestResolve_MissingCredentialsErrors(t *testing.T) {
	_, err := Resolve(&setting.DatasourceSecretSettings{})
	if err == nil {
		t.Fatal("expected error for empty settings, got nil")
	}
}

func TestResolve_InvalidLegacyJSONErrors(t *testing.T) {
	_, err := Resolve(&setting.DatasourceSecretSettings{JWT: "not json"})
	if err == nil {
		t.Fatal("expected error for malformed JWT JSON, got nil")
	}
}

func TestResolve_LegacyJSONMissingFieldsErrors(t *testing.T) {
	_, err := Resolve(&setting.DatasourceSecretSettings{JWT: `{"client_email":"a"}`})
	if err == nil {
		t.Fatal("expected error when JWT JSON omits private_key/token_uri, got nil")
	}
}

func TestResolve_NewFieldsRequireClientEmailAndTokenURI(t *testing.T) {
	// PrivateKey alone is not enough; explicit-fields path needs all three.
	_, err := Resolve(&setting.DatasourceSecretSettings{
		AuthenticationType: TypeJWT,
		PrivateKey:         "PEM",
	})
	if err == nil {
		t.Fatal("expected error when explicit fields are partial, got nil")
	}
}

func TestResolve_UnknownAuthTypeErrors(t *testing.T) {
	_, err := Resolve(&setting.DatasourceSecretSettings{AuthenticationType: "oauth"})
	if err == nil {
		t.Fatal("expected error for unknown auth type, got nil")
	}
}

func TestResolve_GCEAndWIFAreNotYetImplemented(t *testing.T) {
	for _, at := range []string{TypeGCE, TypeWIF} {
		_, err := Resolve(&setting.DatasourceSecretSettings{AuthenticationType: at})
		if err == nil {
			t.Errorf("auth type %q should return not-implemented error in PR-A", at)
		}
	}
}
