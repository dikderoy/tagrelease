package tagrelease

import (
	"fmt"
	"strconv"
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

//a X.Y.Z (Major.Minor.Patch) formatted version according to SEMVER spec
func (f *Formatter) SemVer() string {
	return fmt.Sprintf("%d.%d.%d", f.version.major, f.version.minor, f.version.patch)
}

//a X.Y (Major.Minor) formatted version
func (f *Formatter) Short() string {
	return fmt.Sprintf("%d.%d", f.version.major, f.version.minor)
}

func (f *Formatter) Major() string {
	return strconv.Itoa(f.version.major)
}
func (f *Formatter) Minor() string {
	return strconv.Itoa(f.version.minor)
}
func (f *Formatter) Patch() string {
	return strconv.Itoa(f.version.patch)
}

func (f *Formatter) ReleaseKind() string {
	return f.converter.ReleaseKind()
}

func (f *Formatter) Revision() string {
	return f.converter.Revision()
}

func (f *Formatter) RevisionShort() string {
	return f.converter.Revision()[:7]
}

//a PEP440 compatible release identifier
func (f *Formatter) PEP440() string {
	var kind, diff string
	if f.version.diff > 0 {
		kind = f.ReleaseKind()
		diff = strconv.Itoa(f.version.diff)
	} else {
		kind = ""
		diff = ""
	}
	return fmt.Sprintf(
		"%s%s%s+%s",
		f.SemVer(),
		kind,
		diff,
		f.RevisionShort(),
	)
}

const (
	FormatRelease  = "release"
	FormatPEP440   = "pep440"
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
		return nil
	}
}

//if __name__ == '__main__':
//    parser = argparse.ArgumentParser(
//        description='''
//        GitAppVersionFetcher
//        uses git info to create a version/release identifier of form: x.y.z[(a|b|rc)N+gitref]
//        '''
//    )
//    parser.add_argument('--debug', action='store_true')
//    parser.add_argument('-d', '--dir', type=str, dest='work_dir', help='pick git info from this dir')
//    parser.add_argument('-o', '--out', type=str, dest='output', default='version.info',
//                        help='where to output version identifier, use `-` for stdout')
//    parser.add_argument('-f', '--format', type=str, dest='format', default='release',
//                        choices=('semver', 'major', 'short', 'release', 'revision'),
//                        help='version format')
//    args = parser.parse_args()
//    if args.debug:
//        logging.basicConfig(stream=sys.stderr, level=logging.DEBUG)
//
//    if args.work_dir is not None:
//        logger.info('using cwd [%s]', args.work_dir)
//        os.chdir(args.work_dir)
//    forge = VersionForger()
//    v = getattr(forge, args.format)
//    if args.output == '-':
//        sys.stdout.write(v)
//    else:
//        logger.info('compiled file [%s] with [%s]', args.output, v)
//        with open(args.output, 'w') as fd:
//            fd.write(v)
