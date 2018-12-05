package crypto

import (
	"DNA/common"
	"github.com/magiconair/properties/assert"
	"math/big"
	"testing"
)

func TestPubKey_Serialization(t *testing.T) {
	b := new(big.Int)
	b.SetInt64(int64(100))
	p := &PubKey{
		X:b,
		Y:b,
	}
	sink := new(common.ZeroCopySink)
	p.Serialization(sink)

	p2 := &PubKey{}
	source := common.NewZeroCopySource(sink.Bytes())
	p2.DeSerialization(source)
	assert.Equal(t, p, p2)
}
