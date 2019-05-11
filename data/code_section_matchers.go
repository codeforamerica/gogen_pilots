package data

import "regexp"

func isSuperstrike(codeSection string) bool {
	for _, pattern := range superstrikesPatterns {
		if pattern == codeSection {
			return true
		}
	}
	return false
}

func isPC290(codeSection string) bool {
	for _, pattern := range pc290Patterns {
		if pattern.MatchString(codeSection) {
			return true
		}
	}
	return false
}

var superstrikesPatterns = []string{
	"187 PC",
	"191.5 PC",
	"187-664 PC",
	"191.5-664 PC",
	"209 PC",
	"220 PC",
	"245(D)(3) PC",
	"261(A)(2) PC",
	"261(A)(6) PC",
	"262(A)(2) PC",
	"262(A)(4) PC",
	"264.1 PC",
	"269 PC",
	"286(C)(1) PC",
	"286(C)(2)(A) PC",
	"286(C)(2)(B) PC",
	"286(C)(2)(C) PC",
	"286(C)(3) PC",
	"286(D)(1) PC",
	"286(D)(2) PC",
	"286(D)(3) PC",
	"288(A) PC",
	"288(B)(1) PC",
	"288(B)(2) PC",
	"288A(C)(1) PC",
	"288A(C)(2)(A) PC",
	"288A(C)(2)(B) PC",
	"288A(C)(2)(C) PC",
	"288A(D) PC",
	"288.5(A) PC",
	"289(A)(1)(A) PC",
	"289(A)(1)(B) PC",
	"289(A)(1)(C) PC",
	"289(A)(2)(C) PC",
	"289(J) PC",
	"653F PC",
	"11418(A)(1) PC",
}

var pc290Patterns = []*regexp.Regexp{
	regexp.MustCompile(`290 PC`),
	regexp.MustCompile(`290(\(|\.|[a-zA-Z])+(.*) PC`),
	regexp.MustCompile(`236\.1\([BC]\)(.*) PC`),
	regexp.MustCompile(`243\.4(.*) PC`),
	regexp.MustCompile(`261 PC`),
	regexp.MustCompile(`261(\(|\.|[a-zA-Z])+(.*) PC`),
	regexp.MustCompile(`262\(A\)\(1\) PC`),
	regexp.MustCompile(`264\.1(.*) PC`),
	regexp.MustCompile(`266 PC`),
	regexp.MustCompile(`266C PC`),
	regexp.MustCompile(`266H\(B\)(.*) PC`),
	regexp.MustCompile(`266I\(B\)(.*) PC`),
	regexp.MustCompile(`266J(.*) PC`),
	regexp.MustCompile(`267 PC`),
	regexp.MustCompile(`269 PC`),
	regexp.MustCompile(`269(\(|\.|[a-zA-Z])+(.*) PC`),
	regexp.MustCompile(`285 PC`),
	regexp.MustCompile(`286([^\.]*) PC`),
	regexp.MustCompile(`288 PC`),
	regexp.MustCompile(`288\(([^\.]*) PC`),
	regexp.MustCompile(`288A(.*) PC`),
	regexp.MustCompile(`288\.[23457](.*) PC`),
	regexp.MustCompile(`289([^\.]*) PC`),
	regexp.MustCompile(`311\.1(.*) PC`),
	regexp.MustCompile(`311\.2\([BCD]\) PC`),
	regexp.MustCompile(`311\.([34]|10|11)(.*) PC`),
	regexp.MustCompile(`314 PC`),
	regexp.MustCompile(`314(\(|\.|[a-zA-Z])+(.*) PC`),
	regexp.MustCompile(`647\.6(.*) PC`),
	regexp.MustCompile(`647A(.*) PC`),
	regexp.MustCompile(`653F\(C\) PC`),
}
