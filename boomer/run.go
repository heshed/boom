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
	"crypto/tls"
	"sync"
	"net"
	"net/http"
	"time"
)

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func (b *Boomer) Run() {
	b.run()
}

func (b *Boomer) worker(wg *sync.WaitGroup, ch chan *http.Request, reports *report) {
	host, _, _ := net.SplitHostPort(b.Req.OriginalHost)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: b.AllowInsecure,
			ServerName:         host,
		},
		DisableCompression: b.DisableCompression,
		DisableKeepAlives:  b.DisableKeepAlives,
		// TODO(jbd): Add dial timeout.
		TLSHandshakeTimeout: time.Duration(b.Timeout) * time.Millisecond,
		Proxy:               http.ProxyURL(b.ProxyAddr),
	}
	client := &http.Client{Transport: tr}
	_ = client
	for req := range ch {
		s := time.Now()
		code := 0
		size := int64(0)
		b.SentCounter.Incr(1)
		resp, err := client.Do(req)
		if err == nil {
			size = resp.ContentLength
			code = resp.StatusCode
			resp.Body.Close()
		}

		res := result{
			statusCode:    code,
			duration:      time.Now().Sub(s),
			err:           err,
			contentLength: size,
		}

		reports.update(reports.second, &res)
		reports.update(reports.half, &res)
		reports.update(reports.total, &res)
	}
}

func (b *Boomer) run() {
	rep := NewReport()

	var wg sync.WaitGroup
	wg.Add(2)

	shutdown := make(chan bool)
	go rep.Run(shutdown, &wg)

	var throttle <-chan time.Time
	if b.Qps > 0 {
		throttle = time.Tick(time.Duration(1e6/(b.Qps)) * time.Microsecond)
	}
	jobs := make(chan *http.Request, 1)
	for i := 0; i < b.C; i++ {
		go func() {
			b.worker(&wg, jobs, rep)
		}()
	}

	i := 0
	for {
		i++
		if b.Qps > 0 {
			<-throttle
		}

		body, err := b.BodyReader.Next()
		if err != true {
			break
		}
		if b.N > 0 && i > b.N {
			break
		}
		jobs <- b.Req.Request(&body)
	}
	close(jobs)
	close(shutdown)

	wg.Done()
	wg.Wait()
}
