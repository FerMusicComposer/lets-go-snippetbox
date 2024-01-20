package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/FerMusicComposer/lets-go-snippetbox.git/internal/models"
	"github.com/FerMusicComposer/lets-go-snippetbox.git/ui"
)

// templateData is a struct that has two fields: Snippet and Snippets.
// The Snippet field is a pointer to a models.Snippet object, and the Snippets field
// is a slice of models.Snippet objects.
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
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
// It uses fs.Glob() to get a slice of all filepaths in the ui.Files embedded
// filesystem which match the pattern 'html/pages/*.tmpl'. This essentially
// gives us a slice of all the 'page' templates for the application, just
// like before.
//
// Parameters:
//
//	None
//
// Returns:
//
//	map[string]*template.Template: The initialized template cache.
//	error: An error if any occurred during initialization.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use fs.Glob() to get a slice of all filepaths in the ui.Files embedded
	// filesystem which match the pattern 'html/pages/*.tmpl'. This essentially
	// gives us a slice of all the 'page' templates for the application, just
	// like before.
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}
		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil

}
