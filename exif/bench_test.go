package exif

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

var (
	sampleDir     = "../../test/img/"
	testFilenames = []string{
		"honor20.jpg",
		"hero6.jpg",
		"1.CR2",
		"3.CR2",
		"350D.CR2",
		"XT1.CR2",
		"60D.CR2",
		"6D.CR2",
		"7D.CR2",
		"90D.cr3",
		"2.CR3",
		"1.CR3",
		"1.jpg",
		"2.jpg",
		"1.NEF",
		"2.NEF",
		"3.NEF",
		"1.ARW",
		"2.ARW",
		"4.RW2",
		"hero6.gpr",
		"4.webp",
		"20.jpg",
	}
	testFilenames2 = []string{
		"2.CR2",
	}
)

func BenchmarkScanExif100(b *testing.B) {
	for _, bm := range testFilenames {
		b.Run(bm, func(b *testing.B) {
			f, err := os.Open(sampleDir + bm)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			buf, _ := ioutil.ReadAll(f)
			cb := bytes.NewReader(buf)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err = ScanExif(cb)
				if err != nil {
					if err != ErrNoExif {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

//BenchmarkScanExif100/honor20.jpg         	   95787	     11793 ns/op	    5488 B/op	      25 allocs/op
//BenchmarkScanExif100/hero6.jpg           	  122072	      8912 ns/op	    5351 B/op	      21 allocs/op
//BenchmarkScanExif100/1.CR2               	   50002	     23182 ns/op	   11129 B/op	      31 allocs/op
//BenchmarkScanExif100/3.CR2               	   74066	     18786 ns/op	    5724 B/op	      29 allocs/op
//BenchmarkScanExif100/350D.CR2            	   73534	     16655 ns/op	    5484 B/op	      25 allocs/op
//BenchmarkScanExif100/XT1.CR2             	   78391	     14130 ns/op	    5484 B/op	      25 allocs/op
//BenchmarkScanExif100/60D.CR2             	   56073	     21547 ns/op	   11083 B/op	      29 allocs/op
//BenchmarkScanExif100/6D.CR2              	   56967	     21669 ns/op	   11116 B/op	      31 allocs/op
//BenchmarkScanExif100/7D.CR2              	   56366	     22410 ns/op	   11131 B/op	      31 allocs/op
//BenchmarkScanExif100/90D.cr3             	  265622	      3936 ns/op	     982 B/op	       9 allocs/op
//BenchmarkScanExif100/2.CR3               	  316366	      3631 ns/op	     976 B/op	       9 allocs/op
//BenchmarkScanExif100/1.CR3               	  315829	      3648 ns/op	     976 B/op	       9 allocs/op
//BenchmarkScanExif100/1.jpg               	  156877	     11899 ns/op	    2579 B/op	      18 allocs/op
//BenchmarkScanExif100/2.jpg               	   99021	     12269 ns/op	    5433 B/op	      23 allocs/op
//BenchmarkScanExif100/1.NEF               	   48734	     23879 ns/op	   11242 B/op	      33 allocs/op
//BenchmarkScanExif100/2.NEF               	   46730	     25353 ns/op	   11243 B/op	      33 allocs/op
//BenchmarkScanExif100/3.NEF               	   43644	     26385 ns/op	   11303 B/op	      35 allocs/op
//BenchmarkScanExif100/1.ARW               	   72991	     15820 ns/op	    5658 B/op	      30 allocs/op
//BenchmarkScanExif100/2.ARW               	   78735	     16046 ns/op	    5658 B/op	      30 allocs/op
//BenchmarkScanExif100/4.RW2               	   65305	     17868 ns/op	    5360 B/op	      21 allocs/op
//BenchmarkScanExif100/hero6.gpr           	   67656	     15147 ns/op	    5434 B/op	      23 allocs/op
//BenchmarkScanExif100/4.webp              	    1120	   1117162 ns/op	     240 B/op	       3 allocs/op
//BenchmarkScanExif100/20.jpg              	     536	   1925329 ns/op	     240 B/op	       3 allocs/op

func BenchmarkParseExif100(b *testing.B) {
	for _, bm := range testFilenames2 {
		b.Run(bm, func(b *testing.B) {
			f, err := os.Open(sampleDir + bm)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			buf, _ := ioutil.ReadAll(f)
			cb := bytes.NewReader(buf)
			e, err := ScanExif(cb)
			if err != nil {
				if err != ErrNoExif {
					b.Fatal(err)
				}
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				b.StartTimer()
				e.Aperture()
			}
		})
	}
}

// GPSDate
// BenchmarkParseExif100/2.CR2         	 1000000	      1044 ns/op	     176 B/op	       4 allocs/op
// BenchmarkParseExif100/2.CR2         	 4983446	       272 ns/op	       0 B/op	       0 allocs/op

// GPSCoords
// BenchmarkParseExif100/2.CR2         	 1340478	       850 ns/op	     224 B/op	       4 allocs/op
// BenchmarkParseExif100/2.CR2         	 3894754	       303 ns/op	       0 B/op	       0 allocs/op

// Artist
// BenchmarkParseExif100/0.CR2         	 4855602	       241 ns/op	      16 B/op	       1 allocs/op

// Aperture
// BenchmarkParseExif100/0.CR2         	10543807	       123 ns/op	       0 B/op	       0 allocs/op
// BenchmarkParseExif100/0.CR2         	10259781	       116 ns/op	       0 B/op	       0 allocs/op
// BenchmarkParseExif100/0.CR2         	 3556357	       303 ns/op	      80 B/op	       1 allocs/op

// DateTime
// BenchmarkParseExif100/0.CR2         	 1000000	      1043 ns/op	     208 B/op	       5 allocs/op
// BenchmarkParseExif100/0.CR2         	 6169940	       230 ns/op	       0 B/op	       0 allocs/op

// BenchmarkParseExif100/1.CR2         	 5274560	       256 ns/op	      16 B/op	       2 allocs/op

// BenchmarkParseExif100/1.CR2         	12420832	        97.1 ns/op	       0 B/op	       0 allocs/op

// BenchmarkParseExif100/1.CR2         	 5502931	       235 ns/op	      16 B/op	       1 allocs/op

// BenchmarkParseExif100/1.CR2         	   36753	     31716 ns/op	   13120 B/op	      36 allocs/op

// BenchmarkParseExif100/1.CR2         	   31756	     34805 ns/op	   13706 B/op	      38 allocs/op

// BenchmarkParseExif100/1.CR2         	   32517	     37431 ns/op	   13706 B/op	      38 allocs/op

// BenchmarkParseExif100/1.CR2         	   34333	     35437 ns/op	    9148 B/op	      56 allocs/op

// BenchmarkParseExif100/1.CR2         	   31680	     39503 ns/op	   13214 B/op	      57 allocs/op
