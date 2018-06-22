package slack

import "testing"

func TestHistoryMethods(t *testing.T) {
	for name, tc := range map[string]struct {
		ch   string
		want string
	}{
		"c": {"c584845", "channels"},
		"g": {"gsadfdf", "groups"},
		"i": {"isadfdf", "im"},
		"d": {"dsadfdf", "im"},
		"C": {"C584845", "channels"},
		"G": {"Gsadfdf", "groups"},
		"I": {"Isadfdf", "im"},
		"D": {"Dsadfdf", "im"},
	} {
		t.Run(name, func(t *testing.T) {
			want := tc.want
			have, err := kindof(tc.ch)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			if have != want {
				t.Fatalf("%q: have %q, want %q", tc.ch, have, want)
			}

		})
	}

}
