package vantrue

import (
	"strings"
	"testing"
	"time"
)

func TestParseLineUTCAndKnotSpeed(t *testing.T) {
	const line = "20260421145509,51.311962,N,12.323399,E,10.5,131.9"
	opt := ParseOptions{SpeedColumn: SpeedColumnKnots}
	p, err := parseRecord(strings.Split(line, ","), opt)
	if err != nil {
		t.Fatal(err)
	}
	if !p.Time.Equal(time.Date(2026, 4, 21, 14, 55, 9, 0, time.UTC)) {
		t.Fatalf("time %v", p.Time)
	}
	if p.Lat != 51.311962 || p.Lon != 12.323399 {
		t.Fatalf("lat/lon")
	}
	want := 10.5 * knotsToKMH
	if p.Speed != want || p.Course != 131.9 {
		t.Fatalf("speed %v want %v, course %v", p.Speed, want, p.Course)
	}
}

func TestParseKmhColumn(t *testing.T) {
	const line = "20260421145509,51.311962,N,12.323399,E,10,131.9"
	opt := ParseOptions{SpeedColumn: SpeedColumnKmh}
	p, err := parseRecord(strings.Split(line, ","), opt)
	if err != nil {
		t.Fatal(err)
	}
	if p.Speed != 10 {
		t.Fatalf("speed %v", p.Speed)
	}
}
