package browser

import (
	"log"
	"os"
	"path/filepath"

	"linkedin-automation/internal/stealth"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Browser struct {
	Instance *rod.Browser
	Stealth  *stealth.Config
}

func New() *Browser {
	// Create persistent profile folder
	profileDir := filepath.Join(".", "chrome-profile")
	os.MkdirAll(profileDir, 0755)

	log.Println("Using persistent profile:", profileDir)

	// Launch Rod-managed Chromium (NOT Chrome)
	u := launcher.New().
		Leakless(false).         // avoid antivirus false positives
		UserDataDir(profileDir). // persistent session
		Headless(false).         // show Chromium window
		Set("disable-gpu").
		Set("no-sandbox").
		Set("disable-dev-shm-usage").
		Set("disable-blink-features", "AutomationControlled").
		MustLaunch()

	log.Println("Chromium launched at:", u)

	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		log.Fatalf("Failed to connect to Chromium: %v", err)
	}

	return &Browser{
		Instance: browser,
		Stealth:  stealth.NewConfig(),
	}
}

func (b *Browser) NewPage(url string) *rod.Page {
	page := b.Instance.MustPage()
	stealth.Apply(page, b.Stealth)

	if url != "" {
		page.MustNavigate(url).MustWaitLoad()
	}
	return page
}

func (b *Browser) Close() {
	_ = b.Instance.Close()
}
