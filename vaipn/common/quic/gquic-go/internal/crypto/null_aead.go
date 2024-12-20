package crypto

import "github.com/payske-dev/vaipn-tunnel-core/vaipn/common/quic/gquic-go/internal/protocol"

// NewNullAEAD creates a NullAEAD
func NewNullAEAD(p protocol.Perspective, connID protocol.ConnectionID, v protocol.VersionNumber) (AEAD, error) {
	if v.UsesTLS() {
		return newNullAEADAESGCM(connID, p)
	}
	return &nullAEADFNV128a{perspective: p}, nil
}
