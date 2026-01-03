package search

import (
	"net/url"
	"strconv"
)

const peopleSearchBase = "https://www.linkedin.com/search/results/people/"

// BuildPeopleSearchURL returns a LinkedIn people search URL for a given page
func BuildPeopleSearchURL(input *Input, page int) string {
	params := url.Values{}
	params.Set("keywords", input.Keywords)

	if page > 1 {
		params.Set("page", strconv.Itoa(page))
	}

	return peopleSearchBase + "?" + params.Encode()
}
