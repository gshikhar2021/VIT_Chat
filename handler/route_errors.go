package handler

import (
	"fmt"
	"net/http"
	"strings"
)

type errHandler func(http.ResponseWriter, *http.Request) error

func (fn errHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		if apierr, ok := err.(*data.APIError); ok {
			w.Header().Set("Content-Type", "application/json")
			apierr.SetMsg()
			Warning("API error:", apierr.Error())
			if apierr.Code == 101 || apierr.Code == 201 {
				notFound(w, r)
			} else if apierr.Code == 102 || apierr.Code == 202 || apierr.Code == 303 || apierr.Code == 105 {
				badRequest(w, r)
			} else if apierr.Code == 104 || apierr.Code == 204 || apierr.Code == 304 || apierr.Code == 401 || apierr.Code == 402 {
				unauthorized(w, r)
			} else if apierr.Code == 403 {
				forbidden(w, r)
			} else {
				badRequest(w, r)
			}
			ReportStatus(w, false, apierr)
		} else {
			Danger("Server error", err.Error())
			http.Error(w, err.Error(), 500)
		}
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	Info("Not found request:", r.RequestURI)
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(401)
	Info("forbidden:", r.RequestURI, r.Body)
}

func forbidden(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	Warning("forbidden:", r.RequestURI, r.Body)
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	Info("Bad request:", r.RequestURI, r.Body)
}

// Convenience function to redirect to the error message page
func errorMessage(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), 302)
}

// GET /err?msg=
// shows the error message page
func handleError(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	fmt.Fprintf(writer, "Error: %s!", vals.Get("msg"))
	Warning(fmt.Sprintf("Error: %s!", vals.Get("msg")))
}
