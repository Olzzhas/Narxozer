package main

import (
	"expvar"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/julienschmidt/httprouter"
	"github.com/olzzhas/narxozer/graph/generated"
	"github.com/olzzhas/narxozer/graph/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// all methods
	protected := app.authenticate(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: app.resolver})))

	router.Handler(http.MethodPost, "/protected", protected)

	router.Handler(http.MethodPost, "/query", handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: app.resolver})))
	router.Handler(http.MethodGet, "/", playground.Handler("GraphQL playground", "/query"))

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	router.HandlerFunc(http.MethodPost, "/v1/login", app.loginUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/register", app.registerUserHandler)

	//return app.metrics(app.recoverPanic(app.rateLimit(router)))
	return app.metrics(app.recoverPanic(app.rateLimit(router)))

}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение UserID и Role из контекста
	userID := middleware.GetUserIDFromContext(r.Context())
	userRole := middleware.GetUserRoleFromContext(r.Context())

	// Обработка запроса
	w.Write([]byte("User ID: " + fmt.Sprint(userID) + " Role: " + userRole))
}
