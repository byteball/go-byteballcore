package types

import(
	"errors"
	"fmt"
)

const(
	HASHB64_SIZE	= 44
	CHASH_SIZE	= 32
)

type(
	HashBase64T	string
	HashHexT	string

	SHA1HexT	string

	CHash160T	string


	CHashT		= CHash160T
	CHashesT	= []CHashT

	refCHashesT	= *CHashesT
)

//

func (unit *CHashT) UnmarshalText(text []byte) error {
	if len(text) != CHASH_SIZE {
		return errors.New(fmt.Sprintf("invalid c-hash size %d", len(text)))
	}
	*unit = CHashT(text)
	return nil
}

//

const(
//	CHashT_Undefined CHashT = CHashT("")
	CHashT_Null CHashT = CHashT("")
)

//func (hash CHashT) Undefined() bool { return len(hash) == 0 }
func (hash CHashT) IsNull() bool { return len(hash) == 0 }

//

func (hash CHashT) OrNull() *CHashT {
	if hash.IsNull() { return nil }
	return &hash
}

//

func (hash CHashT) IndexOf(hashes refCHashesT) int {
	for k, hash_ := range *hashes {
		if hash_ == hash { return k }
	}
	return -1
}
