package processor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gogen/processor"

	. "gogen/data"

	"io/ioutil"
	path "path/filepath"
)

var _ = Describe("csvWriter", func() {
	var writer CMSWriter
	var entry CMSEntry
	var outputDir string
	var err error
	var info EligibilityInfo

	BeforeEach(func() {
		entry = CMSEntry{
			RawRow: []string{"A", "Slice", "Of", "Strings"},
		}
		info = EligibilityInfo{
			Over1Lb:   "a eligibility value",
			QFinalSum: "999.9",
		}

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())
		writer = NewCMSWriter(path.Join(outputDir, "written_csv.csv"))
	})

	It("Writes CMS Entries to a csv file", func() {
		writer.WriteEntry(entry, info)
		writer.Flush()

		expectedResultsBody := `Court Number,Ind,Incident Number,Truename,Case Level,Case Dispo,Case Disposition Description ,Dispo Date,Action Number,1st Filed,Charge Level,Charge Date,Current Charge,Current Level,Current Charge Description,Chg Disp,Charge Disposition Description,Ch Dispo Date,Race,Sex,DOB,SFNO,CII,FBI,SSN,DL Number,EOR,PRI_NAME,PRI_DOB,SUBJECT_ID,CII_NUMBER,PRI_SSN,Superstrikes,Superstrike Code Section(s),PC290 Charges,PC290 Code Section(s),PC290 Registration,Two Priors,Over 1lb,Q_final_sum,Age at Conviction,Years Since Event,Years Since Most Recent Conviction,Final Recommendation
A,Slice,Of,Strings,a eligibility value,999.9
`

		pathToOutput, err := path.Abs(path.Join(outputDir, "written_csv.csv"))
		Expect(err).ToNot(HaveOccurred())
		outputBody, err := ioutil.ReadFile(pathToOutput)
		Expect(err).ToNot(HaveOccurred())

		Expect(string(outputBody)).To(Equal(string(expectedResultsBody)))
	})
})
