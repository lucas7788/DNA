package program

import (
	"DNA/common"
	"DNA/common/serialization"
	. "DNA/errors"
	"io"
)

type Program struct {

	//the contract program code,which will be run on VM or specific envrionment
	Code []byte

	//the program code's parameter
	Parameter []byte
}

//Serialize the Program
func (p *Program) Serialize(w io.Writer) error {
	err := serialization.WriteVarBytes(w, p.Parameter)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Execute Program Serialize Code failed.")
	}
	err = serialization.WriteVarBytes(w, p.Code)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Execute Program Serialize Parameter failed.")
	}

	return nil
}

func (p *Program) Serialization(sink *common.ZeroCopySink) error {
	sink.WriteVarBytes(p.Parameter)
	sink.WriteVarBytes(p.Code)
	return nil
}

func (p *Program) Deserialization(source *common.ZeroCopySource) error {
	var irregular, eof bool
	var data []byte
	data, _, irregular, eof = source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	p.Parameter = data
	data, _, irregular, eof = source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	p.Code = data
	if eof {
		return io.ErrUnexpectedEOF
	}
	return nil
}

//Deserialize the Program
func (p *Program) Deserialize(w io.Reader) error {
	val, err := serialization.ReadVarBytes(w)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Execute Program Deserialize Parameter failed.")
	}
	p.Parameter = val
	p.Code, err = serialization.ReadVarBytes(w)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Execute Program Deserialize Code failed.")
	}
	return nil
}
