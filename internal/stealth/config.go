package stealth

type Config struct {
	UserAgent string
	Width     int
	Height    int
}

func NewConfig() *Config {
	return &Config{
		UserAgent: randomUserAgent(),
		Width:     randInt(1280, 1920),
		Height:    randInt(720, 1080),
	}
}
