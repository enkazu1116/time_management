package handlers

import (
	"net/http"
	"time"
	trackingports "time_management/domain/time_tracking/ports"
	"time_management/infrastructure/http/requests"
	"time_management/infrastructure/http/responses"
	"time_management/infrastructure/util/decoder"
	"time_management/infrastructure/util/id"
	"time_management/infrastructure/util/parser"
	"time_management/infrastructure/util/response"
)

type TimeTrackingHandler struct {
	usecase trackingports.TimeTrackingUsecase
}

func NewTimeTrackingHandler(u trackingports.TimeTrackingUsecase, _ ...any) *TimeTrackingHandler {
	return &TimeTrackingHandler{usecase: u}
}

func (h *TimeTrackingHandler) CreateWorkLog(w http.ResponseWriter, r *http.Request) {
	var q requests.CreateWorkLogRequest
	if !decode(w, r, &q) || !valid(w, q.Validate()) {
		return
	}
	d, e := parser.ParseDate(q.WorkDate)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	s, e := parser.ParseDateTime(q.StartedAt)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	f, e := parser.ParseDateTime(q.EndedAt)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	if q.WorkLogID == "" {
		q.WorkLogID = id.NewID()
	}
	v, e := h.usecase.CreateWorkLog(trackingports.CreateWorkLogInput{q.WorkLogID, q.UserID, q.WorkTypeID, q.WorkTypeName, d, s, f, time.Duration(q.WorkDurationSeconds) * time.Second})
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	response.WriteJSON(w, 201, responses.WorkLog(v))
}
func (h *TimeTrackingHandler) CreateBreakLog(w http.ResponseWriter, r *http.Request) {
	var q requests.CreateBreakLogRequest
	if !decode(w, r, &q) || !valid(w, q.Validate()) {
		return
	}
	s, e := parser.ParseDateTime(q.StartedAt)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	f, e := parser.ParseDateTime(q.EndedAt)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	var ra *time.Time
	if q.ResumedAt != "" {
		v, e := parser.ParseDateTime(q.ResumedAt)
		if e != nil {
			response.WriteError(w, 400, e)
			return
		}
		ra = &v
	}
	var rd *time.Duration
	if q.ResumeDurationSeconds != nil {
		v := time.Duration(*q.ResumeDurationSeconds) * time.Second
		rd = &v
	}
	if q.BreakLogID == "" {
		q.BreakLogID = id.NewID()
	}
	v, e := h.usecase.CreateBreakLog(trackingports.CreateBreakLogInput{q.BreakLogID, q.WorkLogID, q.UserID, s, f, time.Duration(q.BreakDurationSeconds) * time.Second, ra, rd})
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	response.WriteJSON(w, 201, responses.BreakLog(v))
}
func (h *TimeTrackingHandler) CreateWeeklyTask(w http.ResponseWriter, r *http.Request) {
	var q requests.CreateWeeklyTaskRequest
	if !decode(w, r, &q) || !valid(w, q.Validate()) {
		return
	}
	d, e := parser.ParseDate(q.WeekStartDate)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	if q.WeeklyTaskID == "" {
		q.WeeklyTaskID = id.NewID()
	}
	v, e := h.usecase.CreateWeeklyTask(trackingports.CreateWeeklyTaskInput{q.WeeklyTaskID, q.UserID, q.WorkTypeID, d, time.Duration(q.TargetDurationSeconds) * time.Second, q.GoalDescription})
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	response.WriteJSON(w, 201, responses.WeeklyTask(v))
}
func (h *TimeTrackingHandler) CreateDailyTask(w http.ResponseWriter, r *http.Request) {
	var q requests.CreateDailyTaskRequest
	if !decode(w, r, &q) || !valid(w, q.Validate()) {
		return
	}
	d, e := parser.ParseDate(q.WorkDate)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	if q.DailyTaskID == "" {
		q.DailyTaskID = id.NewID()
	}
	v, e := h.usecase.CreateDailyTask(trackingports.CreateDailyTaskInput{q.DailyTaskID, q.WeeklyTaskID, q.UserID, d, q.ProgressDescription, q.ProgressRate})
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	response.WriteJSON(w, 201, responses.DailyTask(v))
}
func (h *TimeTrackingHandler) CreateDailyLog(w http.ResponseWriter, r *http.Request) {
	var q requests.CreateDailyLogRequest
	if !decode(w, r, &q) || !valid(w, q.Validate()) {
		return
	}
	d, e := parser.ParseDate(q.WorkDate)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	if q.DailyLogID == "" {
		q.DailyLogID = id.NewID()
	}
	v, e := h.usecase.CreateDailyLog(trackingports.CreateDailyLogInput{q.DailyLogID, q.UserID, d, time.Duration(q.TotalWorkDurationSeconds) * time.Second})
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	response.WriteJSON(w, 201, responses.DailyLog(v))
}
func (h *TimeTrackingHandler) CreateDailySummary(w http.ResponseWriter, r *http.Request) {
	var q requests.CreateDailySummaryRequest
	if !decode(w, r, &q) || !valid(w, q.Validate()) {
		return
	}
	d, e := parser.ParseDate(q.WorkDate)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	if q.DailyAchievementID == "" {
		q.DailyAchievementID = id.NewID()
	}
	v, e := h.usecase.CreateDailySummary(trackingports.CreateDailySummaryInput{q.DailyAchievementID, q.UserID, d, time.Now()})
	if e != nil {
		response.WriteError(w, 500, e)
		return
	}
	response.WriteJSON(w, 201, responses.Summary(v))
}
func (h *TimeTrackingHandler) CreateWeeklyLog(w http.ResponseWriter, r *http.Request) {
	var q requests.CreateWeeklyLogRequest
	if !decode(w, r, &q) || !valid(w, q.Validate()) {
		return
	}
	d, e := parser.ParseDate(q.WeekStartDate)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	if q.WeeklyLogIDPrefix == "" {
		q.WeeklyLogIDPrefix = id.NewID()
	}
	v, e := h.usecase.CreateWeeklyLog(trackingports.CreateWeeklyLogInput{q.WeeklyLogIDPrefix, q.UserID, d, time.Now()})
	if e != nil {
		response.WriteError(w, 500, e)
		return
	}
	response.WriteJSON(w, 201, responses.WeeklyLogSet(v))
}
func (h *TimeTrackingHandler) CreateMonthlySummary(w http.ResponseWriter, r *http.Request) {
	var q requests.CreateMonthlySummaryRequest
	if !decode(w, r, &q) || !valid(w, q.Validate()) {
		return
	}
	m, e := parser.ParseYearMonth(q.YearMonth)
	if e != nil {
		response.WriteError(w, 400, e)
		return
	}
	if q.MonthlyAchievementID == "" {
		q.MonthlyAchievementID = id.NewID()
	}
	v, e := h.usecase.CreateMonthlySummary(trackingports.CreateMonthlySummaryInput{q.MonthlyAchievementID, q.UserID, m, time.Now()}, q.Source == "work_logs")
	if e != nil {
		response.WriteError(w, 500, e)
		return
	}
	response.WriteJSON(w, 201, responses.Summary(v))
}
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if e := decoder.DecodeJSON(r, v); e != nil {
		response.WriteError(w, 400, e)
		return false
	}
	return true
}
func valid(w http.ResponseWriter, e error) bool {
	if e != nil {
		response.WriteError(w, 400, e)
		return false
	}
	return true
}
