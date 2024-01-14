package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/FerMusicComposer/lets-go-snippetbox.git/internal/models"
)

// templateData is a struct that has two fields: Snippet and Snippets.
// The Snippet field is a pointer to a models.Snippet object, and the Snippets field
// is a slice of models.Snippet objects.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
}

// humanDate formats a given time.Time value into a string representation.
//
// It takes a time.Time parameter and returns a string.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// newTemplateCache initializes a new template cache.
//
// It uses the filepath.Glob() function to get a slice of all filepaths that
// match the pattern "./ui/html/pages/*.tmpl". This will essentially give us
// a slice of all the filepaths for our application 'page' templates like:
// [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
//
// Returns:
//   - map[string]*template.Template: A map of template names to template sets.
//   - error: An error if there was a problem initializing the template cache.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.html". This will essentially give
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.html ui/html/pages/view.html]
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		//Creates a template set beginning with the base template, to which all partials and pages will be added dynamically
		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		tmplSet, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		tmplSet, err = tmplSet.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		tmplSet, err = tmplSet.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = tmplSet
	}

	return cache, nil

}
