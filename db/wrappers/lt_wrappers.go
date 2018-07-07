// [fyi] wrappers for custom/unsupported by sqlite3 driver types

package wrappers

import(
	_ "database/sql/driver"

//	"log"
//	"math/big"
//	"errors"
//	"fmt"
	"time"

	"github.com/byteball/go-byteballcore/types"
)

const(
	UNIT_BYTE_SIZE		= types.UNIT_SIZE
	ADDRESS_BYTE_SIZE	= types.ADDRESS_SIZE
)

//

// [tbd]  Address Value

// [tbd]  Unit Value

// [tbd]  Asset Value

/**
type Asset struct{ V *types.AssetT }

func (asset Asset) Value() (driver.Value, error) {
	log.Printf("db.Asset.Value %#v", asset)
	if (asset.V.IsNull()) { return nil, nil }
	return string(*asset.V), nil
}
 **/

//

type Time struct{ V *time.Time }

func (time_ *Time) Scan(src interface{}) error {
//	panic("db.Time.Scan " + string(src.([]byte)))
	// [fyi] src = []byte("<time string>")
//	err := (time_.V).UnmarshalText(src.([]byte))
	t, err := time.Parse("2006-01-02 15:04:05", string(src.([]byte)))
	if err == nil { *time_.V = t }
	return err
/**
	if src == nil {
		*hash.V = types.CHashT_Undefined
		return nil
	}
	err := (hash.V).UnmarshalText(src.([]byte))
	return err
 **/
}

//

type CHash struct{ V *types.CHashT }

func (hash *CHash) Scan(src interface{}) error {
	// [fyi] src = []byte("<c-hash string>")
	if src == nil {
		*hash.V = types.CHashT_Null
		return nil
	}
	err := (hash.V).UnmarshalText(src.([]byte))
	return err
}

type Unit struct{ V *types.UnitT }

func (unit *Unit) Scan(src interface{}) error {
	// [fyi] src = []byte("<unit string>")
	// [tbd] if (src == nil) ...
//	if src == nil {
//		*unit.V = types.UnitT_Null
//		return nil
//	}
	err := (unit.V).UnmarshalText(src.([]byte))
	return err
}

type Address struct{ V *types.AddressT }

func (address *Address) Scan(src interface{}) error {
	// [fyi] src = []byte("<address string>")
	if src == nil {
		*address.V = types.AddressT_Null
		return nil
	}
	err := (address.V).UnmarshalText(src.([]byte))
	return err
}

type Ball struct{ V *types.BallT }

func (ball *Ball) Scan(src interface{}) error {
	// [fyi] src = []byte("<ball string>")
	if src == nil {
		*ball.V = types.BallT_Null
		return nil
	}
	err := (ball.V).UnmarshalText(src.([]byte))
	return err
}

type Asset struct{ V *types.AssetT }

func (asset *Asset) Scan(src interface{}) error {
	// [fyi] src = []byte("<asset string>")
	if src == nil {
		*asset.V = types.AssetT_Null
		return nil
	}
	err := (asset.V).UnmarshalText(src.([]byte))
	return err
}

type Level struct{ V *types.LevelT }

func (level *Level) Scan(src interface{}) error {
	// [fyi] src = int64(<int value>)
	if src == nil {
		*level.V = types.LevelT_Null
		return nil
	}
	*level.V = types.LevelT(src.(int64))
	return nil
}

type MCIndex struct{ V *types.MCIndexT }

func (mcindex *MCIndex) Scan(src interface{}) error {
	// [fyi] src = int64(<int value>)
	if src == nil {
//		*mcindex.V = types.MCIndexT(-1)
		*mcindex.V = types.MCIndexT_Null
	} else {
		*mcindex.V = types.MCIndexT(src.(int64))
	}
	return nil
}

/**
type UUID struct{ V *uuid.UUID }

func (uuid UUID) Value() (driver.Value, error) {
	bs, err := (uuid.V).MarshalText()
	return bs, err
}

func (uuid *UUID) Scan(src interface{}) error {
	// [fyi] src = []byte("<UUID string>")
	err := (uuid.V).UnmarshalText(src.([]byte))
	return err
}

type Numeric struct{ V *big.Int }

func (num Numeric) Value() (driver.Value, error) {
	bs, err := (num.V).MarshalText()
	return bs, err
}

func (num *Numeric) Scan(src interface{}) error {
	// [fyi] src = []byte("<numeric, decimal>")
	err := (num.V).UnmarshalText(src.([]byte))
	return err
}

type Address struct { V *common.Address }

func (addr Address) Value() (driver.Value, error) {
	bs, err := (addr.V).MarshalText()
	if bs != nil {
		// [fyi] drop leading "0x"
		bs = bs[2:]
	}
	return bs, err
}

func (addr *Address) Scan(src interface{}) error {
	// [fyi] src = []byte("<address, hex>")
	err := (addr.V).UnmarshalText(append([]byte("0x"),src.([]byte)...))
	return err
}
 **/
