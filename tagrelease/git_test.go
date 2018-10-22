package tagrelease

import (
	"regexp"
	"testing"
)

func TestGitAdapter_Version(t *testing.T) {
	a := GitAdapter{}

	variants := map[string]*Version{
		"":                  {},
		"0.8.1-84-g57a182a": {0, 8, 1, 84, "g57a182a", ""},
		"0.0.0":             {0, 0, 0, 0, "", ""},
		"0.0":               {0, 0, -1, 0, "", ""},
		"1.0-1-g57a182a":    {1, 0, -1, 1, "g57a182a", ""},
		"1.0.0-0-g57a182a":  {1, 0, 0, 0, "g57a182a", ""},
		"v1.0.0-0-g57a182a": {1, 0, 0, 0, "g57a182a", ""},
	}

	for k, expected := range variants {
		t.Run(k, func(t *testing.T) {
			received := a.evaluate(k)
			if *received != *expected {
				t.Fail()
				t.Logf("expected: %v received: %v", expected, received)
			}
		})
	}

}

func TestGitAdapter_Branch(t *testing.T) {
	a := GitAdapter{}
	v, _ := a.Branch()
	t.Log(v)
}

var refRe = regexp.MustCompile(`[a-f0-9]+`)

func TestGitAdapter_Revision(t *testing.T) {
	a := GitAdapter{}
	v, _ := a.Revision()
	t.Log(v)
	if !refRe.MatchString(v) {
		t.Fail()
	}
}
