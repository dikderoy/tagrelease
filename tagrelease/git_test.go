package tagrelease

import "testing"

func TestGitAdapter_Version(t *testing.T) {
	a := GitAdapter{}
	v := a.Version()

	t.Log(v)
}

func TestGitAdapter_Branch(t *testing.T) {
	a := GitAdapter{}
	v, _ := a.Branch()
	t.Log(v)
}

func TestGitAdapter_Revision(t *testing.T) {
	a := GitAdapter{}
	v, _ := a.Revision()
	t.Log(v)
}
