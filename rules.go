package rules

import (
	"bufio"
	"io"
	"strings"

	"github.com/pkg/errors"
)

// Rule contains a single rule
type Rule struct {
	Interval Duration
	Keep     Duration
}

// ParseRule parses a single rule line
func ParseRule(rule string) (Rule, error) {
	tokens := strings.Split(rule, " ")

	more := func(i int) bool {
		return len(tokens) > i+1
	}
	r := Rule{}
	var err error

	for i, token := range tokens {
		switch token {
		case "interval":
			if more(i) {
				r.Interval, err = ParseDuration(tokens[i+1])
				if err != nil {
					return r, errors.Wrap(err, "failed to parse interval")
				}
			}
		case "keep":
			if more(i) {
				r.Keep, err = ParseDuration(tokens[i+1])
				if err != nil {
					return r, errors.Wrap(err, "failed to parse keep")
				}
			}
		}
	}
	return r, nil
}

// RuleSet contains a series of rules
type RuleSet []Rule

// Parse parses a rule set.
// It just addresses the structure of the rule file, not the behaviour.
// Use RuleSet.Validate() to validate the behaviour of the rule set.
func Parse(r io.Reader) (RuleSet, error) {
	sc := bufio.NewScanner(r)
	rs := make(RuleSet, 0, 3)
	for sc.Scan() {
		txt := strings.TrimSpace(sc.Text())
		if len(txt) == 0 {
			continue
		}
		// ignore comments
		if txt[0] == '#' {
			continue
		}
		rule, err := ParseRule(txt)
		if err != nil {
			return rs, errors.Wrap(err, "failed to parse rule")
		}
		rs = append(rs, rule)
	}
	return rs, nil
}
