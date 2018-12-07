package ledger

import (
	"DNA/common"
	"DNA/core/contract/program"
	tx "DNA/core/transaction"
	"DNA/crypto"
	"bytes"
	"github.com/magiconair/properties/assert"
	"testing"
)

func GenerateBlock() *Block {
	bs, _ := common.HexToBytes("dded139727732e13640c0acba31c459b4cf1e96720fa286ec62946f9b11c36a1")
	u, _ := common.Uint256ParseFromBytes(bs)
	pu, _ := common.ToScriptHash("AbPRaepcpBAFHz9zCj4619qch4Aq5hJARA")
	pp := &program.Program{
		Code:      []byte("1111"),
		Parameter: []byte("1111"),
	}
	txs := make([]*tx.Transaction, 0)
	tx := &tx.Transaction{}
	txbytes, _ := common.HexToBytes("8100047465737404737373730120047465737401dded139727732e13640c0acba31c459b4cf1e96720fa286ec62946f9b11c36a1010001dded139727732e13640c0acba31c459b4cf1e96720fa286ec62946f9b11c36a1010000000000000000000000000000000000000000000000000000000104313131310431313131")
	source := common.NewZeroCopySource(txbytes)
	tx.Deserialization(source)
	var tharray []common.Uint256
	for i := 0; i < 10000; i++ {
		txs = append(txs, tx)
		tharray = append(tharray, tx.Hash())
	}
	txRoot, _ := crypto.ComputeRoot(tharray)
	blockdata := &Blockdata{
		Version:          uint32(1),
		PrevBlockHash:    u,
		TransactionsRoot: txRoot,
		Timestamp:        uint32(1),
		Height:           uint32(1),
		ConsensusData:    uint64(1),
		NextBookKeeper:   pu,
		Program:          pp,
	}

	block := &Block{
		Blockdata:    blockdata,
		Transactions: txs,
	}
	return block
}

func TestBlock_Serialize(t *testing.T) {
	b := GenerateBlock()
	var buffer bytes.Buffer
	b.Serialize(&buffer)
	sink := new(common.ZeroCopySink)
	b.Serialization(sink)

	assert.Equal(t, buffer.Bytes(), sink.Bytes())

	source := common.NewZeroCopySource(sink.Bytes())
	block := &Block{}
	block.Deserialization(source)
	blockhash := block.Hash()
	bhash := b.Hash()
	assert.Equal(t, blockhash.ToArray(), bhash.ToArray())
}

func BenchmarkBlockSerialize(b *testing.B) {
	block := GenerateBlock()
	var buffer bytes.Buffer
	block.Serialize(&buffer)

	for i := 0; i < b.N; i++ {
		var buffer bytes.Buffer
		block.Serialize(&buffer)
	}
}
func BenchmarkBlockSerialization(b *testing.B) {
	block := GenerateBlock()

	for i := 0; i < b.N; i++ {
		sink := new(common.ZeroCopySink)
		block.Serialization(sink)
	}
}

func BenchmarkBlockDeserialize(b *testing.B) {
	block := GenerateBlock()
	sink := new(common.ZeroCopySink)
	block.Serialization(sink)
	block2 := &Block{}
	for i := 0; i < b.N; i++ {
		block2.Deserialize(bytes.NewReader(sink.Bytes()))
	}
}

func BenchmarkBlockDeserialization(b *testing.B) {

	block := GenerateBlock()
	sink := new(common.ZeroCopySink)
	block.Serialization(sink)
	block2 := &Block{}
	for i := 0; i < b.N; i++ {
		source := common.NewZeroCopySource(sink.Bytes())
		block2.Deserialization(source)
	}
}
