package main

import (
	"log"
	"strings"
	"time"

	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/mousemovement"
	"linkedin-automation/internal/search"
	"linkedin-automation/internal/state"
	"linkedin-automation/internal/stealth"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting LinkedIn Automation Bot (Stage 4 PoC)")

	// Load environment variables (used for first-time authentication only)
	_ = godotenv.Load()

	// Initialize persistent state store (SQLite)
	store, err := state.NewStore("state.db")
	if err != nil {
		log.Fatalf("failed to initialize state store: %v", err)
	}
	defer store.Close()

	// Launch browser with persistent Chromium profile
	b := browser.New()
	defer b.Close()

	// Open LinkedIn entry page
	page := b.NewPage("https://www.linkedin.com")
	time.Sleep(2 * time.Second)

	// Verify authenticated session
	if strings.Contains(page.MustInfo().URL, "/feed") {
		log.Println("✅ Already logged in via persistent browser profile")
	} else {
		log.Println("No active session detected, performing login")

		if !auth.Login(page) {
			log.Fatal("Login failed")
		}

		log.Println("Login successful (session persisted)")
	}

	// ---- Stage 4A: Discover & Persist Profiles ----
	discoverProfiles(b, store)

	// ---- Stage 4B: Visit Persisted Profiles ----
	visitProfiles(b, store)

	log.Println("Stage 4 execution complete")
	time.Sleep(10 * time.Second)
}

/*
discoverProfiles

Responsibility:
- Navigate LinkedIn People Search pages
- Simulate light human browsing behavior
- Parse visible profile URLs
- Persist discovered profiles into state store
*/
func discoverProfiles(b *browser.Browser, store *state.Store) {
	log.Println("\n=== Stage 4A | Profile Discovery & Persistence ===")

	input := &search.Input{
		Keywords:  "software engineer",
		PageLimit: 3,
	}

	if err := input.Validate(); err != nil {
		log.Fatalf("invalid search input: %v", err)
	}

	for pageNum := 1; pageNum <= input.PageLimit; pageNum++ {
		log.Printf("\n[Discovery] Navigating to search page %d", pageNum)

		searchURL := search.BuildPeopleSearchURL(input, pageNum)
		page := b.NewPage(searchURL)

		// Mouse controller used only for subtle human-like behavior
		mouse := mousemovement.New(page)

		// Allow page to load naturally
		stealth.RandomDelay(1200, 2000)

		// Light hover over results container (optional realism)
		if container, err := page.Element("div.search-results-container"); err == nil {
			_ = mouse.Hover(container)
			stealth.RandomDelay(400, 700)
		}

		// Light scroll to hydrate lazy-loaded content
		search.LightScroll(page)
		mouse.Wait()

		// Parse visible profile results (read-only)
		profiles := search.ParseVisibleResults(page)

		// Persist discovered profiles
		for _, p := range profiles {
			if err := store.EnsureProfile(p.URL); err != nil {
				log.Printf("[State] Failed to persist profile %s: %v", p.URL, err)
				continue
			}
			log.Printf("[State] Profile persisted: %s", p.URL)
		}

		log.Printf(
			"[Discovery] Page %d complete — profiles discovered: %d",
			pageNum,
			len(profiles),
		)

		// Close page to avoid tab accumulation
		_ = page.Close()
	}

	log.Println("=== Profile discovery complete ===")
}

/*
visitProfiles

Responsibility:
- Visit previously discovered profiles
- Simulate natural reading behavior
- Update visit status in state store
*/
func visitProfiles(b *browser.Browser, store *state.Store) {
	log.Println("\n=== Stage 4B | Profile Visit Loop ===")

	// Fetch a limited number of unvisited profiles
	profiles, err := store.GetPendingProfiles(5)
	if err != nil {
		log.Fatalf("failed to fetch pending profiles: %v", err)
	}

	for _, p := range profiles {
		url := p.ProfileURL
		log.Printf("[Visit] Opening profile: %s", url)

		page := b.NewPage(url)
		mouse := mousemovement.New(page)

		// Allow profile to load fully
		stealth.RandomDelay(2500, 3500)

		// Simulate reading behavior
		search.LightScroll(page)
		mouse.Wait()

		// Mark profile as visited
		if err := store.MarkVisited(url); err != nil {
			log.Printf("[State] Failed to mark visited: %s", url)
		} else {
			log.Printf("[State] Marked visited: %s", url)
		}

		_ = page.Close()
		stealth.RandomDelay(1500, 2500)
	}

	log.Println("=== Profile visit loop complete ===")
}
