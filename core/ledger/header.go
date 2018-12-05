package ledger

import (
	"DNA/common"
	"io"
)

type Header struct {
	Blockdata *Blockdata
}

//Serialize the blockheader
func (h *Header) Serialize(w io.Writer) {
	h.Blockdata.Serialize(w)
	w.Write([]byte{'0'})

}

func (h *Header) Serialization(sink *common.ZeroCopySink) error {
	h.Blockdata.Serialization(sink)
	sink.WriteByte(byte(0))
	return nil
}

func (h *Header) Deserialization(source *common.ZeroCopySource) error {
	h.Blockdata.Deserialization(source)
	_, eof := source.NextByte()
    if eof {
    	return io.ErrUnexpectedEOF
	}
	return nil
}

func (h *Header) Deserialize(r io.Reader) error {
	header := new(Blockdata)
	header.Deserialize(r)
	h.Blockdata = header
	var headerFlag [1]byte
	_, err := io.ReadFull(r, headerFlag[:])
	if err != nil {
		return err
	}
	return nil
}
