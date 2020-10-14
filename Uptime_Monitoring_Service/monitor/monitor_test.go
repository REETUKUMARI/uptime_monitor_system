package monitor

import (
	"testing"
	"time"
)

func TestCheckurl(t *testing.T) {
	link := "http://google.com"
	var tm time.Duration = 20
	status := Checkurl(link, tm)
	if status != "active" {
		t.Errorf("expected %v and got %v", "active", status)
	}
	link = "http://abc.com"
	status = Checkurl(link, tm)
	if status != "active" {
		t.Errorf("expected %v and got %v", "inactive", status)
	}
}
