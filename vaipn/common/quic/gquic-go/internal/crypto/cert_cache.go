package crypto

import (
	"hash/fnv"
	"sync"

	"github.com/payske-dev/vaipn-tunnel-core/vaipn/common/quic/gquic-go/internal/protocol"
	"github.com/golang/groupcache/lru"
)

var (
	// [Vaipn]
	// Replace github.com/hashicorp/golang-lru with github.com/golang/groupcache/lru,
	// adding mutex for safe concurrent access.
	compressedCertsCacheMutex sync.Mutex
	compressedCertsCache      *lru.Cache
)

func getCompressedCert(chain [][]byte, pCommonSetHashes, pCachedHashes []byte) ([]byte, error) {
	// Hash all inputs
	hasher := fnv.New64a()
	for _, v := range chain {
		hasher.Write(v)
	}
	hasher.Write(pCommonSetHashes)
	hasher.Write(pCachedHashes)
	hash := hasher.Sum64()

	var result []byte

	compressedCertsCacheMutex.Lock()
	resultI, isCached := compressedCertsCache.Get(hash)
	compressedCertsCacheMutex.Unlock()
	if isCached {
		result = resultI.([]byte)
	} else {
		var err error
		result, err = compressChain(chain, pCommonSetHashes, pCachedHashes)
		if err != nil {
			return nil, err
		}
		compressedCertsCacheMutex.Lock()
		compressedCertsCache.Add(hash, result)
		compressedCertsCacheMutex.Unlock()
	}

	return result, nil
}

func init() {
	compressedCertsCache = lru.New(protocol.NumCachedCertificates)
}
