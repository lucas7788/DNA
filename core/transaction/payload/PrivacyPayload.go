package payload

import (
	. "DNA/common"
	"DNA/common/serialization"
	"DNA/crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"github.com/ontio/ontology/common"
	"io"
	mrand "math/rand"
	"time"
)

const PrivacyPayloadVersion byte = 0x00

type EncryptedPayloadType byte

type EncryptedPayload []byte

type PayloadEncryptType byte

type PayloadEncryptAttr interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
	Serialization(sink *ZeroCopySink) error
	Deserialization(source *ZeroCopySource) error
	Encrypt(msg []byte, keys interface{}) ([]byte, error)
	Decrypt(msg []byte, keys interface{}) ([]byte, error)
}

const (
	ECDH_AES256 PayloadEncryptType = 0x01
)
const (
	RawPayload EncryptedPayloadType = 0x01
)

type PrivacyPayload struct {
	PayloadType EncryptedPayloadType
	Payload     EncryptedPayload
	EncryptType PayloadEncryptType
	EncryptAttr PayloadEncryptAttr
}

func (pp *PrivacyPayload) Data(version byte) []byte {
	//TODO: implement PrivacyPayload.Data()
	return []byte{0}
}

func (pp *PrivacyPayload) Serialize(w io.Writer, version byte) error {
	w.Write([]byte{byte(pp.PayloadType)})
	err := serialization.WriteVarBytes(w, pp.Payload)
	if err != nil {
		return err
	}
	w.Write([]byte{byte(pp.EncryptType)})
	err = pp.EncryptAttr.Serialize(w)

	return err
}

func (pp *PrivacyPayload) Serialization(sink *ZeroCopySink, version byte) error {
	sink.WriteByte(byte(pp.PayloadType))
	sink.WriteVarBytes(pp.Payload)
	sink.WriteByte(byte(pp.EncryptType))
	err := pp.EncryptAttr.Serialization(sink)
	if err != nil {
		return err
	}
	return nil
}

func (pp *PrivacyPayload) Deserialization(source *ZeroCopySource, version byte) error {
	//TODO
	//payloadType,_,irregular,eof := source.NextVarBytes()
	//if irregular {
	//	return common.ErrIrregularData
	//}
	payloadType,eof := source.NextByte()
	if eof {
		return io.ErrUnexpectedEOF
	}

	pp.PayloadType = EncryptedPayloadType(payloadType)
	p, _, irregular,eof := source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.EOF
	}
	pp.Payload = EncryptedPayload(p)
	t, eof := source.NextByte()
	if eof {
		return io.EOF
	}
    pp.EncryptType = PayloadEncryptType(t)
	switch pp.EncryptType {
	case ECDH_AES256:
		pp.EncryptAttr = new(EcdhAes256)
	default:
		return errors.New("unknown EncryptType")
	}
	err := pp.EncryptAttr.Deserialization(source)
	return err
}

func (pp *PrivacyPayload) Deserialize(r io.Reader, version byte) error {
	var PayloadType [1]byte
	_, err := io.ReadFull(r, PayloadType[:])
	if err != nil {
		return err
	}
	pp.PayloadType = EncryptedPayloadType(PayloadType[0])

	Payload, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	pp.Payload = Payload

	var encryptType [1]byte
	_, err = io.ReadFull(r, encryptType[:])
	if err != nil {
		return err
	}
	pp.EncryptType = PayloadEncryptType(encryptType[0])

	switch pp.EncryptType {
	case ECDH_AES256:
		pp.EncryptAttr = new(EcdhAes256)
	default:
		return errors.New("unknown EncryptType")
	}
	err = pp.EncryptAttr.Deserialize(r)

	return err
}

type EcdhAes256 struct {
	FromPubkey *crypto.PubKey
	ToPubkey   *crypto.PubKey
	Nonce      []byte
}

func (ea *EcdhAes256) Serialize(w io.Writer) error {
	err := ea.FromPubkey.Serialize(w)
	if err != nil {
		return err
	}
	err = ea.ToPubkey.Serialize(w)
	if err != nil {
		return err
	}
	err = serialization.WriteVarBytes(w, ea.Nonce)
	return err
}
func (ea *EcdhAes256) Serialization(sink *ZeroCopySink) error {
	err := ea.FromPubkey.Serialization(sink)
	if err != nil {
		return err
	}
	err = ea.ToPubkey.Serialization(sink)
	if err != nil {
		return err
	}
	sink.WriteVarBytes(ea.Nonce)
	return nil
}
func (ea *EcdhAes256) Deserialize(r io.Reader) error {
	ea.FromPubkey = new(crypto.PubKey)
	err := ea.FromPubkey.DeSerialize(r)
	if err != nil {
		return err
	}

	ea.ToPubkey = new(crypto.PubKey)
	err = ea.ToPubkey.DeSerialize(r)
	if err != nil {
		return err
	}

	nonce, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	ea.Nonce = nonce
	return nil
}

func (ea *EcdhAes256) Deserialization(source *ZeroCopySource) error {
	ea.FromPubkey = new(crypto.PubKey)
	err := ea.FromPubkey.DeSerialization(source)
	if err != nil {
		return err
	}
	ea.ToPubkey = new(crypto.PubKey)
	err = ea.ToPubkey.DeSerialization(source)
	if err != nil {
		return err
	}
	data, _, irregular, eof := source.NextVarBytes()
	if irregular {
		return ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	ea.Nonce = data
	return nil
}

func (ea *EcdhAes256) Encrypt(msg []byte, keys interface{}) ([]byte, error) {
	var key []byte
	switch keys.(type) {
	case []byte:
		key = keys.([]byte)
	default:
		return []byte{}, errors.New("The keys error")
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return []byte{}, err
	}
	x, _ := priv.Curve.ScalarMult(ea.ToPubkey.X, ea.ToPubkey.Y, key)
	aesKey := make([]byte, 32)
	copy(aesKey[32-len(x.Bytes()):], x.Bytes())

	iv := make([]byte, 16)
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		iv[i] = byte(r.Intn(256))
	}

	paddingData := crypto.PKCS5Padding(msg, 16)
	encryption, err := crypto.AesEncrypt(paddingData, aesKey, iv)
	if err != nil {
		return []byte{}, err
	}
	ea.Nonce = iv

	return encryption, nil
}

func (ea *EcdhAes256) Decrypt(msg []byte, keys interface{}) ([]byte, error) {
	var key []byte
	switch keys.(type) {
	case []byte:
		key = keys.([]byte)
	default:
		return []byte{}, errors.New("The keys error")
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return []byte{}, err
	}
	x, _ := priv.Curve.ScalarMult(ea.FromPubkey.X, ea.FromPubkey.Y, key)
	aesKey := make([]byte, 32)
	copy(aesKey[32-len(x.Bytes()):], x.Bytes())

	decryption, _ := crypto.AesDecrypt(msg, aesKey, ea.Nonce)
	if len(decryption) < int(decryption[len(decryption)-1]) {
		return []byte{}, errors.New("decryption error")
	}
	result := crypto.PKCS5UnPadding(decryption)

	return result, nil
}
