package auth

// Authenticator is used to authenticate our users.
type Authenticator struct {
	auth_provider string
}

func (a *Authenticator) GetProvider() string {
	return a.auth_provider
}

// New instantiates the *Authenticator.
func New(provider string) (*Authenticator, error) {
	return &Authenticator{auth_provider: provider}, nil
}
