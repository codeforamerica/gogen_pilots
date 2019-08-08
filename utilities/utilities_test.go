package utilities_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"gogen/utilities"
)

var _ = Describe("Utilities", func() {
	Describe("AddMaps", func() {
		It("adds the value of each key from both maps and returns a new map", func() {
			map1 := map[string]int{"foo": 1, "bar": 4}
			map2 := map[string]int{"baz": 2, "bar": 5}

			result := utilities.AddMaps(map1, map2)

			Expect(result).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo": Equal(1),
				"bar": Equal(9),
				"baz": Equal(2),
			}))
		})

		It("can handle nil maps", func() {
			var map1 map[string]int
			map2 := map[string]int{"baz": 2, "bar": 5}

			result := utilities.AddMaps(map1, map2)

			Expect(result).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bar": Equal(5),
				"baz": Equal(2),
			}))
		})
	})
})
