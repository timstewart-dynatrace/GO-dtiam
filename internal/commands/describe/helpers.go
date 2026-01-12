package describe

import (
	"net/http"

	"github.com/jtimothystewart/dtiam/internal/auth"
	"github.com/jtimothystewart/dtiam/internal/client"
)

// tokenProviderAdapter adapts auth.TokenProvider to client.TokenProvider
type tokenProviderAdapter struct {
	provider auth.TokenProvider
}

func (a *tokenProviderAdapter) GetHeaders() (http.Header, error) {
	return a.provider.GetHeaders()
}

func (a *tokenProviderAdapter) IsValid() bool {
	return a.provider.IsValid()
}

func (a *tokenProviderAdapter) Close() error {
	return a.provider.Close()
}

// newOAuthProvider creates a new OAuth token provider.
func newOAuthProvider(clientID, clientSecret, accountUUID string) client.TokenProvider {
	return &tokenProviderAdapter{
		provider: auth.NewOAuthTokenManager(auth.OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			AccountUUID:  accountUUID,
		}),
	}
}

// newBearerProvider creates a new static bearer token provider.
func newBearerProvider(token string) client.TokenProvider {
	return &tokenProviderAdapter{
		provider: auth.NewStaticTokenManager(token, ""),
	}
}
