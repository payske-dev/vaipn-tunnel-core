package gquic

import (
	"sync"

	"github.com/payske-dev/vaipn-tunnel-core/vaipn/common/quic/gquic-go/internal/protocol"
)

var bufferPool sync.Pool

func getPacketBuffer() *[]byte {
	return bufferPool.Get().(*[]byte)
}

func putPacketBuffer(buf *[]byte) {
	if cap(*buf) != int(protocol.MaxReceivePacketSize) {
		panic("putPacketBuffer called with packet of wrong size!")
	}
	bufferPool.Put(buf)
}

func init() {
	bufferPool.New = func() interface{} {
		b := make([]byte, 0, protocol.MaxReceivePacketSize)
		return &b
	}
}
