package model

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/database/model/file_type"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TuneFile", func() {
	var tf *TuneFile
	var err error
	var tune *music_model.Tune
	var gotTune *music_model.Tune

	Context("having an empty tune file with correct type", func() {
		BeforeEach(func() {
			tf = &TuneFile{
				Type: file_type.MusicModelTune,
				Data: nil,
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
			tune = TestMusicModelTune("tune 1")
			tf, err = TuneFileFromTune(tune)
		})

		It("should have the correct file type", func() {
			Expect(tf.Type).Should(Equal(file_type.MusicModelTune))
		})

		It("should not return an error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		When("getting the MusicModelTune from the tune file", func() {
			BeforeEach(func() {
				gotTune, err = tf.MusicModelTune()
			})

			It("should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should return the same MusicModelTune", func() {
				Expect(gotTune).To(Equal(tune))
			})
		})

		Context("the tune file has the wrong file type", func() {
			BeforeEach(func() {
				tf.Type = file_type.Bww
			})

			When("getting the MusicModelTune from tune file", func() {
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
