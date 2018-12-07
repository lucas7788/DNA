package payload

import (
	"DNA/common"
	"DNA/core/asset"
	"DNA/crypto"
	"bytes"
	"github.com/magiconair/properties/assert"
	"math/big"
	"testing"
	"DNA/core/code"
	. "DNA/core/contract"
)

func TestBookKeepingSerialize(t *testing.T) {
	b := &BookKeeping{
		Nonce:uint64(10),
	}
	var buffer bytes.Buffer
	b.Serialize(&buffer,byte(0))
	sink := new(common.ZeroCopySink)
	b.Serialization(sink, byte(0))
	assert.Equal(t, buffer.Bytes(), sink.Bytes())

	source := common.NewZeroCopySource(sink.Bytes())

	b2 := &BookKeeping{}
    b2.Deserialization(source,byte(0))
	assert.Equal(t, b, b2)
}
func TestDataFile(t *testing.T) {
	b:=new(big.Int)
	b.SetInt64(int64(100))
	d := &DataFile{
		IPFSPath:"test",
		Filename:"test",
		Note:"test",
		Issuer:&crypto.PubKey{
			X:b,
			Y:b,
		},
	}
	var buffer bytes.Buffer
	d.Serialize(&buffer, byte(0))
	sink := new(common.ZeroCopySink)
	d.Serialization(sink, byte(0))
	assert.Equal(t, buffer.Bytes(), sink.Bytes())

	d2 := &DataFile{}
	source := common.NewZeroCopySource(sink.Bytes())
	d2.Deserialization(source, byte(0))
	assert.Equal(t, d, d2)
}

func TestDeployCode(t *testing.T) {
	f := &code.FunctionCode{
		Code: []byte("test"),

		// Contract parameter type list
		ParameterTypes: []ContractParameterType{Boolean},
	}
	d := &DeployCode{
		Code:        f,
		Name:        "test",
		CodeVersion: "test",
		Author:      "test",
		Email:       "test",
		Description: "test",
	}

	var buffer bytes.Buffer
	d.Serialize(&buffer, byte(0))

	sink := new(common.ZeroCopySink)
	d.Serialization(sink, byte(0))
	assert.Equal(t, buffer.Bytes(), sink.Bytes())


	d2 := &DeployCode{}

	source := common.NewZeroCopySource(sink.Bytes())
	d2.Deserialization(source, byte(0))

	assert.Equal(t, d, d2)

}

func TestPrivacyPayload(t *testing.T) {
	b := new(big.Int)
	b.SetInt64(int64(1))
	e := &EcdhAes256{
		FromPubkey: &crypto.PubKey{b,b},
		ToPubkey:   &crypto.PubKey{b, b},
		Nonce:      []byte{1},
	}
	p := &PrivacyPayload{
		PayloadType: EncryptedPayloadType(byte(20)),
		Payload:     EncryptedPayload([]byte{10}),
		EncryptType: PayloadEncryptType(byte(1)),
		EncryptAttr: e,
	}

	var buffer bytes.Buffer
	p.Serialize(&buffer, byte(0))

	sink := new(common.ZeroCopySink)
	p.Serialization(sink, byte(0))
	assert.Equal(t, buffer.Bytes(), sink.Bytes())

	p2 := &PrivacyPayload{}
	source := common.NewZeroCopySource(sink.Bytes())
	p2.Deserialization(source, byte(0))

	assert.Equal(t, p, p2)
}

func TestRecord(t *testing.T) {
	r := &Record{
		RecordType: "test",
		RecordData: []byte("111"),
	}

	var buffer bytes.Buffer
	r.Serialize(&buffer, byte(0))
	sink := new(common.ZeroCopySink)
	r.Serialization(sink, byte(0))

	r2 := &Record{}
	source := common.NewZeroCopySource(sink.Bytes())
	r2.Deserialization(source, byte(0))
	assert.Equal(t, r, r2)
}

func TestRegisterAsset(t *testing.T) {
	b := new(big.Int)
	b.SetInt64(int64(1))
	a := &asset.Asset{
		Name:        "test",
		Description: "test",
		Precision:   byte(1),
		AssetType:   asset.Token,
		RecordType:  asset.Balance,
	}
	pu,_ := common.ToScriptHash("d7239affb684c3c224476eb7bd52d9b2cb5e2aab")
	r := &RegisterAsset{
		Asset:  a,
		Amount: common.Fixed64(int64(100)),
		//Precision  byte
		Issuer:     &crypto.PubKey{b,b},
		Controller: pu,
	}

	var buffer bytes.Buffer
	r.Serialize(&buffer, byte(1))

	sink := new(common.ZeroCopySink)
	r.Serialization(sink, byte(1))

	assert.Equal(t, buffer.Bytes(), sink.Bytes())

    r2 := &RegisterAsset{}
	source := common.NewZeroCopySource(sink.Bytes())
	r2.Deserialization(source, byte(1))

	assert.Equal(t, r, r2)
}