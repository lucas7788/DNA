package transaction

import (
	"DNA/common"
	"DNA/core/contract/program"
	"DNA/core/transaction/payload"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialization(t *testing.T) {
	attributes := make([]*TxAttribute, 0)
	attribute := &TxAttribute{
		Usage:Script,
		Data:[]byte("test"),
		Size:uint32(1),
		}
	attributes = append(attributes, attribute)
	bs, _ := common.HexToBytes("dded139727732e13640c0acba31c459b4cf1e96720fa286ec62946f9b11c36a1")
	u, _ :=common.Uint256ParseFromBytes(bs)
	inputs := make([]*UTXOTxInput, 0)
	input := &UTXOTxInput{
		ReferTxID:u,
		ReferTxOutputIndex: uint16(1),
	}
	inputs = append(inputs, input)
	pu,_ := common.ToScriptHash("d7239affb684c3c224476eb7bd52d9b2cb5e2aab")
	binputs := make([]*BalanceTxInput, 0)
	binput := &BalanceTxInput{
		AssetID:u,
		Value:common.Fixed64(int64(1)),
		ProgramHash:pu,
	}
	binputs = append(binputs, binput)
	outputs := make([]*TxOutput, 0)
	output := &TxOutput{
		AssetID:u,
		Value:common.Fixed64(int64(1)),
		ProgramHash:pu,
	}
	outputs = append(outputs, output)
	p := make([]*program.Program, 0)
	pp := &program.Program{
		Code:[]byte("1111"),
		Parameter:[]byte("1111"),
	}
	p = append(p, pp)
	mapoutput := make(map[common.Uint256][]*TxOutput)
	mapoutput[u] = outputs
	mapfixed := make(map[common.Uint256]common.Fixed64)
	mapfixed[u] = common.Fixed64(int64(1))
	mapoutputamount := make(map[common.Uint256] common.Fixed64)
	mapoutputamount[u] = common.Fixed64(int64(1))
	tx := &Transaction {
		TxType:Record,
		PayloadVersion: byte(0),
		Payload:&payload.Record{
			RecordType:"test",
			RecordData:[]byte("ssss"),
		},
		Attributes:attributes,
		UTXOInputs:inputs,
		BalanceInputs:binputs,
		Outputs:outputs,
		Programs:p,
		//AssetOutputs:mapoutput,
		//AssetInputAmount:mapfixed,
		//AssetOutputAmount:mapoutputamount,
	}
	var buffer bytes.Buffer
	if err := tx.Serialize(&buffer); err != nil {
		fmt.Println("serialization of issue transaction failed")
	}
	res := hex.EncodeToString(buffer.Bytes())
	fmt.Print(res)
	sink := new(common.ZeroCopySink)
	tx.Serialization(sink)
	res2 := hex.EncodeToString(sink.Bytes())
	assert.Equal(t, res, res2)

	txhash := tx.Hash()
	txraw := &Transaction{}
	txraw.Deserialize(bytes.NewReader(buffer.Bytes()))
	txrawhash := txraw.Hash()

	var buffer2 bytes.Buffer
	txraw.Serialize(&buffer2)
	assert.Equal(t, buffer2.Bytes(), buffer.Bytes())
	assert.Equal(t, txrawhash.ToString(), txhash.ToString())

	source := common.NewZeroCopySource(sink.Bytes())
	tx2 := &Transaction{}
	tx2.Deserialization(source)
	txhash2 := tx2.Hash()
	assert.Equal(t, txhash.ToString(), txhash2.ToString())

	assert.Equal(t, txrawhash.ToString(), txhash2.ToString())

}