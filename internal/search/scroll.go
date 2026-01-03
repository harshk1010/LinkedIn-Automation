package search

import (
	"github.com/go-rod/rod"

	"linkedin-automation/internal/stealth"
)

// LightScroll performs minimal scrolling to mimic human behavior
func LightScroll(page *rod.Page) {
	// Scroll down about one viewport
	page.MustEval(`() => {
		window.scrollBy(0, window.innerHeight);
	}`)

	stealth.RandomDelay(800, 1500)

	// Optional slight scroll back up
	page.MustEval(`() => {
		window.scrollBy(0, -200);
	}`)

	stealth.RandomDelay(500, 1000)
}
