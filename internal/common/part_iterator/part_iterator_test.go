package part_iterator_test

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/expander"
	"banduslib/internal/common/music_model/symbols/time_line"
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
			Expect(parts.Nr(1).Measures).Should(HaveLen(3))
		})

		It("should have set the barlines for the split up measures correctly", func() {
			sourceMeasure := tune.Measures[0]
			splitMeasures := parts.Nr(1).Measures
			Expect(splitMeasures[0].LeftBarline).To(Equal(sourceMeasure.LeftBarline))
			Expect(splitMeasures[len(splitMeasures)-1].RightBarline).To(Equal(sourceMeasure.RightBarline))
		})

		It("should have added the symbols from the source measure", func() {
			sourceMeasure := tune.Measures[0]
			splitMeasures := parts.Nr(1).Measures
			Expect(splitMeasures[0].Symbols).NotTo(BeNil())
			Expect(splitMeasures[0].Symbols[0]).To(Equal(sourceMeasure.Symbols[0]))

			Expect(splitMeasures[1].Symbols).NotTo(BeNil())
			Expect(splitMeasures[1].Symbols[0]).To(Equal(sourceMeasure.Symbols[2]))

			Expect(splitMeasures[2].Symbols).NotTo(BeNil())
			Expect(splitMeasures[2].Symbols[0]).To(Equal(sourceMeasure.Symbols[5]))
		})
	})

	Context("having a tune with 2-2 time line", func() {
		BeforeEach(func() {
			//createYamlFromBww("./testfiles/tune_with_2_of_2_part.bww", parser)
			//createYamlFromBww("./testfiles/tune_with_2_of_2_part_normalized.bww", parser)
			musicModel = test.ImportFromYaml("./testfiles/tune_with_2_of_2_part.yaml", embExpander)
			Expect(musicModel).To(HaveLen(1))
			tune = musicModel[0]
		})

		It("should have two parts", func() {
			Expect(parts.Count()).To(Equal(2))
		})

		It("should have added the second time to part 2", func() {
			Expect(parts.Nr(2).Measures).Should(HaveLen(5))
			Expect(parts.Nr(1).Measures[2]).
				To(Equal(parts.Nr(2).Measures[3]))
		})

		It("should have modified the time line type to 'second'", func() {
			Expect(parts.Nr(2).Measures[3].Symbols).ShouldNot(BeEmpty())
			Expect(parts.Nr(2).Measures[3].Symbols[0].TimeLine).To(Equal(
				time_line.TimeLine{
					BoundaryType: time_line.Start,
					Type:         time_line.Second,
				}))
		})

		It("should not have modified the time line in part 1", func() {
			Expect(parts.Nr(1).Measures[2].Symbols[0].TimeLine.Type).
				To(Equal(time_line.SecondOf2))
		})
	})
})
