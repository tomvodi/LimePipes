package database

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes/internal/api_gen/apimodel"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/common/music_model"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database/model"
	"github.com/tomvodi/limepipes/internal/database/model/file_type"
	"github.com/tomvodi/limepipes/internal/interfaces/mocks"
	"gorm.io/gorm"
)

var _ = Describe("DbDataService", func() {
	var err error
	var service *dbService
	var gormDb *gorm.DB
	var validator *mocks.ApiModelValidator

	BeforeEach(func() {
		cfg, err := config.InitTest()
		Expect(err).ShouldNot(HaveOccurred())
		gormDb, err = GetInitTestPostgreSQLDB(cfg.DbConfig(), "testdb")
		validator = mocks.NewApiModelValidator(GinkgoT())
		Expect(err).ShouldNot(HaveOccurred())

		service = &dbService{
			db:        gormDb,
			validator: validator,
		}
	})

	AfterEach(func() {
		db, err := gormDb.DB()
		Expect(err).ShouldNot(HaveOccurred())
		err = db.Close()
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

	Context("creating a tune type", func() {
		var marchType, returnedType *model.TuneType
		BeforeEach(func() {
			marchType, err = service.createTuneType("March")
			Expect(err).ShouldNot(HaveOccurred())
		})

		When("calling get or create march with other lowercase letters", func() {
			BeforeEach(func() {
				returnedType, err = service.getOrCreateTuneType("march")
			})

			It("should return the march type", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnedType).Should(BeComparableTo(marchType))
			})
		})

		When("calling get or create a new tune type", func() {
			BeforeEach(func() {
				returnedType, err = service.getOrCreateTuneType("slow march")
			})

			It("should return the new type with capitalized letters", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnedType.Name).Should(Equal("slow march"))
			})
		})

		When("creating the exactly same tune type again", func() {
			BeforeEach(func() {
				returnedType, err = service.createTuneType(marchType.Name)
			})

			It("should fail", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("creating a tune type with empty name", func() {
			BeforeEach(func() {
				returnedType, err = service.createTuneType("")
			})

			It("should fail", func() {
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("creating a tune only with a title", func() {
		var tune *apimodel.Tune
		BeforeEach(func() {
			tune, err = service.CreateTune(apimodel.CreateTune{
				Title: "title",
			}, nil)
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tune).Should(Equal(
				&apimodel.Tune{
					Id:    1,
					Title: "title",
				}))
		})

		When("getting it again by title", func() {
			var returnedTune *apimodel.Tune
			BeforeEach(func() {
				returnedTune, err = service.getTuneByTitle(tune.Title)
			})

			It("should return the same tune", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnedTune).To(Equal(tune))
			})
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
					Id:       1,
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
				returnedTune, err = service.GetTune(tune.Id)
			})

			It("should return the same tune", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(returnedTune).To(Equal(tune))
			})
		})

		When("getting it again by title", func() {
			var returnedTune *apimodel.Tune
			BeforeEach(func() {
				returnedTune, err = service.getTuneByTitle(tune.Title)
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
				validator.EXPECT().ValidateUpdateTune(update).Return(nil)
				tune, err = service.UpdateTune(tune.Id, update)
			})

			It("should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(tune).To(Equal(&apimodel.Tune{
					Id:       1,
					Title:    "new title",
					Type:     "new type",
					TimeSig:  "new time signature",
					Composer: "new composer",
					Arranger: "new arranger",
				}))
			})

			When("retrieving that updated tune", func() {
				BeforeEach(func() {
					tune, err = service.GetTune(tune.Id)
				})

				It("should return the same updated tune", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(tune).To(Equal(&apimodel.Tune{
						Id:       1,
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
				validator.EXPECT().ValidateUpdateTune(update).
					Return(fmt.Errorf("missing title"))
				tune, err = service.UpdateTune(tune.Id, update)
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
				err = service.AddFileToTune(tune.Id, tuneFile)
			})

			It("should add that file", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("retrieving that tune file again", func() {
				BeforeEach(func() {
					returnTuneFile, err = service.GetTuneFile(tune.Id, file_type.MusicModelTune)
				})

				It("should contain that same music model tune", func() {
					returnTune, err := returnTuneFile.MusicModelTune()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(returnTune).Should(Equal(testTune))
				})
			})

			When("deleting that file", func() {
				BeforeEach(func() {
					err = service.DeleteFileFromTune(tune.Id, file_type.MusicModelTune)
				})

				It("should succeed", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})

				When("retrieving all tune files", func() {
					BeforeEach(func() {
						tuneFiles, err = service.GetTuneFiles(tune.Id)
					})

					It("should have no tune files again", func() {
						Expect(err).ShouldNot(HaveOccurred())
						Expect(tuneFiles).To(BeEmpty())
					})
				})
			})

			When("deleting that tune", func() {
				BeforeEach(func() {
					err = service.DeleteTune(tune.Id)
				})

				It("should have deleted that tune", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})

				When("retrieving all tune files", func() {
					BeforeEach(func() {
						tuneFiles, err = service.GetTuneFiles(tune.Id)
					})

					It("should have removed all tune files", func() {
						Expect(tuneFiles).To(BeEmpty())
					})
				})
			})
		})

		When("deleting that tune", func() {
			BeforeEach(func() {
				err = service.DeleteTune(tune.Id)
			})

			It("should have removed that tune", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("retrieving that tune again", func() {
				BeforeEach(func() {
					tune, err = service.GetTune(tune.Id)
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
			tune1.Id = 1
			tune2.Id = 2
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
					Id:          1,
					Title:       "title",
					Description: "desc",
					Creator:     "creator",
				}))
		})

		When("getting it again from service", func() {
			var returnedSet *apimodel.MusicSet
			BeforeEach(func() {
				returnedSet, err = service.GetMusicSet(musicSet.Id)
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
				validator.EXPECT().ValidateUpdateSet(update).Return(nil)
				musicSet, err = service.UpdateMusicSet(musicSet.Id, update)
			})

			It("should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(musicSet).To(Equal(&apimodel.MusicSet{
					Id:          1,
					Title:       "new title",
					Description: "new desc",
					Creator:     "new creator",
				}))
			})

			When("retrieving that updated set", func() {
				BeforeEach(func() {
					musicSet, err = service.GetMusicSet(musicSet.Id)
				})

				It("should return the same updated tune", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(musicSet).To(Equal(&apimodel.MusicSet{
						Id:          1,
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
				validator.EXPECT().ValidateUpdateSet(update).
					Return(fmt.Errorf("missing title"))
				musicSet, err = service.UpdateMusicSet(musicSet.Id, update)
			})

			It("should fail", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("deleting that music set", func() {
			BeforeEach(func() {
				err = service.DeleteMusicSet(musicSet.Id)
			})

			It("should have removed that music set", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("retrieving that music set again", func() {
				BeforeEach(func() {
					musicSet, err = service.GetMusicSet(musicSet.Id)
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
			set1.Id = 1
			set2.Id = 2
			Expect(sets).To(Equal([]*apimodel.MusicSet{
				set1,
				set2,
			}))
		})
	})
})
