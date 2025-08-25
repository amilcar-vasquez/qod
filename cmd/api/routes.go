package main

import (
	"github.com/julienschmidt/httprouter"
   "net/http"
)

func (app *applicationDependencies) routes() http.Handler {
   router := httprouter.New()
   router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

   return router
}