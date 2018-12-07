package transaction

import (
	"DNA/common"
	"DNA/common/serialization"
	"fmt"
	"io"
)

type UTXOTxInput struct {

	//Indicate the previous Tx which include the UTXO output for usage
	ReferTxID common.Uint256

	//The index of output in the referTx output list
	ReferTxOutputIndex uint16
}

func (ui *UTXOTxInput) Serialize(w io.Writer) {
	ui.ReferTxID.Serialize(w)
	serialization.WriteUint16(w, ui.ReferTxOutputIndex)
}

func (ui *UTXOTxInput) Serialization(sink *common.ZeroCopySink) error {
	sink.WriteBytes(ui.ReferTxID.ToArray())
	sink.WriteUint16(ui.ReferTxOutputIndex)
	return nil
}

func (ui *UTXOTxInput) Deserialize(r io.Reader) error {
	//referTxID
	err := ui.ReferTxID.Deserialize(r)
	if err != nil {
		return err
	}

	//Output Index
	temp, err := serialization.ReadUint16(r)
	ui.ReferTxOutputIndex = uint16(temp)
	if err != nil {
		return err
	}

	return nil
}

func (ui *UTXOTxInput) Deserialization(source *common.ZeroCopySource) error {
	val, eof := source.NextBytes(common.UINT256SIZE)
	if eof {
		return io.ErrUnexpectedEOF
	}
	copy(ui.ReferTxID[:], val)
	if eof {
		return io.ErrUnexpectedEOF
	}
	ui.ReferTxOutputIndex, eof = source.NextUint16()
	if eof {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (ui *UTXOTxInput) ToString() string {
	return fmt.Sprintf("%x%x", ui.ReferTxID.ToString(), ui.ReferTxOutputIndex)
}

func (ui *UTXOTxInput) Equals(other *UTXOTxInput) bool {
	if ui == other {
		return true
	}
	if other == nil {
		return false
	}
	if ui.ReferTxID == other.ReferTxID && ui.ReferTxOutputIndex == other.ReferTxOutputIndex {
		return true
	} else {
		return false
	}
}
