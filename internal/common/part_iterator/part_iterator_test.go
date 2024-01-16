package part_iterator_test

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/expander"
	"banduslib/internal/common/test"
	"banduslib/internal/interfaces"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"path/filepath"
	"strings"

	"banduslib/internal/common/part_iterator"
)

var embExpander = expander.NewEmbellishmentExpander()

func createYamlFromBww(fpath string, parser interfaces.BwwParser) {
	fileDir := filepath.Dir(fpath)
	filenameWithExt := filepath.Base(fpath)
	filenameWithoutExt := strings.TrimSuffix(filenameWithExt, filepath.Ext(filenameWithExt))
	yamlFileName := filenameWithoutExt + ".yaml"
	yamlFilePath := filepath.Join(fileDir, yamlFileName)

	bwwData := test.DataFromFile(fpath)
	musicTunesBww, err := parser.ParseBwwData(bwwData)
	Expect(err).ShouldNot(HaveOccurred())

	test.ExportToYaml(musicTunesBww, yamlFilePath)
}

var _ = Describe("PartIterator", func() {
	var tune *music_model.Tune
	var parts interfaces.MusicPartIterator
	//var parser interfaces.BwwParser
	var musicModel music_model.MusicModel

	BeforeEach(func() {
		//parser = bww.NewBwwParser(embExpander)
	})

	JustBeforeEach(func() {
		parts = part_iterator.New(tune)
	})

	Context("having a tune with two parts", func() {
		BeforeEach(func() {
			//createYamlFromBww("./testfiles/tune_with_two_parts.bww", parser)
			musicModel = test.ImportFromYaml("./testfiles/tune_with_two_parts.yaml", embExpander)
			Expect(musicModel).To(HaveLen(1))
			tune = musicModel[0]
		})

		It("should have two parts", func() {
			Expect(parts.Count()).To(Equal(2))
			Expect(parts.Nr(1)).ShouldNot(BeNil())
			Expect(parts.Nr(2)).ShouldNot(BeNil())
			Expect(parts.Nr(0)).Should(BeNil())
			Expect(parts.Nr(3)).Should(BeNil())
		})

		It("should have parts with the right amount of measures", func() {
			Expect(parts.Count()).To(Equal(2))
			Expect(parts.Nr(1).Measures).Should(HaveLen(2))
			Expect(parts.Nr(2).Measures).Should(HaveLen(2))
		})

		It("parts should have repeat property set", func() {
			Expect(parts.Nr(1).WithRepeat).To(BeTrue())
			Expect(parts.Nr(2).WithRepeat).To(BeFalse())
		})

		When("iterating over parts", func() {
			It("should have returned all parts", func() {
				var iteratedParts []*music_model.MusicPart
				for parts.HasNext() {
					iteratedParts = append(iteratedParts, parts.GetNext())
				}

				Expect(iteratedParts).To(HaveLen(2))
				Expect(iteratedParts[0]).ShouldNot(BeNil())
				Expect(iteratedParts[1]).ShouldNot(BeNil())
				Expect(iteratedParts[0].WithRepeat).To(BeTrue())
				Expect(iteratedParts[1].WithRepeat).To(BeFalse())
			})
		})
	})

	Context("having a tune with 1-2 time line in a single measure", func() {
		BeforeEach(func() {
			//createYamlFromBww("./testfiles/tune_with_normal_1_2_part.bww", parser)
			musicModel = test.ImportFromYaml("./testfiles/tune_with_normal_1_2_part.yaml", embExpander)
			Expect(musicModel).To(HaveLen(1))
			tune = musicModel[0]
		})

		It("should have one parts", func() {
			Expect(parts.Count()).To(Equal(1))
			Expect(parts.Nr(1)).ShouldNot(BeNil())
		})

		It("should have split up the measure with 1-2 time line", func() {
			Expect(parts.Count()).To(Equal(1))
			Expect(parts.Nr(1).Measures).Should(HaveLen(4))
		})
	})
})
