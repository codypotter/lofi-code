package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"loficode/internal/config"
	"loficode/internal/db"
	"loficode/internal/logger"
	"loficode/internal/model"
	errorpage "loficode/internal/templates/pages/error"
	"loficode/internal/templates/pages/home"
	"loficode/internal/templates/pages/notfound"
	"loficode/internal/templates/pages/post"
	"loficode/internal/templates/pages/posts"
	"loficode/internal/templates/pages/privacypolicy"
	"loficode/internal/templates/pages/tos"
	"loficode/internal/templates/pages/unsubscribe"
	"loficode/internal/templates/pages/verified"
	"sort"

	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	ctx := context.Background()
	cfg := config.New(ctx)
	logger.Configure(cfg.LogLevel, logger.LogFormatConsole)

	var database *db.Db
	if cfg.Environment == "development" {
		log.Debug().Msg("Running in development mode")
		database = db.NewDevelopment(ctx, cfg)
	} else {
		log.Debug().Msg("Running in production mode")
		database = db.New(ctx, cfg)
	}
	database.CreateTable()

	ps, err := parseMarkdownFiles()
	if err != nil {
		log.Error().Err(err).Msg("Error parsing markdown files")
		return
	}
	log.Debug().Msg("Parsed markdown files")

	tags := extractTags(ps)
	log.Debug().Msg("Extracted tags")

	recentPosts := extractRecentPosts(ps)

	err = renderStaticPages(ps, tags, recentPosts)
	if err != nil {
		log.Error().Err(err).Msg("Error rendering static pages")
		return
	}

	database.UpsertPosts(ctx, ps)
}

func renderStaticPages(ps []model.Post, tags []string, recentPosts []model.Post) error {
	staticPages := map[string]templ.Component{
		"index.html":          home.Home(tags, recentPosts),
		"posts.html":          posts.Posts(tags),
		"tos.html":            tos.TermsOfService(),
		"privacy-policy.html": privacypolicy.PrivacyPolicy(),
		"404.html":            notfound.NotFound(),
		"unsubscribe.html":    unsubscribe.Unsubscribe(),
		"verified.html":       verified.Verified(),
		"error.html":          errorpage.Error(),
	}
	baseUrl := config.New(context.Background()).BaseUrl
	for _, p := range ps {
		relatedPosts := getRelatedPosts(p, ps)
		htmlOut := filepath.Join("posts", p.Slug+".html")
		staticPages[htmlOut] = post.Post(p, baseUrl, relatedPosts)
	}

	for name, component := range staticPages {
		err := renderStaticPage(name, component)
		if err != nil {
			return fmt.Errorf("failed to render static page: %w", err)
		}
	}
	return nil
}

func getRelatedPosts(p model.Post, ps []model.Post) []model.Post {
	var relatedPosts []model.Post
	for _, post := range ps {
		var intersection []string
		for _, t1 := range p.Tags {
			for _, t2 := range post.Tags {
				if t1 == t2 {
					intersection = append(intersection, t1)
				}
			}
		}
		if post.Slug != p.Slug && len(relatedPosts) < 3 && len(intersection) > 0 {
			relatedPosts = append(relatedPosts, post)
		}

		sort.Slice(relatedPosts, func(i, j int) bool {
			return len(relatedPosts[i].Tags) > len(relatedPosts[j].Tags)
		})
	}
	return relatedPosts
}

func parseMarkdownFiles() ([]model.Post, error) {
	var ps []model.Post

	err := os.MkdirAll("public/posts", 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create public/posts directory: %w", err)
	}

	err = filepath.WalkDir("cms/_posts", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			log.Info().Str("path", path).Msgf("Skipping")
			return nil
		}
		p, err := parseMarkdownFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse markdown file: %w", err)
		}
		ps = append(ps, *p)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return ps, nil
}

func renderStaticPage(name string, component templ.Component) error {
	path := filepath.Join("public", name)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}
	defer f.Close()

	if err := component.Render(context.Background(), f); err != nil {
		return fmt.Errorf("failed to render %s: %w", path, err)
	}
	log.Debug().Str("path", path).Msg("Rendered static page")
	return nil
}

func parseMarkdownFile(path string) (*model.Post, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	context := parser.NewContext()
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("xcode-dark"),
			),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	err = md.Convert(input, &buf, parser.WithContext(context))
	if err != nil {
		return nil, err
	}

	var post model.Post

	metaData := meta.Get(context)
	metaJson, err := json.Marshal(metaData)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(metaJson, &post)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("slug", post.Slug).Msg("Parsed post slug")

	post.Content = buf.String()
	return &post, nil
}

func extractTags(ps []model.Post) []string {
	var tags []string
	tagMap := make(map[string]bool)
	for _, p := range ps {
		for _, t := range p.Tags {
			if _, ok := tagMap[t]; !ok {
				tagMap[t] = true
				tags = append(tags, t)
			}
		}
	}
	return tags
}

func extractRecentPosts(ps []model.Post) []model.Post {
	sorted := make([]model.Post, len(ps))
	copy(sorted, ps)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.After(sorted[j].Date)
	})

	if len(sorted) > 3 {
		return sorted[:3]
	}
	return sorted
}
