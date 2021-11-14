package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/did/exported"
)

var (
	ValidDid      = regexp.MustCompile(`^did:(stateset:)([a-zA-Z0-9]){21,22}([/][a-zA-Z0-9:]+|)$`)
	ValidPubKey   = regexp.MustCompile(`^[123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ]{43,44}$`)
	IsValidDid    = ValidDid.MatchString
	IsValidPubKey = ValidPubKey.MatchString
)

var _ exported.DidDoc = &BaseDidDoc{}

func (dd BaseDidDoc) MarshalYAML() (interface{}, error) {
	bz, err := codec.MarshalYAML(codec.NewProtoCodec(codectypes.NewInterfaceRegistry()), &dd)
	if err != nil {
		return nil, err
	}
	return string(bz), err
}

func (dd BaseDidDoc) String() string {
	out, _ := dd.MarshalYAML()
	return out.(string)
}

func NewBaseDidDoc(did exported.Did, pubKey string) BaseDidDoc {
	return BaseDidDoc{
		Did:         did,
		PubKey:      pubKey,
		Credentials: []*DidCredential{},
	}
}

func (dd BaseDidDoc) GetDid() exported.Did { return dd.Did }
func (dd BaseDidDoc) GetPubKey() string    { return dd.PubKey }
func (dd BaseDidDoc) GetCredentials() []DidCredential {
	lstToRet := make([]DidCredential, 0)
	credentials := dd.Credentials

	for _, cred := range credentials {
		lstToRet = append(lstToRet, *cred)
	}

	return lstToRet
}

func (dd *BaseDidDoc) SetDid(did exported.Did) error {
	if len(dd.Did) != 0 {
		return errors.New("cannot override BaseDidDoc did")
	}

	dd.Did = did

	return nil
}

func (dd *BaseDidDoc) SetPubKey(pubKey string) error {
	if len(dd.PubKey) != 0 {
		return errors.New("cannot override BaseDidDoc pubKey")
	}

	dd.PubKey = pubKey

	return nil
}

func (dd BaseDidDoc) Address() sdk.AccAddress {
	return exported.VerifyKeyToAddr(dd.GetPubKey())
}

func (dd *BaseDidDoc) AddCredential(cred *DidCredential) {
	if dd.Credentials == nil {
		dd.Credentials = make([]*DidCredential, 0)
	}

	dd.Credentials = append(dd.Credentials, cred)
}

func fromJsonString(jsonStatesetDid string) (exported.StatesetDid, error) {
	var did exported.StatesetDid
	err := json.Unmarshal([]byte(jsonStatesetDid), &did)
	if err != nil {
		err := fmt.Errorf("could not unmarshal did into struct due to error: %s", err.Error())
		return exported.StatesetDid{}, err
	}

	return did, nil
}

func UnmarshalStatesetDid(jsonStatesetDid string) (exported.StatesetDid, error) {
	return fromJsonString(jsonStatesetDid)
}