package config

import "regexp"

type Config struct {
	ProtectedRepoRegExps map[string][]*regexp.Regexp
}
