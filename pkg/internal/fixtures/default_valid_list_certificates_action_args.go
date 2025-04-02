package fixtures

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/utils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

func DefaultValidListCertificatesArgs() *wdk.ListCertificatesArgs {
	return &wdk.ListCertificatesArgs{
		Partial: &wdk.ListCertificatesArgsPartial{
			Type:               utils.Ptr(wdk.Base64String(TypeField)),
			Certifier:          utils.Ptr(wdk.PubKeyHex(Certifier)),
			SerialNumber:       utils.Ptr(wdk.Base64String(SerialNumber)),
			Subject:            utils.Ptr(wdk.PubKeyHex(SubjectPubKey)),
			RevocationOutpoint: utils.Ptr(wdk.OutpointString(RevocationOutpoint)),
			Signature:          utils.Ptr(wdk.HexString(Signature)),
		},
		Certifiers: []wdk.PubKeyHex{Certifier},
		Types:      []wdk.Base64String{TypeField},
		Limit:      wdk.PositiveIntegerDefault10Max10000(4),
		Offset:     wdk.PositiveIntegerOrZero(5),
	}
}
