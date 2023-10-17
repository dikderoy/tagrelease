package tagrelease

import (
	"strings"
	"testing"
	"text/template"

	log "github.com/sirupsen/logrus"
)

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

func expected(data *vbr, tpl string) string {
	t, err := template.New("t").Parse(tpl)
	if err != nil {
		log.WithError(err).Error("parsing expectation template")
	}
	sb := strings.Builder{}
	err = t.Execute(&sb, data)
	if err != nil {
		log.WithError(err).Error("executing expectation template")
	}
	return sb.String()
}

type vbr struct {
	Ver    *Version // version variant
	Branch string   // branch variant
	Rev    string   // revision variant
	RKind  string   // expected release kind
}

func TestFormatter(t *testing.T) {
	GlobalConfig.Branches.Master = []string{"master"}
	GlobalConfig.Branches.Trunk = []string{"trunk"}

	strategy, _ := StrategyFactory(StrategyNever)

	versionVariants := map[string]*vbr{
		"1.2.3-4/57a182a/master": {
			Ver: &Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Diff:  4,
				Rev:   "g57a182a",
			},
			Branch: "master",
			RKind:  "rc",
			Rev:    "57a182a57a182a",
		},
		"1.2.3-4/57a182a/trunk": {
			Ver: &Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Diff:  4,
				Rev:   "g57a182a",
			},
			Branch: "trunk",
			RKind:  "b",
			Rev:    "57a182a57a182a",
		},
		"1.2.3-0/57a182a/trunk": {
			Ver: &Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Diff:  0,
				Rev:   "g57a182a",
			},
			Branch: "trunk",
			RKind:  "",
			Rev:    "57a182a57a182a",
		},
		"1.2-9/57a182a/hotfix": {
			Ver: &Version{
				Major: 1,
				Minor: 2,
				Patch: -1,
				Diff:  9,
				Rev:   "g57a182a",
			},
			Branch: "hotfix",
			RKind:  "a",
			Rev:    "57a182a57a182a",
		},
		"1.2/master/57a182a": {
			Ver: &Version{
				Major: 1,
				Minor: 2,
				Patch: -1,
				Diff:  0,
				Rev:   "g57a182a",
			},
			Branch: "master",
			Rev:    "57a182a57a182a",
		},
	}

	for vbrK, varData := range versionVariants {
		f := NewFormatter(&Converter{
			adapter: &mockAdapter{
				version:  varData.Ver,
				branch:   varData.Branch,
				revision: varData.Rev,
			},
			strategy: strategy,
		})

		pep440Expected := expected(varData,
			"{{.Ver.Major}}.{{.Ver.Minor}}.{{.Ver.Patch}}"+
				"{{if ne .RKind \"\" -}} {{.RKind}}{{.Ver.Diff}} {{- end -}}"+
				"+57a182a")

		semverExpected := expected(varData,
			"{{.Ver.Major}}.{{.Ver.Minor}}.{{.Ver.Patch}}"+
				"{{if ne .RKind \"\" -}} -{{.RKind}}.{{.Ver.Diff}} {{- end -}}"+
				"+57a182a")
		variants := map[string]string{
			FormatRelease:  pep440Expected,
			FormatPEP440:   pep440Expected,
			FormatSemver:   semverExpected,
			FormatShort:    expected(varData, "{{.Ver.Major}}.{{.Ver.Minor}}"),
			FormatMajor:    expected(varData, "{{.Ver.Major}}"),
			FormatMinor:    expected(varData, "{{.Ver.Minor}}"),
			FormatPatch:    expected(varData, "{{.Ver.Patch}}"),
			FormatRevision: varData.Rev,
			FormatRevShort: varData.Rev[:7],

			"{{.Major}}+{{.Diff}}": expected(varData, "{{.Ver.Major}}+{{.Ver.Diff}}"),
		}

		for k, expected := range variants {
			t.Run(vbrK+"|"+k, func(t *testing.T) {
				format := FormatFactory(f, k)
				received := format()
				if received != expected {
					t.Fail()
					t.Logf("expected: %v received: %v", expected, received)
				}
			})
		}
	}
}

func TestEscapeSensitiveChars(t *testing.T) {
	s := "1.2.3+abc"
	r := EscapeSensitiveChars(s, "_")
	if r != "1.2.3_abc" {
		t.Fail()
	}
}
