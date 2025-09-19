// file: cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {

	// setup a new router
	router := httprouter.New()
	// handle 404
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	// handle 405
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)
	// setup routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/quotes", a.createQuoteHandler)
	router.HandlerFunc(http.MethodGet, "/v1/quotes/:id", a.displayQuoteHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/quotes/:id", a.updateQuoteHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/quotes/:id", a.deleteQuoteHandler)
	router.HandlerFunc(http.MethodGet, "/v1/quotes", a.listQuotesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", a.registerUserHandler)

	// Chain: CORS -> RateLimit -> RecoverPanic
	return a.recoverPanic(a.rateLimit(a.enableCORS(router)))
}
