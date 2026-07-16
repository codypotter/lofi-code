// Package editorial reviews a draft post using Bedrock, suggesting a
// summary and tags and giving structural/clarity feedback. It never edits
// the post itself - callers surface its Feedback for a human to act on.
package editorial

import (
	"context"
	"encoding/json"
	"fmt"
	"loficode/internal/config"
	"loficode/internal/model"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/rs/zerolog/log"
)

type Feedback struct {
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`
	Notes   []string `json:"notes"`
}

type Editor interface {
	Review(ctx context.Context, post model.Post) (*Feedback, error)
}

const systemPrompt = `You are an editorial assistant reviewing a draft blog post before it's published.

Respond with ONLY a single JSON object, no other text, matching this shape:
{"summary": "...", "tags": ["...", "..."], "notes": ["...", "..."]}

- summary: a 1-2 sentence suggested value for the post's frontmatter "summary" field, written in the author's voice.
- tags: 1 to 5 suggested tags for the post, in kebab-case (e.g. "system-design", not "System Design" or "system_design") to match this blog's existing tag convention.
- notes: specific, actionable editorial feedback about clarity, structure, or flow issues you actually noticed in this post. Do not include generic praise or filler. If there is nothing worth flagging, return an empty array.`

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

	var feedback Feedback
	if err := json.Unmarshal([]byte(stripCodeFence(text)), &feedback); err != nil {
		return nil, fmt.Errorf("failed to parse model response as JSON: %w (raw: %s)", err, text)
	}
	return &feedback, nil
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

// stripCodeFence removes a leading/trailing markdown code fence in case the
// model wrapped its JSON response in one despite instructions not to.
func stripCodeFence(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	return strings.TrimSpace(text)
}
