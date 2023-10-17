package tagrelease

import log "github.com/sirupsen/logrus"

// Adapter interface represents a VCS or other source of information
// which provides initial version and revision data on which
// the rest of the app operates
type Adapter interface {
	Version() *Version
	Revision() (string, error)
	Branch() (string, error)
}

// Version holds normalized version and revision information
// required to generate targeted output
type Version struct {
	Major int
	Minor int
	Patch int
	Diff  int
	Rev   string
}

// Converter is a container type holding methods of accessing underlying Adapter
// and Strategy to retrieve and process Version
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

// Detect method calls underlying Adapter to retrieve normalized Version from Adapter
func (c *Converter) Detect() (v *Version) {
	v = c.adapter.Version()
	log.WithField("version", v).Debug("version detected")
	if *v == empty {
		log.Debug("empty version detected, use first release strategy")
		v.Minor = 1
		return
	}
	if v.Diff == 0 {
		log.Debug("release detected, skip application of increment strategy")
		return
	}
	// use increment strategy
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

// ReleaseKind queries underlying Adapter for branch info to produce release modifier
// (release-candidate/rc, beta/b or alpha/a
func (c *Converter) ReleaseKind() string {
	branch, _ := c.adapter.Branch()
	var kind string
	switch {
	case among(branch, GlobalConfig.Branches.Master):
		kind = "rc"
	case among(branch, GlobalConfig.Branches.Trunk):
		kind = "b"
	default:
		kind = "a"
	}
	log.WithField("kind", kind).Debug("calculated release kind")
	return kind
}

// Revision queries underlying Adapter for raw revision identifier and returns it unchanged,
// catching errors, and returning empty string if any.
func (c *Converter) Revision() string {
	r, err := c.adapter.Revision()
	if err != nil {
		log.WithError(err).Debug("failed to detect revision")
	}
	return r
}
