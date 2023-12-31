package main

import "github.com/FerMusicComposer/lets-go-snippetbox.git/internal/models"

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}
