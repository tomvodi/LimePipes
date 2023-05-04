package bww

import (
	"banduslib/internal/common"
	"banduslib/internal/interfaces"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"os"
)

var _ = Describe("BWW Parser", func() {
	var err error
	var parser interfaces.BwwParser
	var bwwDoc *common.BwwDocument
	var data []byte

	BeforeEach(func() {
		parser = NewBwwParser()
	})

	When("parsing the file with all bww symbols in it", func() {
		BeforeEach(func() {
			var bwwFile *os.File
			bwwFile, err = os.Open("./testfiles/all_symbols.bww")
			Expect(err).ShouldNot(HaveOccurred())
			data, err = io.ReadAll(bwwFile)
			Expect(err).ShouldNot(HaveOccurred())
			bwwDoc, err = parser.ParseBwwData(data)
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwDoc.Tunes).To(HaveLen(2))
		})
	})
})
