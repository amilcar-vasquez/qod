// Filename: cmd/api/healthcheck.go
package main

import (
	"net/http"
)

func (a *applicationDependencies) healthcheckHandler(w http.ResponseWriter,
	r *http.Request) {
	panic("Apples & Oranges") // deliberate panic
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": a.config.environment,
			"version":     appVersion,
		},
	}
	err := a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
