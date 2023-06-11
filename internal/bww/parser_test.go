package bww

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/expander"
	"banduslib/internal/common/music_model/import_message"
	"banduslib/internal/interfaces"
	"banduslib/internal/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

var embExpander = expander.NewEmbellishmentExpander()

func dataFromFile(filePath string) []byte {
	bwwFile, err := os.Open(filePath)
	Expect(err).ShouldNot(HaveOccurred())
	var data []byte
	data, err = io.ReadAll(bwwFile)
	Expect(err).ShouldNot(HaveOccurred())

	return data
}

func exportToYaml(muMo music_model.MusicModel, filePath string) {
	data, err := yaml.Marshal(muMo)
	Expect(err).ShouldNot(HaveOccurred())
	err = os.WriteFile(filePath, data, 0664)
	Expect(err).ShouldNot(HaveOccurred())
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

func nilAllMeasureMessages(muMo music_model.MusicModel) {
	for _, tune := range muMo {
		for _, measure := range tune.Measures {
			measure.ImportMessages = nil
		}
	}
}

var _ = Describe("BWW Parser", func() {
	utils.SetupConsoleLogger()
	var err error
	var parser interfaces.BwwParser
	var musicTunesBww music_model.MusicModel
	var musicTunesExpect music_model.MusicModel

	BeforeEach(func() {
		parser = NewBwwParser(embExpander)
	})

	When("parsing a file with a staff with 4 measures in it", func() {
		BeforeEach(func() {
			data := dataFromFile("./testfiles/four_measures.bww")
			musicTunesBww, err = parser.ParseBwwData(data)
			//exportToYaml(musicTunesBww, "../testfiles/four_measures.yaml")
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
			Expect(musicTunesBww[0].Footer).To(Equal([]string{"Tune Footer"}))
		})
	})

	When("having all possible time signatures", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/time_signatures.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/time_signatures.yaml")
		})

		It("should have parsed all measures", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having a tune with all kinds of melody notes", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/all_melody_notes.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("../testfiles/all_melody_notes.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/all_melody_notes.yaml")
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
			//exportToYaml(musicTunesBww, "./testfiles/embellishment_without_following_note.yaml")
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
			musicTunesExpect = importFromYaml("../testfiles/single_graces.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/single_graces.yaml")
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
			musicTunesExpect = importFromYaml("../testfiles/dots.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/dots.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having fermatas for melody notes", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/fermatas.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/fermatas.yaml")
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
			musicTunesExpect = importFromYaml("../testfiles/rests.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/rests.yaml")
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
			musicTunesExpect = importFromYaml("../testfiles/accidentals.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/accidentals.yaml")
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
			musicTunesExpect = importFromYaml("../testfiles/doublings.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/doublings.yaml")
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
			musicTunesExpect = importFromYaml("../testfiles/grips.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/grips.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having taorluaths", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/taorluaths.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("../testfiles/taorluaths.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/taorluaths.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having bubblys", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/bubblys.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("../testfiles/bubblys.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/bubblys.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having throw on d", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/throwds.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("../testfiles/throwds.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/throwds.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having birls", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/birls.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("../testfiles/birls.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/birls.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having strikes", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/strikes.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("../testfiles/strikes.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/strikes.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having peles", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/peles.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("../testfiles/peles.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/peles.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having double strikes", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/double_strikes.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("../testfiles/double_strikes.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/double_strikes.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having triple strikes", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/triple_strikes.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("../testfiles/triple_strikes.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/triple_strikes.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having double graces", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/double_grace.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/double_grace.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/double_grace.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having ties", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/ties.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/ties.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/ties.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having ties in old format with error messages", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/ties_old_with_errors.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/ties_old_with_errors.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/ties_old_with_errors.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww[0].Measures[0].ImportMessages[0]).Should(
				Equal(&import_message.ImportMessage{
					Symbol: "^tla",
					Type:   import_message.Warning,
					Text:   "tie in old format (^tla) must follow a note and can't be the first symbol in a measure",
					Fix:    import_message.SkipSymbol,
				}))
			Expect(musicTunesBww[0].Measures[1].ImportMessages[0]).Should(
				Equal(&import_message.ImportMessage{
					Symbol: "^tlg",
					Type:   import_message.Error,
					Text:   "tie in old format (^tlg) must follow a note with pitch and length",
					Fix:    import_message.SkipSymbol,
				}))
			nilAllMeasureMessages(musicTunesBww)
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having ties in old format", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/ties_old.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/ties_old.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/ties_old.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having irregular groups", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/irregular_groups.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/irregular_groups.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/irregular_groups.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having triplets", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/triplets.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/triplets.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/triplets.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having time lines", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/time_lines.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/time_lines.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/time_lines.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having space symbols", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/space.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/space.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/space.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having a file with a tune containing inline text and comments", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/tune_with_inline_comments.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/tune_with_inline_comments.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_inline_comments.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having a file with two tunes", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/two_tunes.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/two_tunes.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/two_tunes.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having a file with a tune with comments, the comment should not be propagated to first measure", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/single_tune_comment_does_not_appear_in_first_measure.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/single_tune_comment_does_not_appear_in_first_measure.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/single_tune_comment_does_not_appear_in_first_measure.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having a file with the first tune without a title", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/first_tune_no_title.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/first_tune_no_title.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/first_tune_no_title.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having a tune with no proper staff ending before next staff starts", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/tune_with_no_staff_ending_before_next_staff.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/tune_with_no_staff_ending_before_next_staff.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_no_staff_ending_before_next_staff.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having a tune staff that ends with EOF", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/tune_staff_ends_with_eof.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/tune_staff_ends_with_eof.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/tune_staff_ends_with_eof.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having tune title and config with missing parameter in list", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/tune_with_missing_parameter_in_list.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/tune_with_missing_parameter_in_list.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_missing_parameter_in_list.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with multiple bagpipe reader version definitions", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/tune_with_multiple_bagpipe_reader_version_definitions.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/tune_with_multiple_bagpipe_reader_version_definitions.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_multiple_bagpipe_reader_version_definitions.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having tune with symbol and measure comments", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/tune_with_symbol_comment.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/tune_with_symbol_comment.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_symbol_comment.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having tune with time line end after staff end", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/tune_with_time_line_end_after_staff_end.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/tune_with_time_line_end_after_staff_end.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_time_line_end_after_staff_end.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having tune with inline tune tempo", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/tunetempo_inline.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/tunetempo_inline.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/tunetempo_inline.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with all cadences in it", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/cadences.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/cadences.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/cadences.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with piobairached throws and doublings in it", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/pio_throws_and_doublings.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/pio_throws_and_doublings.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/pio_throws_and_doublings.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with piobairached grips in it", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/pio_grips.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/pio_grips.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/pio_grips.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with piobairached echo beats in it", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/pio_echo_beats.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/pio_echo_beats.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/pio_echo_beats.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with piobairached darodos in it", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/pio_darodo.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/pio_darodo.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/pio_darodo.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with piobairached lemluaths in it", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/pio_lemluaths.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/pio_lemluaths.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/pio_lemluaths.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with piobairached taorluaths in it", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/pio_taorluaths.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/pio_taorluaths.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/pio_taorluaths.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with piobairached crunluaths in it", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/pio_crunluaths.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/pio_crunluaths.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/pio_crunluaths.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with piobairached triplings in it", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/pio_triplings.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/pio_triplings.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/pio_triplings.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with misc movements in it", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/pio_misc.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/pio_misc.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/pio_misc.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with segno and dalsegno", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/segno_dalsegno.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/segno_dalsegno.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/segno_dalsegno.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having file with fine and dacapoalfine", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/fine_dacapoalfine.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/fine_dacapoalfine.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/fine_dacapoalfine.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having inline comment shouldn't remove measures", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/inline_comment_removes_first_staff_measures.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("./testfiles/inline_comment_removes_first_staff_measures.yaml")
			//exportToYaml(musicTunesBww, "./testfiles/inline_comment_removes_first_staff_measures.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(BeComparableTo(musicTunesExpect))
		})
	})

	When("having a tune with repeats", func() {
		BeforeEach(func() {
			bwwData := dataFromFile("./testfiles/tune_with_repeats.bww")
			musicTunesBww, err = parser.ParseBwwData(bwwData)
			musicTunesExpect = importFromYaml("../testfiles/tune_with_repeats.yaml")
			//exportToYaml(musicTunesBww, "../testfiles/tune_with_repeats.yaml")
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

	When("parsing the file with all piobaireached symbols in it", func() {
		BeforeEach(func() {
			data := dataFromFile("./testfiles/all_piobaireached_symbols.bww")
			musicTunesBww, err = parser.ParseBwwData(data)
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(HaveLen(11))
		})
	})

})
