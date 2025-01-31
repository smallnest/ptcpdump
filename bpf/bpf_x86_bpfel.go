// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64

package bpf

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type BpfExecEventT struct {
	Pid               uint32
	FilenameTruncated uint8
	ArgsTruncated     uint8
	_                 [2]byte
	ArgsSize          uint32
	Filename          [512]int8
	Args              [4096]int8
}

type BpfFlowPidKeyT struct {
	Saddr [4]uint32
	Sport uint16
	_     [2]byte
}

type BpfFlowPidValueT struct{ Pid uint32 }

type BpfPacketEventT struct {
	Meta struct {
		Timestamp  uint64
		PacketType uint8
		_          [3]byte
		Ifindex    uint32
		Pid        uint32
		_          [4]byte
		PayloadLen uint64
		PacketSize uint64
	}
	Payload [1500]uint8
	_       [4]byte
}

// LoadBpf returns the embedded CollectionSpec for Bpf.
func LoadBpf() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_BpfBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load Bpf: %w", err)
	}

	return spec, err
}

// LoadBpfObjects loads Bpf and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*BpfObjects
//	*BpfPrograms
//	*BpfMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func LoadBpfObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := LoadBpf()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// BpfSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type BpfSpecs struct {
	BpfProgramSpecs
	BpfMapSpecs
}

// BpfSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type BpfProgramSpecs struct {
	KprobeSecuritySkClassifyFlow  *ebpf.ProgramSpec `ebpf:"kprobe__security_sk_classify_flow"`
	RawTracepointSchedProcessExec *ebpf.ProgramSpec `ebpf:"raw_tracepoint__sched_process_exec"`
	RawTracepointSchedProcessExit *ebpf.ProgramSpec `ebpf:"raw_tracepoint__sched_process_exit"`
	RawTracepointSchedProcessFork *ebpf.ProgramSpec `ebpf:"raw_tracepoint__sched_process_fork"`
	TcEgress                      *ebpf.ProgramSpec `ebpf:"tc_egress"`
	TcIngress                     *ebpf.ProgramSpec `ebpf:"tc_ingress"`
}

// BpfMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type BpfMapSpecs struct {
	ExecEventStack   *ebpf.MapSpec `ebpf:"exec_event_stack"`
	ExecEvents       *ebpf.MapSpec `ebpf:"exec_events"`
	FilterPidMap     *ebpf.MapSpec `ebpf:"filter_pid_map"`
	FlowPidMap       *ebpf.MapSpec `ebpf:"flow_pid_map"`
	PacketEventStack *ebpf.MapSpec `ebpf:"packet_event_stack"`
	PacketEvents     *ebpf.MapSpec `ebpf:"packet_events"`
}

// BpfObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to LoadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type BpfObjects struct {
	BpfPrograms
	BpfMaps
}

func (o *BpfObjects) Close() error {
	return _BpfClose(
		&o.BpfPrograms,
		&o.BpfMaps,
	)
}

// BpfMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to LoadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type BpfMaps struct {
	ExecEventStack   *ebpf.Map `ebpf:"exec_event_stack"`
	ExecEvents       *ebpf.Map `ebpf:"exec_events"`
	FilterPidMap     *ebpf.Map `ebpf:"filter_pid_map"`
	FlowPidMap       *ebpf.Map `ebpf:"flow_pid_map"`
	PacketEventStack *ebpf.Map `ebpf:"packet_event_stack"`
	PacketEvents     *ebpf.Map `ebpf:"packet_events"`
}

func (m *BpfMaps) Close() error {
	return _BpfClose(
		m.ExecEventStack,
		m.ExecEvents,
		m.FilterPidMap,
		m.FlowPidMap,
		m.PacketEventStack,
		m.PacketEvents,
	)
}

// BpfPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to LoadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type BpfPrograms struct {
	KprobeSecuritySkClassifyFlow  *ebpf.Program `ebpf:"kprobe__security_sk_classify_flow"`
	RawTracepointSchedProcessExec *ebpf.Program `ebpf:"raw_tracepoint__sched_process_exec"`
	RawTracepointSchedProcessExit *ebpf.Program `ebpf:"raw_tracepoint__sched_process_exit"`
	RawTracepointSchedProcessFork *ebpf.Program `ebpf:"raw_tracepoint__sched_process_fork"`
	TcEgress                      *ebpf.Program `ebpf:"tc_egress"`
	TcIngress                     *ebpf.Program `ebpf:"tc_ingress"`
}

func (p *BpfPrograms) Close() error {
	return _BpfClose(
		p.KprobeSecuritySkClassifyFlow,
		p.RawTracepointSchedProcessExec,
		p.RawTracepointSchedProcessExit,
		p.RawTracepointSchedProcessFork,
		p.TcEgress,
		p.TcIngress,
	)
}

func _BpfClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed bpf_x86_bpfel.o
var _BpfBytes []byte
