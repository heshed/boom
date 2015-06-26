// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package boomer

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
	"github.com/labstack/gommon/color"
)

const (
	barChar = "∎"
)

type stat struct {
	avgTotal float64
	fastest  float64
	slowest  float64
	average  float64
	rps      float64

	total time.Duration

	errorDist      map[string]int
	errorCount     Counter
	statusCodeDist map[int]int
	lats           []float64
	sizeTotal      int64
	tick           int
}

func NewStat(tick int) *stat {
	return &stat{
		statusCodeDist: make(map[int]int),
		errorDist:      make(map[string]int),
		tick:           tick,
	}
}

type report struct {
	second *stat
	half   *stat
	total  *stat

	output string
}

func NewReport() *report {
	return &report{
		second: NewStat(1),
		half:   NewStat(30),
		total:  NewStat(-1),
	}
}

func (r *report) clear(s *stat) {
	s.avgTotal = 0
	s.fastest = 0
	s.slowest = 0
	s.average = 0
	s.rps = 0

	s.total = 0

	s.lats = []float64{}
	s.sizeTotal = 0
}

func (r *report) Run(shutdown <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Second)
	ticker30 := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	r.printColumns()
	for {
		select {
		case <-shutdown:
			r.printStatusCodes()
			r.printHistogram()
			r.printLatencies()
			if r.total.errorCount.Val() > 0 {
				r.printErrors()
			}
			return
		case <-ticker.C:
			r.printStat(r.second)
			r.clear(r.second)

		case <-ticker30.C:
			r.printStat(r.half)
			r.clear(r.half)
			r.printColumns()
		}
	}
}

func (r *report) printColumns() {
	fmt.Println("\n :..:..:..:   Avg.Count   Avg.Res(ms)   fails   Fastest(ms)   Slowest(ms)   Avg.Res.Size")
}

func (r *report) update(s *stat, res *result) {
	if res.err != nil {
		s.errorDist[res.err.Error()]++
		s.errorCount.Incr(1)
	} else {
		s.lats = append(s.lats, res.duration.Seconds())
		s.avgTotal += res.duration.Seconds()
		s.statusCodeDist[res.statusCode]++
		if res.contentLength > 0 {
			s.sizeTotal += res.contentLength
		}
	}
}

func (r *report) printStat(s *stat) {
	sort.Float64s(s.lats)

	var avgSizeTotal int64 = 0
	if len(s.lats) > 0 {
		s.rps = float64(len(s.lats)) / (float64)(s.tick)
		s.average = s.avgTotal / float64(len(s.lats))
		s.fastest = s.lats[0] * 1000
		s.slowest = s.lats[len(s.lats)-1] * 1000
		avgSizeTotal = s.sizeTotal / int64(len(s.lats))
	} else {
		s.fastest = 0
		s.slowest = 0
	}

	var now string
	if s.tick == 30 {
		now = " [ 30Sec ]"
	} else {
		now = time.Now().Local().Round(time.Second).Format("  15:04:05")
	}

	line := fmt.Sprintf("%s   %6.0f     %6.0f       %6d      %6.0f       %6.0f           %d",
		now,
		s.rps,
		s.average*1000,
		s.errorCount.Val(),
		s.fastest,
		s.slowest,
		avgSizeTotal,
	)

	if s.tick == 30 {
		fmt.Println(color.Blue(line))
	} else {
		fmt.Println(line)
	}
}

// Prints percentile latencies.
func (r *report) printLatencies() {
	s := r.total
	pctls := []int{10, 25, 50, 75, 90, 95, 99}
	data := make([]float64, len(pctls))
	j := 0
	for i := 0; i < len(s.lats) && j < len(pctls); i++ {
		current := i * 100 / len(s.lats)
		if current >= pctls[j] {
			data[j] = s.lats[i]
			j++
		}
	}
	fmt.Printf("\nLatency distribution:\n")
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			fmt.Printf("  %v%% in %4.4f secs.\n", pctls[i], data[i])
		}
	}
}

func (r *report) printHistogram() {
	fmt.Printf("\nResponse       time   histogram:\n")
	s := r.total
	if len(s.lats) <= 0 {
		fmt.Println("  no data")
		return
	}
	sort.Float64s(s.lats)
	s.fastest = s.lats[0]
	s.slowest = s.lats[len(s.lats)-1]
	bc := 12
	buckets := make([]float64, bc+1)
	counts := make([]int, bc+1)
	//bs := (s.slowest - s.fastest) / float64(bc)
	buckets[0] = 0.01
	buckets[1] = 0.05
	for i := 2; i < bc; i++ {
		//buckets[i] = s.fastest + bs*float64(i)
		buckets[i] = (float64)(i+1) * 0.1
	}
	buckets[bc] = s.slowest
	last_index := len(buckets) - 1
	var bi int
	var max int
	for i := 0; i < len(s.lats); {
		if s.lats[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		} else {
			i++
			counts[last_index]++
			if max < counts[last_index] {
				max = counts[last_index]
			}
		}
	}
	for i := 0; i < len(buckets); i++ {
		// Normalize bar lengths.
		var barLen int
		if max > 0 {
			barLen = counts[i] * 40 / max
		}
		fmt.Printf("  %6.0f [%10d]\t|%v\n", buckets[i]*1000, counts[i], strings.Repeat(barChar, barLen))
	}
}

// Prints status code distribution.
func (r *report) printStatusCodes() {
	s := r.total
	fmt.Printf("\nStatus code distribution:\n")
	for code, num := range s.statusCodeDist {
		fmt.Printf("  [%d]\t%d responses\n", code, num)
	}
}

func (r *report) printErrors() {
	s := r.total
	fmt.Printf("\nError distribution: %d\n", s.errorCount.Val())
	for err, num := range s.errorDist {
		fmt.Printf("  [%d]\t%s\n", num, err)
	}
}
