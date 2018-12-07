package payload

import (
	"DNA/common"
	"DNA/crypto"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestBookKeeper_Serialize(t *testing.T) {
    bint := new(big.Int)
    bint.SetInt64(int64(100))
	p := &crypto.PubKey{
		X:bint,
		Y:bint,
	}
	a := &BookKeeper{
		PubKey:p,
		Action:BookKeeperAction_ADD,
		Cert:[]byte("test"),
		Issuer:p,
	}
	var buffer bytes.Buffer
	a.Serialize(&buffer, byte(0))
	raw:=common.ToHexString(buffer.Bytes())
	sink := &common.ZeroCopySink{}
	a.Serialization(sink,byte(0))
	newstr := common.ToHexString(sink.Bytes())
	assert.Equal(t, raw, newstr)

	book := &crypto.PubKey{}

	book.DeSerialize(bytes.NewReader(buffer.Bytes()))
	fmt.Print(book)
}
