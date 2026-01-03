package search

import (
	"log"
	"strings"

	"github.com/go-rod/rod"
)

type Profile struct {
	URL  string
	Name string // optional
}

// ParseVisibleResults - Stage 4 (Phase 3)
// Read-only parsing of visible LinkedIn search results
func ParseVisibleResults(page *rod.Page) []Profile {
	log.Println("[Phase 3] Parsing visible search results...")

	var profiles []Profile
	seen := make(map[string]bool)

	// Collect anchors that look like profile links
	links, err := page.Elements(`a[href*="/in/"]`)
	if err != nil {
		log.Println("[Parser] No profile links found")
		return profiles
	}

	log.Printf("[Parser] Found %d raw profile links", len(links))

	for _, link := range links {
		href, err := link.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		url := cleanProfileURL(*href)
		if url == "" {
			continue
		}

		// Deduplication (per page)
		if seen[url] {
			log.Printf("[Duplicate] Skipping: %s", url)
			continue
		}
		seen[url] = true

		// Best-effort name extraction (very conservative)
		name := extractName(link)

		profiles = append(profiles, Profile{
			URL:  url,
			Name: name,
		})

		log.Printf("[Parsed] %s", url)
	}

	log.Printf("[Phase 3 Complete] Found %d unique profiles", len(profiles))
	return profiles
}

// cleanProfileURL keeps only real-looking LinkedIn profile URLs
func cleanProfileURL(raw string) string {
	if !strings.Contains(raw, "/in/") {
		return ""
	}

	// Make absolute
	if strings.HasPrefix(raw, "/") {
		raw = "https://www.linkedin.com" + raw
	}

	// Remove tracking params
	if idx := strings.Index(raw, "?"); idx != -1 {
		raw = raw[:idx]
	}

	// Reject obviously non-human or system-generated profiles
	// (service cards, internal URNs, etc.)
	if strings.Contains(raw, "ACoAA") {
		return ""
	}

	if !strings.HasPrefix(raw, "https://www.linkedin.com/in/") {
		return ""
	}

	return raw
}

// extractName tries to read visible link text if it looks human
func extractName(link *rod.Element) string {
	text, err := link.Text()
	if err != nil {
		return ""
	}

	text = strings.TrimSpace(text)

	// Very conservative: reject long promotional strings
	if text == "" || len(text) > 50 {
		return ""
	}

	// Remove common UI noise
	if strings.Contains(strings.ToLower(text), "provides services") {
		return ""
	}

	return text
}
