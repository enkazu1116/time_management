package messages

import "net/http"

const (
	ContentTypeHeader = "Content-Type"
	ContentTypeJSON   = "application/json"

	ResponseKeyStatus = "status"
	ResponseKeyError  = "error"

	StatusOK = "正常です"

	DateRequired      = "日付は必須です"
	YearMonthRequired = "年月は必須です"
	DateTimeRequired  = "日時は必須です"

	RequestNameRequired                = "名前は必須です"
	RequestPasswordRequired            = "パスワードは必須です"
	RequestEmailRequired               = "メールアドレスは必須です"
	RequestUserIDRequired              = "ユーザーIDは必須です"
	RequestWorkTypeIDRequired          = "作業種類IDは必須です"
	RequestWorkTypeNameRequired        = "作業種類名は必須です"
	RequestWorkLogIDRequired           = "作業ログIDは必須です"
	RequestWorkDateRequired            = "作業日は必須です"
	RequestStartedAtRequired           = "開始日時は必須です"
	RequestEndedAtRequired             = "終了日時は必須です"
	RequestWeekStartDateRequired       = "週開始日は必須です"
	RequestWeeklyTaskIDRequired        = "週次タスクIDは必須です"
	RequestNegativeDuration            = "時間に負の値は指定できません"
	RequestProgressRateOutOfRange      = "進捗率は0以上1以下で指定してください"
	RequestResumeDurationWithoutResume = "再開までの時間を指定する場合は再開日時も指定してください"
	RequestUnsupportedMonthlySource    = "集計元は未指定またはwork_logsを指定してください"

	RecordNotFound = "対象のデータが見つかりません"
)

type HTTPResultMessage struct {
	StatusCode  int
	ContentType string
	Key         string
	Message     string
}

func OK() HTTPResultMessage {
	return HTTPResultMessage{
		StatusCode:  http.StatusOK,
		ContentType: ContentTypeJSON,
		Key:         ResponseKeyStatus,
		Message:     StatusOK,
	}
}

func Error(statusCode int, message string) HTTPResultMessage {
	return HTTPResultMessage{
		StatusCode:  statusCode,
		ContentType: ContentTypeJSON,
		Key:         ResponseKeyError,
		Message:     message,
	}
}
