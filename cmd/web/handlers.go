package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FerMusicComposer/lets-go-snippetbox.git/internal/models"
	"github.com/FerMusicComposer/lets-go-snippetbox.git/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// This struct represents form data and errors. All fields are exported so they
// can be read by the HTML template.
type snippetCreateForm struct {
	Title               string
	Content             string
	Expires             int
	validator.Validator // Embed a validator
}

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

// snippetCreate initializes a new createSnippetForm instance and passes it to the template. It also sets the initial value for the snippet expiry to 365 days.
//
// Parameters:
// - w: an http.ResponseWriter object that provides methods for building a HTTP response.
// - r: an *http.Request object that represents the incoming HTTP request.
//
// Return:
// None.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.html", data)
}

// snippetCreatePost handles the creation of a new snippet.
//
// It takes in an http.ResponseWriter and an http.Request as parameters.
// After creating the snippet, it redirects the user to the snippet view page.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError() helper to
	// send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)

	}

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an integer. So we need to manually covert
	// the form data to an integer using strconv.Atoi(), and we send a 400 Bad
	// Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	// Because the Validator type is embedded by the snippetCreateForm struct,
	// we can call CheckField() directly on it to execute our validation checks.
	// CheckField() will add the provided key and error message to the
	// FieldErrors map if the check does not evaluate to true. For example, in
	// the first line here we "check that the form.Title field is not blank". In
	// the second, we "check that the form.Title field has a maximum character
	// length of 100" and so on.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// Use the Valid() method to see if any of the checks failed. If they did,
	// then re-render the template passing in the form in the same way as
	// before.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
