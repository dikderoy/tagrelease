package tagrelease

import "testing"

type mockAdapter struct {
	version  *Version
	branch   string
	revision string
}

func (m *mockAdapter) Version() *Version {
	return m.version
}

func (m *mockAdapter) Revision() (string, error) {
	return m.revision, nil
}

func (m *mockAdapter) Branch() (string, error) {
	return m.branch, nil
}

func TestFormatter(t *testing.T) {
	GlobalConfig.Branches.Master = []string{"master"}
	strategy, _ := StrategyFactory(StrategyNever)
	f := NewFormatter(&Converter{
		adapter: &mockAdapter{
			version: &Version{
				major: 1,
				minor: 2,
				patch: 3,
				diff:  4,
				rev:   "g57a182a",
			},
			branch:   "master",
			revision: "57a182a57a182a",
		},
		strategy: strategy,
	})

	variants := map[string]string{
		FormatRelease:  "1.2.3rc4+57a182a",
		FormatPEP440:   "1.2.3rc4+57a182a",
		FormatSemver:   "1.2.3",
		FormatShort:    "1.2",
		FormatMajor:    "1",
		FormatMinor:    "2",
		FormatPatch:    "3",
		FormatRevision: "57a182a57a182a",
		FormatRevShort: "57a182a",
	}

	for k, expected := range variants {
		t.Run(k, func(t *testing.T) {
			format := FormatFactory(f, k)
			received := format()
			if received != expected {
				t.Fail()
				t.Logf("expected: %v received: %v", expected, received)
			}
		})
	}
}
