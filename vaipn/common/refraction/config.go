/*
 * Copyright (c) 2021, Vaipn Inc.
 * All rights reserved.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package refraction

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/payske-dev/vaipn-tunnel-core/vaipn/common"
)

// ConjureConfig specifies the additional configuration for a Conjure dial.
type ConjureConfig struct {

	// RegistrationCacheTTL specifies how long to retain a successful Conjure
	// registration for reuse in a subsequent dial. This value should be
	// synchronized with the Conjure station configuration. When
	// RegistrationCacheTTL is 0, registrations are not cached.
	RegistrationCacheTTL time.Duration

	// RegistrationCacheKey defines a scope or affinity for cached Conjure
	// registrations. For example, the key can reflect the target Vaipn server
	// as well as the current network ID. This ensures that any replay will
	// always use the same cached registration, including its phantom IP(s). And
	// ensures that the cache scope is restricted to the current network: when
	// the network changes, the client's public IP changes, and previous
	// registrations will become invalid. When the client returns to the original
	// network, the previous registrations may be valid once again (assuming
	// the client reverts back to its original public IP).
	RegistrationCacheKey string

	// APIRegistrarBidirectionalURL specifies the bidirectional API
	// registration endpoint. Setting APIRegistrarBidirectionalURL enables
	// API registration. The domain fronting configuration provided by
	// APIRegistrarHTTPClient may ignore the host portion of this URL,
	// implicitly providing another value; the path portion is always used in
	// the request. Only one of API registration or decoy registration can be
	// enabled for a single dial.
	APIRegistrarBidirectionalURL string

	// APIRegistrarHTTPClient specifies a custom HTTP client (and underlying
	// dialers) to be used for Conjure API registration. The
	// APIRegistrarHTTPClient enables domain fronting of API registration web
	// requests. This parameter is required when API registration is enabled.
	APIRegistrarHTTPClient *http.Client

	// APIRegistrarDelay specifies how long to wait after a successful API
	// registration before initiating the phantom dial(s), as required by the
	// Conjure protocol. This value depends on Conjure station operations and
	// should be synchronized with the Conjure station configuration.
	APIRegistrarDelay time.Duration

	// DoDecoyRegistration indicates to use decoy registration instead of API
	// registration. Only one of API registration or decoy registration can
	// be enabled for a single dial.
	DoDecoyRegistration bool

	// DecoyRegistrarWidth specifies how many decoys to use per registration.
	DecoyRegistrarWidth int

	// DecoyRegistrarDelay specifies how long to wait after a successful API
	// registration before initiating the phantom dial(s), as required by the
	// Conjure protocol.
	//
	// Limitation: this value is not exposed by gotapdance and is currently
	// ignored.
	DecoyRegistrarDelay time.Duration

	// EnableIPv6Dials specifies whether to attempt to dial an IPv6 phantom in
	// addition to and concurrent with an IPv4 phantom dial.
	EnableIPv6Dials bool

	// EnablePortRandomization specifies whether to enable destination port
	// randomization.
	EnablePortRandomization bool

	// EnableRegistrationOverrides specifies whether to allow the Conjure
	// system to provide parameter overrides, such as alternative prefixes,
	// in the registration response.
	EnableRegistrationOverrides bool

	// Transport may be protocol.CONJURE_TRANSPORT_MIN_OSSH,
	// protocol.CONJURE_TRANSPORT_PREFIX_OSSH, or
	// protocol.CONJURE_TRANSPORT_DTLS_OSSH.
	Transport string

	// STUNServerAddress specifies the STUN server to use with
	// protcol.CONJURE_TRANSPORT_DTLS_OSSH.
	STUNServerAddress string

	// DTLSEmptyInitialPacket specifies whether to prefix the DTLS flow with
	// an initial empty packet. Used only for
	// protocol.CONJURE_TRANSPORT_DTLS_OSSH.
	DTLSEmptyInitialPacket bool

	// DiagnosticID identifies this dial in diagnostics.
	DiagnosticID string

	// Logger is used for logging diagnostics.
	Logger common.Logger
}

// Dialer is the dialer function type expected by gotapdance.
type Dialer func(ctx context.Context, network, laddr, raddr string) (net.Conn, error)
