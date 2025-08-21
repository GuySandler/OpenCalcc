# OpenClacc

If the program runs with more than 125mb of memory, close and reopen it (sometimes it starts with extra but this should not happen normally)

## custom fork of Mathcat:
added csc, sec, and cot

Xor bitwise from ^ to ^^

changed power from ** to ^

added the ability to make funcs with no set args to allow optional args

deg2rad()

rad2deg()

degree funcs with all 6 trig funcs


## main features:
Graphing

## features:
history

custom functions in math library

tracing x and y on graphs

graph intersections


## genuine memory optimization techniques:
optimized build command

scaling point generation for graphing

hard limiting memory

## build commands:

### basic build:
go build && ./opencalcc

### optimized build history:
GOMAXPROCS=4 GOMEMLIMIT=1GiB CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o openCalcc main.go && strip openCalcc

GOMAXPROCS=2 GOMEMLIMIT=512MiB GOGC=50 CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o openCalcc main.go && strip openCalcc && upx --best --lzma openCalcc

GOMAXPROCS=2 GOMEMLIMIT=120MiB GOGC=50 CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -o openCalcc main.go && strip openCalcc && upx --best --lzma openCalcc

GOMAXPROCS=1 GOMEMLIMIT=124MiB GOGC=25 CGO_ENABLED=1 go build -ldflags="-s -w" -trimpath -gcflags="-m=2" -o openCalcc main.go && strip -s openCalcc && upx --ultra-brute --best --lzma openCalcc

## known bugs:
doing float mode trig gives wrong answer

## Architecture decisions
I decided to make the main file in main.go using fyne becuase it's a well made desktop app library.

I tried to make my own math tool but I quickly found mathcat online so I forked it and added a few more features

## iamges
<img width="1458" height="1037" alt="Screenshot from 2025-08-20 18-44-33" src="https://github.com/user-attachments/assets/55abc4c0-f1be-47ac-afdd-9f2be7656611" />
<img width="1458" height="1037" alt="Screenshot from 2025-08-20 18-42-59" src="https://github.com/user-attachments/assets/a1c9207b-e697-4280-bf73-d268765a8ae8" />
