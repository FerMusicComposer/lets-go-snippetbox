package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
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

// newTemplateData creates a new instance of the templateData struct.
//
// It takes a pointer to an http.Request as its parameter and returns a pointer to a templateData struct.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

// decodePostForm decodes the form data from the HTTP request and populates the
// given destination struct.
//
// It calls ParseForm() on the request to parse the form data. If there is an
// error during parsing, it returns the error.
//
// Then, it calls Decode() on the form decoder instance, passing the target
// destination struct and the request's PostForm as parameters. If there is an
// error during decoding, it checks if the error is of type *form.InvalidDecoderError.
// If it is, it panics. Otherwise, it returns the error.
//
// If everything is successful, it returns nil.
func (app *application) decodePostForm(r *http.Request, destination any) error {
	// Call ParseForm() on the request, in the same way that we did in our
	// createSnippetPost handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = app.formDecoder.Decode(destination, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var InvalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &InvalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal.
		return err

	}

	return nil
}

// isAuthenticated checks if the user is authenticated.
//
// It takes a *http.Request as a parameter.
// Returns a bool indicating whether the user is authenticated or not.
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
