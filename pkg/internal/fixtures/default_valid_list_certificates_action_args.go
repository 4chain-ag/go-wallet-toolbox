package fixtures

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/pkg/to"
)

func DefaultValidListCertificatesArgs() *wdk.ListCertificatesArgs {
	return &wdk.ListCertificatesArgs{
		Partial: &wdk.ListCertificatesArgsPartial{
			Type:               to.Ptr(primitives.Base64String(TypeField)),
			Certifier:          to.Ptr(primitives.PubKeyHex(Certifier)),
			SerialNumber:       to.Ptr(primitives.Base64String(SerialNumber)),
			Subject:            to.Ptr(primitives.PubKeyHex(SubjectPubKey)),
			RevocationOutpoint: to.Ptr(primitives.OutpointString(RevocationOutpoint)),
			Signature:          to.Ptr(primitives.HexString(Signature)),
		},
		Certifiers: []primitives.PubKeyHex{Certifier},
		Types:      []primitives.Base64String{TypeField},
		Limit:      primitives.PositiveIntegerDefault10Max10000(4),
		Offset:     primitives.PositiveIntegerOrZero(5),
	}
}
