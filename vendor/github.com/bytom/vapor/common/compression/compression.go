package compression

import "fmt"

// Compression is intterface
type Compression interface {
	CompressBytes(data []byte) []byte
	DecompressBytes(data []byte) ([]byte, error)
}

const (
	SnappyBackendStr = "snappy" // legacy, defaults to SnappyBackendStr.
)

type compressionCreator func() Compression

var backends = map[string]compressionCreator{}

func registerCompressionCreator(backend string, creator compressionCreator, force bool) {
	_, ok := backends[backend]
	if !force && ok {
		return
	}
	backends[backend] = creator
}

func NewCompression(backend string) Compression {
	compression, ok := backends[backend]
	if !ok {
		panic(fmt.Sprintf("Cannot find compression algorithm:[%s]", backend))

	}
	return compression()
}
