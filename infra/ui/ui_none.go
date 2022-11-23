package ui

type None struct{}

func (None) UpdateStatus(s Status, msg string) {}

func (None) Display(level MessageLevel, msg string) {}

func (None) ClearStatus() {}

var _ UI = (*None)(nil)
