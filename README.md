# OpenClacc

custom fork of Mathcat:
added csc, sec, and cot
Xor bitwise from ^ to ^^
changed power from ** to ^
added the ability to make funcs with no set args to allow optional args
deg2rad()
rad2deg()

main features:
Graphing

features:
history
custom functions in math library
tracing x and y on graphs
graph intersections


genuine memory optimization techniques:
optimized build command
scaling point generation for graphing


build commands:

basic build:
go build && ./opencalcc

optimized build:
GOMAXPROCS=4 GOMEMLIMIT=1GiB CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o openCalcc main.go && strip openCalcc


