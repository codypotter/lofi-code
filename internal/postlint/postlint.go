// Package postlint runs cheap, deterministic checks against a post that
// don't need an LLM: missing alt text, dead internal links, and gaps in
// the frontmatter fields the Decap CMS config requires.
package postlint

import (
	"fmt"
	"loficode/internal/model"
	"regexp"
	"strings"
)

type Finding struct {
	Message string
}

var (
	imgMissingAltRe = regexp.MustCompile(`<img\s+src="([^"]*)"\s+alt=""`)
	linkHrefRe      = regexp.MustCompile(`<a\s+href="([^"]+)"`)
	internalPostRe  = regexp.MustCompile(`^(?:https?://(?:www\.)?loficode\.com)?/posts/([^/?#]+)`)
	kebabCaseRe     = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
)

// Check runs all deterministic checks against post. allPosts is the full set
// of known posts (including post itself), used to validate internal links.
func Check(post model.Post, allPosts []model.Post) []Finding {
	var findings []Finding
	findings = append(findings, checkFrontmatter(post)...)
	findings = append(findings, checkImageAltText(post)...)
	findings = append(findings, checkInternalLinks(post, allPosts)...)
	return findings
}

func checkFrontmatter(post model.Post) []Finding {
	var findings []Finding

	if strings.TrimSpace(post.Title) == "" {
		findings = append(findings, Finding{Message: "Missing `title` in frontmatter."})
	}
	if strings.TrimSpace(post.Slug) == "" {
		findings = append(findings, Finding{Message: "Missing `slug` in frontmatter."})
	}
	if strings.TrimSpace(post.Summary) == "" {
		findings = append(findings, Finding{Message: "Missing `summary` in frontmatter."})
	}
	if strings.TrimSpace(post.OpenGraphImage) == "" {
		findings = append(findings, Finding{Message: "Missing `openGraphImage` in frontmatter."})
	}
	if len(post.Tags) == 0 {
		findings = append(findings, Finding{Message: "No `tags` set; the CMS requires between 1 and 5."})
	} else if len(post.Tags) > 5 {
		findings = append(findings, Finding{Message: fmt.Sprintf("%d tags set; the CMS only allows up to 5.", len(post.Tags))})
	}
	for _, tag := range post.Tags {
		if !kebabCaseRe.MatchString(tag) {
			findings = append(findings, Finding{
				Message: fmt.Sprintf("Tag `%s` isn't kebab-case; tags are used as exact-match DynamoDB keys, so inconsistent casing/formatting splits a topic across multiple tag pages.", tag),
			})
		}
	}

	return findings
}

func checkImageAltText(post model.Post) []Finding {
	var findings []Finding
	for _, match := range imgMissingAltRe.FindAllStringSubmatch(post.Content, -1) {
		findings = append(findings, Finding{
			Message: fmt.Sprintf("Image `%s` has no alt text.", match[1]),
		})
	}
	return findings
}

func checkInternalLinks(post model.Post, allPosts []model.Post) []Finding {
	knownSlugs := make(map[string]bool, len(allPosts))
	for _, p := range allPosts {
		knownSlugs[p.Slug] = true
	}

	var findings []Finding
	seen := make(map[string]bool)
	for _, match := range linkHrefRe.FindAllStringSubmatch(post.Content, -1) {
		href := match[1]
		linkMatch := internalPostRe.FindStringSubmatch(href)
		if linkMatch == nil {
			continue
		}
		slug := linkMatch[1]
		if knownSlugs[slug] || seen[slug] {
			continue
		}
		seen[slug] = true
		findings = append(findings, Finding{
			Message: fmt.Sprintf("Link to `/posts/%s` doesn't match any known post slug.", slug),
		})
	}
	return findings
}
