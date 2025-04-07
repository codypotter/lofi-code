package application

import (
	"loficode/internal/templates/components"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

func (a application) SearchResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()
	cursor := query.Get("cursor")
	log := log.With().Str("handler", "SearchResults").Str("cursor", cursor).Str("method", r.Method).Logger()
	currentUrl := r.Header.Get("HX-Current-URL")
	if currentUrl != "" {
		log.Debug().Str("currentUrl", currentUrl).Msg("Current URL")
		parsedURL, err := url.Parse(currentUrl)
		if err == nil {
			log.Debug().Stringer("parsedURL", parsedURL).Msg("Parsed URL")
			query = parsedURL.Query()
		}
	}

	tag := query.Get("tag")
	if tag == "" {
		tag = "all"
	}
	log = log.With().Str("tag", tag).Logger()

	posts, nextCursor, err := a.db.GetPostsByTag(ctx, tag, cursor)
	if err != nil {
		log.Error().Err(err).Msg("Error getting posts by tag")
		components.Notification("is-danger", "Error getting posts").Render(r.Context(), w)
		return
	}

	log.Debug().Int("posts", len(posts)).Msg("Got posts")
	components.SearchResults(posts, nextCursor).Render(r.Context(), w)
}
