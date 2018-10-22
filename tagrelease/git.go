package tagrelease

import (
	"github.com/sirupsen/logrus"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type GitAdapter struct {
}

func (git *GitAdapter) Describe() (o string) {
	bo, err := exec.Command("git", "describe", "--tags").Output()
	if err != nil {
		logrus.WithError(err).Debug("no output from describe")
		o = "0.0.0"
	} else {
		o = string(bo)
	}
	logrus.WithField("out", o).Debug("git describe output")
	return o
}

var reGitDescription = regexp.MustCompile(
	`^(v?(?P<major>\d+).(?P<minor>\d+)(.(?P<patch>\d+))?)(-(?P<diff>\d+)-(?P<ref>[0-9A-Za-z]+))?$`,
)

func (git *GitAdapter) evaluate(desc string) *Version {
	matches := reGitDescription.FindAllStringSubmatch(desc, -1)
	logrus.WithField("matches", matches).Debug("matches")
	if matches == nil {
		return &Version{}
	}

	var major, minor, patch, diff int
	major, _ = strconv.Atoi(matches[0][2])
	minor, _ = strconv.Atoi(matches[0][3])
	if matches[0][4] != "" { // we need to know if patch is not defined
		patch, _ = strconv.Atoi(matches[0][5])
	} else {
		patch = -1
	}
	diff, _ = strconv.Atoi(matches[0][7]) //will return 0 on error (we okay with that)

	return &Version{
		major:  major,
		minor:  minor,
		patch:  patch,
		diff:   diff,
		rev:    matches[0][8], //empty value is "" which is exactly what we need
		suffix: "",
	}
}

func (git *GitAdapter) Version() *Version {
	return git.evaluate(git.Describe())
}

func (git *GitAdapter) Revision() (o string, err error) {
	bo, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bo)), nil
}

func (git *GitAdapter) Branch() (o string, err error) {
	bo, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bo)), nil
}
