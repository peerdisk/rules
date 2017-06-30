package rules

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"math"

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

// String returns the string representation of the rule.
func (r Rule) String() string {
	return fmt.Sprintf("interval %v keep %v", r.Interval, r.Keep)
}

// RuleSet contains a series of rules
type RuleSet []Rule

// Validate validates the behaviour of a rule set.
func (rs RuleSet) Validate() error {
	lowest := Duration(math.MaxUint64)
	for _, rule := range rs {
		if rule.Keep < rule.Interval {
			return errors.Errorf("%q has a keep lower than it's interval", rule.String())
		}
		if rule.Interval < lowest {
			lowest = rule.Interval
		}
	}

	for _, rule := range rs {
		if rule.Interval%lowest != 0 {
			return errors.Errorf("%v doesn't divide into lowest interval %v", rule.Interval, lowest)
		}
	}

	return nil
}

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
