package model

type ScanMode string

const (
	ScanModeSource ScanMode = "source"
	ScanModeBinary ScanMode = "binary"
	ScanModeIot    ScanMode = "iot"
)
