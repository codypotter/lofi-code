// cmd/prreview builds an editorial PR comment for changed blog posts. It
// takes changed markdown file paths as args, runs deterministic checks plus
// a Bedrock editorial review on each, and prints one combined Markdown
// comment to stdout. It never modifies post files - a human decides what to
// take from its suggestions.
package main

import (
	"context"
	"fmt"
	"loficode/internal/config"
	"loficode/internal/editorial"
	"loficode/internal/logger"
	"loficode/internal/model"
	"loficode/internal/postlint"
	"loficode/internal/postparser"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const commentMarker = "<!-- lofi-code-editorial-bot -->"

func main() {
	ctx := context.Background()
	cfg := config.New(ctx)
	logger.Configure(cfg.LogLevel, logger.LogFormatConsole)

	changedFiles := os.Args[1:]
	if len(changedFiles) == 0 {
		log.Info().Msg("No changed post files to review")
		return
	}

	allPosts, err := postparser.ParseDir("cms/_posts")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse existing posts")
	}

	var editor editorial.Editor
	if cfg.Environment == "development" {
		editor = editorial.NewNoop()
	} else {
		editor = editorial.New(cfg)
	}

	var sb strings.Builder
	sb.WriteString(commentMarker + "\n")
	sb.WriteString("## Editorial review\n\n")

	for _, path := range changedFiles {
		if !strings.HasSuffix(path, ".md") {
			continue
		}

		post, err := postparser.ParseFile(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Failed to parse post")
			fmt.Fprintf(&sb, "### `%s`\n\nFailed to parse this post: %s\n\n", path, err)
			continue
		}

		renderPostReview(ctx, &sb, path, *post, allPosts, editor)
	}

	fmt.Print(sb.String())
}

func renderPostReview(ctx context.Context, sb *strings.Builder, path string, post model.Post, allPosts []model.Post, editor editorial.Editor) {
	fmt.Fprintf(sb, "### `%s`\n\n", path)

	findings := postlint.Check(post, allPosts)
	if len(findings) > 0 {
		sb.WriteString("**Checks:**\n")
		for _, f := range findings {
			fmt.Fprintf(sb, "- %s\n", f.Message)
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("**Checks:** all good\n\n")
	}

	reviewCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	feedback, err := editor.Review(reviewCtx, post)
	if err != nil {
		log.Error().Err(err).Str("slug", post.Slug).Msg("Editorial review failed")
		sb.WriteString("_Editorial (LLM) review failed - see workflow logs._\n\n")
		return
	}

	if feedback.Summary != "" && feedback.Summary != post.Summary {
		fmt.Fprintf(sb, "**Suggested summary:** %s\n\n", feedback.Summary)
	}
	if len(feedback.Tags) > 0 {
		fmt.Fprintf(sb, "**Suggested tags:** %s\n\n", strings.Join(feedback.Tags, ", "))
	}
	if len(feedback.Notes) > 0 {
		sb.WriteString("**Editorial notes:**\n")
		for _, n := range feedback.Notes {
			fmt.Fprintf(sb, "- %s\n", n)
		}
		sb.WriteString("\n")
	}
}
