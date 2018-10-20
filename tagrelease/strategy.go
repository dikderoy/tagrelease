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

func StrategyFactory(name string) (Strategy, error) {
	switch name {
	case StrategyMajor:
		return func(version *Version) {
			version.patch = 0
			version.minor = 0
			version.major++
		}, nil
	case StrategyMinor:
		return func(version *Version) {
			version.patch = 0
			version.minor++
		}, nil
	case StrategyPatch:
		return func(version *Version) {
			if version.patch == -1 {
				version.patch = 1
			} else {
				version.patch++
			}
		}, nil
	case StrategyUpstream:
		//if patch present - increment patch,
		//if not - increment minor,
		return func(version *Version) {
			switch {
			case version.patch == -1:
				version.patch = 0
				version.minor++
			case version.patch > -1:
				version.patch++
			}
		}, nil
	case StrategyNever:
		//do not alter semver in any way
		return func(version *Version) {
			if version.patch == -1 {
				version.patch = 0
			}
		}, nil
	default:
		return nil, errors.New("unknown strategy")
	}
}
