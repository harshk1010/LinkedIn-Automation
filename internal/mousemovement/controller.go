package mousemovement

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// MouseController provides human-like mouse interactions using Rod
type MouseController struct {
	page *rod.Page
	rand *rand.Rand
}

// New creates a new mouse controller bound to a page
func New(page *rod.Page) *MouseController {
	return &MouseController{
		page: page,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Hover hovers over an element with a human-like delay
func (mc *MouseController) Hover(el *rod.Element) error {
	time.Sleep(time.Duration(120+mc.rand.Intn(200)) * time.Millisecond)
	return el.Hover()
}

// ClickWithDelay clicks an element with human-like hesitation
func (mc *MouseController) ClickWithDelay(el *rod.Element) error {
	if err := mc.Hover(el); err != nil {
		return err
	}

	time.Sleep(time.Duration(150+mc.rand.Intn(150)) * time.Millisecond)
	return el.Click(proto.InputMouseButtonLeft, 1)
}

// ScrollAndHover scrolls an element into view and hovers
func (mc *MouseController) ScrollAndHover(el *rod.Element) error {
	el.ScrollIntoView()
	time.Sleep(time.Duration(300+mc.rand.Intn(400)) * time.Millisecond)
	return mc.Hover(el)
}

// Wait creates a small human-like pause
func (mc *MouseController) Wait() {
	time.Sleep(time.Duration(200+mc.rand.Intn(300)) * time.Millisecond)
}
