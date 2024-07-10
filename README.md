grpcx-benchmark

	goos: linux
	goarch: amd64
	pkg: go-driver/grpcx
	cpu: 12th Gen Intel(R) Core(TM) i3-12100F
	Memory: 16G
	Go: 1.21
	异步调用（不需要返回值的调用）qps
	PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
	73684 ubuntu    20   0 1246124  26344  16500 R 155.5   0.3   0:12.43 main
	2024-07-07 20:32:40 handle.go:73 DEBUG 874325 request/s NumGoroutine 22
	2024-07-07 20:32:46 handle.go:73 DEBUG 882882 request/s NumGoroutine 22
	2024-07-07 20:32:47 handle.go:73 DEBUG 892185 request/s NumGoroutine 22
	2024-07-07 20:32:48 handle.go:73 DEBUG 884763 request/s NumGoroutine 22
	2024-07-07 20:32:49 handle.go:73 DEBUG 899013 request/s NumGoroutine 22
	2024-07-07 20:32:50 handle.go:73 DEBUG 892694 request/s NumGoroutine 22
	2024-07-07 20:32:51 handle.go:73 DEBUG 886786 request/s NumGoroutine 22
	2024-07-07 20:32:52 handle.go:73 DEBUG 887195 request/s NumGoroutine 22
	2024-07-07 20:32:53 handle.go:73 DEBUG 877819 request/s NumGoroutine 22
	2024-07-07 20:32:54 handle.go:73 DEBUG 882005 request/s NumGoroutine 22
	2024-07-07 20:32:55 handle.go:73 DEBUG 878792 request/s NumGoroutine 22
	同步调用（需要返回值的调用）qps
	PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
	46916 ubuntu    20   0 1246124  29052  16816 R 177.3   0.4 124:27.72 main
	2024-07-07 20:10:13 handle.go:73 DEBUG 770865 request/s NumGoroutine 19
	2024-07-07 20:10:14 handle.go:73 DEBUG 793968 request/s NumGoroutine 19
	2024-07-07 20:10:15 handle.go:73 DEBUG 786259 request/s NumGoroutine 19
	2024-07-07 20:10:16 handle.go:73 DEBUG 806728 request/s NumGoroutine 19
	2024-07-07 20:10:17 handle.go:73 DEBUG 775267 request/s NumGoroutine 19
	2024-07-07 20:10:18 handle.go:73 DEBUG 769429 request/s NumGoroutine 19
	2024-07-07 20:10:19 handle.go:73 DEBUG 795161 request/s NumGoroutine 19
	2024-07-07 20:10:20 handle.go:73 DEBUG 774706 request/s NumGoroutine 19
	2024-07-07 20:10:21 handle.go:73 DEBUG 781837 request/s NumGoroutine 19
	2024-07-07 20:10:22 handle.go:73 DEBUG 817479 request/s NumGoroutine 19