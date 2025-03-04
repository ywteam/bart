package scan

import "context"


type IScanner interface {
	Scan(ctx context.Context, assets ...interface{}) ([]*ScannerResult, error)	
}
type Scanner struct {	
	IScanner
}

type IScannerResult interface {
	GetMetrics() map[string]interface{}
	GetAsset() interface{}
	GetError() error
	HasError() bool
}

type ScannerResult struct {
	IScannerResult
}