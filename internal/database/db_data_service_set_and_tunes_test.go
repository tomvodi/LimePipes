package database

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes/internal/api/apimodel"
	"gorm.io/gorm"
	"os"
)

var _ = Describe("DbDataService", func() {
	var err error
	var service *dbService
	var gormDb *gorm.DB
	var tune1 *apimodel.Tune
	var tune2 *apimodel.Tune
	var tune3 *apimodel.Tune
	var musicSet *apimodel.MusicSet

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

	Context("having some tunes created", func() {
		BeforeEach(func() {
			tune1, err = service.CreateTune(apimodel.CreateTune{Title: "tune 1"}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			tune2, err = service.CreateTune(apimodel.CreateTune{Title: "tune 2"}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			tune3, err = service.CreateTune(apimodel.CreateTune{Title: "tune 3"}, nil)
			Expect(err).ShouldNot(HaveOccurred())
		})

		When("creating a set with these tunes", func() {
			var expectedTuneOrder []apimodel.Tune

			BeforeEach(func() {
				expectedTuneOrder = []apimodel.Tune{
					*tune2,
					*tune3,
					*tune1,
				}
				var tuneIds []uint64
				for _, tune := range expectedTuneOrder {
					tuneIds = append(tuneIds, tune.ID)
				}

				musicSet, err = service.CreateMusicSet(
					apimodel.CreateSet{
						Title: "test music set",
						Tunes: tuneIds,
					},
					nil,
				)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should have created that set with these tunes", func() {
				Expect(musicSet.Tunes).To(Equal(expectedTuneOrder))
			})

			When("retrieving that music set from database", func() {
				BeforeEach(func() {
					musicSet, err = service.GetMusicSet(musicSet.ID)
				})

				It("should have the tunes in correct order", func() {
					Expect(err).ShouldNot(HaveOccurred())
					Expect(musicSet.Tunes).To(Equal(expectedTuneOrder))
				})
			})

			When("retrieving the music set by tune ids", func() {
				var foundMusicSet *apimodel.MusicSet
				BeforeEach(func() {
					foundMusicSet, err = service.getMusicSetByTuneIds([]uint64{
						expectedTuneOrder[0].ID,
						expectedTuneOrder[1].ID,
						expectedTuneOrder[2].ID,
					})
				})

				It("should get the music set", func() {
					Expect(foundMusicSet).To(Equal(musicSet))
				})
			})

			When("updating the music set with another order of tunes", func() {
				BeforeEach(func() {
					expectedTuneOrder = []apimodel.Tune{
						*tune3,
						*tune2,
						*tune1,
					}
					var tuneIds []uint64
					for _, tune := range expectedTuneOrder {
						tuneIds = append(tuneIds, tune.ID)
					}

					musicSet, err = service.UpdateMusicSet(
						musicSet.ID,
						apimodel.UpdateSet{
							Title: "test music set",
							Tunes: tuneIds,
						},
					)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("should have updated that set with these new tune order", func() {
					Expect(musicSet.Tunes).To(Equal(expectedTuneOrder))
				})
			})
		})

		Context("creating a set without tunes", func() {
			var musicSetAfterAssignment *apimodel.MusicSet
			var musicSets []*apimodel.MusicSet

			BeforeEach(func() {
				musicSet, err = service.CreateMusicSet(
					apimodel.CreateSet{Title: "test music set"},
					nil,
				)
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("assigning tunes in an arbitrary order "+
				"to the music set", func() {
				BeforeEach(func() {
					musicSetAfterAssignment, err = service.AssignTunesToMusicSet(
						musicSet.ID,
						[]uint64{tune2.ID, tune1.ID, tune3.ID},
					)
				})

				It("should succeed", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("should have the tunes assigned in the same order", func() {
					Expect(musicSetAfterAssignment.Tunes).To(Equal(
						[]apimodel.Tune{
							*tune2,
							*tune1,
							*tune3,
						}))
				})

				When("getting the same set from service", func() {
					BeforeEach(func() {
						musicSetAfterAssignment, err = service.GetMusicSet(musicSetAfterAssignment.ID)
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should also have the tunes in the same order", func() {
						Expect(musicSetAfterAssignment.Tunes).To(Equal(
							[]apimodel.Tune{
								*tune2,
								*tune1,
								*tune3,
							}))
					})
				})

				When("getting the list of sets", func() {
					BeforeEach(func() {
						musicSets, err = service.MusicSets()
					})

					It("should also have the tunes in the same order", func() {
						Expect(err).ShouldNot(HaveOccurred())
						Expect(musicSets[0].Tunes).To(Equal(
							[]apimodel.Tune{
								*tune2,
								*tune1,
								*tune3,
							}))
					})
				})

				When("trying to delete a tune that is assigned to the set", func() {
					BeforeEach(func() {
						err = service.DeleteTune(musicSetAfterAssignment.Tunes[0].ID)
					})

					It("should not be possible", func() {
						Expect(err).Should(HaveOccurred())
					})
				})

				When("deleting that set", func() {
					BeforeEach(func() {
						err = service.DeleteMusicSet(musicSetAfterAssignment.ID)
					})

					It("should get deleted", func() {
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

		})
	})

})
