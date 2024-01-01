package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// The notFound helper is simply a convenience wrapper around clientError which sends
// a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// render retrieves the appropriate template set from the cache based on the page
// name (like 'home.tmpl'). If no entry exists in the cache with the provided name,
// then it creates a new error, calls the serverError() helper method and then returns.
//
// If the entry exists, then it is written first into a bytes.Buffer just in case
// any errors occur at runtime. If there is no problem with the template, then it is
// written to the http.ResponseWriter via buffer.WriteTo().
//
// Parameters:
//   - w: an http.ResponseWriter object used to write the HTTP response.
//   - status: an integer representing the HTTP status code.
//   - page: a string representing the name of the page.
//   - td: a *templateData object containing the data to be passed to the template.
//
// Return type(s):
//   - None
func (app *application) render(w http.ResponseWriter, status int, page string, td *templateData) {
	tmplSet, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := tmplSet.ExecuteTemplate(buf, "base", td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)

	buf.WriteTo(w)
}
