package database

import (
	"banduslib/internal/api/apimodel"
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/database/model"
	"banduslib/internal/database/model/file_type"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"os"
)

var _ = Describe("DbDataService", func() {
	var err error
	var service *dbService
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

	Context("creating a tune without a title", func() {
		BeforeEach(func() {
			_, err = service.CreateTune(apimodel.CreateTune{
				Title: "",
			}, nil)
		})

		It("should return an error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("creating a valid tune with all fields", func() {
		var tune *apimodel.Tune
		BeforeEach(func() {
			tune, err = service.CreateTune(apimodel.CreateTune{
				Title:    "title",
				Type:     "march",
				TimeSig:  "2/4",
				Composer: "mr. x",
				Arranger: "mr. y",
			}, nil)
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tune).Should(Equal(
				&apimodel.Tune{
					ID:       1,
					Title:    "title",
					Type:     "march",
					TimeSig:  "2/4",
					Composer: "mr. x",
					Arranger: "mr. y",
				}))
		})

		When("getting it again from service", func() {
			var returnedTune *apimodel.Tune
			BeforeEach(func() {
				returnedTune, err = service.GetTune(tune.ID)
			})

			It("should return the same tune", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnedTune).To(Equal(tune))
			})
		})

		When("updating that tune", func() {
			BeforeEach(func() {
				update := apimodel.UpdateTune{
					Title:    "new title",
					Type:     "new type",
					TimeSig:  "new time signature",
					Composer: "new composer",
					Arranger: "new arranger",
				}
				tune, err = service.UpdateTune(tune.ID, update)
			})

			It("should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(tune).To(Equal(&apimodel.Tune{
					ID:       1,
					Title:    "new title",
					Type:     "new type",
					TimeSig:  "new time signature",
					Composer: "new composer",
					Arranger: "new arranger",
				}))
			})

			When("retrieving that updated tune", func() {
				BeforeEach(func() {
					tune, err = service.GetTune(tune.ID)
				})

				It("should return the same updated tune", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(tune).To(Equal(&apimodel.Tune{
						ID:       1,
						Title:    "new title",
						Type:     "new type",
						TimeSig:  "new time signature",
						Composer: "new composer",
						Arranger: "new arranger",
					}))
				})
			})
		})

		When("updating that tune with an empty title", func() {
			BeforeEach(func() {
				update := apimodel.UpdateTune{
					Title:    "",
					Type:     "new type",
					TimeSig:  "new time signature",
					Composer: "new composer",
					Arranger: "new arranger",
				}
				tune, err = service.UpdateTune(tune.ID, update)
			})

			It("should fail", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("adding a file to that tune", func() {
			var testTune *music_model.Tune
			var tuneFile *model.TuneFile
			var tuneFiles []*model.TuneFile
			var returnTuneFile *model.TuneFile

			BeforeEach(func() {
				testTune = model.TestMusicModelTune("test tune")
				tuneFile, err = model.TuneFileFromTune(testTune)
				Expect(err).ShouldNot(HaveOccurred())
				err = service.AddFileToTune(tune.ID, tuneFile)
			})

			It("should add that file", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("retrieving that tune file again", func() {
				BeforeEach(func() {
					returnTuneFile, err = service.GetTuneFile(tune.ID, file_type.MusicModelTune)
				})

				It("should contain that same music model tune", func() {
					returnTune, err := returnTuneFile.MusicModelTune()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(returnTune).Should(Equal(testTune))
				})
			})

			When("deleting that file", func() {
				BeforeEach(func() {
					err = service.DeleteFileFromTune(tune.ID, file_type.MusicModelTune)
				})

				It("should succeed", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})

				When("retrieving all tune files", func() {
					BeforeEach(func() {
						tuneFiles, err = service.GetTuneFiles(tune.ID)
					})

					It("should have no tune files again", func() {
						Expect(err).ShouldNot(HaveOccurred())
						Expect(tuneFiles).To(BeEmpty())
					})
				})
			})

			When("deleting that tune", func() {
				BeforeEach(func() {
					err = service.DeleteTune(tune.ID)
				})

				It("should have deleted that tune", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})

				When("retrieving all tune files", func() {
					BeforeEach(func() {
						tuneFiles, err = service.GetTuneFiles(tune.ID)
					})

					It("should have removed all tune files", func() {
						Expect(tuneFiles).To(BeEmpty())
					})
				})
			})
		})

		When("deleting that tune", func() {
			BeforeEach(func() {
				err = service.DeleteTune(tune.ID)
			})

			It("should have removed that tune", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("retrieving that tune again", func() {
				BeforeEach(func() {
					tune, err = service.GetTune(tune.ID)
				})

				It("should return a not found error", func() {
					Expect(err).To(Equal(common.NotFound))
				})
			})
		})
	})

	Context("creating two tunes", func() {
		var tune1 *apimodel.Tune
		var tune2 *apimodel.Tune
		var tunes []*apimodel.Tune

		BeforeEach(func() {
			tune1, err = service.CreateTune(apimodel.CreateTune{
				Title: "tune1",
			}, nil)
			tune2, err = service.CreateTune(apimodel.CreateTune{
				Title: "tune2",
			}, nil)
		})

		It("should return both tunes", func() {
			tunes, err = service.Tunes()
			Expect(err).ShouldNot(HaveOccurred())
			tune1.ID = 1
			tune2.ID = 2
			Expect(tunes).To(Equal([]*apimodel.Tune{
				tune1,
				tune2,
			}))
		})
	})

	// Sets
	Context("creating a set without a title", func() {
		BeforeEach(func() {
			_, err = service.CreateMusicSet(apimodel.CreateSet{
				Title: "",
			}, nil)
		})

		It("should return an error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("creating a valid set with all fields", func() {
		var musicSet *apimodel.MusicSet
		BeforeEach(func() {
			musicSet, err = service.CreateMusicSet(apimodel.CreateSet{
				Title:       "title",
				Description: "desc",
				Creator:     "creator",
			}, nil)
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicSet).Should(Equal(
				&apimodel.MusicSet{
					ID:          1,
					Title:       "title",
					Description: "desc",
					Creator:     "creator",
				}))
		})

		When("getting it again from service", func() {
			var returnedSet *apimodel.MusicSet
			BeforeEach(func() {
				returnedSet, err = service.GetMusicSet(musicSet.ID)
			})

			It("should return the same musicSet", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnedSet).To(Equal(musicSet))
			})
		})

		When("updating that music set", func() {
			BeforeEach(func() {
				update := apimodel.UpdateSet{
					Title:       "new title",
					Description: "new desc",
					Creator:     "new creator",
				}
				musicSet, err = service.UpdateMusicSet(musicSet.ID, update)
			})

			It("should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(musicSet).To(Equal(&apimodel.MusicSet{
					ID:          1,
					Title:       "new title",
					Description: "new desc",
					Creator:     "new creator",
				}))
			})

			When("retrieving that updated set", func() {
				BeforeEach(func() {
					musicSet, err = service.GetMusicSet(musicSet.ID)
				})

				It("should return the same updated tune", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(musicSet).To(Equal(&apimodel.MusicSet{
						ID:          1,
						Title:       "new title",
						Description: "new desc",
						Creator:     "new creator",
					}))
				})
			})
		})

		When("updating that music set with an empty title", func() {
			BeforeEach(func() {
				update := apimodel.UpdateSet{
					Title:       "",
					Description: "new desc",
					Creator:     "new creator",
				}
				musicSet, err = service.UpdateMusicSet(musicSet.ID, update)
			})

			It("should fail", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("deleting that music set", func() {
			BeforeEach(func() {
				err = service.DeleteMusicSet(musicSet.ID)
			})

			It("should have removed that music set", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("retrieving that music set again", func() {
				BeforeEach(func() {
					musicSet, err = service.GetMusicSet(musicSet.ID)
				})

				It("should return a not found error", func() {
					Expect(err).To(Equal(common.NotFound))
				})
			})
		})
	})

	Context("creating two music sets", func() {
		var set1 *apimodel.MusicSet
		var set2 *apimodel.MusicSet
		var sets []*apimodel.MusicSet

		BeforeEach(func() {
			set1, err = service.CreateMusicSet(apimodel.CreateSet{
				Title: "set1",
			}, nil)
			set2, err = service.CreateMusicSet(apimodel.CreateSet{
				Title: "set2",
			}, nil)
		})

		It("should return both sets", func() {
			sets, err = service.MusicSets()
			Expect(err).ShouldNot(HaveOccurred())
			set1.ID = 1
			set2.ID = 2
			Expect(sets).To(Equal([]*apimodel.MusicSet{
				set1,
				set2,
			}))
		})
	})
})
