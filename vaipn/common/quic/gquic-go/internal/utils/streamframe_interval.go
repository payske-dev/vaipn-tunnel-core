package utils

import "github.com/payske-dev/vaipn-tunnel-core/vaipn/common/quic/gquic-go/internal/protocol"

// ByteInterval is an interval from one ByteCount to the other
type ByteInterval struct {
	Start protocol.ByteCount
	End   protocol.ByteCount
}
