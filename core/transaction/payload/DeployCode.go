package payload

import (
	"DNA/common/serialization"
	. "DNA/common"
	. "DNA/core/code"
	"io"
)

const DeployCodePayloadVersion byte = 0x00

type DeployCode struct {
	Code        *FunctionCode
	Name        string
	CodeVersion string
	Author      string
	Email       string
	Description string
}

func (dc *DeployCode) Data(version byte) []byte {
	// TODO: Data()

	return []byte{0}
}

func (dc *DeployCode) Serialize(w io.Writer, version byte) error {

	err := dc.Code.Serialize(w)
	if err != nil {
		return err
	}

	err = serialization.WriteVarString(w, dc.Name)
	if err != nil {
		return err
	}

	err = serialization.WriteVarString(w, dc.CodeVersion)
	if err != nil {
		return err
	}

	err = serialization.WriteVarString(w, dc.Author)
	if err != nil {
		return err
	}

	err = serialization.WriteVarString(w, dc.Email)
	if err != nil {
		return err
	}

	err = serialization.WriteVarString(w, dc.Description)
	if err != nil {
		return err
	}

	return nil
}

func (dc *DeployCode) Deserialize(r io.Reader, version byte) error {
	err := dc.Code.Deserialize(r)
	if err != nil {
		return err
	}

	dc.Name, err = serialization.ReadVarString(r)
	if err != nil {
		return err
	}

	dc.CodeVersion, err = serialization.ReadVarString(r)
	if err != nil {
		return err
	}

	dc.Author, err = serialization.ReadVarString(r)
	if err != nil {
		return err
	}

	dc.Email, err = serialization.ReadVarString(r)
	if err != nil {
		return err
	}

	dc.Description, err = serialization.ReadVarString(r)
	if err != nil {
		return err
	}

	return nil
}

func (dc *DeployCode) Serialization(sink *ZeroCopySink, version byte) error {
	err := dc.Code.Serialization(sink)
	if err != nil {
		return err
	}
	sink.WriteString(dc.Name)
	sink.WriteString(dc.CodeVersion)
	sink.WriteString(dc.Author)
	sink.WriteString(dc.Email)
	sink.WriteString(dc.Description)
	return nil
}

//note: DeployCode.Code has data reference of param source
func (dc *DeployCode) Deserialization(source *ZeroCopySource, version byte) error {
	err := dc.Code.Deserialization(source)
	if err != nil {
		return err
	}
	var eof, irregular bool
	dc.Name, _, irregular, eof = source.NextString()
	if irregular {
		return ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	dc.CodeVersion, _, irregular, eof = source.NextString()
	if irregular {
		return ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	dc.Author, _, irregular, eof = source.NextString()
	if irregular {
		return ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	dc.Email, _, irregular, eof = source.NextString()
	if irregular {
		return ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	dc.Description, _, irregular, eof = source.NextString()
	if irregular {
		return ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}