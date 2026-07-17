// Package editorial reviews a draft post using Bedrock, suggesting a
// summary and tags and giving structural/clarity feedback. It never edits
// the post itself - callers surface its Feedback for a human to act on.
package editorial

import (
	"context"
	"fmt"
	"loficode/internal/config"
	"loficode/internal/model"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/rs/zerolog/log"
)

type Feedback struct {
	Summary string
	Tags    []string
	Notes   []string
}

type Editor interface {
	Review(ctx context.Context, post model.Post) (*Feedback, error)
}

const systemPrompt = `You are an editorial assistant reviewing a draft blog post before it's published.

Respond using exactly these three tags, in this order, with no other text before, between, or after them:

<summary>A 1-2 sentence suggested value for the post's frontmatter "summary" field. Lead with the point — no "This post..." or "In this post..." phrasing. Write in the author's voice, direct and opinionated, as if the author is telling a friend what the post is about. Look at the existing summary for tone reference.</summary>
<tags>1 to 5 comma-separated tags for the post, in kebab-case (e.g. "system-design", not "System Design" or "system_design") to match this blog's existing tag convention.</tags>
<notes>
Specific, actionable editorial feedback about clarity, structure, or flow issues you actually noticed in this post, one per line prefixed with "- ". Do not include generic praise or filler. If there is nothing worth flagging, leave this empty.
</notes>`

func New(cfg *config.Config) Editor {
	return &bedrockEditor{
		client:  bedrockruntime.NewFromConfig(cfg.AwsConfig),
		modelId: cfg.BedrockModelId,
	}
}

func NewNoop() Editor {
	return noopEditor{}
}

type noopEditor struct{}

func (noopEditor) Review(_ context.Context, post model.Post) (*Feedback, error) {
	log.Debug().Str("slug", post.Slug).Msg("noopEditor: skipping Bedrock review")
	return &Feedback{}, nil
}

type bedrockEditor struct {
	client  *bedrockruntime.Client
	modelId string
}

func (e *bedrockEditor) Review(ctx context.Context, post model.Post) (*Feedback, error) {
	prompt := fmt.Sprintf(
		"Title: %s\nCurrent tags: %s\nCurrent summary: %s\n\nBody (rendered HTML):\n%s",
		post.Title, strings.Join(post.Tags, ", "), post.Summary, post.Content,
	)

	resp, err := e.client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId: aws.String(e.modelId),
		System: []types.SystemContentBlock{
			&types.SystemContentBlockMemberText{Value: systemPrompt},
		},
		Messages: []types.Message{
			{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{Value: prompt},
				},
			},
		},
		InferenceConfig: &types.InferenceConfiguration{
			MaxTokens:   aws.Int32(1024),
			Temperature: aws.Float32(0.2),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("bedrock converse failed: %w", err)
	}

	text, err := extractText(resp.Output)
	if err != nil {
		return nil, err
	}

	return parseFeedback(text)
}

func extractText(output types.ConverseOutput) (string, error) {
	msgOutput, ok := output.(*types.ConverseOutputMemberMessage)
	if !ok {
		return "", fmt.Errorf("unexpected converse output type %T", output)
	}
	var sb strings.Builder
	for _, block := range msgOutput.Value.Content {
		if textBlock, ok := block.(*types.ContentBlockMemberText); ok {
			sb.WriteString(textBlock.Value)
		}
	}
	if sb.Len() == 0 {
		return "", fmt.Errorf("model returned no text content")
	}
	return sb.String(), nil
}

var (
	summaryTagRe = regexp.MustCompile(`(?s)<summary>(.*?)</summary>`)
	tagsTagRe    = regexp.MustCompile(`(?s)<tags>(.*?)</tags>`)
	notesTagRe   = regexp.MustCompile(`(?s)<notes>(.*?)</notes>`)
)

// parseFeedback extracts the <summary>/<tags>/<notes> tags from the model's
// response. Plain, tag-delimited text is more resilient to parse here than
// asking the model for JSON: free-form prose in notes (quotes, code
// snippets, backticks) is prone to breaking JSON string escaping, and a
// missing or malformed tag just leaves that one field empty instead of
// failing the whole response.
func parseFeedback(text string) (*Feedback, error) {
	summaryMatch := summaryTagRe.FindStringSubmatch(text)
	tagsMatch := tagsTagRe.FindStringSubmatch(text)
	notesMatch := notesTagRe.FindStringSubmatch(text)

	if summaryMatch == nil && tagsMatch == nil && notesMatch == nil {
		return nil, fmt.Errorf("model response had none of the expected <summary>/<tags>/<notes> tags: %s", text)
	}

	feedback := &Feedback{}
	if summaryMatch != nil {
		feedback.Summary = strings.TrimSpace(summaryMatch[1])
	}
	if tagsMatch != nil {
		feedback.Tags = splitTags(tagsMatch[1])
	}
	if notesMatch != nil {
		feedback.Notes = splitNotes(notesMatch[1])
	}
	return feedback, nil
}

func splitTags(raw string) []string {
	var tags []string
	for _, t := range strings.Split(raw, ",") {
		if t = strings.TrimSpace(t); t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func splitNotes(raw string) []string {
	var notes []string
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "-")
		if line = strings.TrimSpace(line); line != "" {
			notes = append(notes, line)
		}
	}
	return notes
}
