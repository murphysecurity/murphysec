package api

type TaskKind string

const (
	TaskKindNormal  TaskKind = "Normal"
	TaskKindBinary  TaskKind = "Binary"
	TaskKindIotScan TaskKind = "IotScan"
)
