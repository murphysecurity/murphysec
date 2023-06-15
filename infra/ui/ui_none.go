package ui

type none struct{}

func (none) UpdateStatus(s Status, msg string) {}

func (none) Display(level MessageLevel, msg string) {}

func (none) ClearStatus() {}

var _ UI = (*none)(nil)

var None UI = &none{}
