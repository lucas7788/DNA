package payload

import (
	"DNA/common"
	"io"
)

const TransferAssetayloadVersion byte = 0x00

type TransferAsset struct {
}

func (a *TransferAsset) Data(version byte) []byte {
	//TODO: implement TransferAsset.Data()
	return []byte{0}

}

func (a *TransferAsset) Serialize(w io.Writer, version byte) error {
	return nil
}

func (a *TransferAsset) Deserialize(r io.Reader, version byte) error {
	return nil
}

func (a *TransferAsset) Serialization(sink *common.ZeroCopySink, version byte) error  {
	return nil
}

func (a *TransferAsset) Deserialization(source *common.ZeroCopySource, version byte) error {
	return nil
}