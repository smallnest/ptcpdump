package event

import (
	"github.com/mozillazg/ptcpdump/bpf"
	"strings"
)

type ProcessExec struct {
	Pid int

	Filename          string
	FilenameTruncated bool

	Args          []string
	ArgsTruncated bool
}

func ParseProcessExecEvent(event bpf.BpfExecEventT) (*ProcessExec, error) {
	var p ProcessExec
	if event.ArgsTruncated == 1 {
		p.ArgsTruncated = true
	}
	if event.FilenameTruncated == 1 {
		p.FilenameTruncated = true
	}
	p.Pid = int(event.Pid)
	bs := strings.Builder{}
	for i := 0; i < int(event.ArgsSize); i++ {
		b := byte(event.Args[i])
		if b == '\x00' {
			p.Args = append(p.Args, bs.String())
			bs.Reset()
		} else {
			bs.WriteByte(b)
		}
	}

	bs.Reset()
	for _, i := range event.Filename {
		b := byte(i)
		if b == '\x00' {
			break
		}
		bs.WriteByte(b)
	}
	p.Filename = bs.String()

	return &p, nil
}

func (p ProcessExec) FilenameStr() string {
	s := string(p.Filename)
	if p.FilenameTruncated {
		s += "..."
	}
	return s
}

func (p ProcessExec) ArgsStr() string {
	s := strings.Join(p.Args, " ")
	if p.ArgsTruncated {
		s += "..."
	}
	return s
}
