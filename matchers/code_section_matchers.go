package matchers

import (
	"regexp"
)

var prop64matcher = regexp.MustCompile(`(11357|11358|11359|11360)`)
var relatedChargeMatcher = regexp.MustCompile(`(647\(f\)\s*PC|602\s*PC|466\s*PC|148\.9\s*PC|148\s*PC|11364\s*HS|11550\s*HS|4140\s*BP|4149\s*BP|4060\s*BP|40508\s*VC|1320[^\d\.][^\.]*PC)`)
var Prop64MatchersByCodeSection = map[string]*regexp.Regexp{
	"11357":                 regexp.MustCompile(`11357.*`),
	"11358":                 regexp.MustCompile(`11358.*`),
	"11359":                 regexp.MustCompile(`11359.*`),
	"11360":                 regexp.MustCompile(`11360.*`),
}

func ExtractProp64Section(codeSection string) (bool, string) {
	if IsProp64Charge(codeSection) {
		return true, prop64matcher.FindStringSubmatch(codeSection)[1]
	} else {
		return false, ""
	}
}

func ExtractRelatedChargeSection(codeSection string) (bool, string) {
	if IsRelatedCharge(codeSection) {
		return true, relatedChargeMatcher.FindStringSubmatch(codeSection)[1]
	} else {
		return false, ""
	}
}

func IsProp64Charge(codeSection string) bool {
	return prop64matcher.Match([]byte(codeSection))
}

func IsRelatedCharge(codeSection string) bool {
	return relatedChargeMatcher.Match([]byte(codeSection))
}
