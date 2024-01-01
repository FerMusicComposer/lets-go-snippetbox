package main

import (
	"html/template"
	"path/filepath"

	"github.com/FerMusicComposer/lets-go-snippetbox.git/internal/models"
)

// templateData is a struct that has two fields: Snippet and Snippets.
// The Snippet field is a pointer to a models.Snippet object, and the Snippets field
// is a slice of models.Snippet objects.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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
	pages, err := filepath.Glob("../../ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		//Creates a template set beginning with the base templte, to which all partials and pages will be added dynamically
		tmplSet, err := template.ParseFiles("../../ui/html/base.html")
		if err != nil {
			return nil, err
		}

		tmplSet, err = tmplSet.ParseGlob("../../ui/html/partials/*.html")
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
