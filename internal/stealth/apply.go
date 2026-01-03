package stealth

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func Apply(page *rod.Page, cfg *Config) {
	// 1. Set realistic viewport
	page.MustSetViewport(cfg.Width, cfg.Height, 1, false)

	// 2. Override User-Agent (correct Rod v0.114.6 API)
	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: cfg.UserAgent,
	})

	// 3. Hide webdriver flag before any site JS runs
	page.MustEvalOnNewDocument(`
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined
		});
	`)
}
