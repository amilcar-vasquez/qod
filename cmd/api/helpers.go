package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "io"
  "net/http"
  "strings"
)

//create an envelope type
type envelope map[string]any

func (a *applicationDependencies)writeJSON(w http.ResponseWriter,
                                           status int, data envelope,
                                           headers http.Header) error  {
    jsResponse, err := json.MarshalIndent(data, "", "\t")
    if err != nil {
        return err
    }
jsResponse = append(jsResponse, '\n')
    // additional headers to be set
    for key, value := range headers {
        w.Header()[key] = value
    }
    // set content type header
    w.Header().Set("Content-Type", "application/json")
    // explicitly set the response status code
    w.WriteHeader(status) 
    _, err = w.Write(jsResponse)
    if err != nil {
        return err
    }

    return nil


}

//readJSON function
func (a *applicationDependencies)readJSON(w http.ResponseWriter, r *http.Request, destination any) error{
    // what is the max size of the request body (250KB seems reasonable)
    maxBytes := 256_000
    r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
    // our decoder will check for unknown fields
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    // let start the decoding
    err := dec.Decode(destination)

    // check for the different errors
    var syntaxError *json.SyntaxError
    var unmarshalTypeError *json.UnmarshalTypeError
    var invalidUnmarshalError *json.InvalidUnmarshalError
    var maxBytesError *http.MaxBytesError

    switch {
    case errors.As(err, &syntaxError):
        return fmt.Errorf("the body contains badly-formed JSON (at character %d)", syntaxError.Offset)
    case errors.Is(err, io.ErrUnexpectedEOF):
        return errors.New("the body contains badly-formed JSON")
    case strings.HasPrefix(err.Error(), "json: unknown field "):
        fieldName := strings.TrimPrefix(err.Error(), 
                                        "json: unknown field ")
        return fmt.Errorf("body contains unknown key %s", fieldName)
    case errors.As(err, &maxBytesError):
        return fmt.Errorf("the body must not be larger than %d bytes", maxBytesError.Limit)
    case errors.As(err, &unmarshalTypeError):
        if unmarshalTypeError.Field != "" {
            return fmt.Errorf("the body contains the incorrect JSON type for field %q", unmarshalTypeError.Field)
        }
        return fmt.Errorf("the body contains the incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
    case errors.Is(err, io.EOF):
        return errors.New("the body must not be empty")
    case errors.As(err, &invalidUnmarshalError):
        panic(err)
    case err != nil:
        return err
    }

    // call decode again to check if there is only one JSON value in the body
    err = dec.Decode(&struct{}{})
    if err != io.EOF {
        return errors.New("the body must contain only one JSON value")
    }

    return nil
}