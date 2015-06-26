# boom

Boom is a tiny program that sends some load to a web application. It's similar to Apache Bench ([ab](http://httpd.apache.org/docs/2.2/programs/ab.html)), but with better availability across different platforms and a less troubling installation experience.

fork from [boom](https://github.com/rakyll/boom)

## Installation

Simple as it takes to type the following command:

    go get github.com/heshed/boom

## Usage

Boom supports custom headers, request body and basic authentication. It runs provided number of requests in the provided concurrency level, and prints stats.

~~~
Usage: boom [options...] <url>

Options:
  -n  Number of requests to run.
  -c  Number of requests to run concurrently. Total number of requests cannot
      be smaller than the concurency level.
  -q  Rate limit, in seconds (QPS).
  -o  Output type. If none provided, a summary is printed.
      "csv" is the only supported alternative. Dumps the response
      metrics in comma-seperated values format.

  -i  Request body from file
  -m  HTTP method, one of GET, POST, PUT, DELETE, HEAD, OPTIONS.
  -h  Custom HTTP headers, name1:value1;name2:value2.
  -t  Timeout in ms.
  -A  HTTP Accept header.
  -d  HTTP request body.
  -T  Content-type, defaults to "text/html".
  -a  Basic authentication, username:password.
  -x  HTTP Proxy address as host:port.

  -allow-insecure       Allow bad/expired TLS/SSL certificates.
  -disable-compression  Disable compression.
  -disable-keepalive    Disable keep-alive, prevents re-use of TCP
                        connections between different HTTP requests.
  -cpus                 Number of used cpu cores.
                        (default for current machine is 4 cores)
~~~

This is what happens when you run Boom:

    boom -i q -q 500 -c 100    http://192.168.59.103:18000/  ✚ ✱ ◼
    :......................................:
    URL: http://192.168.59.103:18000/
    Method: GET
    Body:
    Header:
    Username:
    Password:
    OriginalHost: 192.168.59.103:18000
    Number of requests to run: -1
    Number of concurrency: 100
    QPS: 500
    Timeout: 0
    AllowInsecure: false
    DisableCompression: false
    DisableKeepAlives: false
    ProxyAddr:
    Output:
    BodyFile: /Users/md/.go/src/boom/q
    :......................................:

     :..:..:..:   Avg.Count   Avg.Res(ms)   fails   Fastest(ms)   Slowest(ms)   Avg.Res.Size
      03:35:52      494          4            0           1           11           166
      03:35:53      503         14            0           1         1253           168
      03:35:54      498          4            0           1           12           168
      03:35:55      502          4            0           1           12           168
      03:35:56      501          4            0           1           11           168
      03:35:57      491          4            0           1           21           168
      03:35:58      503          4            0           1           10           168
      03:35:59      500          4            0           2           13           168
      03:36:00      498          5            0           2           14           168
      03:36:01      503          4            0           2           12           168
      03:36:02      499          5            0           2           16           168
      03:36:03      489          5            0           2           18           168
      03:36:04      492          5            0           2           27           168
      03:36:05      498          5            0           2           19           168
      03:36:06      496          5            0           2           17           168
      03:36:07      494          9            0           2           83           168
      03:36:08      497          6            0           2           20           168
      03:36:09      497          6            0           2           24           168
      03:36:10      495          6            0           3           22           168
      03:36:11      497          7            0           3           25           168
      03:36:12      493          7            0           3           22           168
      03:36:13      497          8            0           3           26           168
      03:36:14      495          8            0           3           26           168
      03:36:15      492         10            0           3           33           168
      03:36:16      492         10            0           4           32           168
      03:36:17      478         11            0           4           29           168
      03:36:18      488         11            0           4           32           168
      03:36:19      482         13            0           4           33           168
      03:36:20      475         14            0           4           40           168
     [ 30Sec ]      493          7            0           1         1253           167

     :..:..:..:   Avg.Count   Avg.Res(ms)   fails   Fastest(ms)   Slowest(ms)   Avg.Res.Size
      03:36:21      464         15            0           5           40           168
      03:36:22      454         18            0           4           40           168
      03:36:23      426         20            0           8           45           168
      03:36:24      441         20            0           4           37           168
      03:36:25      235         20           90           5           42           168
      03:36:26        0          0          293           0            0           0
      03:36:27        0          0          500           0            0           0
      03:36:28        0          0          705           0            0           0
      03:36:29        0          0          911           0            0           0
      03:36:30        0          0         1114           0            0           0
      03:36:31        0          0         1311           0            0           0
      03:36:32        2         36         1501          27           44           168
      03:36:33        0          0         1696           0            0           0
      03:36:34        0          0         1892           0            0           0
      03:36:37        0          0         1916           0            0           0
      03:36:37       28        113         1916           4         2775           168
      03:36:41      433         14         1916           4           37           168
      03:36:41      154         98         1916           3         2785           168
      03:36:42      298          8         1916           3           25           168
      03:36:43      381         15         1916           3          546           168
      03:36:44      244         10         1916           2          518           168
      03:36:45      299          7         1916           3           29           168
      03:36:46      486          8         1916           3          424           168
      03:36:47      322          9         1916           3          364           168
      03:36:48      303         10         1916           2          386           168
      03:36:49      360          8         1916           2          286           168
      03:36:51      480          5         1916           3           22           168
      03:36:51      319         12         1916           2          411           168
     [ 30Sec ]      189         15         1916           2         2785           168

     :..:..:..:   Avg.Count   Avg.Res(ms)   fails   Fastest(ms)   Slowest(ms)   Avg.Res.Size
      03:36:52      318          9         1916           2          365           168
      03:36:53      372          7         1916           3          262           168
      03:36:55      490          5         1916           3           22           168
      03:36:55      306          8         1916           2          402           168
      03:36:56      373          7         1916           2          263           168
      03:36:57      300          8         1916           2          406           168
      03:36:58      482          5         1916           2           18           168
      03:36:59      384          7         1916           2          271           168

    Status code distribution:
      [404] 23552 responses
      [200] 8 responses

    Response       time   histogram:
          10 [     18608]   |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
          50 [      4882]   |∎∎∎∎∎∎∎∎∎∎
         300 [        34]   |
         400 [         9]   |
         500 [        12]   |
         600 [         5]   |
         700 [         0]   |
         800 [         0]   |
         900 [         0]   |
        1000 [         0]   |
        1100 [         0]   |
        1200 [         0]   |
        2785 [        10]   |

    Latency distribution:
      10% in 0.0029 secs.
      25% in 0.0036 secs.
      50% in 0.0056 secs.
      75% in 0.0093 secs.
      90% in 0.0157 secs.
      95% in 0.0205 secs.
      99% in 0.0278 secs.

    Error distribution: 1916
      [1916]    Get http://192.168.59.103:18000/index.htmla: dial tcp 192.168.59.103:18000: can't assign requested address

## License

Copyright 2014 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

