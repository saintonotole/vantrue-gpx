package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vantrue-gpx/internal/gpx"
	"vantrue-gpx/internal/vantrue"
)

func main() {
	in := flag.String("i", "", "input .dat file or directory of .dat files")
	out := flag.String("o", "", "output .gpx path, or output directory when -i is a directory (default: next to each .dat)")
	speedCol := flag.String("speed-column", "knots", `speed field: "knots" (×1.852→km/h) or "kmh"`)
	flag.Parse()
	if *in == "" {
		flag.Usage()
		os.Exit(2)
	}

	opts := vantrue.ParseOptions{}
	switch *speedCol {
	case "knots":
		opts.SpeedColumn = vantrue.SpeedColumnKnots
	case "kmh":
		opts.SpeedColumn = vantrue.SpeedColumnKmh
	default:
		fmt.Fprintf(os.Stderr, "unknown -speed-column %q (use knots or kmh)\n", *speedCol)
		os.Exit(2)
	}

	st, err := os.Stat(*in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if st.IsDir() {
		outDir := strings.TrimSpace(*out)
		if outDir != "" {
			if err := os.MkdirAll(outDir, 0o755); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		entries, err := os.ReadDir(*in)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".dat") {
				continue
			}
			inPath := filepath.Join(*in, e.Name())
			base := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
			outPath := filepath.Join(*in, base+".gpx")
			if outDir != "" {
				outPath = filepath.Join(outDir, base+".gpx")
			}
			if err := convertOne(inPath, outPath, opts); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", inPath, err)
				os.Exit(1)
			}
		}
		return
	}

	outPath := resolveOutputFilePath(*in, *out)
	if err := convertOne(*in, outPath, opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func resolveOutputFilePath(inPath, outSpec string) string {
	outSpec = strings.TrimSpace(outSpec)
	if outSpec == "" {
		return strings.TrimSuffix(inPath, filepath.Ext(inPath)) + ".gpx"
	}
	if st, err := os.Stat(outSpec); err == nil && st.IsDir() {
		base := filepath.Base(inPath)
		return filepath.Join(outSpec, strings.TrimSuffix(base, filepath.Ext(base))+".gpx")
	}
	return outSpec
}

func convertOne(inPath, outPath string, popts vantrue.ParseOptions) error {
	f, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pts, err := vantrue.ParseFile(f, popts)
	if err != nil {
		return err
	}

	w, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer w.Close()

	base := filepath.Base(inPath)
	return gpx.WriteGPX(w, pts, gpx.EncodeOptions{
		TrackName: strings.TrimSuffix(base, filepath.Ext(base)),
		Creator:   "vantrue-gpx",
	})
}
