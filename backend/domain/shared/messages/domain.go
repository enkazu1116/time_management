package messages

const (
	UUIDRequired = "UUIDは必須です"
	InvalidUUID  = "UUIDの形式が不正です"

	UserIDRequired     = "ユーザーIDは必須です"
	UserNameRequired   = "ユーザー名は必須です"
	UserNameTooLong    = "ユーザー名は100文字以内で入力してください"
	EmailRequired      = "メールアドレスは必須です"
	InvalidEmail       = "メールアドレスの形式が不正です"
	InvalidEmailDomain = "メールアドレスのドメインが不正です"
	RoleRequired       = "権限は必須です"

	WorkTypeRequired       = "作業種類は必須です"
	WorkLogIDRequired      = "作業ログIDは必須です"
	BreakLogIDRequired     = "休憩ログIDは必須です"
	WorkDateRequired       = "作業日は必須です"
	StartedAtRequired      = "開始日時は必須です"
	EndedAtRequired        = "終了日時は必須です"
	StartedAtAfterEndedAt  = "開始日時は終了日時より前である必要があります"
	NegativeWorkDuration   = "作業時間に負の値は指定できません"
	NegativeBreakDuration  = "休憩時間に負の値は指定できません"
	NegativeResumeDuration = "再開までの時間に負の値は指定できません"

	WeeklyTaskIDRequired       = "週次タスクIDは必須です"
	DailyTaskIDRequired        = "日次タスクIDは必須です"
	DailyLogIDRequired         = "日次ログIDは必須です"
	WeeklyLogIDRequired        = "週次ログIDは必須です"
	WeekStartDateRequired      = "週開始日は必須です"
	TargetDurationNegative     = "目標時間に負の値は指定できません"
	TotalWorkDurationNegative  = "合計作業時間に負の値は指定できません"
	ActualDurationNegative     = "実績時間に負の値は指定できません"
	ProgressRateOutOfRange     = "進捗率は0以上1以下で指定してください"
	AchievementRateOutOfRange  = "達成率に負の値は指定できません"
	AverageProgressOutOfRange  = "平均進捗率は0以上1以下で指定してください"
	WeeklyLogUpdatedAtRequired = "週次ログの更新日時は必須です"

	DailyAchievementIDRequired   = "日次実績IDは必須です"
	WeeklyLogIDPrefixRequired    = "週次ログIDプレフィックスは必須です"
	MonthlyAchievementIDRequired = "月次実績IDは必須です"
	YearMonthRequired            = "年月は必須です"
	UpdatedAtRequired            = "更新日時は必須です"
)
