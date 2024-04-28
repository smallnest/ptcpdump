package consumer

import (
	"context"
	"github.com/mozillazg/ptcpdump/bpf"
	"github.com/mozillazg/ptcpdump/internal/dev"
	"github.com/mozillazg/ptcpdump/internal/event"
	"github.com/mozillazg/ptcpdump/internal/writer"
	"log"
)

type PacketEventConsumer struct {
	writers []writer.PacketWriter
	devices map[int]dev.Device
}

func NewPacketEventConsumer(writers []writer.PacketWriter, devices map[int]dev.Device) *PacketEventConsumer {
	return &PacketEventConsumer{
		writers: writers,
		devices: devices,
	}
}

func (c *PacketEventConsumer) Start(ctx context.Context, ch <-chan bpf.BpfPacketEventT, maxPacketCount uint) {
	var n uint
	for {
		select {
		case <-ctx.Done():
			return
		case pt := <-ch:
			c.parsePacketEvent(pt)
			n++
			if maxPacketCount > 0 && n == maxPacketCount {
				log.Printf("%d packets captured", n)
				return
			}
		}
	}
}

func (c *PacketEventConsumer) parsePacketEvent(pt bpf.BpfPacketEventT) {
	pevent, err := event.ParsePacketEvent(c.devices, pt)
	if err != nil {
		log.Printf("[PacketEventConsumer] parse event failed: %s", err)
		return
	}

	for _, w := range c.writers {
		if err := w.Write(pevent); err != nil {
			log.Printf("[PacketEventConsumer] write packet failed: %s", err)
		}
		w.Flush()
	}
}

func (c *PacketEventConsumer) Stop() {

}
