package log_test

import (
	"testing"

	"github.com/oa-pass/pass-tools/lib/log"
)

func TestDefaultLogging(t *testing.T) {
	log.Instance{}.Printf("testing")
}

func TestNilDebug(t *testing.T) {
	var l log.Instance
	if l.Debug != nil {
		t.Fatalf("Debug logger should be nil")
	}
}
