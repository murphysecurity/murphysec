package v5

import _ "embed"

//go:embed testdata/1.yaml
var testData1 string

//go:embed testdata/2.yaml
var testData2 string

//go:embed testdata/3.yaml
var testData3 string

//go:embed testdata/4.yaml
var testData4 string

var testDataList = []string{
	testData1,
	testData2,
	testData3,
	testData4,
}
