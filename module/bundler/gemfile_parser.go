package bundler

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type gem struct {
	Name    string
	Git     string
	Version string
	Tag     string
	Groups  string
	Require string
}
type gemfile struct {
	Source string
	Ruby   string
	Gems   []*gem
}

var (
	// Gem attributes
	REGEX_GEM            = regexp.MustCompile(`gem\s*['|"](.*?)['|"]`)
	REGEX_VERSION        = regexp.MustCompile(`,\s*["|'](.*?)['|"]`)
	REGEX_GIT            = regexp.MustCompile(`git:\s*["|'](.*?)['|"]`)
	REGEX_TAG            = regexp.MustCompile(`tag:\s*["|'](.*?)['|"]`)
	REGEX_REQUIRE_STRING = regexp.MustCompile(`require:\s*["|'](.*?)['|"]`)
	REGEX_REQUIRE_BOOL   = regexp.MustCompile(`require:\s*(false|true)`)
	// Gemfile header
	REGEX_SOURCE = regexp.MustCompile(`source\s*['|"](.*?)['|"]`)
	REGEX_RUBY   = regexp.MustCompile(`ruby\s*['|"](.*?)['|"]`)
	// Gem Groups
	REGEX_GROUP     = regexp.MustCompile(`group\s(.*?)\sdo`)
	REGEX_END_GROUP = regexp.MustCompile(`^end$`)
)

func (gf *gemfile) Parse(file io.Reader) {
	scanner := bufio.NewScanner(file)
	gf.parseGemfile(scanner, "")
}

func (gf *gemfile) parseGemfile(scanner *bufio.Scanner, groups string) {
	for scanner.Scan() {
		line := scanner.Text()
		gem := gem{}

		if REGEX_RUBY.MatchString(line) {
			gf.Ruby = REGEX_RUBY.FindStringSubmatch(line)[1]
		}
		if REGEX_SOURCE.MatchString(line) {
			gf.Source = REGEX_SOURCE.FindStringSubmatch(line)[1]
		}
		if REGEX_GEM.MatchString(line) {
			gem.Name = REGEX_GEM.FindStringSubmatch(line)[1]
		}
		if REGEX_VERSION.MatchString(line) {
			gem.Version = REGEX_VERSION.FindStringSubmatch(line)[1]
		}
		if REGEX_GIT.MatchString(line) {
			gem.Git = REGEX_GIT.FindStringSubmatch(line)[1]
		}
		if REGEX_TAG.MatchString(line) {
			gem.Tag = REGEX_TAG.FindStringSubmatch(line)[1]
		}
		if REGEX_GROUP.MatchString(line) {
			gf.parseGemfile(scanner, REGEX_GROUP.FindStringSubmatch(line)[1])
		}
		if REGEX_REQUIRE_STRING.MatchString(line) {
			gem.Require = fmt.Sprintf(`"%s"`, REGEX_REQUIRE_STRING.FindStringSubmatch(line)[1])
		}
		if REGEX_REQUIRE_BOOL.MatchString(line) {
			gem.Require = REGEX_REQUIRE_BOOL.FindStringSubmatch(line)[1]
		}
		if REGEX_END_GROUP.MatchString(line) && groups != "" {
			gf.parseGemfile(scanner, "")
		}
		if gem.Name != "" {
			gem.Groups = groups
			gf.Gems = append(gf.Gems, &gem)
		}
	}
}
