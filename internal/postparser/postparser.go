package postparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"loficode/internal/model"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// ParseDir walks dir for markdown post files and parses each into a model.Post.
func ParseDir(dir string) ([]model.Post, error) {
	var ps []model.Post

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			log.Info().Str("path", path).Msgf("Skipping")
			return nil
		}
		p, err := ParseFile(path)
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

// ParseFile parses a single markdown post file, extracting its frontmatter
// metadata into a model.Post and rendering its body to HTML in post.Content.
func ParseFile(path string) (*model.Post, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	context := parser.NewContext()
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
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
