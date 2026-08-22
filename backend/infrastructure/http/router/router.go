package router

import (
	"net/http"

	"time_management/infrastructure/http/handlers"
)

type Handlers struct {
	Health       *handlers.HealthHandler
	User         *handlers.UserHandler
	TimeTracking *handlers.TimeTrackingHandler
}

func New(handlers Handlers) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Health.Health)
	mux.HandleFunc("POST /users", handlers.User.RegisterUser)
	mux.HandleFunc("POST /work-logs", handlers.TimeTracking.CreateWorkLog)
	mux.HandleFunc("POST /break-logs", handlers.TimeTracking.CreateBreakLog)
	mux.HandleFunc("POST /weekly-tasks", handlers.TimeTracking.CreateWeeklyTask)
	mux.HandleFunc("POST /daily-tasks", handlers.TimeTracking.CreateDailyTask)
	mux.HandleFunc("POST /daily-logs", handlers.TimeTracking.CreateDailyLog)
	mux.HandleFunc("POST /summaries/daily", handlers.TimeTracking.CreateDailySummary)
	mux.HandleFunc("POST /weekly-logs", handlers.TimeTracking.CreateWeeklyLog)
	mux.HandleFunc("POST /summaries/monthly", handlers.TimeTracking.CreateMonthlySummary)
	return mux
}
