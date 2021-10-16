package config

import (
	"io"
	"os"
	"regexp"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open config file")
	}

	configBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	var config struct {
		ProtectedRepoRegExps map[string][]string `json:"protectedReposRegExps"`
	}

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse config file")
	}

	protectedRepoRegExps := map[string][]*regexp.Regexp{}

	for remote, regexps := range config.ProtectedRepoRegExps {
		resultRegexps := make([]*regexp.Regexp, 0, len(regexps))
		for _, regexpStr := range regexps {
			reg, err2 := regexp.Compile(regexpStr)
			if err2 != nil {
				return nil, errors.Wrapf(err2, "failed to parse regexp for remote %s", remote)
			}

			resultRegexps = append(resultRegexps, reg)
		}
		protectedRepoRegExps[remote] = resultRegexps
	}

	return &Config{
		ProtectedRepoRegExps: protectedRepoRegExps,
	}, nil
}
