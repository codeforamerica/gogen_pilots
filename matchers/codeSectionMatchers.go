package matchers

import (

	"regexp"
)

func Prop64Matcher(codeSection string) (bool, string) {
	pat := regexp.MustCompile(`(11357|11358|11359|11360).*`)
	if pat.Match([]byte(codeSection)) {
		return true, pat.FindStringSubmatch(codeSection)[1]
	} else {
		return false, ""
	}
}

func RelatedChargeMatcher(codeSection string) (bool, string) {
	pat := regexp.MustCompile(`(647\(f\)\s*PC|602\s*PC|466\s*PC|148\.9\s*PC|148\s*PC|11364\s*HS|11550\s*HS|4140\s*BP|4149\s*BP|4060\s*BP|40508\s*VC|1320[^\d\.][^\.]*PC).*`)
	if pat.Match([]byte(codeSection)) {
		return true, pat.FindStringSubmatch(codeSection)[1]
	} else {
		return false, ""
	}
}