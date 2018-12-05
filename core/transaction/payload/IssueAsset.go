package payload

import (
	"DNA/common"
	"io"
)


const IssueAssetPayloadVersion byte = 0x00

type IssueAsset struct {
}

func (a *IssueAsset) Data(version byte) []byte {
	//TODO: implement IssueAsset.Data()
	return []byte{0}

}

func (a *IssueAsset) Serialize(w io.Writer, version byte) error {
	return nil
}

func (a *IssueAsset) Deserialize(r io.Reader, version byte) error {
	return nil
}

func (a *IssueAsset) Serialization(sink *common.ZeroCopySink, version byte) error {
	return nil
}

func (a *IssueAsset) Deserialization(source *common.ZeroCopySource, version byte) error {

	return nil
}
