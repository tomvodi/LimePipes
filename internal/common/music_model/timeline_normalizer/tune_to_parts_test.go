package timeline_normalizer

import (
	"banduslib/internal/bww"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/expander"
	"banduslib/internal/common/test"
	"banduslib/internal/interfaces"
	"banduslib/internal/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"path/filepath"
	"strings"
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

var _ = Describe("tuneToParts", func() {
	var tune *music_model.Tune
	var partedTune PartedTune
	var parser interfaces.BwwParser
	var musicModel music_model.MusicModel

	BeforeEach(func() {
		parser = bww.NewBwwParser(embExpander)
	})

	JustBeforeEach(func() {
		partedTune = tuneToParts(tune)
	})

	Context("splitting a tune with to parts", func() {
		BeforeEach(func() {
			//createYamlFromBww("./testfiles/tune_with_two_parts.bww", parser)
			musicModel = test.ImportFromYaml("./testfiles/tune_with_two_parts.yaml", embExpander)
			Expect(musicModel).To(HaveLen(1))
			tune = musicModel[0]
		})

		It("should have two parts", func() {
			Expect(parser).ShouldNot(BeNil())
			Expect(partedTune.Parts).To(HaveLen(2))
			Expect(partedTune.Parts[0].WithRepeat).To(BeTrue())
			Expect(partedTune.Parts[1].WithRepeat).To(BeFalse())
		})
	})
})

var _ = Describe("NormalizeTimeline", func() {
	utils.SetupConsoleLogger()
	var err error
	var musicModelBefore music_model.MusicModel
	var musicModelNormalized music_model.MusicModel
	var tuneAfter *music_model.Tune
	var tuneExpected *music_model.Tune
	var normalizer *timelineNormalize
	var parser interfaces.BwwParser

	BeforeEach(func() {
		parser = bww.NewBwwParser(embExpander)
		normalizer = &timelineNormalize{}
	})

	When("having normal 1 - 2 time line", func() {
		BeforeEach(func() {
			//createYamlFromBww("./testfiles/tune_with_normal_1_2_part.bww", parser)
			musicModelBefore = test.ImportFromYaml("./testfiles/tune_with_normal_1_2_part.yaml", embExpander)
			Expect(musicModelBefore).To(HaveLen(1))

			tuneExpected = musicModelBefore[0]
			tuneAfter, err = normalizer.NormalizeTimeline(musicModelBefore[0])
		})

		It("should not change anything", func() {
			Expect(parser).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tuneAfter).Should(BeComparableTo(tuneExpected))
		})
	})

	When("having a 2 - 2 time line", func() {
		BeforeEach(func() {
			//createYamlFromBww("./testfiles/tune_with_2_of_2_part.bww", parser)
			//createYamlFromBww("./testfiles/tune_with_2_of_2_part_normalized.bww", parser)
			musicModelBefore = test.ImportFromYaml("./testfiles/tune_with_2_of_2_part.yaml", embExpander)
			Expect(musicModelBefore).To(HaveLen(1))
			musicModelNormalized = test.ImportFromYaml("./testfiles/tune_with_2_of_2_part_normalized.yaml", embExpander)
			Expect(musicModelNormalized).To(HaveLen(1))

			tuneExpected = musicModelNormalized[0]
			tuneAfter, err = normalizer.NormalizeTimeline(musicModelBefore[0])
		})

		It("should have normalized tune correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tuneAfter).Should(BeComparableTo(tuneExpected))
		})
	})

	When("having a 2 - 2 time line spanning over measures", func() {
		BeforeEach(func() {
			//createYamlFromBww("./testfiles/tune_with_2_of_2_spanning_over_measures.bww", parser)
			//createYamlFromBww("./testfiles/tune_with_2_of_2_spanning_over_measures_normalized.bww", parser)
			musicModelBefore = test.ImportFromYaml("./testfiles/tune_with_2_of_2_spanning_over_measures.yaml", embExpander)
			Expect(musicModelBefore).To(HaveLen(1))
			musicModelNormalized = test.ImportFromYaml("./testfiles/tune_with_2_of_2_spanning_over_measures_normalized.yaml", embExpander)
			Expect(musicModelNormalized).To(HaveLen(1))

			tuneExpected = musicModelNormalized[0]
			tuneAfter, err = normalizer.NormalizeTimeline(musicModelBefore[0])
		})

		It("should have normalized tune correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tuneAfter).Should(BeComparableTo(tuneExpected))
		})
	})
})
