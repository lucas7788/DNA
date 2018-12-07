package transaction

import (
	"DNA/common"
	"bytes"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestTxOutput(t *testing.T) {
	bs, _ := common.HexToBytes("dded139727732e13640c0acba31c459b4cf1e96720fa286ec62946f9b11c36a1")
	u, _ := common.Uint256ParseFromBytes(bs)
	pu, _ := common.ToScriptHash("d7239affb684c3c224476eb7bd52d9b2cb5e2aab")
	tx := &TxOutput{
		AssetID:     u,
		Value:       common.Fixed64(int64(100)),
		ProgramHash: pu,
	}

	var buffer bytes.Buffer
	tx.Serialize(&buffer)

	sink := new(common.ZeroCopySink)
	tx.Serialization(sink)
	assert.Equal(t, buffer.Bytes(), sink.Bytes())

	tx2 := &TxOutput{}
	source := common.NewZeroCopySource(sink.Bytes())
	tx2.Deserialization(source)
	assert.Equal(t, tx, tx2)
}

func TestUTXOTxInput_Deserialization(t *testing.T) {
	bs, _ := common.HexToBytes("dded139727732e13640c0acba31c459b4cf1e96720fa286ec62946f9b11c36a1")
	u, _ := common.Uint256ParseFromBytes(bs)
	utxo := &UTXOTxInput{
		//Indicate the previous Tx which include the UTXO output for usage
		ReferTxID: u,

		//The index of output in the referTx output list
		ReferTxOutputIndex: uint16(1),
	}

	var buffer bytes.Buffer
	utxo.Serialize(&buffer)

	sink := new(common.ZeroCopySink)
	utxo.Serialization(sink)
	assert.Equal(t, buffer.Bytes(), sink.Bytes())

	utxo2 := &UTXOTxInput{}

	source := common.NewZeroCopySource(sink.Bytes())
	utxo2.Deserialization(source)
	assert.Equal(t, utxo2, utxo)
}
