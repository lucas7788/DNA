package ledger

import (
	"DNA/common"
	"bytes"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestSerialize(t *testing.T) {
	block := GenerateBlock()
	header := &Header{
		Blockdata: block.Blockdata,
	}
	var buffer bytes.Buffer
	header.Serialize(&buffer)

	sink := new(common.ZeroCopySink)
	header.Serialization(sink)
	assert.Equal(t, buffer.Bytes(), sink.Bytes())

	header2 := &Header{}
	source := common.NewZeroCopySource(sink.Bytes())
	header2.Deserialization(source)
	assert.Equal(t, header, header2)
}
