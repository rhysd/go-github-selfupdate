package selfupdate

import (
	"fmt"
	"regexp"
)

type settings struct {
	filters []*regexp.Regexp
}

// Option defines an optional feature for the self updater
type Option func(*settings)

func defaultSettings() *settings {
	return &settings{}
}

// AssetFilter sets a filter to select the proper released asset
// from releases with multiple artifacts. The filter is regexp string.
func AssetFilter(filter string) Option {
	return func(s *settings) {
		rex, err := regexp.Compile(filter)
		if err != nil {
			panic(fmt.Sprintf("invalid regexp passed as option: %v", err))
		}
		s.filters = append(s.filters, rex)
	}
}
