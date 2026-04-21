package gpx

import (
	"encoding/xml"
	"io"
	"math"
	"time"

	"vantrue-gpx/internal/vantrue"
)

const (
	nsGPX = "http://www.topografix.com/GPX/1/1"
	nsTPX = "http://www.garmin.com/xmlschemas/TrackPointExtension/v2"
	nsVT  = "https://www.vantrue.net/gpx-extensions/1"
	nsXSI = "http://www.w3.org/2001/XMLSchema-instance"
)

type EncodeOptions struct {
	TrackName string
	Creator   string
}

func WriteGPX(w io.Writer, pts []vantrue.Point, opt EncodeOptions) error {
	if opt.Creator == "" {
		opt.Creator = "vantrue-gpx"
	}
	name := opt.TrackName
	if name == "" {
		name = "Vantrue track"
	}

	doc := gpxDoc{
		XMLName:     xml.Name{Space: nsGPX, Local: "gpx"},
		Version:     "1.1",
		Creator:     opt.Creator,
		XmlnsXsi:    nsXSI,
		XmlnsGpxtpx: nsTPX,
		Metadata: metadata{
			Time: time.Now().UTC().Format(time.RFC3339),
		},
		Trk: trk{
			Name: name,
			Seg: trkseg{
				Points: make([]trkpt, 0, len(pts)),
			},
		},
	}

	for _, p := range pts {
		tp := trkpt{
			Lat:  p.Lat,
			Lon:  p.Lon,
			Time: p.Time.UTC().Format(time.RFC3339),
		}
		ext := extensions{
			Tpx: tpxExt{
				Speed:  p.Speed / 3.6,
				Course: normalizeCourseDegrees(p.Course),
			},
		}
		if p.AccelX != nil && p.AccelY != nil && p.AccelZ != nil {
			ext.VtAccel = &vtAccel{
				X: *p.AccelX,
				Y: *p.AccelY,
				Z: *p.AccelZ,
			}
			doc.XmlnsVt = nsVT
		}
		tp.Ext = &ext
		doc.Trk.Seg.Points = append(doc.Trk.Seg.Points, tp)
	}

	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return err
	}
	_, err := w.Write([]byte("\n"))
	return err
}

type gpxDoc struct {
	XMLName     xml.Name `xml:"http://www.topografix.com/GPX/1/1 gpx"`
	Version     string   `xml:"version,attr"`
	Creator     string   `xml:"creator,attr"`
	XmlnsXsi    string   `xml:"xmlns:xsi,attr,omitempty"`
	XmlnsGpxtpx string   `xml:"xmlns:gpxtpx,attr,omitempty"`
	XmlnsVt     string   `xml:"xmlns:vt,attr,omitempty"`
	Metadata    metadata `xml:"metadata"`
	Trk         trk      `xml:"trk"`
}

type metadata struct {
	Time string `xml:"time,omitempty"`
}

type trk struct {
	Name string `xml:"name"`
	Seg  trkseg `xml:"trkseg"`
}

type trkseg struct {
	Points []trkpt `xml:"trkpt"`
}

type trkpt struct {
	Lat  float64     `xml:"lat,attr"`
	Lon  float64     `xml:"lon,attr"`
	Time string      `xml:"time"`
	Ext  *extensions `xml:"extensions,omitempty"`
}

type extensions struct {
	Tpx     tpxExt   `xml:"http://www.garmin.com/xmlschemas/TrackPointExtension/v2 TrackPointExtension"`
	VtAccel *vtAccel `xml:"https://www.vantrue.net/gpx-extensions/1 imu,omitempty"`
}

type tpxExt struct {
	Speed  float64 `xml:"http://www.garmin.com/xmlschemas/TrackPointExtension/v2 speed,omitempty"`
	Course float64 `xml:"http://www.garmin.com/xmlschemas/TrackPointExtension/v2 course,omitempty"`
}

type vtAccel struct {
	X float64 `xml:"https://www.vantrue.net/gpx-extensions/1 accelX"`
	Y float64 `xml:"https://www.vantrue.net/gpx-extensions/1 accelY"`
	Z float64 `xml:"https://www.vantrue.net/gpx-extensions/1 accelZ"`
}

func normalizeCourseDegrees(c float64) float64 {
	c = math.Mod(c, 360)
	if c < 0 {
		c += 360
	}
	if c >= 360 {
		return 0
	}
	return c
}
