package model

type ScanStrategy string

const (
	ScanStrategyBackup ScanStrategy = "backup"
	ScanStrategyNormal ScanStrategy = "normal"
)

func (s ScanStrategy) MarshalText() ([]byte, error) {
	return []byte(ScanStrategyNormal), nil
}
