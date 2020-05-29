package main

import "testing"

func Test1(t *testing.T) {
	idata := NewIdata()
	idata.iMin = 0
	err := fileLoad("testdata/interests1.tsv", idata)
	if err != nil {
		t.Errorf("fileLoad failed: %v", err)
	}
	if idata.ignorePeople() != 1 {
		t.Error("expected 1 ignored person")
	} else {
		for id := range idata.pIgnore {
			if idata.pToID["user5"] != id {
				t.Error("user5 should be ignored")
			}
		}
	}
	idata.ignoreInterests()
	if len(idata.iIgnore) != 2 {
		t.Error("should be 2")
	} else {
		for id := range idata.iIgnore {
			if idata.toStr[id] != "bats" && idata.toStr[id] != "animals" {
				t.Error("bats or animals allowed")
			}
		}
	}
	return
}
