package rules

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestParseRule(t *testing.T) {
	type args struct {
		rule string
	}
	tests := []struct {
		name    string
		args    args
		want    Rule
		wantErr bool
	}{
		{"simple", args{"interval 6h keep 1w"}, Rule{Interval: MustParseDuration("6h"), Keep: MustParseDuration("1w")}, false},
		{"ordering", args{"keep 3w interval 3h"}, Rule{Interval: MustParseDuration("3h"), Keep: MustParseDuration("3w")}, false},
		{"bad duration", args{"keep 3wt interval 3h"}, Rule{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRule(tt.args.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    RuleSet
		wantErr bool
	}{
		{
			"simple",
			args{strings.NewReader("interval 6h keep 7w\ninterval 10h keep 10w")},
			RuleSet{Rule{Interval: MustParseDuration("6h"), Keep: MustParseDuration("7w")}, Rule{Interval: MustParseDuration("10h"), Keep: MustParseDuration("10w")}},
			false,
		},
		{
			"comment",
			args{strings.NewReader("# hello\ninterval 6h keep 7w\ninterval 10h keep 10w")},
			RuleSet{Rule{Interval: MustParseDuration("6h"), Keep: MustParseDuration("7w")}, Rule{Interval: MustParseDuration("10h"), Keep: MustParseDuration("10w")}},
			false,
		},
		{
			"empty lines",
			args{strings.NewReader("\n\n\ninterval 6h keep 7w\ninterval 10h keep 10w")},
			RuleSet{Rule{Interval: MustParseDuration("6h"), Keep: MustParseDuration("7w")}, Rule{Interval: MustParseDuration("10h"), Keep: MustParseDuration("10w")}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
