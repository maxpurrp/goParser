module parser

go 1.21.5

require github.com/PuerkitoBio/goquery v1.8.1

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	golang.org/x/net v0.17.0 // indirect
)

replace goParser/src/web/sendReq => ./web

replace goParser/src/carProcess/carParser => ./carProcess
