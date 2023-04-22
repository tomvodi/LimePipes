package database

import (
	"banduslib/internal/api/apimodel"
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

	Context("having some tunes and a set", func() {
		var tune1 *apimodel.Tune
		var tune2 *apimodel.Tune
		var tune3 *apimodel.Tune
		var initialMusicSet *apimodel.MusicSet
		var musicSetAfterAssignment *apimodel.MusicSet

		BeforeEach(func() {
			tune1, err = service.CreateTune(apimodel.CreateTune{Title: "tune 1"})
			Expect(err).ShouldNot(HaveOccurred())
			tune2, err = service.CreateTune(apimodel.CreateTune{Title: "tune 2"})
			Expect(err).ShouldNot(HaveOccurred())
			tune3, err = service.CreateTune(apimodel.CreateTune{Title: "tune 3"})
			Expect(err).ShouldNot(HaveOccurred())
			initialMusicSet, err = service.CreateMusicSet(
				apimodel.CreateSet{Title: "test music set"},
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		When("assigning tunes in an arbitrary order "+
			"to the music set", func() {
			BeforeEach(func() {
				musicSetAfterAssignment, err = service.AssignTunesToMusicSet(
					initialMusicSet.ID,
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
		})

	})

})
