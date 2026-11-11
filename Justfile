build:
    go build -o renderer .

run: build
    ./renderer

cpu: build
    go tool pprof -http=:8080 cpu.prof

mem: build
    go tool pprof -http=:8081 mem.prof

bench:
    go test -bench=. -benchmem ./...

build-prof:
    go build -pgo=cpu.prof -o renderer-prof .

run-prof: build-prof
    ./renderer-prof