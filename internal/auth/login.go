// package auth

// import (
// 	"log"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/go-rod/rod"
// )

// // Login performs a demo-safe LinkedIn login using a persistent browser profile.
// // - If already logged in, it exits early.
// // - If CAPTCHA/checkpoint appears, it waits for manual resolution.
// // - Fails gracefully after limited retries.
// func Login(page *rod.Page) bool {
// 	// 0. Early check: already logged in via persistent profile
// 	if isLoggedIn(page) {
// 		log.Println("Already logged in (persistent profile)")
// 		return true
// 	}

// 	email := os.Getenv("LINKEDIN_EMAIL")
// 	password := os.Getenv("LINKEDIN_PASSWORD")

// 	if email == "" || password == "" {
// 		log.Println("Missing LinkedIn credentials in env")
// 		return false
// 	}

// 	// 1. Navigate to login page
// 	page.MustNavigate("https://www.linkedin.com/login")
// 	page.MustWaitLoad()

// 	// 2. Fill credentials
// 	emailInput := page.Timeout(15 * time.Second).
// 		MustElement(`input[name="session_key"]`)
// 	passInput := page.MustElement(`input[name="session_password"]`)

// 	emailInput.MustClick()
// 	emailInput.MustInput(email)
// 	time.Sleep(700 * time.Millisecond)

// 	passInput.MustClick()
// 	passInput.MustInput(password)
// 	time.Sleep(700 * time.Millisecond)

// 	// 3. Submit login form
// 	page.MustElement(`button[type="submit"]`).MustClick()

// 	// 4. Graceful post-login check loop
// 	const maxChecks = 20
// 	for i := 1; i <= maxChecks; i++ {
// 		time.Sleep(5 * time.Second)

// 		info, err := page.Info()
// 		if err != nil {
// 			log.Println("Page unavailable (possibly closed by user)")
// 			return false
// 		}

// 		url := info.URL
// 		log.Printf("Post-login check %d/%d: %s\n", i, maxChecks, url)

// 		// Success
// 		if strings.Contains(url, "/feed") {
// 			return true
// 		}

// 		// CAPTCHA / checkpoint → wait for human
// 		if strings.Contains(url, "/checkpoint") {
// 			log.Println("Checkpoint/CAPTCHA detected. Waiting for manual completion...")
// 			continue
// 		}

// 		// Hard failure
// 		if strings.Contains(url, "/login") {
// 			log.Println("Redirected back to login page. Login failed.")
// 			return false
// 		}
// 	}

// 	log.Println("Login not completed after maximum wait time")
// 	return false
// }

// // isLoggedIn checks if the session is already authenticated
// func isLoggedIn(page *rod.Page) bool {
// 	info, err := page.Info()
// 	if err != nil {
// 		return false
// 	}
// 	return strings.Contains(info.URL, "/feed")
// }

package auth

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// Login performs LinkedIn login using environment variables and persistent session.
// If already authenticated, login is skipped.
func Login(page *rod.Page) bool {

	// 0: Fast check — already logged in?
	if isLoggedIn(page) {
		log.Println("Session already active (persistent Chrome profile). Skipping login.")
		return true
	}

	// 1: Read credentials
	email := os.Getenv("LINKEDIN_EMAIL")
	password := os.Getenv("LINKEDIN_PASSWORD")

	if email == "" || password == "" {
		log.Println("Missing LinkedIn credentials in env")
		return false
	}

	log.Println("Performing LinkedIn login...")

	// 2: Navigate to login page
	page.MustNavigate("https://www.linkedin.com/login")
	page.MustWaitLoad()

	// 3: Fill in email
	emailField := page.Timeout(15 * time.Second).MustElement(`input[name="session_key"]`)
	emailField.MustClick()
	emailField.MustInput(email)
	time.Sleep(400 * time.Millisecond)

	// 4: Fill password
	passField := page.MustElement(`input[name="session_password"]`)
	passField.MustClick()
	passField.MustInput(password)
	time.Sleep(400 * time.Millisecond)

	// 5: Submit form
	page.MustElement(`button[type="submit"]`).MustClick()

	// 6: Post-login flow validation
	const maxChecks = 25

	for i := 1; i <= maxChecks; i++ {
		time.Sleep(3 * time.Second)

		info, err := page.Info()
		if err != nil {
			log.Println("Login interrupted: page closed or unreachable")
			return false
		}

		url := info.URL
		log.Printf("Login progress check %d/%d → %s\n", i, maxChecks, url)

		// -----------------------------------------
		// SUCCESS CONDITION
		// -----------------------------------------
		if strings.Contains(url, "/feed") {
			log.Println("Login successful ✔")
			return true
		}

		// -----------------------------------------
		// CAPTCHA or Security Checkpoint
		// -----------------------------------------
		if strings.Contains(url, "/checkpoint") {
			log.Println("⚠ CAPTCHA or LinkedIn Checkpoint detected.")
			log.Println("Please complete it manually in the browser...")
			continue
		}

		// -----------------------------------------
		// Login failure — returned back to login page
		// -----------------------------------------
		if strings.Contains(url, "/login") {
			if i > 5 { // Give LinkedIn time to redirect
				log.Println("❌ Login failed — redirected back to login page.")
				return false
			}
		}
	}

	log.Println("❌ Login timed out after waiting too long.")
	return false
}

// isLoggedIn determines if LinkedIn is already authenticated
func isLoggedIn(page *rod.Page) bool {
	info, err := page.Info()
	if err != nil {
		return false
	}
	return strings.Contains(info.URL, "/feed")
}
