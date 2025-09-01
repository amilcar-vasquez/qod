package main

import (
	"github.com/julienschmidt/httprouter"
   "net/http"
)

func (a *applicationDependencies)routes() http.Handler  {

   // setup a new router
   router := httprouter.New()
   // handle 404
   router.NotFound = http.HandlerFunc(a.notFoundResponse)
  // handle 405
   router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)
   // setup routes
   router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthCheckHandler)
   return router
  
}

