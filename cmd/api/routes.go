package main

import (
	"expvar"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	//Health Check
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

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
