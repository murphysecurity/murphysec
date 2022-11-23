package display

type _NONE struct{}

func (_ _NONE) ClearStatus() {}

func (_ _NONE) UpdateStatus(s Status, msg string) {}

func (_ _NONE) Display(level MsgLevel, msg string) {}
