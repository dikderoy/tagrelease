package tagrelease

import "testing"

func TestStrategy(t *testing.T) {
	proto := Version{}

	variants := map[string]*Version{
		StrategyMajor: {
			major: 1,
			minor: 0,
			patch: 0,
		},
		StrategyMinor: {
			major: 0,
			minor: 1,
			patch: 0,
		},
		StrategyPatch: {
			major: 0,
			minor: 0,
			patch: 1,
		},
		StrategyUpstream: {
			major: 0,
			minor: 0,
			patch: 1,
		},
		StrategyNever: {},
	}

	for strategy, expected := range variants {
		t.Run(strategy, func(t *testing.T) {
			s, err := StrategyFactory(strategy)
			if err != nil {
				t.Fail()
			}
			monkey := proto
			s(&monkey)
			if *expected != monkey {
				t.Fail()
				t.Logf("expected: %v received: %v", expected, monkey)
			}
		})
	}

}

func TestStrategyWithEmptyPatch(t *testing.T) {
	proto := Version{
		patch: -1,
	}

	variants := map[string]*Version{
		StrategyMajor: {
			major: 1,
			minor: 0,
			patch: 0,
		},
		StrategyMinor: {
			major: 0,
			minor: 1,
			patch: 0,
		},
		StrategyPatch: {
			major: 0,
			minor: 0,
			patch: 1,
		},
		StrategyUpstream: {
			major: 0,
			minor: 1,
			patch: 0,
		},
		StrategyNever: {},
	}

	for strategy, expected := range variants {
		t.Run(strategy, func(t *testing.T) {
			s, err := StrategyFactory(strategy)
			if err != nil {
				t.Fail()
			}
			monkey := proto
			s(&monkey)
			if *expected != monkey {
				t.Fail()
				t.Logf("expected: %v received: %v", expected, monkey)
			}
		})
	}

}
