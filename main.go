package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"time_management/domain/time_tracking/attestation"
	trackingusecases "time_management/domain/time_tracking/usecases"
	userusecases "time_management/domain/user/usecases"
	easchain "time_management/infrastructure/blockchain/eas"
	httphandlers "time_management/infrastructure/http/handlers"
	httprouter "time_management/infrastructure/http/router"
	memory "time_management/infrastructure/persistence/memory"
	memoryachievement "time_management/infrastructure/persistence/memory/time_tracking/achievement"
	memorytaskprogress "time_management/infrastructure/persistence/memory/time_tracking/task_progress"
	memoryworkrecord "time_management/infrastructure/persistence/memory/time_tracking/work_record"
)

func main() {
	store := memory.NewStore()

	workLogRepository := memoryworkrecord.NewWorkLogRepository(store)
	breakLogRepository := memoryworkrecord.NewBreakLogRepository(store)
	weeklyTaskRepository := memorytaskprogress.NewWeeklyTaskRepository(store)
	dailyTaskRepository := memorytaskprogress.NewDailyTaskRepository(store)
	dailyLogRepository := memorytaskprogress.NewDailyLogRepository(store)
	weeklyLogRepository := memorytaskprogress.NewWeeklyLogRepository(store)
	dailyAchievementRepository := memoryachievement.NewDailyAchievementRepository(store)
	monthlyAchievementRepository := memoryachievement.NewMonthlyAchievementRepository(store)

	chain, err := easchain.NewBaseSepoliaFromEnv(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer chain.Close()
	attestationWriter, err := attestation.NewWriterFromEnv(chain)
	if err != nil {
		log.Fatal(err)
	}
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	attestationService, err := attestation.NewService(
		attestationWriter,
		workLogRepository,
		breakLogRepository,
		dailyTaskRepository,
		weeklyTaskRepository,
		dailyAchievementRepository,
		store,
		tokyo,
	)
	if err != nil {
		log.Fatal(err)
	}

	userUsecase := userusecases.NewUserUsecase(store)
	timeTrackingUsecase := trackingusecases.NewService(
		workLogRepository,
		breakLogRepository,
		weeklyTaskRepository,
		dailyTaskRepository,
		dailyLogRepository,
		weeklyLogRepository,
		dailyAchievementRepository,
		monthlyAchievementRepository,
	)
	_ = dailyTaskRepository

	handlerSet := httprouter.Handlers{
		Health: httphandlers.NewHealthHandler(),
		User:   httphandlers.NewUserHandler(userUsecase),
		TimeTracking: httphandlers.NewTimeTrackingHandler(
			timeTrackingUsecase,
		),
		Attestation: httphandlers.NewAttestationHandler(attestationService),
	}

	addr := ":" + env("PORT", "8080")
	log.Printf("time-management-api listening on %s", addr)
	if err := http.ListenAndServe(addr, httprouter.New(handlerSet)); err != nil {
		log.Fatal(err)
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
