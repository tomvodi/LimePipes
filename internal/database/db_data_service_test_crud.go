package database

import (
	"banduslib/internal/api/apimodel"
	"banduslib/internal/common"
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
			})
		})

		It("should return an error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("creating two tunes with same title", func() {
		BeforeEach(func() {
			title := "test title"
			_, err = service.CreateTune(apimodel.CreateTune{
				Title: title,
			})
			Expect(err).ShouldNot(HaveOccurred())
			_, err = service.CreateTune(apimodel.CreateTune{
				Title: title,
			})
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
			})
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
			})
			tune2, err = service.CreateTune(apimodel.CreateTune{
				Title: "tune2",
			})
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
			})
		})

		It("should return an error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("creating two sets with same title", func() {
		BeforeEach(func() {
			title := "test title"
			_, err = service.CreateMusicSet(apimodel.CreateSet{
				Title: title,
			})
			Expect(err).ShouldNot(HaveOccurred())
			_, err = service.CreateMusicSet(apimodel.CreateSet{
				Title: title,
			})
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
			})
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
			})
			set2, err = service.CreateMusicSet(apimodel.CreateSet{
				Title: "set2",
			})
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
