package rules

import "testing"

func TestParseDuration(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    Duration
		wantErr bool
	}{
		{"single token", args{"1m"}, Duration(60), false},
		{"double token", args{"1m1s"}, Duration(61), false},
		{"decimal", args{"1.5m"}, Duration(90), false},
		{"double decimal", args{"1.5w1.5h"}, Duration((7 * 1.5 * 86400) + (3600 * 1.5)), false},
		{"maximal", args{"1y1mo1w1d1m1s"}, Duration(1 + 60 + 86400 + (86400 * 7) + (86400 * 7 * 30) + (86400 * 7 * 30 * 12)), false},

		{"bad unit", args{"1x"}, Duration(0), true},
		{"bad number", args{"1..3m"}, Duration(0), true},
		{"bad negative number", args{"-1.3m"}, Duration(0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_String(t *testing.T) {
	cases := [...]string{
		"1m",
		"1h",
		"1m30s",
		"1y3mo1h10s",
	}

	for _, c := range cases {
		d := MustParseDuration(c)
		if d.String() != c {
			t.Errorf("d.String() == %v, wanted %v", d.String(), c)
		}
	}
}

func BenchmarkParseDuration(b *testing.B) {
	bench := func(dur string) {
		b.Run(dur, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				ParseDuration(dur)
			}
		})
	}

	bench("1m")
	bench("1m1s")
	bench("1h1m5s")
}
