package auth

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

// New instantiates the *Authenticator.
func New() (*Authenticator, error) {
	//provider, err := oidc.NewProvider(
	//	context.Background(),
	//	"https://"+os.Getenv("AUTH_DOMAIN")+"/",
	//)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//conf := oauth2.Config{
	//	ClientID:     os.Getenv("AUTH_CLIENT_ID"),
	//	ClientSecret: os.Getenv("AUTH_CLIENT_SECRET"),
	//	RedirectURL:  "http://" + os.Getenv("AUTH_CALLBACK_URL"),
	//	Endpoint:     provider.Endpoint(),
	//	Scopes:       []string{oidc.ScopeOpenID, "profile"},
	//}

	return &Authenticator{}, nil
}
