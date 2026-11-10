package integrity

import (
	"encoding/hex"
	"slices"
	"testing"
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
)

func TestDigestDayIgnoresRecordOrderAndLocalTimeZone(t *testing.T) {
	original := fixture()
	reordered := clone(original)
	slices.Reverse(reordered.WorkLogs)
	slices.Reverse(reordered.BreakLogs)
	slices.Reverse(reordered.DailyTasks)
	slices.Reverse(reordered.WeeklyTasks)
	slices.Reverse(reordered.DailyAchievements)
	for i := range reordered.DailyAchievements {
		slices.Reverse(reordered.DailyAchievements[i].WorkTypeSummaries)
	}
	reordered.Day = reordered.Day.In(time.FixedZone("local", 9*60*60))
	for i := range reordered.WorkLogs {
		reordered.WorkLogs[i].StartedAt = reordered.WorkLogs[i].StartedAt.In(time.FixedZone("local", -7*60*60))
	}

	got := mustDigest(t, reordered)
	want := mustDigest(t, original)
	if got != want {
		t.Fatalf("DigestDay() changed with order or location:\n got %x\nwant %x", got, want)
	}
}

func TestDigestDayUsesRemainingFieldsAsSortTieBreakers(t *testing.T) {
	original := fixture()
	original.WorkLogs[1].WorkLogID = original.WorkLogs[0].WorkLogID
	reordered := clone(original)
	slices.Reverse(reordered.WorkLogs)

	if got, want := mustDigest(t, reordered), mustDigest(t, original); got != want {
		t.Fatalf("DigestDay() changed when equal-ID records were reordered:\n got %x\nwant %x", got, want)
	}
}

func TestDigestDayChangesWhenAFieldChanges(t *testing.T) {
	oneNanosecond := time.Nanosecond
	oneDurationUnit := time.Duration(1)
	changes := []struct {
		name   string
		mutate func(*DayRecords)
	}{
		{"day", func(r *DayRecords) {
			r.Day = r.Day.AddDate(0, 0, 1)
			for i := range r.WorkLogs {
				r.WorkLogs[i].WorkDate = r.WorkLogs[i].WorkDate.AddDate(0, 0, 1)
			}
			for i := range r.BreakLogs {
				r.BreakLogs[i].StartedAt = r.BreakLogs[i].StartedAt.AddDate(0, 0, 1)
				r.BreakLogs[i].EndedAt = r.BreakLogs[i].EndedAt.AddDate(0, 0, 1)
			}
			for i := range r.DailyTasks {
				r.DailyTasks[i].WorkDate = r.DailyTasks[i].WorkDate.AddDate(0, 0, 1)
			}
			for i := range r.DailyAchievements {
				r.DailyAchievements[i].Achievement.WorkDate = r.DailyAchievements[i].Achievement.WorkDate.AddDate(0, 0, 1)
			}
		}},
		{"work log ID", func(r *DayRecords) { r.WorkLogs[0].WorkLogID += "x" }},
		{"work type ID", func(r *DayRecords) { r.WorkLogs[0].WorkType.WorkTypeID += "x" }},
		{"work type name", func(r *DayRecords) { r.WorkLogs[0].WorkType.WorkTypeName += "x" }},
		{"work date timestamp", func(r *DayRecords) { r.WorkLogs[0].WorkDate = r.WorkLogs[0].WorkDate.Add(oneNanosecond) }},
		{"work start timestamp by one nanosecond", func(r *DayRecords) { r.WorkLogs[0].StartedAt = r.WorkLogs[0].StartedAt.Add(oneNanosecond) }},
		{"work end timestamp", func(r *DayRecords) { r.WorkLogs[0].EndedAt = r.WorkLogs[0].EndedAt.Add(oneNanosecond) }},
		{"work duration", func(r *DayRecords) { r.WorkLogs[0].WorkDuration += oneDurationUnit }},
		{"break ID", func(r *DayRecords) { r.BreakLogs[0].BreakLogID += "x" }},
		{"break work log ID", func(r *DayRecords) { r.BreakLogs[0].WorkLogID += "x" }},
		{"break start timestamp", func(r *DayRecords) { r.BreakLogs[0].StartedAt = r.BreakLogs[0].StartedAt.Add(oneNanosecond) }},
		{"break end timestamp", func(r *DayRecords) { r.BreakLogs[0].EndedAt = r.BreakLogs[0].EndedAt.Add(oneNanosecond) }},
		{"break duration", func(r *DayRecords) { r.BreakLogs[0].BreakDuration += oneDurationUnit }},
		{"resumed timestamp", func(r *DayRecords) { *r.BreakLogs[0].ResumedAt = r.BreakLogs[0].ResumedAt.Add(oneNanosecond) }},
		{"resume duration", func(r *DayRecords) { *r.BreakLogs[0].ResumeDuration += oneDurationUnit }},
		{"optional resumed timestamp presence", func(r *DayRecords) { r.BreakLogs[0].ResumedAt = nil }},
		{"optional resume duration presence", func(r *DayRecords) { r.BreakLogs[0].ResumeDuration = nil }},
		{"daily task ID", func(r *DayRecords) { r.DailyTasks[0].DailyTaskID += "x" }},
		{"daily task weekly ID", func(r *DayRecords) { r.DailyTasks[0].WeeklyTaskID += "x" }},
		{"daily task work date", func(r *DayRecords) { r.DailyTasks[0].WorkDate = r.DailyTasks[0].WorkDate.Add(oneNanosecond) }},
		{"progress description", func(r *DayRecords) { r.DailyTasks[0].ProgressDescription += "x" }},
		{"progress rate", func(r *DayRecords) { r.DailyTasks[0].ProgressRate += 0.01 }},
		{"weekly task ID", func(r *DayRecords) { r.WeeklyTasks[0].WeeklyTaskID += "x" }},
		{"weekly work type ID", func(r *DayRecords) { r.WeeklyTasks[0].WorkTypeID += "x" }},
		{"week start timestamp", func(r *DayRecords) {
			r.WeeklyTasks[0].WeekStartDate = r.WeeklyTasks[0].WeekStartDate.Add(oneNanosecond)
		}},
		{"target duration", func(r *DayRecords) { r.WeeklyTasks[0].TargetDuration += oneDurationUnit }},
		{"goal description", func(r *DayRecords) { r.WeeklyTasks[0].GoalDescription += "x" }},
		{"achievement ID", func(r *DayRecords) { r.DailyAchievements[0].Achievement.DailyAchievementID += "x" }},
		{"achievement work date", func(r *DayRecords) {
			r.DailyAchievements[0].Achievement.WorkDate = r.DailyAchievements[0].Achievement.WorkDate.Add(oneNanosecond)
		}},
		{"total work duration", func(r *DayRecords) { r.DailyAchievements[0].Achievement.TotalWorkDuration += oneDurationUnit }},
		{"total break duration", func(r *DayRecords) { r.DailyAchievements[0].Achievement.TotalBreakDuration += oneDurationUnit }},
		{"work count", func(r *DayRecords) { r.DailyAchievements[0].Achievement.WorkCount++ }},
		{"break count", func(r *DayRecords) { r.DailyAchievements[0].Achievement.BreakCount++ }},
		{"achievement update timestamp", func(r *DayRecords) {
			r.DailyAchievements[0].Achievement.UpdatedAt = r.DailyAchievements[0].Achievement.UpdatedAt.Add(oneNanosecond)
		}},
		{"summary achievement ID", func(r *DayRecords) { r.DailyAchievements[0].WorkTypeSummaries[0].DailyAchievementID += "x" }},
		{"summary work type ID", func(r *DayRecords) { r.DailyAchievements[0].WorkTypeSummaries[0].WorkTypeID += "x" }},
		{"summary work count", func(r *DayRecords) { r.DailyAchievements[0].WorkTypeSummaries[0].WorkCount++ }},
		{"summary work duration", func(r *DayRecords) { r.DailyAchievements[0].WorkTypeSummaries[0].WorkDuration += oneDurationUnit }},
		{"summary average work duration", func(r *DayRecords) {
			r.DailyAchievements[0].WorkTypeSummaries[0].AverageWorkDuration += oneDurationUnit
		}},
		{"break summary achievement ID", func(r *DayRecords) { r.DailyAchievements[0].BreakSummary.DailyAchievementID += "x" }},
		{"break summary count", func(r *DayRecords) { r.DailyAchievements[0].BreakSummary.BreakCount++ }},
		{"break summary duration", func(r *DayRecords) { r.DailyAchievements[0].BreakSummary.BreakDuration += oneDurationUnit }},
		{"average break duration", func(r *DayRecords) { r.DailyAchievements[0].BreakSummary.AverageBreakDuration += oneDurationUnit }},
		{"average resume duration", func(r *DayRecords) { r.DailyAchievements[0].BreakSummary.AverageResumeDuration += oneDurationUnit }},
		{"longest resume duration", func(r *DayRecords) { r.DailyAchievements[0].BreakSummary.LongestResumeDuration += oneDurationUnit }},
	}

	want := mustDigest(t, fixture())
	for _, change := range changes {
		t.Run(change.name, func(t *testing.T) {
			changed := clone(fixture())
			change.mutate(&changed)
			if got := mustDigest(t, changed); got == want {
				t.Fatalf("DigestDay() did not change after changing %s", change.name)
			}
		})
	}
}

func TestDigestDayIncludesUserAndCollectionMembership(t *testing.T) {
	emptyA := DayRecords{UserID: "user-a", Day: time.Date(2026, time.November, 10, 0, 0, 0, 0, time.UTC)}
	emptyB := emptyA
	emptyB.UserID = "user-b"
	if mustDigest(t, emptyA) == mustDigest(t, emptyB) {
		t.Fatal("DigestDay() did not change with the user ID")
	}

	withRecords := fixture()
	withoutWorkLog := clone(withRecords)
	withoutWorkLog.WorkLogs = withoutWorkLog.WorkLogs[1:]
	if mustDigest(t, withRecords) == mustDigest(t, withoutWorkLog) {
		t.Fatal("DigestDay() did not change when a record was removed")
	}
}

func TestDigestDayGoldenVector(t *testing.T) {
	const want = "c6e61fb4cdc39e8228a3760ea09fa62498169e88c5841e288c5d654347b4d696"
	digest := mustDigest(t, fixture())
	got := hex.EncodeToString(digest[:])
	if got != want {
		t.Fatalf("DigestDay() = %s, want %s", got, want)
	}
}

func TestDigestDayEmptyDay(t *testing.T) {
	records := DayRecords{
		UserID: "00000000-0000-0000-0000-000000000001",
		Day:    time.Date(2026, time.November, 10, 18, 30, 0, 0, time.FixedZone("JST", 9*60*60)),
	}
	const want = "74b1144d06b9e8bc1440f1e74c0a365f57dc97f6734998c3be7030f7f11f8b55"
	digest := mustDigest(t, records)
	got := hex.EncodeToString(digest[:])
	if got != want {
		t.Fatalf("empty DigestDay() = %s, want %s", got, want)
	}

	serialized, err := SerializeDay(records)
	if err != nil {
		t.Fatal(err)
	}
	if len(serialized) == 0 {
		t.Fatal("SerializeDay() returned no versioned representation for an empty day")
	}
}

func mustDigest(t *testing.T, records DayRecords) [32]byte {
	t.Helper()
	digest, err := DigestDay(records)
	if err != nil {
		t.Fatal(err)
	}
	return digest
}

func fixture() DayRecords {
	day := time.Date(2026, time.November, 10, 0, 0, 0, 0, time.UTC)
	resumedAt := day.Add(11*time.Hour + 15*time.Minute)
	resumeDuration := 15 * time.Minute
	records := DayRecords{
		UserID: "00000000-0000-0000-0000-000000000001",
		Day:    day.Add(12 * time.Hour),
		WorkLogs: []workrecordmodel.WorkLog{{
			WorkLogID: "work-b", UserID: "00000000-0000-0000-0000-000000000001",
			WorkType: workrecordmodel.WorkType{WorkTypeID: "type-b", WorkTypeName: "Build"},
			WorkDate: day, StartedAt: day.Add(9 * time.Hour), EndedAt: day.Add(12 * time.Hour),
			WorkDuration: 3 * time.Hour,
		}, {
			WorkLogID: "work-a", UserID: "00000000-0000-0000-0000-000000000001",
			WorkType: workrecordmodel.WorkType{WorkTypeID: "type-a", WorkTypeName: "Plan"},
			WorkDate: day, StartedAt: day.Add(8 * time.Hour), EndedAt: day.Add(9 * time.Hour),
			WorkDuration: time.Hour,
		}},
		BreakLogs: []workrecordmodel.BreakLog{{
			BreakLogID: "break-b", WorkLogID: "work-b", UserID: "00000000-0000-0000-0000-000000000001",
			StartedAt: day.Add(11 * time.Hour), EndedAt: resumedAt, BreakDuration: 15 * time.Minute,
			ResumedAt: &resumedAt, ResumeDuration: &resumeDuration,
		}, {
			BreakLogID: "break-a", WorkLogID: "work-a", UserID: "00000000-0000-0000-0000-000000000001",
			StartedAt: day.Add(8*time.Hour + 30*time.Minute), EndedAt: day.Add(8*time.Hour + 35*time.Minute),
			BreakDuration: 5 * time.Minute,
		}},
		DailyTasks: []taskprogressmodel.DailyTask{{
			DailyTaskID: "daily-b", WeeklyTaskID: "weekly-b", UserID: "00000000-0000-0000-0000-000000000001",
			WorkDate: day, ProgressDescription: "implemented", ProgressRate: 0.75,
		}, {
			DailyTaskID: "daily-a", WeeklyTaskID: "weekly-a", UserID: "00000000-0000-0000-0000-000000000001",
			WorkDate: day, ProgressDescription: "planned", ProgressRate: 0.25,
		}},
		WeeklyTasks: []taskprogressmodel.WeeklyTask{{
			WeeklyTaskID: "weekly-b", UserID: "00000000-0000-0000-0000-000000000001", WorkTypeID: "type-b",
			WeekStartDate: day.AddDate(0, 0, -1), TargetDuration: 20 * time.Hour, GoalDescription: "ship",
		}, {
			WeeklyTaskID: "weekly-a", UserID: "00000000-0000-0000-0000-000000000001", WorkTypeID: "type-a",
			WeekStartDate: day.AddDate(0, 0, -1), TargetDuration: 10 * time.Hour, GoalDescription: "design",
		}},
		DailyAchievements: []analysismodel.DailySummarySet{{
			Achievement: achievementmodel.DailyAchievement{
				DailyAchievementID: "achievement-a", UserID: "00000000-0000-0000-0000-000000000001",
				WorkDate: day, TotalWorkDuration: 4 * time.Hour, TotalBreakDuration: 20 * time.Minute,
				WorkCount: 2, BreakCount: 2, UpdatedAt: day.Add(20 * time.Hour),
			},
			WorkTypeSummaries: []analysismodel.DailyWorkTypeSummary{{
				DailyAchievementID: "achievement-a", WorkTypeID: "type-b", WorkCount: 1,
				WorkDuration: 3 * time.Hour, AverageWorkDuration: 3 * time.Hour,
			}, {
				DailyAchievementID: "achievement-a", WorkTypeID: "type-a", WorkCount: 1,
				WorkDuration: time.Hour, AverageWorkDuration: time.Hour,
			}},
			BreakSummary: analysismodel.DailyBreakSummary{
				DailyAchievementID: "achievement-a", BreakCount: 2, BreakDuration: 20 * time.Minute,
				AverageBreakDuration: 10 * time.Minute, AverageResumeDuration: 15 * time.Minute,
				LongestResumeDuration: 15 * time.Minute,
			},
		}},
	}
	return records
}

func clone(records DayRecords) DayRecords {
	cloned := records
	cloned.WorkLogs = slices.Clone(records.WorkLogs)
	cloned.BreakLogs = slices.Clone(records.BreakLogs)
	for i := range cloned.BreakLogs {
		if records.BreakLogs[i].ResumedAt != nil {
			value := *records.BreakLogs[i].ResumedAt
			cloned.BreakLogs[i].ResumedAt = &value
		}
		if records.BreakLogs[i].ResumeDuration != nil {
			value := *records.BreakLogs[i].ResumeDuration
			cloned.BreakLogs[i].ResumeDuration = &value
		}
	}
	cloned.DailyTasks = slices.Clone(records.DailyTasks)
	cloned.WeeklyTasks = slices.Clone(records.WeeklyTasks)
	cloned.DailyAchievements = slices.Clone(records.DailyAchievements)
	for i := range cloned.DailyAchievements {
		cloned.DailyAchievements[i].WorkTypeSummaries = slices.Clone(records.DailyAchievements[i].WorkTypeSummaries)
	}
	return cloned
}
