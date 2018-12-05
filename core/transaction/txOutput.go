package transaction

import (
	"DNA/common"
	"io"
)

type TxOutput struct {
	AssetID     common.Uint256
	Value       common.Fixed64
	ProgramHash common.Uint160
}

func (o *TxOutput) Serialize(w io.Writer) {
	o.AssetID.Serialize(w)
	o.Value.Serialize(w)
	o.ProgramHash.Serialize(w)
}

func (o *TxOutput) Deserialize(r io.Reader) {
	o.AssetID.Deserialize(r)
	o.Value.Deserialize(r)
	o.ProgramHash.Deserialize(r)
}
func (o *TxOutput) Deserialization(source *common.ZeroCopySource) error {
	val,eof := source.NextBytes(common.UINT256SIZE)
	if eof {
		return io.ErrUnexpectedEOF
	}
	copy(o.AssetID[:], val)
	if eof {
		return io.ErrUnexpectedEOF
	}
	data, eof := source.NextInt64()
	if eof {
		return io.ErrUnexpectedEOF
	}
	o.Value = common.Fixed64(data)
	val,eof = source.NextBytes(common.UINT160SIZE)
	if eof {
		return io.ErrUnexpectedEOF
	}
	copy(o.ProgramHash[:], val)
	return nil
}

func (o *TxOutput) Serialization(sink *common.ZeroCopySink) error {
	sink.WriteUint256(o.AssetID)
	sink.WriteFixed64(o.Value)
	sink.WriteUint160(o.ProgramHash)
	return nil
}
