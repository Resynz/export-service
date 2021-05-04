package notify

type NotifyResult string

const (
	NotifyResultSuccess NotifyResult = "success"
	NotifyResultFailed  NotifyResult = "failed"
)

type NotifyParam struct {
	NotifyUrl string
	Result    NotifyResult
	Message   string
}
