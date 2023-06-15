package ui

type iDEA struct{}

func (iDEA) UpdateStatus(s Status, msg string) {}

func (iDEA) Display(level MessageLevel, msg string) {}

func (iDEA) ClearStatus() {}

var _ UI = (*iDEA)(nil)

var IDEA UI = &iDEA{}
