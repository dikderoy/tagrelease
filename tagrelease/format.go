package tagrelease

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type Formatter struct {
	converter *Converter
	version   *Version
}

func NewFormatter(converter *Converter) *Formatter {
	return &Formatter{
		converter: converter,
		version:   converter.Detect(),
	}
}

// XYZ - a X.Y.Z (Major.Minor.Patch) formatted version according to SEMVER spec
func (f *Formatter) XYZ() string {
	return fmt.Sprintf("%d.%d.%d", f.version.Major, f.version.Minor, f.version.Patch)
}

// SemVer formats a version as extended identifier according to semver 2.0.0 spec
// a X.Y.Z-kind.diff+commit (ex.: 0.1.1-a.3+b8def90)
func (f *Formatter) SemVer() string {
	var kind, diff, release string
	if f.version.Diff > 0 {
		kind = "-" + f.ReleaseKind()
		diff = "." + strconv.Itoa(f.version.Diff)
	}
	if !GlobalConfig.Strategy.NoReleaseID {
		release = "+" + f.RevisionShort()
	}
	return fmt.Sprintf(
		"%s%s%s%s",
		f.XYZ(),
		kind,
		diff,
		release,
	)
}

// Short - an X.Y (Major.Minor) formatted version
func (f *Formatter) Short() string {
	return fmt.Sprintf("%d.%d", f.version.Major, f.version.Minor)
}

func (f *Formatter) Major() string {
	return strconv.Itoa(f.version.Major)
}
func (f *Formatter) Minor() string {
	return strconv.Itoa(f.version.Minor)
}
func (f *Formatter) Patch() string {
	return strconv.Itoa(f.version.Patch)
}
func (f *Formatter) Diff() string {
	return strconv.Itoa(f.version.Diff)
}

func (f *Formatter) ReleaseKind() string {
	return f.converter.ReleaseKind()
}

func (f *Formatter) Revision() string {
	return f.converter.Revision()
}

func (f *Formatter) RevisionShort() string {
	rev := f.converter.Revision()
	if len(rev) < 7 {
		return rev
	}
	return rev[:7]
}

// PEP440 - a PEP440 compatible release identifier
func (f *Formatter) PEP440() string {
	var kind, diff, release string
	if f.version.Diff > 0 {
		kind = f.ReleaseKind()
		diff = strconv.Itoa(f.version.Diff)
	}
	if !GlobalConfig.Strategy.NoReleaseID {
		release = "+" + f.RevisionShort()
	}
	return fmt.Sprintf(
		"%s%s%s%s",
		f.XYZ(),
		kind,
		diff,
		release,
	)
}

const (
	FormatRelease  = "release"
	FormatPEP440   = "pep440"
	FormatXYZ      = "xyz"
	FormatSemver   = "semver"
	FormatShort    = "short"
	FormatMajor    = StrategyMajor
	FormatMinor    = StrategyMinor
	FormatPatch    = StrategyPatch
	FormatRevision = "revision"
	FormatRevShort = "revshort"
)

var FormatList = []string{
	FormatRelease,
	FormatPEP440,
	FormatXYZ,
	FormatSemver,
	FormatShort,
	FormatMajor,
	FormatMinor,
	FormatPatch,
	FormatRevision,
	FormatRevShort,
}

func FormatFactory(fe *Formatter, format string) func() string {
	switch format {
	case FormatRelease:
		return fe.PEP440
	case FormatPEP440:
		return fe.PEP440
	case FormatSemver:
		return fe.SemVer
	case FormatXYZ:
		return fe.XYZ
	case FormatShort:
		return fe.Short
	case FormatMajor:
		return fe.Major
	case FormatMinor:
		return fe.Minor
	case FormatPatch:
		return fe.Patch
	case FormatRevision:
		return fe.Revision
	case FormatRevShort:
		return fe.RevisionShort
	default:
		// default to custom template (an unknown format is processed as custom template)
		return func() string {
			return FormatTemplate(fe, format)
		}
	}
}

func FormatTemplate(fe *Formatter, tplSource string) string {
	tpl, err := template.New("user-supplied").Parse(tplSource)
	if err != nil {
		log.WithError(err).Fatal("failed to parse template")
	}
	sb := strings.Builder{}
	err = tpl.Execute(&sb, fe)
	if err != nil {
		log.WithError(err).Fatal("failed to execute template")
	}
	return sb.String()
}

var sensitiveChars = []string{"/", "+", "-", "~", "*", "@", "#", "!", "^", "$", "%", "&", "(", ")"}

func EscapeSensitiveChars(out string, escChar string) string {
	var escapee []string
	for _, x := range sensitiveChars {
		escapee = append(escapee, x, escChar)
	}
	r := strings.NewReplacer(escapee...)
	return r.Replace(out)
}
