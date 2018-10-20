package tagrelease

import "github.com/sirupsen/logrus"

type Adapter interface {
	Version() *Version
	Revision() (string, error)
	Branch() (string, error)
}

type Version struct {
	major  int
	minor  int
	patch  int
	diff   int
	rev    string
	suffix string
}

type Converter struct {
	adapter  Adapter
	strategy Strategy
}

func NewConverter(adapter Adapter, strategy Strategy) *Converter {
	return &Converter{
		adapter:  adapter,
		strategy: strategy,
	}
}

var empty = Version{}

func (c *Converter) Detect() (v *Version) {
	v = c.adapter.Version()
	if *v == empty {
		v.minor = 1
		return
	}
	//use increment strategy
	c.strategy(v)
	return
}

func among(elem string, stack []string) bool {
	for i := range stack {
		if stack[i] == elem {
			return true
		}
	}
	return false
}

func (c *Converter) ReleaseKind() string {
	branch, _ := c.adapter.Branch()
	switch {
	case among(branch, GlobalConfig.Branches.Master):
		return "rc"
	case among(branch, GlobalConfig.Branches.Trunk):
		return "b"
	default:
		return "a"
	}
}

func (c *Converter) Revision() string {
	r, err := c.adapter.Revision()
	if err != nil {
		logrus.WithError(err).Debug("failed to detect revision")
	}
	return r
}
