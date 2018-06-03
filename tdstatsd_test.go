package main

import (
	"sort"
	"testing"
)

var mustParse = func(t *testing.T, data []byte) []TDPool {
	pools, err := parse(data)
	if err != nil {
		t.Fatal(err)
	}
	return pools
}

func TestParseSort(t *testing.T) {
	data := []byte(testDataHead + testDataPools + testDataTail)
	pools := mustParse(t, data)
	if len(pools) == 0 {
		// nothing to sort
		t.SkipNow()
	}

	// check sort. last 4 pools must be online
	sort.Sort(byStatus(pools))
	for _, p := range pools[2:] {
		if p.Status != "online" {
			t.Fatal("Wrong pools order")
		}
	}
}

func TestParseEmpty(t *testing.T) {
	t.Log(testDataHead + testDataTail)
	data := []byte(testDataHead + testDataTail)
	if pools := mustParse(t, data); len(pools) != 0 {
		t.Fatal("Expected 0 pools")
	}
}

func TestParseWrongData(t *testing.T) {
	data := []byte("wrong data")
	if pools := mustParse(t, data); len(pools) != 0 {
		t.Fatal("Expected 0 pools")
	}
}

const testDataHead = `trafficd pid: 12345

Shaitan Traffic Detector 11.1.1.1.1 B99/11/2025 01:23 (MarsOS)

Server started Sun Feb 31 23:55:55 2022
Process 12345 started Sun Feb 31 23:55:55 2022

ConnectionQueue:
-----------------------------------------
Current/Peak/Limit Queue Length            100/100500/12345
Total Connections Queued                   23984752345234
Average Queue Length (1, 5, 15 minutes)    0,02, 0,04, 0,16
Average Queueing Delay                     0,03 milliseconds

ListenSocket http-listener-1:
------------------------
Address                   0.0.0.0:10119
Acceptor Threads          148
Default Virtual Server    FOO_BAR

ListenSocket AI-listene:
------------------------
Address                   0.0.0.0:10129
Acceptor Threads          148
Default Virtual Server    JEDI_MASTER

KeepAliveInfo:
--------------------
KeepAliveCount        6/39552
KeepAliveHits         32450924
KeepAliveFlushes      0
KeepAliveRefusals     0
KeepAliveTimeouts     7
KeepAliveTimeout      20 seconds

SessionCreationInfo:
------------------------
Active Sessions           76
Keep-Alive Sessions       6
Total Sessions Created    9873/7832

Proxy Cache:
---------------------------
Proxy Cache Enabled              yes
Object Cache Entries             233
Cache lookup (hits/misses)       2345234/465733
Requests served from Cache       142532
Revalidation (successful/total)  435634/568756 ( 34,12%)
Heap space used                  356342

Native pools:
----------------------------
NativePool:
Idle/Peak/Limit               5/5/512
Work Queue Length/Peak/Limit  0/0/0

DNSCacheInfo:
------------------
enabled             yes
CacheEntries        0/1024
HitRatio            0/0 (  0,00%)

Async DNS disabled

Performance Counters:
------------------------------------------------
                           Average         Total      Percent

Total number of requests:               98734598
Request processing time:    0,0360  2028741,0000

default-bucket (Default bucket)
Number of Requests:                     98734598    (100,00%)
Number of Invocations:                 908327453    (100,00%)
Latency:                    0,0005    23423,4297    (  1,41%)
Function Processing Time:   0,0355  2342323,5000    ( 98,59%)
Total Response Time:        0,0360  2028741,0000    (100,00%)

Origin server statistics (for http):
---------------------------------------------------------------------------------------------------------------------------------------------------------------
Pool-name             Host:Port                    Status  ActiveConn  IdleConn  StickyConn  Timeouts  Aborted  Sticky-Reqs  Total-Reqs  BytesTrans  BytesRecvd

`

const testDataPools = `
fidget-server-pool-1  http://oe-oe-foo-bar2:8899   online  2           2         1           0         12       35463456     35645645    345G        897G
AI-server-pool        http://oe-oe-foo-bar2:9999   online  0           1         0           0         185      34545        13486       576M        87963M
fidget-server-pool-1  http://oe-oe-foo-bar1:8899   offline 1           2         1           0         15       34564562     98372452    876G        765G
fidget-server-pool-1  http://oe-oe-foo-bar2:8889   online  1           3         1           1         156      82345867     23452345    786G        345G
AI-server-pool        http://oe-oe-foo-bar1:9999   fooline 0           1         0           12        154      34564        34523       754M        864M
fidget-server-pool-1  http://oe-oe-foo-bar1:8889   online  0           6         0           0         14       35635645     34565324    654G        234G
`

const testDataTail = `
Origin server statistics (for tcp):
-----------------------------------------------------------------------------------------------
Pool-name  Host:Port  Status  ActiveConn  Timeouts  Aborted  Total-Reqs  BytesTrans  BytesRecvd


Sessions:
--------------------------------------------------------------------------------------------------------------------------------------------------
Process  Status    Client         Age  VS        Method  URI                                          Function         Origin-Server

35011    response  555.23.1.196   123  FOO_BAR  POST    /aaaaaaaaaa/bbbbbbbbb/cccccccccccccc_Service  foxy-retrieve/   http://oe-oe-foo-bar2:0899
35011    response  555.23.1.197   456  FOO_BAR  POST    /aaaaaaaaaa/bbbbbbbbb/cccccccccccccc_Service  foxy-retrieve/   http://oe-oe-foo-bar2:0889
35011    response  555.23.1.197   789  FOO_BAR  POST    /aaaaaaaaaa/bbbbbbbbb/cccccccccccccc_Service  foxy-retrieve/   http://oe-oe-foo-bar1:0899
35011    response  555.23.1.1949  5    FOO_BAR  GET     /.abcd                                        service-pump/    -

TCP Proxy:
------------------------------------
Active Connections                  0
Avg Duration                        0,00 seconds
Requests (timeout/aborted/total)    0/0/0
`
