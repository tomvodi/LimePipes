package health

import (
	"context"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database"
	"sync"
	"testing"
	"time"
)

func TestNewHealthCheck(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(1)

	g := NewGomegaWithT(t)
	cfg, err := config.InitTest()
	g.Expect(err).ShouldNot(HaveOccurred())

	gormDb, err := database.GetInitTestPostgreSQLDB(cfg.DbConfig(), "testdb")
	g.Expect(err).ShouldNot(HaveOccurred())

	check := NewHealthCheck(gormDb)

	go func() {
		defer wg.Done()
		checker, err := check.GetChecker()
		if err != nil {
			log.Error().Err(err)
			return
		}

		log.Info().Msgf("checker status %v",
			checker.Check(context.Background()),
		)

		checker.Start()

		time.Sleep(180 * time.Second)

		log.Info().Msg("finished check routine")
	}()

	wg.Wait()
}
