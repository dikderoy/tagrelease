package tagrelease

import "testing"

func TestStrategy(t *testing.T) {
	proto := Version{}

	variants := map[string]*Version{
		StrategyMajor: {
			Major: 1,
			Minor: 0,
			Patch: 0,
		},
		StrategyMinor: {
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
		StrategyPatch: {
			Major: 0,
			Minor: 0,
			Patch: 1,
		},
		StrategyUpstream: {
			Major: 0,
			Minor: 0,
			Patch: 1,
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
		Patch: -1,
	}

	variants := map[string]*Version{
		StrategyMajor: {
			Major: 1,
			Minor: 0,
			Patch: 0,
		},
		StrategyMinor: {
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
		StrategyPatch: {
			Major: 0,
			Minor: 0,
			Patch: 1,
		},
		StrategyUpstream: {
			Major: 0,
			Minor: 1,
			Patch: 0,
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
