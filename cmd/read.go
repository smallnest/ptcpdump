package cmd

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/mozillazg/ptcpdump/internal/metadata"
	"github.com/mozillazg/ptcpdump/internal/parser"
	"github.com/mozillazg/ptcpdump/internal/writer"
)

func read(ctx context.Context, opts Options) error {
	f, err := os.Open(opts.ReadPath())
	if err != nil {
		return err
	}
	defer f.Close()

	var p parser.Parser
	pcache := metadata.NewProcessCache()
	stdoutWriter := writer.NewStdoutWriter(os.Stdout, pcache)
	ext := filepath.Ext(opts.ReadPath())

	switch ext {
	case extPcap:
		pr, err := parser.NewPcapParser(f)
		if err != nil {
			return err
		}
		stdoutWriter.Decoder = pr.Decoder()
		p = pr
		break
	default:
		p, err = parser.NewPcapNGParser(f, pcache)
		if err != nil {
			return err
		}
	}

	for {
		e, err := p.Parse()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err := stdoutWriter.Write(e); err != nil {
			return err
		}
	}

	return nil
}
