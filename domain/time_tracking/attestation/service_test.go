package attestation

import (
	"context"
	"errors"
	"testing"
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
	"time_management/infrastructure/persistence/memory"
	memoryachievement "time_management/infrastructure/persistence/memory/time_tracking/achievement"
	memorytaskprogress "time_management/infrastructure/persistence/memory/time_tracking/task_progress"
	memoryworkrecord "time_management/infrastructure/persistence/memory/time_tracking/work_record"
)

type recordingChain struct {
	nextUID  [32]byte
	records  map[[32]byte]OnChainRecord
	payloads map[[32]byte]Payload
	readErr  error
}

func (chain *recordingChain) Attest(_ context.Context, payload Payload) ([32]byte, error) {
	chain.nextUID[31]++
	uid := chain.nextUID
	chain.records[uid] = OnChainRecord{
		DayDigest: payload.DayDigest, SubjectCommitment: payload.SubjectCommitment,
	}
	chain.payloads[uid] = payload
	return uid, nil
}

func (chain *recordingChain) Read(_ context.Context, uid [32]byte) (OnChainRecord, error) {
	if chain.readErr != nil {
		return OnChainRecord{}, chain.readErr
	}
	return chain.records[uid], nil
}

func TestServiceVerifyReturnsMismatchAfterRecordEdit(t *testing.T) {
	service, store, chain, date := testService(t)

	attested, err := service.Attest(context.Background(), testUserID, date)
	if err != nil {
		t.Fatal(err)
	}
	verified, err := service.Verify(context.Background(), testUserID, date)
	if err != nil {
		t.Fatal(err)
	}
	if verified.Status != StatusMatch || verified.UID != attested.UID {
		t.Fatalf("Verify() = %#v, want match with UID %x", verified, attested.UID)
	}

	workLog, err := store.FindByWorkLogID(testWorkLogID)
	if err != nil {
		t.Fatal(err)
	}
	workLog.WorkType.WorkTypeName = "edited"
	if err := store.UpdateWorkLog(workLog); err != nil {
		t.Fatal(err)
	}

	verified, err = service.Verify(context.Background(), testUserID, date)
	if err != nil {
		t.Fatal(err)
	}
	if verified.Status != StatusMismatch {
		t.Fatalf("Verify() status = %q, want %q", verified.Status, StatusMismatch)
	}
	_ = chain
}

func TestServiceAttestReferencesPreviousUID(t *testing.T) {
	service, _, chain, date := testService(t)
	first, err := service.Attest(context.Background(), testUserID, date)
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Attest(context.Background(), testUserID, date)
	if err != nil {
		t.Fatal(err)
	}
	if first.UID == second.UID {
		t.Fatal("second Attest() reused the first UID")
	}
	if chain.payloads[second.UID].RefUID != first.UID {
		t.Fatalf("second refUID = %x, want %x", chain.payloads[second.UID].RefUID, first.UID)
	}
}

func TestServiceUsesTokyoBoundaryForBreaks(t *testing.T) {
	service, store, _, date := testService(t)
	breakLog := workrecordmodel.NewBreakLog(
		"00000000-0000-4000-8000-000000000501",
		testWorkLogID,
		testUserID,
		time.Date(2026, time.November, 9, 15, 30, 0, 0, time.UTC), // Nov 10 00:30 JST
		time.Date(2026, time.November, 9, 15, 45, 0, 0, time.UTC),
		15*time.Minute,
		nil,
		nil,
	)
	if err := store.CreateBreakLog(breakLog); err != nil {
		t.Fatal(err)
	}
	records, err := service.dayRecords(testUserID, date)
	if err != nil {
		t.Fatal(err)
	}
	if len(records.BreakLogs) != 1 || records.BreakLogs[0].BreakLogID != breakLog.BreakLogID {
		t.Fatalf("Tokyo day breaks = %#v, want %s", records.BreakLogs, breakLog.BreakLogID)
	}
}

func TestServiceVerifyTreatsDeletedAchievementAsMismatch(t *testing.T) {
	service, store, _, date := testService(t)
	if _, err := service.Attest(context.Background(), testUserID, date); err != nil {
		t.Fatal(err)
	}
	if err := store.DeleteDaily(testAchievementID); err != nil {
		t.Fatal(err)
	}
	result, err := service.Verify(context.Background(), testUserID, date)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusMismatch {
		t.Fatalf("Verify() status = %q, want %q", result.Status, StatusMismatch)
	}
}

func TestServiceVerifyDistinguishesNotAttestedAndRPCFailure(t *testing.T) {
	service, _, chain, date := testService(t)
	result, err := service.Verify(context.Background(), testUserID, date)
	if err != nil || result.Status != StatusNotAttested {
		t.Fatalf("Verify() = %#v, %v; want not-attested", result, err)
	}

	if _, err := service.Attest(context.Background(), testUserID, date); err != nil {
		t.Fatal(err)
	}
	chain.readErr = errors.New("RPC unavailable")
	if _, err := service.Verify(context.Background(), testUserID, date); !errors.Is(err, ErrAttestationFailed) {
		t.Fatalf("Verify() error = %v, want %v", err, ErrAttestationFailed)
	}
}

const (
	testUserID        = "00000000-0000-4000-8000-000000000101"
	testWorkLogID     = "00000000-0000-4000-8000-000000000201"
	testAchievementID = "00000000-0000-4000-8000-000000000401"
)

func testService(t *testing.T) (*Service, *memory.Store, *recordingChain, time.Time) {
	t.Helper()
	t.Setenv(HMACSecretEnv, "0123456789abcdef0123456789abcdef")
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	date := time.Date(2026, time.November, 10, 0, 0, 0, 0, time.UTC)
	store := memory.NewStore()
	workLogs := memoryworkrecord.NewWorkLogRepository(store)
	breakLogs := memoryworkrecord.NewBreakLogRepository(store)
	dailyTasks := memorytaskprogress.NewDailyTaskRepository(store)
	weeklyTasks := memorytaskprogress.NewWeeklyTaskRepository(store)
	achievements := memoryachievement.NewDailyAchievementRepository(store)
	workLog := workrecordmodel.NewWorkLog(
		testWorkLogID, testUserID,
		workrecordmodel.NewWorkType("00000000-0000-4000-8000-000000000301", "build"),
		date, date.Add(time.Hour), date.Add(2*time.Hour), time.Hour,
	)
	if err := workLogs.Create(workLog); err != nil {
		t.Fatal(err)
	}
	achievement := achievementmodel.NewDailyAchievement(
		testAchievementID, testUserID, date,
		time.Hour, 0, 1, 0, date.Add(3*time.Hour),
	)
	if err := achievements.Save(
		achievement,
		[]analysismodel.DailyWorkTypeSummary{analysismodel.NewDailyWorkTypeSummary(
			achievement.DailyAchievementID, workLog.WorkType.WorkTypeID, 1, time.Hour, time.Hour,
		)},
		analysismodel.NewDailyBreakSummary(achievement.DailyAchievementID, 0, 0, 0, 0, 0),
	); err != nil {
		t.Fatal(err)
	}
	chain := &recordingChain{
		records:  map[[32]byte]OnChainRecord{},
		payloads: map[[32]byte]Payload{},
	}
	writer, err := NewWriterFromEnv(chain)
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewService(writer, workLogs, breakLogs, dailyTasks, weeklyTasks, achievements, store, tokyo)
	if err != nil {
		t.Fatal(err)
	}
	return service, store, chain, date
}
