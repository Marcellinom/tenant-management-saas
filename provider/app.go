package provider

import (
	"fmt"
	"github.com/samber/do"
)

type Application struct {
	engine   *WebEngine
	database *Connection
	i        *do.Injector
}

func (a Application) Engine() *WebEngine {
	return a.engine
}

func (a Application) Database() *Connection {
	return a.database
}

func NewApplication(web_engine *WebEngine, db_connections *Connection) *Application {
	return &Application{engine: web_engine, database: db_connections, i: do.DefaultInjector}
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
