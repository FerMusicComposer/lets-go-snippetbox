package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FerMusicComposer/lets-go-snippetbox.git/internal/models"
	"github.com/julienschmidt/httprouter"
)

// home handles the HTTP request for the home page.
//
// It checks if the current request URL path exactly matches "/". If it doesn't, it uses
// the http.NotFound() function to send a 404 response to the client.
// Importantly, it then returns from the handler. If it doesn't return, the handler
// would keep executing and also render the contents of the home.html file.
//
// Parameters:
// - w: an http.ResponseWriter object used to write the HTTP response.
// - r: an *http.Request object representing the HTTP request.
//
// Returns:
// - None.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because httprouter matches the "/" path exactly, we can now remove the
	// manual check of r.URL.Path != "/" from this handler.

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)
}

// snippetView handles the HTTP request for viewing a snippet.
//
// It extracts the value of the id parameter from the query string and tries to convert it
// to an integer using the strconv.Atoi() function. If it cannot be converted to an integer,
// or the value is less than 1, it returns a 404 page not found response.
//
// Parameters:
// - w: http.ResponseWriter: the response writer that will be used to write the HTTP response.
// - r: *http.Request: the HTTP request object that contains the request information.
//
// Return:
// - None.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context. We'll talk about request context
	// in detail later in the book, but for now it's enough to know that you can
	// use the ParamsFromContext() function to retrieve a slice containing these
	// parameter names and values like so:
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" named
	// parameter from the slice and validate it as normal.
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}

// snippetCreate handles the creation of a new snippet.
//
// It takes in an http.ResponseWriter and an http.Request as parameters.
// After creating the snippet, it redirects the user to the snippet view page.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Checking if the request method is a POST is now superfluous and can be
	// removed, because this is done automatically by httprouter.

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

}
