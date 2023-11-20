module github.com/netsys-lab/scion-host

go 1.20

require (
	github.com/jessevdk/go-flags v1.5.0
	github.com/netsec-ethz/bootstrapper v0.0.7
)

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/grandcat/zeroconf v1.0.0 // indirect
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac // indirect
	github.com/insomniacslk/dhcp v0.0.0-20211209223715-7d93572ebe8e // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mdlayher/ndp v0.10.0 // indirect
	github.com/miekg/dns v1.1.27 // indirect
	github.com/pelletier/go-toml v1.8.1-0.20200708110244-34de94e6a887 // indirect
	github.com/u-root/uio v0.0.0-20210528114334-82958018845c // indirect
	gitlab.com/golang-commonmark/puny v0.0.0-20191124015043-9f83538fa04f // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.6.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)

replace github.com/netsec-ethz/bootstrapper => ./dev/bootstrapper
