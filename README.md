# sdc

Package sdc implements simple dense coding which compresses an array of integers
and allows random access in constant time.

The algorithm is taken from *Simple dense coding* described in [Simple
Random Access Compression][1] by Kimmo Fredriksson and Fedor Nikitin. It
uses an auxiliary structure to support fast select operations from [Fast,
Small, Simple Rank/Select on Bitmaps][2].

[1]: http://cs.uef.fi/~fredriks/pub/papers/fi09.pdf
[2]: http://dcc.uchile.cl/~gnavarro/ps/sea12.1.pdf

## Installation
```sh
go get github.com/robskie/sdc
```

## API Reference

Godoc documentation can be found [here][3].

[3]:https://godoc.org/github.com/robskie/sdc

## Benchmarks

These benchmarks are done on a Core i5 at 2.3GHz. You can run these benchmarks
by typing ```go test github.com/robskie/sdc -bench=.*``` from terminal.

```
BenchmarkAdd-4  20000000            87.5 ns/op
BenchmarkGet-4   5000000             264 ns/op
```
