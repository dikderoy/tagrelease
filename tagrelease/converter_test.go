package tagrelease

import "testing"

func prepareConverterWithMockAdapter() Converter {
	mock := mockAdapter{
		version:  &Version{},
		branch:   "master",
		revision: "57a182a57a182a",
	}
	strategy, _ := StrategyFactory(StrategyNever)
	c := Converter{
		adapter:  &mock,
		strategy: strategy,
	}
	return c
}

func TestConverter_Detect(t *testing.T) {
	c := prepareConverterWithMockAdapter()

	v := Version{
		Minor: 1,
	}

	t.Run("empty", func(t *testing.T) {
		r := c.Detect()
		if v != *r {
			t.Fail()
			t.Logf("expected: %v received: %v", v, r)
		}
	})
	t.Run("not-empty", func(t *testing.T) {
		c.adapter.(*mockAdapter).version.Minor = 1
		r := c.Detect()
		if v != *r {
			t.Fail()
			t.Logf("expected: %v received: %v", empty, r)
		}
	})
}

func TestConverter_ReleaseKind(t *testing.T) {
	GlobalConfig.Branches.Master = []string{"master"}
	GlobalConfig.Branches.Trunk = []string{"trunk"}
	c := prepareConverterWithMockAdapter()
	mock := c.adapter.(*mockAdapter)
	t.Run("release-candidate", func(t *testing.T) {
		if c.ReleaseKind() != "rc" {
			t.Fail()
		}
	})
	t.Run("beta", func(t *testing.T) {
		mock.branch = "trunk"
		if c.ReleaseKind() != "b" {
			t.Fail()
		}
	})
	t.Run("alpha", func(t *testing.T) {
		mock.branch = "hot-fix"
		if c.ReleaseKind() != "a" {
			t.Fail()
		}
	})
}

func TestConverter_Revision(t *testing.T) {
	c := prepareConverterWithMockAdapter()
	if c.Revision() != "57a182a57a182a" {
		t.Fail()
	}
}
