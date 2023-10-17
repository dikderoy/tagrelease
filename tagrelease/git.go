package tagrelease

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type GitAdapter struct {
}

// Describe executes and returns output from `git describe --tags`
func (git *GitAdapter) Describe() (o string) {
	bo, err := exec.Command("git", "describe", "--tags").Output()
	if err != nil {
		log.WithError(err).Debug("no output from describe")
		o = "0.0.0"
	} else {
		o = strings.TrimSpace(string(bo))
	}
	log.WithField("out", o).Debug("git describe output")
	return o
}

var reGitDescription = regexp.MustCompile(
	`^(v?(?P<major>\d+).(?P<minor>\d+)(.(?P<patch>\d+))?)(-(?P<diff>\d+)-(?P<ref>[0-9A-Za-z]+))?$`,
)

func (git *GitAdapter) evaluate(desc string) *Version {
	matches := reGitDescription.FindAllStringSubmatch(desc, -1)
	log.WithField("matches", matches).Debug("matches")
	if matches == nil {
		return &Version{}
	}

	var major, minor, patch, diff int
	major, _ = strconv.Atoi(matches[0][2])
	minor, _ = strconv.Atoi(matches[0][3])
	if matches[0][4] != "" { // we need to know if a patch is not defined
		patch, _ = strconv.Atoi(matches[0][5])
	} else {
		patch = -1
	}
	diff, _ = strconv.Atoi(matches[0][7]) // will return 0 on error (we are okay with that)

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
		Diff:  diff,
		Rev:   matches[0][8], // empty value is "" which is exactly what we need
	}
}

// Version calls Describe and then evaluate results to populate Version struct
func (git *GitAdapter) Version() *Version {
	return git.evaluate(git.Describe())
}

// Revision executes and returns output from `git rev-parse HEAD`
func (git *GitAdapter) Revision() (o string, err error) {
	bo, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}
	o = string(bo)
	log.WithField("out", o).Debug("git rev-parse HEAD output")
	return strings.TrimSpace(o), nil
}

// Branch executes and returns output from `git rev-parse --abbrev-ref HEAD`
func (git *GitAdapter) Branch() (o string, err error) {
	bo, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	o = string(bo)
	log.WithField("out", o).Debug("git rev-parse --abbrev-ref HEAD output")
	return strings.TrimSpace(o), nil
}
