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

func Bind[T any](app *Application, label string, dependency T) {
	do.ProvideNamed[T](app.i, label, func(injector *do.Injector) (T, error) {
		return dependency, nil
	})
}

func Make[T any](app *Application, label string) T {
	d, err := do.InvokeNamed[T](app.i, label)
	if err != nil {
		panic(fmt.Errorf("error when creating object %s: %w", label, err))
	}
	return d
}
