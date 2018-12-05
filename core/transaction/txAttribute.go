package transaction

import (
	"DNA/common"
	"DNA/common/serialization"
	. "DNA/errors"
	"errors"
	"io"
)

type TransactionAttributeUsage byte

const (
	Nonce          TransactionAttributeUsage = 0x00
	Script         TransactionAttributeUsage = 0x20
	DescriptionUrl TransactionAttributeUsage = 0x81
	Description    TransactionAttributeUsage = 0x90
)

func IsValidAttributeType(usage TransactionAttributeUsage) bool {
	return usage == Nonce || usage == Script ||
		usage == DescriptionUrl || usage == Description
}

type TxAttribute struct {
	Usage TransactionAttributeUsage
	Data  []byte
	Size  uint32
}

func NewTxAttribute(u TransactionAttributeUsage, d []byte) TxAttribute {
	tx := TxAttribute{u, d, 0}
	tx.Size = tx.GetSize()
	return tx
}

func (u *TxAttribute) GetSize() uint32 {
	if u.Usage == DescriptionUrl {
		return uint32(len([]byte{(byte(0xff))}) + len([]byte{(byte(0xff))}) + len(u.Data))
	}
	return 0
}

func (tx *TxAttribute) Serialize(w io.Writer) error {
	if err := serialization.WriteUint8(w, byte(tx.Usage)); err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction attribute Usage serialization error.")
	}
	if !IsValidAttributeType(tx.Usage) {
		return NewDetailErr(errors.New("[TxAttribute] error"), ErrNoCode, "Unsupported attribute Description.")
	}
	if err := serialization.WriteVarBytes(w, tx.Data); err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction attribute Data serialization error.")
	}
	return nil
}

func (tx *TxAttribute) Deserialize(r io.Reader) error {
	val, err := serialization.ReadBytes(r, 1)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction attribute Usage deserialization error.")
	}
	tx.Usage = TransactionAttributeUsage(val[0])
	if !IsValidAttributeType(tx.Usage) {
		return NewDetailErr(errors.New("[TxAttribute] error"), ErrNoCode, "Unsupported attribute Description.")
	}
	tx.Data, err = serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction attribute Data deserialization error.")
	}
	return nil

}

func (tx *TxAttribute) Deserialization(source *common.ZeroCopySource) error {
	val,eof := source.NextBytes(1)
	if eof {
		return io.ErrUnexpectedEOF
	}
	tx.Usage = TransactionAttributeUsage(val[0])
	if !IsValidAttributeType(tx.Usage) {
		return NewDetailErr(errors.New("[TxAttribute] error"), ErrNoCode, "Unsupported attribute Description.")
	}
	data, _, irregular, eof := source.NextVarBytes()
	tx.Data = data
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}
	return nil
}

func (tx *TxAttribute) Serialization(sink *common.ZeroCopySink) error {
	sink.WriteByte(byte(tx.Usage))
	if !IsValidAttributeType(tx.Usage) {
		return NewDetailErr(errors.New("[TxAttribute] error"), ErrNoCode, "Unsupported attribute Description.")
	}
	sink.WriteVarBytes(tx.Data)
	return nil
}
