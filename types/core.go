package types

import(
	"errors"
	"fmt"

	"encoding/json"
	"math/big"
	"time"
)

const(
	ADDRESS_SIZE	= CHASH_SIZE

	UNIT_SIZE	= HASHB64_SIZE
	BALL_SIZE	= HASHB64_SIZE
	ASSET_SIZE	= HASHB64_SIZE
)

type(
//	AddressT	= CHashT
	AddressT	CHashT
	AddressesT	= []AddressT

	UnitT		HashBase64T
//	UnitsT		= []UnitT
	UnitsT		[]UnitT

	BallT		HashBase64T
	BallsT		= []BallT

	AssetT		HashBase64T

	LevelT		int

	MCIndexT	int
	MCIndexesT	= []MCIndexT

	ContentHashT	= CHashT

//	TimeT		time.Time
	TimeT		= time.Time

	BigIntT		= int64

	// authors

	AuthorT struct{
		Address		AddressT		`json:"address"`
		Definition	DefinitionT		`json:"definition,omitempty"`
		Authentifiers	AuthentifiersT		`json:"authentifiers,omitempty"`
	}
	AuthorsT	= []AuthorT

	// [fyi] high order typed json content
	DefinitionT	interface{}

	AuthentifiersT	= map[PathT] AuthentifierT

	PathT		string
	AuthentifierT	string

	// messages

	MessageT	struct{
		App		string			`json:"app"`
		Payload_hash	string			`json:"payload_hash"`
		Payload_location string			`json:"payload_location"`
		Payload		PayloadT		`json:"payload"`

		// [fyi] "text" and other
		PayloadText	string			`json:"-"`

		Payload_uri	*string			`json:"payload_uri,omitempty""`
		Payload_uri_hash *string		`json:"payload_uri_hash,omitempty"`

		Spend_proofs	SpendProofsT		`json:"spend_proofs,omitempty"`
	}
	MessagesT	= []MessageT

	// payload

	PayloadT struct {
		Denomination	DenominationT		`json:"denomination,omitempty"`
		Inputs		InputsT			`json:"inputs,omitempty"`
		Outputs		OutputsT		`json:"outputs,omitempty"`
		Asset		AssetT			`json:"asset,omitempty"`
		Timestamp	int64			`json:"timestamp,omitempty"`

		// [fyi] "address_definition_change"
		Definition_chash  CHashT		`json:"defnition_chash,omitempty"`
		Address		AddressT		`json:"address,omitempty"`

		// [fyi] "poll"
		Question	string			`json:"question,omitempty"`
		Choices		[]string		`json:"choices,omitempty"`

		// [fyi] "vote"
		Unit		UnitT			`json:"unit,omitempty"`
		Choice		string			`json:"choice,omitempty"`

		// [fyi] "attestation"
		Profile		map[string] string	`json:"profile,omitempty"`

		// [fyi] "asset"

		Cap		int64			`json:"cap,omitempty"`

		Is_private		bool		`json:"is_private,omitempty"`
		Is_transferrable	bool		`json:"is_transferrable,omitempty"`
		Auto_destroy		bool		`json:"auto_destroy,omitempty"`
		Fixed_denominations	bool		`json:"fixed_denominations,omitempty"`
		Issued_by_definer_only	bool		`json:"issued_by_definer_only,omitempty"`
		Cosigned_by_definer	bool		`json:"cosigned_by_definer,omitempty"`
		Spender_attested	bool		`json:"spender_attested,omitempty"`

		Issue_condition		interface{}	`json:"issue_condition,omitempty"`
		Transfer_condition	interface{}	`json:"transfer_condition,omitempty"`

		Attestors	AddressesT		`json:"attestors,omitempty"`
		Denominations	[]struct{
			Denomination	DenominationT	`json:"denomination,omitempty"`
			Count_coins	BigIntT		`json:"count_coins"`
		}					`json:"denominations,omitempty"`

		// [fyi] "data_feed"
		DataFeed	DataFeedT		`json:"datafeed,omitempty"`

	}

	DenominationT		int64

	SpendProofT struct{
		Address		AddressT		`json:"address"`
		Spend_proof	string			`json:"spend_proof"`
	}
	SpendProofsT	= []SpendProofT

	InputT struct{
		Type		string			`json:"type,omitempty"`
		Unit		UnitT			`json:"unit,omitempty"`
		Message_index	MessageIndexT		`json:"message_index,omitempty"`
		Output_index	OutputIndexT		`json:"output_index,omitempty"`
		From_main_chain_index  MCIndexT		`json:"from_main_chain_index,omitempty"`
		To_main_chain_index  MCIndexT		`json:"to_main_chain_index,omitempty"`
		Address		AddressT		`json:"address,omitempty"`
		Amount		AmountT			`json:"amount,omitempty"`
		Serial_number	SerialNumberT		`json:"serial_number,omitempty"`
	}
	InputsT		= []InputT

	MessageIndexT	int
	InputIndexT	int
	OutputIndexT	int

	AmountT		BigIntT
	SerialNumberT	BigIntT
	// [fyi] in percents
	ShareT		int64

	OutputT struct{
		Address		AddressT		`json:"address"`
		Amount		AmountT			`json:"amount"`
	}
	OutputsT	= []OutputT

	DataFeedT	= map[string] interface{}

	// unit properties

	UnitPropsT struct{
		Unit		UnitT
		PropsT
	}
	UnitPropssT	= []UnitPropsT

	UnitPropsBallT struct{
		UnitPropsT
		Ball		BallT
	}
	UnitPropsBallsT	= []UnitPropsBallT

	// [fyi] per DB table units
	PropsT  struct{
		Creation_date		TimeT
		Version			string
		Alt			string
		Witness_list_unit	UnitT
		Last_ball_unit		UnitT
		Content_hash		ContentHashT
		Headers_commission	int
		Payload_commission	int
		Is_free			int
		Is_on_main_chain	int
		Main_chain_index	MCIndexT
		Latest_included_mc_index  MCIndexT
		Level			LevelT
		Witnessed_level		LevelT
		Is_stable		int
		Sequence		string
		Best_parent_unit	UnitT
	}
	PropssT		= []PropsT

	JointT struct{
		Unit		UnitT
	}
)


type(
//	refUnitsT	= *UnitsT
//	refAddressesT	= *AddressesT
)

//

func (unit UnitT) MarshalText() ([]byte, error) {
	bs := ([]byte)(unit)
	return bs, nil
}

func (unit *UnitT) UnmarshalText(text []byte) error {
	if len(text) != UNIT_SIZE {
		return errors.New(fmt.Sprintf("invalid unit size %d", len(text)))
	}
	*unit = UnitT(text)
	return nil
}

func (unit *AssetT) UnmarshalText(text []byte) error {
	if len(text) != ASSET_SIZE {
		return errors.New(fmt.Sprintf("invalid asset size %d", len(text)))
	}
	*unit = AssetT(text)
	return nil
}

func (address *AddressT) UnmarshalText(text []byte) error {
	if len(text) != ADDRESS_SIZE {
		return errors.New(fmt.Sprintf("invalid address size %d", len(text)))
	}
	*address = AddressT(text)
	return nil
}

func (ball *BallT) UnmarshalText(text []byte) error {
	if len(text) != BALL_SIZE {
		return errors.New(fmt.Sprintf("invalid ball size %d", len(text)))
	}
	*ball = BallT(text)
	return nil
}

func (level *LevelT) UnmarshalText(text []byte) error {
	int := big.NewInt(0)
	err := int.UnmarshalText(text)
	if err != nil { return err }
	if !int.IsInt64() {
		return errors.New(fmt.Sprintf("Level.Scan: %s out of range ", string(text)))
	}
	*level = LevelT(int.Int64())
	return nil
}

//

func (level *LevelT) UnmarshalJSON(json []byte) error {
	return level.UnmarshalText(json)
}

func (mcindex *MCIndexT) MarshalJSON() ([]byte, error) {
	if mcindex.IsNull() { return json.Marshal(nil) }
	return json.Marshal(int(*mcindex))
}

//

const(
	UnitT_Null UnitT = UnitT("")

	BallT_Null BallT = BallT("")

	AssetT_Null AssetT = AssetT("")

	AddressT_Null AddressT = AddressT("")

	LevelT_Null LevelT = LevelT(-1)

	MCIndexT_Null MCIndexT = MCIndexT(-1)

	AmountT_Null AmountT = AmountT(-666)
)

func (unit UnitT) IsNull() bool { return len(unit) == 0 }

func (ball BallT) IsNull() bool { return len(ball) == 0 }

func (asset AssetT) IsNull() bool { return len(asset) == 0 }

func (address AddressT) IsNull() bool { return len(address) == 0 }

func (level LevelT) IsNull() bool { return int(level) < 0 }

func (mcindex MCIndexT) IsNull() bool { return int(mcindex) < 0 }

// [tbd] verify this
func (amount AmountT) IsNull() bool { return int(amount) < 0 }

func (denomination DenominationT) IsNull() bool { return int(denomination) <= 0 }

func (serno SerialNumberT) IsNull() bool { return int(serno) <= 0 }

//

func (unit UnitT) OrNull() *UnitT {
	if unit.IsNull() { return nil }
	return &unit
}

func (address AddressT) OrNull() *AddressT {
	if address.IsNull() { return nil }
	return &address
}

func (mcindex MCIndexT) OrNull() *MCIndexT {
	if mcindex.IsNull() { return nil }
	return &mcindex
}

func (asset AssetT) OrNull() *AssetT {
	if asset.IsNull() { return nil }
	return &asset
}

func (serno SerialNumberT) OrNull() *SerialNumberT {
	if serno.IsNull() { return nil }
	return &serno
}

//

func (unit UnitT) IndexOf(units UnitsT) int {
	for k, unit_ := range units {
		if unit_ == unit { return k }
	}
	return -1
}

func (address AddressT) IndexOf(addresses AddressesT) int {
	for k, address_ := range addresses {
		if address_ == address { return k }
	}
	return -1
}

//

func (units UnitsT) Compare (units_ UnitsT) int {
	d := len(units) - len(units_)
	if d != 0 { return d }
	for _k, unit := range units {
		if unit < units_[_k] { return -1 }
		if unit > units_[_k] { return +1 }
	}
	return 0
}

//

func (units *UnitsT) Join(sep string) string {
	panic("[tbd] Units.Join")
	return ""
}

//

func (units UnitsT) Uniq () UnitsT {
	uniq := make(UnitsT, 0, len(units))
	umap := make(map[UnitT] bool)
	for _, unit := range units {
		if _, _exists := umap[unit]; ! _exists {
			uniq = append(uniq, unit)
			umap[unit] = true
		}
	}
	return uniq
}
