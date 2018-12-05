package payload

import (
	. "DNA/common"
	"DNA/common/serialization"
	"io"
)

const BookKeepingPayloadVersion byte = 0x03
const BookKeepingPayloadVersionBase byte = 0x02

type BookKeeping struct {
	Nonce uint64
}

func (a *BookKeeping) Data(version byte) []byte {
	return []byte{0}
}

func (a *BookKeeping) Serialize(w io.Writer, version byte) error {
	if version == BookKeepingPayloadVersionBase {
		return nil
	}
	err := serialization.WriteUint64(w, a.Nonce)
	if err != nil {
		return err
	}
	return nil
}

func (a *BookKeeping) Serialization(sink *ZeroCopySink, version byte) error {
	if version == BookKeepingPayloadVersionBase {
		return nil
	}
	sink.WriteUint64(uint64(a.Nonce))
	return nil
}

func (a *BookKeeping) Deserialize(r io.Reader, version byte) error {
	if version == BookKeepingPayloadVersionBase {
		return nil
	}
	var err error
	a.Nonce, err = serialization.ReadUint64(r)
	if err != nil {
		return err
	}
	return nil
}

func (a *BookKeeping) Deserialization(source *ZeroCopySource,version byte) error {
	if version == BookKeepingPayloadVersionBase {
		return nil
	}
	data, eof := source.NextUint64()
	if eof {
		return io.ErrUnexpectedEOF
	}
	a.Nonce = data
	return nil
}


