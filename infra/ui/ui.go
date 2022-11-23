package ui

type UI interface {
	UpdateStatus(s Status, msg string)
	Display(level MessageLevel, msg string)
	ClearStatus()
}
