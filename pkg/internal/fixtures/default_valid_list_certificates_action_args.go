package fixtures

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/go-softwarelab/common/pkg/to"
)

func DefaultValidListCertificatesArgs() *wdk.ListCertificatesArgs {
	return &wdk.ListCertificatesArgs{
		Partial: &wdk.ListCertificatesArgsPartial{
			Type:               to.Ptr(wdk.Base64String(TypeField)),
			Certifier:          to.Ptr(wdk.PubKeyHex(Certifier)),
			SerialNumber:       to.Ptr(wdk.Base64String(SerialNumber)),
			Subject:            to.Ptr(wdk.PubKeyHex(SubjectPubKey)),
			RevocationOutpoint: to.Ptr(wdk.OutpointString(RevocationOutpoint)),
			Signature:          to.Ptr(wdk.HexString(Signature)),
		},
		Certifiers: []wdk.PubKeyHex{Certifier},
		Types:      []wdk.Base64String{TypeField},
		Limit:      wdk.PositiveIntegerDefault10Max10000(4),
		Offset:     wdk.PositiveIntegerOrZero(5),
	}
}
