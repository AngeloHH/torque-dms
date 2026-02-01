package domain

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type RulePattern struct {
	Pattern string `yaml:"pattern"`
	Message string `yaml:"message"`
}

type FieldRule struct {
	MinLength int           `yaml:"min_length"`
	MaxLength int           `yaml:"max_length"`
	Blacklist []string      `yaml:"blacklist"`
	Rules     []RulePattern `yaml:"rules"`
}

var validationRules map[string]FieldRule

func LoadValidationRules(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read validation rules: %w", err)
	}

	err = yaml.Unmarshal(data, &validationRules)
	if err != nil {
		return fmt.Errorf("failed to parse validation rules: %w", err)
	}

	return nil
}

func Validate(field string, value string) error {
	rule, exists := validationRules[field]
	if !exists {
		return nil
	}

	if rule.MinLength > 0 && len(value) < rule.MinLength {
		return fmt.Errorf("%s must be at least %d characters", field, rule.MinLength)
	}

	if rule.MaxLength > 0 && len(value) > rule.MaxLength {
		return fmt.Errorf("%s must be at most %d characters", field, rule.MaxLength)
	}

	for _, pattern := range rule.Blacklist {
		regex := regexp.MustCompile(pattern)
		if regex.MatchString(value) {
			return fmt.Errorf("%s is not allowed", field)
		}
	}

	for _, r := range rule.Rules {
		regex := regexp.MustCompile(r.Pattern)
		if !regex.MatchString(value) {
			return fmt.Errorf("%s %s", field, r.Message)
		}
	}

	return nil
}