package code

import (
	. "DNA/common"
	"DNA/common/log"
	"DNA/common/serialization"
	. "DNA/core/contract"
	"fmt"
	"io"
)

type FunctionCode struct {
	// Contract Code
	Code []byte

	// Contract parameter type list
	ParameterTypes []ContractParameterType

	// Contract return type list
	ReturnTypes []ContractParameterType
}

// method of SerializableData
func (fc *FunctionCode) Serialize(w io.Writer) error {
	err := serialization.WriteVarBytes(w, ContractParameterTypeToByte(fc.ParameterTypes))
	if err != nil {
		return err
	}

	err = serialization.WriteVarBytes(w, fc.Code)
	if err != nil {
		return err
	}

	return nil
}

func (fc *FunctionCode) Serialization(sink *ZeroCopySink) error {
	sink.WriteVarBytes(ContractParameterTypeToByte(fc.ParameterTypes))
	sink.WriteVarBytes(fc.Code)
	return nil
}

func (fc *FunctionCode) Deserialization(source *ZeroCopySource) error {
	var eof, irregular bool
	var data []byte
	data, _, irregular, eof = source.NextVarBytes()
	if irregular {
		return ErrIrregularData
	}
	fc.ParameterTypes = ByteToContractParameterType(data)
	data, _, irregular, eof = source.NextVarBytes()
	if irregular {
		return ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	fc.Code = data
	return nil
}

// method of SerializableData
func (fc *FunctionCode) Deserialize(r io.Reader) error {
	p, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	fc.ParameterTypes = ByteToContractParameterType(p)

	fc.Code, err = serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}

	return nil
}

// method of ICode
// Get code
func (fc *FunctionCode) GetCode() []byte {
	return fc.Code
}

// method of ICode
// Get the list of parameter value
func (fc *FunctionCode) GetParameterTypes() []ContractParameterType {
	return fc.ParameterTypes
}

// method of ICode
// Get the list of return value
func (fc *FunctionCode) GetReturnTypes() []ContractParameterType {
	return fc.ReturnTypes
}

// method of ICode
// Get the hash of the smart contract
func (fc *FunctionCode) CodeHash() Uint160 {
	hash, err := ToCodeHash(fc.Code)
	if err != nil {
		log.Debug(fmt.Sprintf("[FunctionCode] ToCodeHash err=%s", err))
		return Uint160{0}
	}

	return hash
}
