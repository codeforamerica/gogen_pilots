package data

import (
	"regexp"
	"sort"
	"strings"
	"time"
)

type DOJHistory struct {
	SubjectID         string
	Name              string
	WeakName          string
	CII               string
	DOB               time.Time
	SSN               string
	CDL               string
	PC290Registration bool
	Convictions       []*DOJRow
	seenConvictions   map[string]bool
	OriginalCII       string
}

func (history *DOJHistory) PushRow(row DOJRow) {
	if history.SubjectID == "" {
		history.SubjectID = row.SubjectID
		history.Name = row.Name
		history.WeakName = strings.Split(row.Name, " ")[0]
		history.CII = formatCII(row.CII)
		history.OriginalCII = row.CII
		history.DOB = row.DOB
		history.SSN = row.SSN
		history.CDL = row.CDL
		history.seenConvictions = make(map[string]bool)
	}

	if row.Convicted && !history.seenConvictions[row.CountOrder] {
		history.Convictions = append(history.Convictions, &row)
		history.seenConvictions[row.CountOrder] = true
	}

	if row.PC290Registration {
		history.PC290Registration = true
	}
}

func (history *DOJHistory) Match(entry CMSEntry) MatchData {
	var results = make(map[string]bool)

	results["cii"] = entry.CII != "" && entry.CII == history.CII
	results["ssn"] = entry.SSN != "" && entry.SSN == history.SSN
	results["cdl"] = entry.CDL != "" && entry.CDL == history.CDL

	if entry.CourtNumber != "" {
		matched := false
		for _, row := range history.Convictions {
			if row.County == "SAN FRANCISCO" && row.MatchingCourtNumber(entry.CourtNumber) {
				matched = true
				break
			}
		}
		results["courtno"] = matched
	}

	name := entry.FormattedName
	dateOfBirth := entry.DateOfBirth

	if (name != "" && dateOfBirth != time.Time{}) {
		results["nameAndDob"] = name == history.Name && dateOfBirth == history.DOB
		results["weakNameAndDob"] = entry.WeakName == history.WeakName && dateOfBirth == history.DOB
	}

	if (entry.WeakName != "" && entry.BookingDate != time.Time{}) {
		matched := false
		if entry.WeakName == history.WeakName {

			for _, row := range history.Convictions {
				if row.County == "SAN FRANCISCO" && row.CycleDate == entry.BookingDate {
					matched = true
					break
				}
			}
		}
		results["weakNameAndCycleDate"] = matched
	}

	matchStrength := 0
	for _, val := range results {
		if val {
			matchStrength++
		}
	}

	return MatchData{
		History:       history,
		MatchResults:  results,
		MatchStrength: matchStrength,
	}
}

func (history *DOJHistory) OnlyProp64Misdemeanors() bool {
	for _, row := range history.Convictions {
		matcher := regexp.MustCompile(`(11357|11358|11359|11360).*`)
		if !matcher.Match([]byte(row.CodeSection)) || row.Felony {
			return false
		}
	}
	return true
}

func (history *DOJHistory) PC290CodeSections() []string {
	var result []string

	for _, row := range history.Convictions {
		for _, pattern := range pc290Patterns {
			if pattern.MatchString(row.CodeSection) {
				result = append(result, row.CodeSection)
			}
		}
	}
	return result
}

func (history *DOJHistory) SuperstrikesCodeSections() []string {
	var result []string
	for _, row := range history.Convictions {
		for _, pattern := range superstrikesPatterns {
			if pattern == row.CodeSection {
				result = append(result, row.CodeSection)
			}
		}
	}
	return result
}

func (history *DOJHistory) ThreeConvictionsSameCode(codeSection string) bool {
	matchingCycles := make(map[time.Time]bool)
	for _, row := range history.Convictions {
		if codeSection == strings.Replace(row.CodeSection, " ", "", -1) {
			matchingCycles[row.CycleDate] = true
		}
	}
	return len(matchingCycles) > 2
}

func (history *DOJHistory) MostRecentConvictionDate() time.Time {
	if len(history.Convictions) == 0 {
		return time.Time{}
	}
	convictions := history.Convictions
	sort.Slice(convictions, func(i, j int) bool {
		return convictions[i].DispositionDate.Before(convictions[j].DispositionDate)
	})
	return convictions[len(convictions)-1].DispositionDate
}

var pc290Patterns = []*regexp.Regexp{
	regexp.MustCompile(`290(.*) PC`),
	regexp.MustCompile(`236\.1\([BC]\)(.*) PC`),
	regexp.MustCompile(`243\.4(.*) PC`),
	regexp.MustCompile(`261(.*) PC`),
	regexp.MustCompile(`262\(A\)\(1\) PC`),
	regexp.MustCompile(`264\.1(.*) PC`),
	regexp.MustCompile(`266 PC`),
	regexp.MustCompile(`266C PC`),
	regexp.MustCompile(`266H\(B\)(.*) PC`),
	regexp.MustCompile(`266I\(B\)(.*) PC`),
	regexp.MustCompile(`266J(.*) PC`),
	regexp.MustCompile(`267 PC`),
	regexp.MustCompile(`269(.*) PC`),
	regexp.MustCompile(`285 PC`),
	regexp.MustCompile(`286([^\.]*) PC`),
	regexp.MustCompile(`288([^\.]*) PC`),
	regexp.MustCompile(`288A(.*) PC`),
	regexp.MustCompile(`288\.[23457](.*) PC`),
	regexp.MustCompile(`289([^\.]*) PC`),
	regexp.MustCompile(`311\.1(.*) PC`),
	regexp.MustCompile(`311\.2\([BCD]\) PC`),
	regexp.MustCompile(`311\.([34]|10|11)(.*) PC`),
	regexp.MustCompile(`314(.*) PC`),
	regexp.MustCompile(`647\.6(.*) PC`),
	regexp.MustCompile(`647A(.*) PC`),
	regexp.MustCompile(`653F\(C\) PC`),
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
