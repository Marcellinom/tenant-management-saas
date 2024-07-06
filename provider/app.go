package provider

import (
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/provider/auth"
	"github.com/samber/do"
	"os"
)

type Application struct {
	engine   *WebEngine
	database *Connection
	auth     *auth.Authenticator
	i        *do.Injector
}

const (
	BILLING    = "BILLING"
	IAM        = "IAM"
	ONBOARDING = "ONBOARDING"
)

func IntegrateWith(module string) bool {
	v := fmt.Sprintf("%s_INTEGRATION_MODE", module)
	env := os.Getenv(v)
	return env == "true" || env == "1"
}

func (a Application) Auth() *auth.Authenticator {
	return a.auth
}

func (a Application) RegisterAuth() {
	a.auth.RegisterCallback(a.engine)
}

func NewApplication(engine *WebEngine, database *Connection, auth *auth.Authenticator) *Application {
	return &Application{engine: engine, database: database, auth: auth}
}

func (a Application) Engine() *WebEngine {
	return a.engine
}

func (a Application) DefaultDatabase() *Database {
	databases := *a.database
	return databases[os.Getenv("DB_CONNECTION")]
}

func (a Application) UseConnection(name string) (*Database, bool) {
	v, exists := (*a.database)[name]
	return v, exists
}

func Bind[T any](app *Application, dependency T) {
	do.Provide[T](app.i, func(injector *do.Injector) (T, error) {
		return dependency, nil
	})
}

func Make[T any](app *Application) T {
	d, err := do.Invoke[T](app.i)
	if err != nil {
		panic(fmt.Errorf("error when creating object %T: %w", new(T), err))
	}
	return d
}
