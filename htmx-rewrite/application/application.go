package application

import (
	"context"
	"fmt"
	"loficode/config"
	"loficode/db"
	"loficode/model"
	"loficode/templates/components"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi"
)

type application struct {
	db *db.Db
}

func New(ctx context.Context, c *config.Config) application {
	return application{
		db: db.New(context.Background(), c),
	}
}

// /api/posts/:slug/comments
func (a *application) Comments(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	comments, err := a.db.GetCommentsBySlug(slug)
	if err != nil {
		log.Printf("Error getting comments: %v\n", err)
		components.Notification("is-warning", "Error getting comments").Render(r.Context(), w)
	}

	components.Comments(comments).Render(r.Context(), w)
}

func (a *application) CommentForm(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	name := r.FormValue("name")
	email := r.FormValue("email")
	comment := r.FormValue("comment")

	if name == "" || email == "" || comment == "" {
		components.CommentForm(components.CommentFormConfig{
			Slug:         slug,
			Name:         name,
			Email:        email,
			Comment:      comment,
			ErrorMessage: "Please enter a valid name, email, and comment.",
		}).Render(context.Background(), w)
		return
	}
	if err := a.db.AddComment(slug, model.Comment{
		Name:  name,
		Email: email,
		Text:  comment,
		Date:  time.Now(),
	}); err != nil {
		log.Printf("Error adding comment: %v\n", err)
		components.Notification("is-danger", "Error adding comment").Render(r.Context(), w)
		return
	}
	components.CommentForm(components.CommentFormConfig{
		Slug:    slug,
		Name:    name,
		Email:   email,
		Comment: comment,
		Success: true,
	}).Render(context.Background(), w)
}

func (a *application) PostPreviews(w http.ResponseWriter, r *http.Request) {
	components.PostPreviews([]model.Post{
		{
			Title:   "Hello, World!",
			Slug:    "hello-world",
			Summary: "This is a description of the post.",
		},
		{
			Title:   "Hello, World!",
			Slug:    "hello-world-1",
			Summary: "This is a description of the post.",
		},
		{
			Title:   "Hello, World!",
			Slug:    "hello-world-2",
			Summary: "This is a description of the post.",
		},
	}).Render(r.Context(), w)
}

func (a *application) SearchResults(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	cursor := query.Get("cursor")
	currentURL := r.Header.Get("HX-Current-URL")
	if currentURL != "" {
		parsedURL, err := url.Parse(currentURL)
		if err == nil {
			query = parsedURL.Query()
		}
	}

	tag := query.Get("tag")
	if tag == "" {
		tag = "all"
	}
	fmt.Printf("SearchResults: %+v\n", tag)
	fmt.Printf("SearchResults: %+v\n", cursor)

	posts, nextCursor, err := a.db.GetPostsByTag(tag, cursor)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}
	components.SearchResults(posts, nextCursor).Render(r.Context(), w)
}

func (a *application) Tags(w http.ResponseWriter, r *http.Request) {
	components.Tags(components.TagsConfig{
		Size:             "is-large",
		Tags:             []string{"tag1", "tag2", "tag3"},
		EnableNavigation: true,
	}).Render(r.Context(), w)
}

func (a *application) Subscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		components.MailingListForm(components.MailingListConfig{
			ErrorMessage: "Please enter a valid email address.",
		})
		return
	}
	components.MailingListForm(components.MailingListConfig{
		SuccessMessage: "Thank you for subscribing!",
	}).Render(r.Context(), w)
}
