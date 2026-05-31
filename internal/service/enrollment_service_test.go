package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
)

func setupPostgres(t *testing.T) *infra.Postgresql {
	cfg := &config.Config{
		DatabaseHost:    "localhost",
		DatabasePort:    "5432",
		DatabaseUser:    "postgres",
		DatabasePass:    "althea",
		DatabaseName:    "healthcare_gov",
		DatabaseSSLMode: "disable",
	}

	pg := infra.NewPostgresql(cfg)
	if pg.Db == nil {
		t.Fatal("Database not initialized")
	}

	return pg
}

func TestEnrollmentService_EnrollPlan_RaceCondition(t *testing.T) {
	pg := setupPostgres(t)
	defer pg.Db.Close()

	repo := repository.NewEnrollmentRepository(pg)
	service := &EnrollmentServiceImpl{
		repo: repo,
	}

	ctx := context.Background()
	input := EnrollPlanInput{
		UserID: "97e7f136-20da-42eb-8f16-2f421a14c5a6",
		PlanID: "01107fe3-791c-48ae-a027-c360ce9448cf",
	}

	const concurrency = 10
	var wg sync.WaitGroup

	startRaceGate := make(chan struct{})
	errorChan := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			<-startRaceGate

			_, err := service.EnrollPlan(ctx, input)
			if err != nil {
				errorChan <- err
			}
		}()
	}

	close(startRaceGate)
	wg.Wait()
	close(errorChan)

	var successCount int
	var conflictCount int

	for err := range errorChan {
		if errors.Is(err, ErrUserHasActivePlan) {
			conflictCount++
		}
	}

	successCount = concurrency - conflictCount

	if successCount > 1 {
		t.Log("Success count", successCount)
	}

	if conflictCount > 1 {
		t.Errorf("Conflict count %d", conflictCount)
	}
}
