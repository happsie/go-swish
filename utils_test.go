package goswish

import "testing"

func TestGetInstructionID(t *testing.T) {
	actual := "https://mss.cpc.getswish.net/swish-cpcapi/api/v1/paymentrequests/0409C36FD37C4DAD838E227C5FFF3859"
	expected := "0409C36FD37C4DAD838E227C5FFF3859"

	ID, err := GetInstructionID(actual)
	if err != nil {
		t.FailNow()
	}
	if ID != expected {
		t.FailNow()
	}
}