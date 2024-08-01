package database

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes/internal/api_gen/apimodel"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/common/music_model"
	"github.com/tomvodi/limepipes/internal/common/music_model/helper"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database/model"
	"github.com/tomvodi/limepipes/internal/database/model/file_type"
	"gorm.io/gorm"
)

var _ = Describe("DbDataService Import", func() {
	var err error
	var returnTunes []*apimodel.ImportTune
	var bwwFileData *common.BwwFileTuneData
	var tuneFile *model.TuneFile
	var tuneFileTune *tune.Tune
	var service *dbService
	var muMo music_model.MusicModel
	var musicSet *apimodel.MusicSet
	var fileInfo *common.ImportFileInfo
	var gormDb *gorm.DB

	BeforeEach(func() {
		cfg, err := config.InitTest()
		Expect(err).ShouldNot(HaveOccurred())
		gormDb, err = GetInitTestPostgreSQLDB(cfg.DbConfig(), "testdb")
		Expect(err).ShouldNot(HaveOccurred())

		service = &dbService{db: gormDb}
	})

	AfterEach(func() {
		db, err := gormDb.DB()
		Expect(err).ShouldNot(HaveOccurred())
		err = db.Close()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("having a music model with two tunes", func() {
		BeforeEach(func() {
			fileInfo, err = common.NewImportFileInfo("testfile.bww", []byte(`BagpipeReader:1.0`))
			Expect(err).ShouldNot(HaveOccurred())
			muMo = music_model.MusicModel{
				model.TestMusicModelTune("tune 1"),
				model.TestMusicModelTune("tune 2"),
			}
		})

		When("importing this music model", func() {
			BeforeEach(func() {
				returnTunes, err = service.ImportMusicModel(muMo, fileInfo, nil)
			})

			It("should return two apimodel tunes", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnTunes).Should(HaveLen(2))
				Expect(returnTunes[0].Set).ShouldNot(BeNil())
				Expect(returnTunes[1].Set).ShouldNot(BeNil())
				setId := returnTunes[0].Set.Id
				Expect(setId).To(Equal(returnTunes[1].Set.Id))
			})

			It("should have imported both tunes into database", func() {
				Expect(returnTunes[0].ImportedToDatabase).To(BeTrue())
				Expect(returnTunes[1].ImportedToDatabase).To(BeTrue())
			})

			When("retrieving tune file for music model", func() {
				BeforeEach(func() {
					tuneFile, err = service.GetTuneFile(
						returnTunes[0].Id,
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
						Expect(tuneFileTune).Should(BeComparableTo(muMo[0], helper.CompareOpts))
						Expect(returnTunes[0].Set).ShouldNot(BeNil())
						Expect(returnTunes[1].Set).ShouldNot(BeNil())
						setId := returnTunes[0].Set.Id
						Expect(setId).To(Equal(returnTunes[1].Set.Id))
					})
				})
			})

			When("importing this music model a second time", func() {
				BeforeEach(func() {
					returnTunes, err = service.ImportMusicModel(muMo, fileInfo, nil)
				})

				It("should return two apimodel tunes again", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(returnTunes).Should(HaveLen(2))
				})

				It("shouldn't have imported both tunes again", func() {
					Expect(returnTunes[0].ImportedToDatabase).To(BeFalse())
					Expect(returnTunes[1].ImportedToDatabase).To(BeFalse())
				})
			})

			When("having a a tune with title of already imported tune but with another arranger", func() {
				BeforeEach(func() {
					tune1 := muMo[0]
					tune1.Arranger = "another arranger"
					muMo = music_model.MusicModel{
						tune1,
					}
				})

				When("importing that tune with different arranger", func() {
					BeforeEach(func() {
						returnTunes, err = service.ImportMusicModel(muMo, fileInfo, nil)
					})

					It("should succeed", func() {
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should return that tune", func() {
						Expect(returnTunes).To(HaveLen(1))
						Expect(returnTunes[0].Arranger).To(Equal("another arranger"))
					})
				})
			})
		})

		Context("having bww tune file data", func() {
			BeforeEach(func() {
				bwwFileData = &common.BwwFileTuneData{}
				bwwFileData.AddTuneData(muMo[0].Title, []byte("& LA_4 !t"))
				bwwFileData.AddTuneData(muMo[1].Title, []byte("& B_4 !t"))
			})

			When("importing this music model", func() {

				BeforeEach(func() {
					returnTunes, err = service.ImportMusicModel(muMo, fileInfo, bwwFileData)
				})

				It("should return two apimodel tunes", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(returnTunes).Should(HaveLen(2))
					Expect(returnTunes[0].Set).ShouldNot(BeNil())
					Expect(returnTunes[1].Set).ShouldNot(BeNil())
					setId := returnTunes[0].Set.Id
					Expect(setId).To(Equal(returnTunes[1].Set.Id))
				})

				When("retrieving the tune file for bww", func() {
					var getTuneFileErr error
					BeforeEach(func() {
						tuneFile, getTuneFileErr = service.GetTuneFile(
							returnTunes[0].Id,
							file_type.Bww,
						)
					})

					It("should return the tune file data", func() {
						Expect(getTuneFileErr).ShouldNot(HaveOccurred())
						bwwData := bwwFileData.Data(0)
						Expect(tuneFile.Data).To(Equal(bwwData))
					})
				})
			})
		})
	})

	Context("having a music model with three tunes, where two of them have the same title", func() {
		BeforeEach(func() {
			fileInfo, err = common.NewImportFileInfo("testfile.bww", []byte(`BagpipeReader:1.0`))
			Expect(err).ShouldNot(HaveOccurred())
			muMo = music_model.MusicModel{
				model.TestMusicModelTune("scotty"),
				model.TestMusicModelTune("wings"),
				model.TestMusicModelTune("scotty"),
			}
		})

		When("importing this music model", func() {
			BeforeEach(func() {
				returnTunes, err = service.ImportMusicModel(muMo, fileInfo, nil)
			})

			It("should return three apimodel tunes", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnTunes).Should(HaveLen(3))
			})

			It("shouldn't have imported the last tune into database", func() {
				Expect(returnTunes[0].ImportedToDatabase).To(BeTrue())
				Expect(returnTunes[1].ImportedToDatabase).To(BeTrue())
				Expect(returnTunes[2].ImportedToDatabase).To(BeFalse())
			})

			When("I retrieve the set", func() {
				BeforeEach(func() {
					setId := returnTunes[0].Set.Id
					musicSet, err = service.GetMusicSet(setId)
				})

				It("should successfully got that set", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("should have three tunes, where the first and last are the same", func() {
					Expect(musicSet.Tunes).To(HaveLen(3))
					Expect(musicSet.Tunes[0]).To(Equal(musicSet.Tunes[2]))
				})
			})
		})
	})
})
