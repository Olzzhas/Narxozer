package main

import (
	"expvar"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/julienschmidt/httprouter"
	"github.com/olzzhas/narxozer/graph/generated"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	//GraphQL
	router.Handler(http.MethodPost, "/query", handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: app.resolver})))
	router.Handler(http.MethodGet, "/", playground.Handler("GraphQL playground", "/query"))

	//User Routes
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.getUserById)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/login", app.loginUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/logout", app.logoutHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/refresh", app.refreshHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/update/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/image/:id", app.setProfileImageHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	//return app.metrics(app.recoverPanic(app.rateLimit(router)))
	return app.metrics(app.recoverPanic(app.rateLimit(app.authenticate(router))))

}
