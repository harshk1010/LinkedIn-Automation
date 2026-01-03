package search

import (
	"log"
	"time"
	"math/rand"

	"github.com/go-rod/rod"

	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/stealth"
)

// Navigator handles safe navigation to LinkedIn search pages
type Navigator struct {
	Browser *browser.Browser
}

// NewNavigator creates a new search navigator
func NewNavigator(b *browser.Browser) *Navigator {
	return &Navigator{
		Browser: b,
	}
}

// OpenSearch navigates to the LinkedIn people search page
func (n *Navigator) OpenSearch(searchURL string) *rod.Page {
	log.Printf("Navigating to search URL")

	page := n.Browser.NewPage(searchURL)

	// Ensure results container exists (basic sanity check)
	page.MustElement("div.search-results-container")

	stealth.RandomDelay(1500, 3000)

	return page
}

// randomDelay pauses execution for a human-like duration
func randomDelay(minMs, maxMs int) {
	delay := rand.Intn(maxMs-minMs+1) + minMs
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
