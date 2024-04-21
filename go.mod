module github.com/mozillazg/ptcpdump

go 1.21.7

require (
	github.com/cilium/ebpf v0.14.0
	github.com/florianl/go-tc v0.4.3
	github.com/gopacket/gopacket v1.2.0
	golang.org/x/sys v0.17.0
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/mdlayher/netlink v1.6.0 // indirect
	github.com/mdlayher/socket v0.1.1 // indirect
	golang.org/x/exp v0.0.0-20230224173230-c95f2b4c22f2 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
)

replace github.com/gopacket/gopacket => github.com/mozillazg/gopacket v0.0.0-20240420072046-71afeafe42df