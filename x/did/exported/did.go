package exported

import (
	"bytes"
	ed25519Local "crypto/ed25519"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/gogo/protobuf/proto"
	naclBox "golang.org/x/crypto/nacl/box"
)

var DidPrefix = "did:stateset:"
var AgreementDidPrefix = "did:stateset:agreement:"
var PurchaseOrderDidPrefix = "did:stateset:purchaseorder:"
var InvoiceDidPrefix = "did:stateset:invoice:"
var LoanDidPrefix = "did:stateset:loan:"

type Did = string

func UnprefixedDid(did Did) string {
	// Assumes that DID is valid (check IsValidDid regex)
	// Removes 8 characters (for did:stateset:)
	return did[8:]
}

func UnprefixedDidFromPubKey(pubKey string) string {
	// Assumes that PubKey is valid (check IsValidPubKey regex)
	// Since result is not prefixed (did:stateset:), string returned rather than DID
	pubKeyBz := base58.Decode(pubKey)
	return base58.Encode(pubKeyBz[:16])
}

type DidDoc interface {
	proto.Message

	SetDid(did Did) error
	GetDid() Did
	SetPubKey(pubkey string) error
	GetPubKey() string
	Address() sdk.AccAddress
}

type Secret struct {
	Seed                 string `json:"seed" yaml:"seed"`
	SignKey              string `json:"signKey" yaml:"signKey"`
	EncryptionPrivateKey string `json:"encryptionPrivateKey" yaml:"encryptionPrivateKey"`
}

func NewSecret(seed, signKey, encryptionPrivateKey string) Secret {
	return Secret{
		Seed:                 seed,
		SignKey:              signKey,
		EncryptionPrivateKey: encryptionPrivateKey,
	}
}

func (s Secret) Equals(other Secret) bool {
	return s.Seed == other.Seed &&
		s.SignKey == other.SignKey &&
		s.EncryptionPrivateKey == other.EncryptionPrivateKey
}

func (s Secret) String() string {
	output, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", string(output))
}

type StatesetDid struct {
	Did                 string `json:"did" yaml:"did"`
	VerifyKey           string `json:"verifyKey" yaml:"verifyKey"`
	EncryptionPublicKey string `json:"encryptionPublicKey" yaml:"encryptionPublicKey"`
	Secret              Secret `json:"secret" yaml:"secret"`
}

func NewStatesetDid(did, verifyKey, encryptionPublicKey string, secret Secret) StatesetDid {
	return StatesetDid{
		Did:                 did,
		VerifyKey:           verifyKey,
		EncryptionPublicKey: encryptionPublicKey,
		Secret:              secret,
	}
}

func (id StatesetDid) Equals(other StatesetDid) bool {
	return id.Did == other.Did &&
		id.VerifyKey == other.VerifyKey &&
		id.EncryptionPublicKey == other.EncryptionPublicKey &&
		id.Secret.Equals(other.Secret)
}

func (id StatesetDid) String() string {
	output, err := json.MarshalIndent(id, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", string(output))
}

func VerifyKeyToAddr(verifyKey string) sdk.AccAddress {
	var pubKey ed25519.PubKey
	pubKey.Key = base58.Decode(verifyKey)
	return sdk.AccAddress(pubKey.Address())
}

func (id StatesetDid) Address() sdk.AccAddress {
	return VerifyKeyToAddr(id.VerifyKey)
}

func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

func FromMnemonic(mnemonic string) (StatesetDid, error) {
	seed := sha256.New()
	seed.Write([]byte(mnemonic))

	var seed32 [32]byte
	copy(seed32[:], seed.Sum(nil)[:32])

	return FromSeed(seed32)
}

func Gen() (StatesetDid, error) {
	var seed [32]byte
	_, err := io.ReadFull(cryptoRand.Reader, seed[:])
	if err != nil {
		return StatesetDid{}, err
	}
	return FromSeed(seed)
}

func FromSeed(seed [32]byte) (StatesetDid, error) {
	publicKeyBytes, privateKeyBytes, err := ed25519Local.GenerateKey(bytes.NewReader(seed[0:32]))
	if err != nil {
		return StatesetDid{}, err
	}
	publicKey := []byte(publicKeyBytes)
	privateKey := []byte(privateKeyBytes)

	signKey := base58.Encode(privateKey[:32])
	keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKey[:]))
	if err != nil {
		return StatesetDid{}, err
	}

	return StatesetDid{
		Did:                 DidPrefix + base58.Encode(publicKey[:16]),
		VerifyKey:           base58.Encode(publicKey),
		EncryptionPublicKey: base58.Encode(keyPairPublicKey[:]),
		Secret: Secret{
			Seed:                 hex.EncodeToString(seed[0:32]),
			SignKey:              signKey,
			EncryptionPrivateKey: base58.Encode(keyPairPrivateKey[:]),
		},
	}, nil
}

func (id StatesetDid) SignMessage(msg []byte) ([]byte, error) {
	var privateKey ed25519.PrivKey
	privateKey.Key = append(base58.Decode(id.Secret.SignKey), base58.Decode(id.VerifyKey)...)

	return privateKey.Sign(msg)
}

func (id StatesetDid) VerifySignedMessage(msg []byte, sig []byte) bool {
	var publicKey ed25519.PubKey
	publicKey.Key = base58.Decode(id.VerifyKey)

	return publicKey.VerifySignature(msg, sig)
}