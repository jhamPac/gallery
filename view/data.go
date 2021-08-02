package view

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	AlertMsgGeneric = "An error occured, please try again. If the problem persists please contact us."
)

type Data struct {
	Alert   *Alert
	Content interface{}
}

type Alert struct {
	Level   string
	Message string
}
