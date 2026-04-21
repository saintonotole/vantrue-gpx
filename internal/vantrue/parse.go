package vantrue

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type SpeedColumn int

const (
	SpeedColumnKnots SpeedColumn = iota
	SpeedColumnKmh
)

type ParseOptions struct {
	SpeedColumn SpeedColumn
}

type Point struct {
	Time                   time.Time
	Lat                    float64
	Lon                    float64
	Speed                  float64
	Course                 float64
	AccelX, AccelY, AccelZ *float64
}

const timeLayout = "20060102150405"
const knotsToKMH = 1.852

func ParseFile(r io.Reader, opt ParseOptions) ([]Point, error) {
	sc := bufio.NewScanner(r)
	const max = 1024 * 1024
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, max)

	var out []Point
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, ",")
		p, err := parseRecord(fields, opt)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}
		out = append(out, p)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func parseRecord(fields []string, opt ParseOptions) (Point, error) {
	var p Point
	if len(fields) < 7 {
		return p, fmt.Errorf("expected at least 7 fields, got %d", len(fields))
	}

	ts, err := time.ParseInLocation(timeLayout, fields[0], time.UTC)
	if err != nil {
		return p, fmt.Errorf("timestamp: %w", err)
	}
	p.Time = ts

	lat, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return p, fmt.Errorf("latitude: %w", err)
	}
	switch strings.ToUpper(strings.TrimSpace(fields[2])) {
	case "N":
		p.Lat = lat
	case "S":
		p.Lat = -lat
	default:
		return p, fmt.Errorf("invalid latitude hemisphere %q", fields[2])
	}

	lon, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return p, fmt.Errorf("longitude: %w", err)
	}
	switch strings.ToUpper(strings.TrimSpace(fields[4])) {
	case "E":
		p.Lon = lon
	case "W":
		p.Lon = -lon
	default:
		return p, fmt.Errorf("invalid longitude hemisphere %q", fields[4])
	}

	speed, err := strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return p, fmt.Errorf("speed: %w", err)
	}
	p.Speed = normalizeSpeedKMH(speed, opt)

	course, err := strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return p, fmt.Errorf("course: %w", err)
	}
	p.Course = course

	if len(fields) > 7 {
		if err := parseOptionalFloats(fields[7:], &p); err != nil {
			return p, err
		}
	}
	return p, nil
}

func parseOptionalFloats(rest []string, p *Point) error {
	trimmed := trimAll(rest)
	switch len(trimmed) {
	case 0:
		return nil
	case 1, 2, 4, 5:
		return fmt.Errorf("unsupported trailing field count %d (expected 3 or 6+)", len(trimmed))
	case 3:
		x, err := strconv.ParseFloat(trimmed[0], 64)
		if err != nil {
			return fmt.Errorf("accel x: %w", err)
		}
		y, err := strconv.ParseFloat(trimmed[1], 64)
		if err != nil {
			return fmt.Errorf("accel y: %w", err)
		}
		z, err := strconv.ParseFloat(trimmed[2], 64)
		if err != nil {
			return fmt.Errorf("accel z: %w", err)
		}
		p.AccelX, p.AccelY, p.AccelZ = fp(x), fp(y), fp(z)
		return nil
	default:
		x, err1 := strconv.ParseFloat(trimmed[0], 64)
		y, err2 := strconv.ParseFloat(trimmed[1], 64)
		z, err3 := strconv.ParseFloat(trimmed[2], 64)
		if err1 != nil || err2 != nil || err3 != nil {
			return fmt.Errorf("trailing fields present but first three are not all floats")
		}
		p.AccelX, p.AccelY, p.AccelZ = fp(x), fp(y), fp(z)
		return nil
	}
}

func trimAll(s []string) []string {
	out := make([]string, 0, len(s))
	for _, v := range s {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func fp(v float64) *float64 {
	return &v
}

func normalizeSpeedKMH(raw float64, opt ParseOptions) float64 {
	if opt.SpeedColumn == SpeedColumnKmh {
		return raw
	}
	return raw * knotsToKMH
}
