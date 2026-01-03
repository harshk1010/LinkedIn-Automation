package search

import "errors"

type Input struct {
	Keywords  string
	PageLimit int
}

func (i *Input) Validate() error {
	if i.Keywords == "" {
		return errors.New("keywords cannot be empty")
	}
	if i.PageLimit <= 0 {
		return errors.New("page limit must be greater than zero")
	}
	return nil
}
