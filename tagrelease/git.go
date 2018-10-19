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
		//o = "0.0.0"
		o = "0.8.1-84-g57a182a"
	} else {
		o = string(bo)
	}
	logrus.WithField("out", o).Debug("git describe output")
	return o
}

var reGitDescription = regexp.MustCompile(
	`^((?P<major>\d+).(?P<minor>\d+)(.(?P<patch>\d+))?)(-(?P<diff>\d+)-(?P<ref>[0-9A-Za-z]+))?$`,
)

func (git *GitAdapter) Version() *Version {
	output := git.Describe()
	matches := reGitDescription.FindAllStringSubmatch(output, -1)
	logrus.WithField("matches", matches).Debug("matches")
	if matches == nil {
		return &Version{}
	}

	major, _ := strconv.Atoi(matches[0][2])
	minor, _ := strconv.Atoi(matches[0][3])
	patch, _ := strconv.Atoi(matches[0][5])
	diff, _ := strconv.Atoi(matches[0][7])

	return &Version{
		major:  major,
		minor:  minor,
		patch:  patch,
		diff:   diff,
		rev:    matches[0][8],
		suffix: "",
	}
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
