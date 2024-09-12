package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/helper"
	mumotune "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
)

var _ = Describe("TuneFile", func() {
	var tf *TuneFile
	var err error
	var parsedTune *messages.ParsedTune
	var gotTune *mumotune.Tune

	Context("having an empty parsedTune file with correct type", func() {
		BeforeEach(func() {
			tf = &TuneFile{
				Format: fileformat.Format_MUSIC_MODEL,
				Data:   nil,
			}
		})

		When("getting the music model from it", func() {
			BeforeEach(func() {
				gotTune, err = tf.MusicModelTune()
			})

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("a TuneFile created from a music model", func() {
		BeforeEach(func() {
			parsedTune = TestParsedTune("parsedTune 1")
			tf, err = TuneFileFromMusicModelTune(parsedTune.Tune)
		})

		It("should have the correct file type", func() {
			Expect(tf.Format).Should(Equal(fileformat.Format_MUSIC_MODEL))
		})

		It("should not return an error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		When("getting the MusicModelTune from the parsedTune file", func() {
			BeforeEach(func() {
				gotTune, err = tf.MusicModelTune()
			})

			It("should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should return the same MusicModelTune", func() {
				Expect(gotTune).To(BeComparableTo(parsedTune.Tune, helper.MusicModelCompareOptions))
			})
		})

		Context("the parsedTune file has the wrong file type", func() {
			BeforeEach(func() {
				tf.Format = fileformat.Format_BWW
			})

			When("getting the MusicModelTune from parsedTune file", func() {
				BeforeEach(func() {
					gotTune, err = tf.MusicModelTune()
				})

				It("should return an error", func() {
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})
})
