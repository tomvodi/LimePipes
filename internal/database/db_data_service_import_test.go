package database

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/helper"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/file_type"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes/internal/apigen/apimodel"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database/model"
	"gorm.io/gorm"
)

var _ = Describe("DbDataService Import", func() {
	var err error
	var returnTunes []*apimodel.ImportTune
	var returnSet *apimodel.BasicMusicSet
	var tuneFile *model.TuneFile
	var tuneFileTune *tune.Tune
	var service *dbService
	var importTunes []*messages.ImportedTune
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

	Context("having a music model with one tune without tune file data", func() {
		BeforeEach(func() {
			fileInfo, err = common.NewImportFileInfo("testfile.bww", file_type.Type_BWW, []byte(`BagpipeReader:1.0`))
			Expect(err).ShouldNot(HaveOccurred())
			importTunes = []*messages.ImportedTune{
				model.TestImportedTune("tune 1"),
			}
		})

		When("importing this music model", func() {
			BeforeEach(func() {
				returnTunes, returnSet, err = service.ImportTunes(importTunes, fileInfo)
			})

			It("should return one apimodel tune", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnTunes).Should(HaveLen(1))
			})
		})
	})

	Context("having a music model with two tunes", func() {
		BeforeEach(func() {
			fileInfo, err = common.NewImportFileInfo("testfile.bww", file_type.Type_BWW, []byte(`BagpipeReader:1.0`))
			Expect(err).ShouldNot(HaveOccurred())
			importTunes = []*messages.ImportedTune{
				model.TestImportedTune("tune 1"),
				model.TestImportedTune("tune 2"),
			}
		})

		When("importing this music model", func() {
			BeforeEach(func() {
				returnTunes, returnSet, err = service.ImportTunes(importTunes, fileInfo)
			})

			It("should return two apimodel tunes", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnTunes).Should(HaveLen(2))
				Expect(returnSet).ShouldNot(BeNil())
			})

			When("retrieving tune file for music model", func() {
				BeforeEach(func() {
					tuneFile, err = service.GetTuneFile(
						returnTunes[0].Id,
						file_type.Type_MUSIC_MODEL,
					)
				})

				It("should return the tune file", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})

				When("getting music model tune from the tune file", func() {
					BeforeEach(func() {
						tuneFileTune, err = tuneFile.MusicModelTune()
					})

					It("should return the same data as for the imported music model tune", func() {
						Expect(err).ShouldNot(HaveOccurred())
						Expect(tuneFileTune).Should(BeComparableTo(importTunes[0].Tune, helper.MusicModelCompareOptions))
						Expect(returnSet).ShouldNot(BeNil())
					})
				})
			})

			When("importing this music model a second time", func() {
				BeforeEach(func() {
					_, _, err = service.ImportTunes(importTunes, fileInfo)
				})

				It("should return an error", func() {
					Expect(err).Should(HaveOccurred())
				})
			})

			When("having a a tune with title of already imported tune but with another arranger", func() {
				BeforeEach(func() {
					tune1 := importTunes[0]
					tune1.Tune.Arranger = "another arranger"
					importTunes = []*messages.ImportedTune{
						tune1,
					}
				})

				When("importing that tune with different arranger", func() {
					BeforeEach(func() {
						fileInfo.Hash = "another hash because of different arranger"
						returnTunes, returnSet, err = service.ImportTunes(importTunes, fileInfo)
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
				importTunes[0].TuneFileData = []byte("& LA_4 !t")
				importTunes[1].TuneFileData = []byte("& B_4 !t")
			})

			When("importing this music model", func() {

				BeforeEach(func() {
					returnTunes, returnSet, err = service.ImportTunes(importTunes, fileInfo)
				})

				It("should return two apimodel tunes", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(returnTunes).Should(HaveLen(2))
					Expect(returnSet).ShouldNot(BeNil())
				})

				When("retrieving the tune file for bww", func() {
					var getTuneFileErr error
					BeforeEach(func() {
						tuneFile, getTuneFileErr = service.GetTuneFile(
							returnTunes[0].Id,
							file_type.Type_BWW,
						)
					})

					It("should return the tune file data", func() {
						Expect(getTuneFileErr).ShouldNot(HaveOccurred())
						Expect(tuneFile.Data).To(Equal(importTunes[0].TuneFileData))
					})
				})
			})
		})
	})

	Context("having a music model with three tunes, where two of them have the same title", func() {
		BeforeEach(func() {
			fileInfo, err = common.NewImportFileInfo("testfile.bww", file_type.Type_BWW, []byte(`BagpipeReader:1.0`))
			Expect(err).ShouldNot(HaveOccurred())
			importTunes = []*messages.ImportedTune{
				model.TestImportedTune("scotty"),
				model.TestImportedTune("wings"),
				model.TestImportedTune("scotty"),
			}
		})

		Context("when tunes with duplicate title have different file data", func() {
			BeforeEach(func() {
				importTunes[0].TuneFileData = []byte("& LA_4 !t")
				importTunes[1].TuneFileData = []byte("& B_4 !t")
				importTunes[2].TuneFileData = []byte("& C_4 !t")
			})

			When("importing this music model", func() {
				BeforeEach(func() {
					returnTunes, returnSet, err = service.ImportTunes(importTunes, fileInfo)
				})

				It("should return three apimodel tunes", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(returnTunes).Should(HaveLen(3))
				})

				It("should have imported all three tunes into database", func() {
					Expect(returnTunes[0].Id).ToNot(Equal(returnTunes[2].Id))
				})

				When("I retrieve the set", func() {
					BeforeEach(func() {
						setID := returnSet.Id
						musicSet, err = service.GetMusicSet(setID)
					})

					It("should successfully got that set", func() {
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should have three tunes, where the tunes with same title are not the same", func() {
						Expect(musicSet.Tunes).To(HaveLen(3))
						Expect(musicSet.Tunes[0]).NotTo(Equal(musicSet.Tunes[2]))
					})
				})
			})
		})

		Context("when tunes with duplicate title have the same file data", func() {
			BeforeEach(func() {
				importTunes[0].TuneFileData = []byte("& LA_4 !t")
				importTunes[1].TuneFileData = []byte("& B_4 !t")
				importTunes[2].TuneFileData = []byte("& LA_4 !t")
			})

			When("importing this music model", func() {
				BeforeEach(func() {
					returnTunes, returnSet, err = service.ImportTunes(importTunes, fileInfo)
				})

				It("should return three apimodel tunes", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(returnTunes).Should(HaveLen(3))
				})

				It("should have imported two tunes and return the duplicate for the third one", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(returnTunes[0].Id).To(Equal(returnTunes[2].Id))
				})

				When("I retrieve the set", func() {
					BeforeEach(func() {
						setID := returnSet.Id
						musicSet, err = service.GetMusicSet(setID)
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
})
