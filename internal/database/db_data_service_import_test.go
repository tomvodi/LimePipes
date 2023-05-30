package database

import (
	"banduslib/internal/api/apimodel"
	"banduslib/internal/common/music_model"
	"banduslib/internal/database/model"
	"banduslib/internal/database/model/file_type"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"os"
)

var _ = Describe("DbDataService Import", func() {
	var err error
	var returnTunes []*apimodel.ImportTune
	var tuneFile *model.TuneFile
	var tuneFileTune *music_model.Tune
	var service *dbService
	var muMo music_model.MusicModel
	var filename string
	var gormDb *gorm.DB

	BeforeEach(func() {
		gormDb, err = GetInitSqliteDb("testing.db")
		Expect(err).ShouldNot(HaveOccurred())

		service = &dbService{db: gormDb}
	})

	AfterEach(func() {
		db, err := gormDb.DB()
		Expect(err).ShouldNot(HaveOccurred())
		err = db.Close()
		Expect(err).ShouldNot(HaveOccurred())
		err = os.Remove("testing.db")
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("having a music model with two tunes", func() {
		BeforeEach(func() {
			filename = "testfile"
			muMo = music_model.MusicModel{
				model.TestMusicModelTune("tune 1"),
				model.TestMusicModelTune("tune 2"),
			}
		})

		When("importing this music model", func() {
			BeforeEach(func() {
				returnTunes, err = service.ImportMusicModel(muMo, filename)
			})

			It("should return two apimodel tunes", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnTunes).Should(HaveLen(2))
				Expect(returnTunes[0].Set).ShouldNot(BeNil())
				Expect(returnTunes[1].Set).ShouldNot(BeNil())
				setId := returnTunes[0].Set.ID
				Expect(setId).To(Equal(returnTunes[1].Set.ID))
			})

			When("retrieving tune file for music model", func() {
				BeforeEach(func() {
					tuneFile, err = service.GetTuneFile(
						returnTunes[0].ID,
						file_type.MusicModelTune,
					)
				})

				It("should return the tune file", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})

				When("getting music model tune from tune file", func() {
					BeforeEach(func() {
						tuneFileTune, err = tuneFile.MusicModelTune()
					})

					It("should return the same data as for the imported music model tune", func() {
						Expect(err).ShouldNot(HaveOccurred())
						Expect(tuneFileTune).Should(Equal(muMo[0]))
					})
				})
			})
		})
	})
})
