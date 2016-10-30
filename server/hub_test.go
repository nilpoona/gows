package main

import (
	"testing"
)

func TestNewHubManager(t *testing.T) {
	hm := newHubManager(5)
	if len(hm.hubs) != 5 {
		t.Errorf("指定したhubの数 %v\n実際に起動した数 %v", 5, len(hm.hubs))
	}
}
