package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"loficode/model"
	"loficode/templates/pages/home"
	"loficode/templates/pages/notfound"
	"loficode/templates/pages/post"
	"loficode/templates/pages/posts"
	"loficode/templates/pages/privacypolicy"
	"loficode/templates/pages/tos"
	"sort"

	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	var ps []model.Post

	err := os.MkdirAll("public/posts", 0755)
	if err != nil {
		return fmt.Errorf("failed to create public/posts directory: %w", err)
	}

	err = filepath.WalkDir("cms/_posts", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			fmt.Printf("Skipping: %s\n", path)
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
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	fmt.Printf("Parsed %d posts\n", len(ps))

	allTags := extractTags(ps)
	fmt.Printf("Found %d tags\n", len(allTags))

	recentPosts := extractRecentPosts(ps)

	staticPages := map[string]templ.Component{
		"index.html":          home.Home(allTags, recentPosts),
		"posts.html":          posts.Posts(allTags),
		"tos.html":            tos.TermsOfService(),
		"privacy-policy.html": privacypolicy.PrivacyPolicy(),
		"404.html":            notfound.NotFound(),
	}
	for _, p := range ps {
		htmlOut := filepath.Join("posts", p.Slug+".html")
		staticPages[htmlOut] = post.Post(p)
	}

	for name, component := range staticPages {
		err := renderStaticPage(name, component)
		if err != nil {
			return fmt.Errorf("failed to render static page: %w", err)
		}
	}

	return nil
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
	fmt.Println("Wrote:", path)
	return nil
}

func parseMarkdownFile(path string) (*model.Post, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	context := parser.NewContext()
	md := goldmark.New(
		goldmark.WithExtensions(meta.Meta),
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

	fmt.Printf("Parsed: %s\n", post.Title)

	post.Content = buf.String()
	return &post, nil
}

// extractTags extracts all the tags from the posts.
// It returns a unique list of tags.
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

// extractRecentPosts extracts the 3 most recent posts.
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
