package code

import (
	"DNA/common"
	. "DNA/core/contract"
	"bytes"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFuncCode(t *testing.T) {
	f := &FunctionCode{
		Code: []byte("test"),

		// Contract parameter type list
		ParameterTypes: []ContractParameterType{Boolean},

	}
	var buffer bytes.Buffer
	f.Serialize(&buffer)
	sink := new(common.ZeroCopySink)
	f.Serialization(sink)

	assert.Equal(t, buffer.Bytes(), sink.Bytes())

	f2 := &FunctionCode{}
	source := common.NewZeroCopySource(sink.Bytes())
	f2.Deserialization(source)
	assert.Equal(t, f, f2)
}
