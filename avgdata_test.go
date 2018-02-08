package main

import "testing"

func TestNewAvgData(t *testing.T) {
	avgdata, err := NewAvgData("testdata/simple.csv")
	if err != nil {
		t.Error(err)
	}
	if len(avgdata.Records) != 1 {
		t.Errorf("expected 1 record")
	}
	if avgdata.Records[0].Load != 1 {
		t.Errorf("unexpected record value")
	}
}
