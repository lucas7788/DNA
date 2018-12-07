package asset

import (
	"DNA/common"
	"bytes"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestAsset_Serialize(t *testing.T) {
	a := &Asset{
		Name:"test",
		Description:"test",
		Precision:byte(1),
		AssetType:Token,
		RecordType:UTXO,
	}
	var buffer bytes.Buffer
	a.Serialize(&buffer)
	raw:=common.ToHexString(buffer.Bytes())
	sink := &common.ZeroCopySink{}
	a.Serialization(sink)
	newstr := common.ToHexString(sink.Bytes())
	assert.Equal(t, raw, newstr)

	adeser := Asset{}

	adeser.Deserialize(bytes.NewReader(buffer.Bytes()))

	adeserNew := Asset{}
	source := common.NewZeroCopySource(buffer.Bytes())
	adeserNew.Deserialization(source)

	assert.Equal(t, adeser, adeserNew)
}
