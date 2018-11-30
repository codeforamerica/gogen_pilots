package data_test

import (
	"gogen/data"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("DOJHistory", func() {
	Describe("PushRow", func() {
		Context("An empty history", func() {
			var history data.DOJHistory

			BeforeEach(func() {
				history = data.DOJHistory{}
			})

			PIt("Sets the name of the history", func() {
				row := []string{}
				history.PushRow(row)
			})
		})
	})
})
