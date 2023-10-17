package tagrelease

import "errors"

type Strategy func(version *Version)

const (
	StrategyMajor    = "major"
	StrategyMinor    = "minor"
	StrategyPatch    = "patch"
	StrategyUpstream = "upstream"
	StrategyNever    = "never"
)

var StrategyList = []string{
	StrategyMajor,
	StrategyMinor,
	StrategyPatch,
	StrategyUpstream,
	StrategyNever,
}

// StrategyFactory creates and returns requested increment strategy function
// or error if such not defined
func StrategyFactory(name string) (Strategy, error) {
	switch name {
	case StrategyMajor:
		return func(version *Version) {
			version.Patch = 0
			version.Minor = 0
			version.Major++
		}, nil
	case StrategyMinor:
		return func(version *Version) {
			version.Patch = 0
			version.Minor++
		}, nil
	case StrategyPatch:
		return func(version *Version) {
			if version.Patch == -1 {
				version.Patch = 1
			} else {
				version.Patch++
			}
		}, nil
	case StrategyUpstream:
		// if patch present - increment patch,
		// if not - increment minor,
		return func(version *Version) {
			switch {
			case version.Patch == -1:
				version.Patch = 0
				version.Minor++
			case version.Patch > -1:
				version.Patch++
			}
		}, nil
	case StrategyNever:
		// do not alter semver in any way
		return func(version *Version) {
			if version.Patch == -1 {
				version.Patch = 0
			}
		}, nil
	default:
		return nil, errors.New("unknown strategy")
	}
}
