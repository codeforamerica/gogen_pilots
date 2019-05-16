package data

import (
	"regexp"
	"strings"
)

func IsSuperstrike(codeSection string) bool {
	for _, pattern := range superstrikesPatterns {
		if pattern.MatchString(codeSection) {
			return true
		}
	}
	return false
}

func IsPC290(codeSection string) bool {
	for _, pattern := range pc290Patterns {
		if pattern.MatchString(codeSection) {
			return true
		}
	}
	return false
}

var gangEnhancement = `186\.22\(B\)\(4\)`
var superstrikesWithGangEnhancement = []string{
	`136\.1`,
	`215`,
	`213\(A\)\(1\)\(A\)`,
	`246`,
	`519`,
	`12022\.55`,
}

var enhanceableOffenses = `((` + strings.Join(superstrikesWithGangEnhancement, `)|(`) + `))`
var superstrikesPatterns = []*regexp.Regexp{
	regexp.MustCompile(`37 PC`),
	regexp.MustCompile(`128 PC`),
	regexp.MustCompile(enhanceableOffenses + `.*` + gangEnhancement + ` PC`),
	regexp.MustCompile(gangEnhancement + `.*` + enhanceableOffenses + ` PC`),
	regexp.MustCompile(`187 PC`),
	regexp.MustCompile(`188 PC`),
	regexp.MustCompile(`189(\.[15])? PC`),
	regexp.MustCompile(`190(\.\d{1,2})?(\(.*\))? PC`),
	regexp.MustCompile(`191(\.5)? PC`),
	regexp.MustCompile(`205 PC`),
	regexp.MustCompile(`207 PC`),
	regexp.MustCompile(`209(\.5)? PC`),
	regexp.MustCompile(`217\.1 PC`),
	regexp.MustCompile(`218 PC`),
	regexp.MustCompile(`219 PC`),
	regexp.MustCompile(`220 PC`),
	regexp.MustCompile(`245\(D\)\(3\) PC`),
	regexp.MustCompile(`261(\(.*\))? PC`),
	regexp.MustCompile(`262(\(.*\))? PC`),
	regexp.MustCompile(`264\.1 PC`),
	regexp.MustCompile(`269 PC`),
	regexp.MustCompile(`273AB(\(.*\))? PC`),
	regexp.MustCompile(`286((\(C\)\([123]\)(\([ABC]\))?)|(\(D\)\([123]\)))? PC`),
	regexp.MustCompile(`287 PC`),
	regexp.MustCompile(`288((\(A\))|(\(B\)\([12]\)))? PC`),
	regexp.MustCompile(`288A((\(D\))|(\(C\)\(1\))|(\(C\)\(2\)\([ABC]\)))? PC`),
	regexp.MustCompile(`288\.5(\(A\))? PC`),
	regexp.MustCompile(`289((\(J\))|(\(A\)\(1\)\([ABC]\))|(\(A\)\(2\)\(C\)))? PC`),
	regexp.MustCompile(`451\.5 PC`),
	regexp.MustCompile(`653F PC`),
	regexp.MustCompile(`667\.(61|7|71) PC`),
	regexp.MustCompile(`4500 PC`),
	regexp.MustCompile(`11418((\(A\)\(1\))|(\(B\)\([12]\))) PC`),
	regexp.MustCompile(`12308 PC`),
	regexp.MustCompile(`12310 PC`),
	regexp.MustCompile(`18745 PC`),
	regexp.MustCompile(`18755 PC`),
	regexp.MustCompile(`1672\(A\) MV`),
}

var pc290Patterns = []*regexp.Regexp{
	regexp.MustCompile(`236\.1\([BC]\)(.*) PC`),
	regexp.MustCompile(`243\.4(.*) PC`),
	regexp.MustCompile(`261 PC`),
	regexp.MustCompile(`261(\(|\.|[a-zA-Z])+(.*) PC`),
	regexp.MustCompile(`262\(A\)(\(1\))? PC`),
	regexp.MustCompile(`264\.1(\((.*))? PC`),
	regexp.MustCompile(`266 PC`),
	regexp.MustCompile(`266C PC`),
	regexp.MustCompile(`266H\(B\)(.*) PC`),
	regexp.MustCompile(`266I\(B\)(.*) PC`),
	regexp.MustCompile(`266J PC`),
	regexp.MustCompile(`267 PC`),
	regexp.MustCompile(`269 PC`),
	regexp.MustCompile(`269(\(|\.|[a-zA-Z])+(.*) PC`),
	regexp.MustCompile(`272 PC`),
	regexp.MustCompile(`285 PC`),
	regexp.MustCompile(`287 PC`),
	regexp.MustCompile(`286([^\.]*) PC`),
	regexp.MustCompile(`288 PC`),
	regexp.MustCompile(`288A(.*) PC`),
	regexp.MustCompile(`288\.[23457](.*) PC`),
	regexp.MustCompile(`289([^\.]*) PC`),
	regexp.MustCompile(`311\.1(.*) PC`),
	regexp.MustCompile(`311\.2\([BCD]\) PC`),
	regexp.MustCompile(`311\.([34]|10|11)(.*) PC`),
	regexp.MustCompile(`314\([12]\) PC`),
	regexp.MustCompile(`451\.5 PC`),
	regexp.MustCompile(`647\.6(.*) PC`),
	regexp.MustCompile(`647A(.*) PC`),
	regexp.MustCompile(`653F\([BC]\) PC`),
}
