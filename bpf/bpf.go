package bpf

import (
	"encoding/binary"
	"log"
	"net"
	"unsafe"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"golang.org/x/sys/unix"
	"golang.org/x/xerrors"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -no-strip -target native -type packet_event_t Bpf ./ptcpdump.c -- -I./headers -I. -Wall

const tcFilterName = "ptcpdump"

type BPF struct {
	spec       *ebpf.CollectionSpec
	objs       *BpfObjects
	links      []link.Link
	closeFuncs []func()
}

func NewBPF() (*BPF, error) {
	spec, err := LoadBpf()
	if err != nil {
		return nil, err
	}
	return &BPF{
		spec: spec,
		objs: &BpfObjects{},
	}, nil
}

func (b *BPF) Load() error {
	err := b.spec.LoadAndAssign(b.objs, &ebpf.CollectionOptions{
		Programs: ebpf.ProgramOptions{
			LogLevel: ebpf.LogLevelInstruction,
			LogSize:  ebpf.DefaultVerifierLogSize * 8,
		},
	})
	return err
}

func (b *BPF) Close() {
	for _, lk := range b.links {
		if err := lk.Close(); err != nil {
			log.Printf("[bpf] close link %v failed: %+v", lk, err)
		}
	}
	for i := len(b.closeFuncs) - 1; i > 0; i-- {
		f := b.closeFuncs[i]
		f()
	}
	if err := b.objs.Close(); err != nil {
		log.Printf("[bpf] close objects failed: %+v", err)
	}
}

func (b *BPF) AttachKprobes() error {
	lk, err := link.Kprobe("security_sk_classify_flow",
		b.objs.KprobeSecuritySkClassifyFlow, &link.KprobeOptions{})
	if err != nil {
		return xerrors.Errorf("attach kprobe/security_sk_classify_flow: %w", err)
	}
	b.links = append(b.links, lk)
	return nil
}

func (b *BPF) AttachTcHooks(dev *net.Interface) error {
	closeFunc, err := ensureTcQdisc(dev)
	if err != nil {
		if closeFunc != nil {
			closeFunc()
		}
		return xerrors.Errorf("attach tc hooks: %w", err)
	}

	c1, err := attachTcHook(dev, b.objs.TcEgress, false)
	if err != nil {
		closeFunc()
		return xerrors.Errorf("attach tc hooks: %w", err)
	}

	c2, err := attachTcHook(dev, b.objs.TcIngress, true)
	if err != nil {
		c1()
		closeFunc()
		return xerrors.Errorf("attach tc hooks: %w", err)
	}

	b.closeFuncs = append(b.closeFuncs, closeFunc, c1, c2)
	return nil
}

func (b *BPF) NewPacketEventReader() (*perf.Reader, error) {
	reader, err := perf.NewReader(b.objs.PacketEvents, 1500*1000)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return reader, nil
}

func attachTcHook(dev *net.Interface, prog *ebpf.Program, ingress bool) (func(), error) {
	tcnl, err := tc.Open(&tc.Config{})
	if err != nil {
		return nil, err
	}
	closeFunc := func() {
		if err := tcnl.Close(); err != nil {
			log.Printf("tcnl.Close() failed: %+v", err)
		}
	}

	fd := uint32(prog.FD())
	name := tcFilterName
	parent := tc.HandleMinEgress
	if ingress {
		parent = tc.HandleMinIngress
	}

	filter := tc.Object{
		tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(dev.Index),
			Handle:  0,
			Parent:  core.BuildHandle(tc.HandleRoot, parent),
			Info:    1<<16 | uint32(htons(unix.ETH_P_ALL)),
		},
		tc.Attribute{
			Kind: "bpf",
			BPF: &tc.Bpf{
				FD:   &fd,
				Name: &name,
			},
		},
	}
	if err := tcnl.Filter().Add(&filter); err != nil {
		return closeFunc, xerrors.Errorf("add filter: %w", err)
	}

	newCloseFunc := func() {
		if err := tcnl.Filter().Delete(&filter); err != nil {
			log.Printf("delete tcnl filter failed: %+v", err)
		}
		closeFunc()
	}
	return newCloseFunc, nil
}

func ensureTcQdisc(dev *net.Interface) (func(), error) {
	tcnl, err := tc.Open(&tc.Config{})
	if err != nil {
		return nil, err
	}
	closeFunc := func() {
		if err := tcnl.Close(); err != nil {
			log.Printf("tcnl.Close() failed: %+v", err)
		}
	}

	qdisc := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(dev.Index),
			Handle:  core.BuildHandle(tc.HandleRoot, 0x0000),
			Parent:  tc.HandleIngress,
		},
		Attribute: tc.Attribute{
			Kind: "clsact",
		},
	}

	if err := tcnl.Qdisc().Replace(&qdisc); err != nil {
		return closeFunc, err
	}

	newCloseFunc := func() {
		if err := tcnl.Qdisc().Delete(&qdisc); err != nil {
			log.Printf("delete tcnl qdisc failed: %+v", err)
		}
		closeFunc()
	}

	return newCloseFunc, nil
}

func htons(n uint16) uint16 {
	b := *(*[2]byte)(unsafe.Pointer(&n))
	return binary.BigEndian.Uint16(b[:])
}
