package bww

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

func dataFromFile(filePath string) []byte {
	bwwFile, err := os.Open(filePath)
	Expect(err).ShouldNot(HaveOccurred())
	var data []byte
	data, err = io.ReadAll(bwwFile)
	Expect(err).ShouldNot(HaveOccurred())

	return data
}

func exportToYaml(tunes []*music_model.Tune, filePath string) {
	data, err := yaml.Marshal(tunes)
	Expect(err).ShouldNot(HaveOccurred())
	err = os.WriteFile(filePath, data, 0664)
	Expect(err).ShouldNot(HaveOccurred())
}

func importFromYaml(filePath string) []*music_model.Tune {
	tunes := make([]*music_model.Tune, 0)
	fileData, err := os.ReadFile(filePath)
	Expect(err).ShouldNot(HaveOccurred())
	err = yaml.Unmarshal(fileData, &tunes)
	Expect(err).ShouldNot(HaveOccurred())

	return tunes
}

var _ = Describe("BWW Parser", func() {
	var err error
	var parser interfaces.BwwParser
	var musicTunesBww []*music_model.Tune
	var musicTunesExpect []*music_model.Tune

	BeforeEach(func() {
		parser = NewBwwParser()
	})

	When("parsing a file with a staff with 4 measures in it", func() {
		BeforeEach(func() {
			data := dataFromFile("./testfiles/four_measures.bww")
			musicTunesBww, err = parser.ParseBwwData(data)
		})

		It("should have parsed 4 measures", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(HaveLen(1))
			Expect(musicTunesBww[0].Measures).To(HaveLen(4))
		})
	})

	When("having a tune with title, composer, type and footer", func() {
		BeforeEach(func() {
			data := dataFromFile("./testfiles/full_tune_header.bww")
			musicTunesBww, err = parser.ParseBwwData(data)
		})

		It("should have parsed 4 measures", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(HaveLen(1))
			Expect(musicTunesBww[0].Title).To(Equal("Tune Title"))
			Expect(musicTunesBww[0].Composer).To(Equal("Composer"))
			Expect(musicTunesBww[0].Type).To(Equal("Tune Type"))
			Expect(musicTunesBww[0].Footer).To(Equal("Tune Footer"))
		})
	})

	When("having a tune with all kinds of melody notes", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/all_melody_notes.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/all_melody_notes.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having only an embellishment without a following melody note", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/embellishment_without_following_note.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/embellishment_without_following_note.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(Equal(musicTunesExpect))
		})
	})

	When("having single grace notes following a melody note", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/single_graces.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/single_graces.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(Equal(musicTunesExpect))
		})
	})

	When("having dots for the melody note", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/dots.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/dots.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having rests", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/rests.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/rests.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having accidentals", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/accidentals.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/accidentals.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having doublings", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/doublings.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/doublings.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having grips", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/grips.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/grips.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("parsing the file with all bww symbols in it", func() {
		BeforeEach(func() {
			data := dataFromFile("./testfiles/all_symbols.bww")
			musicTunesBww, err = parser.ParseBwwData(data)
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(HaveLen(2))
		})
	})
})
