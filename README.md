# cowmap

A lock free single writer multiple reader map using atomic pointer operations.

# Why Cow?
Initials of [Copy-on-write](https://en.wikipedia.org/wiki/Copy-on-write)

# Usecase
cowmap provides a faster concurrent read operation for single writer multiple reader cases where write operations are known to be rare. It creates huge amounts of garbage, so make sure to memory profile your case before using it.

# How
It reads the current map atomically, creates a copy of the map, modifies it, then swaps the map pointers atomically.

# Tests
Unit tests require race detector to be enabled. It only tests if there is a data race occurrence during concurrent read/write operations.

# Benchmarks
```
$ go test -run="#" -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: cowmap
Benchmark_Maps/concurrent_read/copy_on_write_map-10         	413894556	         2.512 ns/op	       0 B/op	       0 allocs/op
Benchmark_Maps/concurrent_read/rwmutex_map-10               	 8570545	       141.0 ns/op	       0 B/op	       0 allocs/op
Benchmark_Maps/single_write/copy_on_write_map-10            	   10000	    662256 ns/op	  705642 B/op	      92 allocs/op
Benchmark_Maps/single_write/rwmutex_map-10                  	 3820129	       403.7 ns/op	     186 B/op	       2 allocs/op
PASS
ok  	cowmap	11.532s
```
