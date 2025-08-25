package main

import (
	"fmt"
	"net/http"
     
)

func (app *application) healthcheckHandler(w http.ResponseWriter, 
                                           r *http.Request) {
   
     js := `{"status": "available", "environment": %q, "version": %q}`
     version := "1.0.0" // Define your version here or get it from config
     js = fmt.Sprintf(js, app.config.env, version)
     // Content-Type is text/plain by default
     w.Header().Set("Content-Type", "application/json")
     // Write the JSON as the HTTP response body.
     w.Write([]byte(js))

}
