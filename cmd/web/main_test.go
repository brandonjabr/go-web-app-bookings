package main

import "testing"

func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Errorf("failed to run application")
	}
}
