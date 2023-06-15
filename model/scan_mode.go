package model

type ScanMode string

const (
	ScanModeSource     ScanMode = "source"
	ScanModeSourceDeep ScanMode = "source_deep"
	ScanModeBinary     ScanMode = "binary"
	ScanModeIot        ScanMode = "iot"
	ScanModeStandard   ScanMode = "standard"
	ScanModeComp       ScanMode = "comp"
)
