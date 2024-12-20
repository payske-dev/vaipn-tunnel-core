package handshake

import (
	"crypto/x509"

	"github.com/payske-dev/vaipn-tunnel-core/vaipn/common/quic/gquic-go/internal/crypto"
	"github.com/payske-dev/vaipn-tunnel-core/vaipn/common/quic/gquic-go/internal/protocol"
	"github.com/bifurcation/mint"
)

// Sealer seals a packet
type Sealer interface {
	Seal(dst, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) []byte
	Overhead() int
}

// mintTLS combines some methods needed to interact with mint.
type mintTLS interface {
	crypto.TLSExporter
	Handshake() mint.Alert
}

// A TLSExtensionHandler sends and received the QUIC TLS extension.
// It provides the parameters sent by the peer on a channel.
type TLSExtensionHandler interface {
	Send(mint.HandshakeType, *mint.ExtensionList) error
	Receive(mint.HandshakeType, *mint.ExtensionList) error
	GetPeerParams() <-chan TransportParameters
}

type baseCryptoSetup interface {
	HandleCryptoStream() error
	ConnectionState() ConnectionState

	GetSealer() (protocol.EncryptionLevel, Sealer)
	GetSealerWithEncryptionLevel(protocol.EncryptionLevel) (Sealer, error)
	GetSealerForCryptoStream() (protocol.EncryptionLevel, Sealer)
}

// CryptoSetup is the crypto setup used by gQUIC
type CryptoSetup interface {
	baseCryptoSetup

	Open(dst, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) ([]byte, protocol.EncryptionLevel, error)
}

// CryptoSetupTLS is the crypto setup used by IETF QUIC
type CryptoSetupTLS interface {
	baseCryptoSetup

	OpenHandshake(dst, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) ([]byte, error)
	Open1RTT(dst, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) ([]byte, error)
}

// ConnectionState records basic details about the QUIC connection.
// Warning: This API should not be considered stable and might change soon.
type ConnectionState struct {
	HandshakeComplete bool                // handshake is complete
	ServerName        string              // server name requested by client, if any (server side only)
	PeerCertificates  []*x509.Certificate // certificate chain presented by remote peer
}
