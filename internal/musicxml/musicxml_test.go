package musicxml

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/expander"
	"banduslib/internal/musicxml/model"
	"banduslib/internal/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	"os"
)

var embExpander = expander.NewEmbellishmentExpander()

func exportToMusicXml(score *model.Score, filePath string) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	Expect(err).ShouldNot(HaveOccurred())
	defer f.Close()

	err = WriteScore(score, f)
	Expect(err).ShouldNot(HaveOccurred())
}

func importFromMusicXml(filePath string) *model.Score {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0660)
	Expect(err).ShouldNot(HaveOccurred())
	defer f.Close()

	score, err := ReadScore(f)
	Expect(err).ShouldNot(HaveOccurred())

	return score
}

func importFromYaml(filePath string) music_model.MusicModel {
	muMo := make(music_model.MusicModel, 0)
	fileData, err := os.ReadFile(filePath)
	Expect(err).ShouldNot(HaveOccurred())
	err = yaml.Unmarshal(fileData, &muMo)
	Expect(err).ShouldNot(HaveOccurred())

	embExpander.ExpandModel(muMo)

	return muMo
}

var _ = Describe("ScoreFromMusicModelTune", func() {
	utils.SetupConsoleLogger()
	var err error
	var score *model.Score
	var readScore *model.Score

	Context("having a tune with four measures", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/four_measures.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/four_measures.musicxml")
			readScore = importFromMusicXml("./testfiles/four_measures.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with all melody notes", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/all_melody_notes.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/all_melody_notes.musicxml")
			readScore = importFromMusicXml("./testfiles/all_melody_notes.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with single grace notes", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/single_graces.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/single_graces.musicxml")
			readScore = importFromMusicXml("./testfiles/single_graces.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with doublings", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/doublings.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/doublings.musicxml")
			readScore = importFromMusicXml("./testfiles/doublings.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with strikes", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/strikes.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/strikes.musicxml")
			readScore = importFromMusicXml("./testfiles/strikes.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with grips", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/grips.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/grips.musicxml")
			readScore = importFromMusicXml("./testfiles/grips.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with taorluaths", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/taorluaths.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/taorluaths.musicxml")
			readScore = importFromMusicXml("./testfiles/taorluaths.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with bubblys", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/bubblys.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/bubblys.musicxml")
			readScore = importFromMusicXml("./testfiles/bubblys.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with birls", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/birls.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/birls.musicxml")
			readScore = importFromMusicXml("./testfiles/birls.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with throw on Ds", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/throwds.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/throwds.musicxml")
			readScore = importFromMusicXml("./testfiles/throwds.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with peles", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/peles.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/peles.musicxml")
			readScore = importFromMusicXml("./testfiles/peles.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with double strikes", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/double_strikes.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/double_strikes.musicxml")
			readScore = importFromMusicXml("./testfiles/double_strikes.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})

	Context("having a file with triple strikes", func() {
		BeforeEach(func() {
			muMo := importFromYaml("../testfiles/triple_strikes.yaml")
			score, err = ScoreFromMusicModelTune(muMo[0])
			//exportToMusicXml(score, "./testfiles/triple_strikes.musicxml")
			readScore = importFromMusicXml("./testfiles/triple_strikes.musicxml")
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(readScore).Should(BeComparableTo(score))
		})
	})
})
