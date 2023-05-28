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
	var muMo music_model.MusicModel
	var gotMuMo music_model.MusicModel

	Context("having an empty tune file with correct type", func() {
		BeforeEach(func() {
			tf = &TuneFile{
				Type: file_type.MusicModel,
				Data: nil,
			}
		})

		When("getting the music model from it", func() {
			BeforeEach(func() {
				gotMuMo, err = tf.MusicModel()
			})

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("a TuneFile created from a music model", func() {
		BeforeEach(func() {
			muMo = testMusicModel()
			tf, err = TuneFileFromMusicModel(muMo)
		})

		It("should have the correct file type", func() {
			Expect(tf.Type).Should(Equal(file_type.MusicModel))
		})

		It("should not return an error", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		When("getting the MusicModel from the tune file", func() {
			BeforeEach(func() {
				gotMuMo, err = tf.MusicModel()
			})

			It("should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should return the same MusicModel", func() {
				Expect(gotMuMo).To(Equal(muMo))
			})
		})

		Context("the tune file has the wrong file type", func() {
			BeforeEach(func() {
				tf.Type = file_type.Bww
			})

			When("getting the MusicModel from tune file", func() {
				BeforeEach(func() {
					gotMuMo, err = tf.MusicModel()
				})

				It("should return an error", func() {
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})
})
