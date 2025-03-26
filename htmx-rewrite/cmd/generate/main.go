package main

import (
	"bytes"
	"context"
	"fmt"
	"loficode/templates/pages/post"

	"io/fs"
	"loficode/templates/components"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	err := os.MkdirAll("public/posts", 0755)
	if err != nil {
		panic(err)
	}
	err = filepath.WalkDir("cms", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		p, err := parseMarkdownFile(path)
		if err != nil {
			return err
		}
		fmt.Printf("Parsed post: %+v\n", p)

		htmlOut := filepath.Join("public/posts", p.Slug+".html")
		w, err := os.Create(htmlOut)
		if err != nil {
			return err
		}
		err = post.Post(*p).Render(context.Background(), w)
		if err != nil {
			return err
		}
		fmt.Println("Wrote:", htmlOut)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func parseMarkdownFile(path string) (*components.Post, error) {
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

	metaData := meta.Get(context)
	post := &components.Post{
		Content: buf.String(),
	}

	if title, ok := metaData["title"].(string); ok {
		post.Name = title
	}
	if slug, ok := metaData["slug"].(string); ok {
		post.Slug = slug
	}

	return post, nil
}
