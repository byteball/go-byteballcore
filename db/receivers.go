package db

import(
	"database/sql"
	"fmt"
	"log"
//	"strings"
	"regexp"

 lt	"github.com/byteball/go-byteballcore/db/wrappers"

	"github.com/byteball/go-byteballcore/types"
)

type(
)


var(
	tagReceiverRe = regexp.MustCompile("Receiver$")
)


func _log(tag string, format string, args... interface{}) {

//	if tagReceiverRe.MatchString(tag) { return }
	return

	s := fmt.Sprintf(format, args...)
	log.Printf(tag + " " + s)
}

//

type UnitRow struct{ Unit types.UnitT }

type UnitsReceiver struct{
	Rows	[]UnitRow
}

func (rcvr *UnitsReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitReceiver", "unit %#v", row)

	return err
}

//

type WitnessListUnitRow struct{
	Witness_list_unit  types.UnitT
}

type WitnessListUnitsReceiver struct{
	Rows	[]WitnessListUnitRow
}

func (rcvr *WitnessListUnitsReceiver) Scan(sqlRows *sql.Rows) error {
	row := WitnessListUnitRow{}

	params := []interface{}{
		&lt.Unit{ &row.Witness_list_unit },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("WitnessListUnitsReceiver", "row %#v", row)

	return err
}

//

type WitnessListUnitAndUnitRow struct{
	Witness_list_unit  types.UnitT
	Unit		types.UnitT
}

type WitnessListUnitAndUnitsReceiver struct{
	Rows	[]WitnessListUnitAndUnitRow
}

func (rcvr *WitnessListUnitAndUnitsReceiver) Scan(sqlRows *sql.Rows) error {
	row := WitnessListUnitAndUnitRow{}

	params := []interface{}{
		&lt.Unit{ &row.Witness_list_unit },
		&lt.Unit{ &row.Unit },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("WitnessListUnitAndUnitsReceiver", "row %#v", row)

	return err
}

//

type ParentChildUnitRow struct{
	Parent_unit	types.UnitT
	Child_unit	types.UnitT
}

type ParentChildUnitsReceiver struct{
	Rows	[]ParentChildUnitRow
}

func (rcvr *ParentChildUnitsReceiver) Scan(sqlRows *sql.Rows) error {
	row := ParentChildUnitRow{}

	params := []interface{}{
		&lt.Unit{ &row.Parent_unit },
		&lt.Unit{ &row.Child_unit },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("ParentChildUnitsReceiver", "row %#v", row)

	return err
}

//

type BestParentUnitRow struct{
	Best_parent_unit  types.UnitT
	Witnessed_level	types.LevelT
}

type BestParentUnitsReceiver struct{
	Rows	[]BestParentUnitRow
}

func (rcvr *BestParentUnitsReceiver) Scan(sqlRows *sql.Rows) error {
	row := BestParentUnitRow{}

	params := []interface{}{
		&lt.Unit{ &row.Best_parent_unit },
		&lt.Level{ &row.Witnessed_level },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("BestParentUnitsReceiver", "row %#v", row)

	return err
}

//

type BestParentUnitCountRow struct{
	Best_parent_unit  types.UnitT
	Witnessed_level	types.LevelT
	Count		int
}

type BestParentUnitCountsReceiver struct{
	Rows	[]BestParentUnitCountRow
}

func (rcvr *BestParentUnitCountsReceiver) Scan(sqlRows *sql.Rows) error {
	row := BestParentUnitCountRow{}

	params := []interface{}{
		&lt.Unit{ &row.Best_parent_unit },
		&lt.Level{ &row.Witnessed_level },
		&row.Count,
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("BestParentUnitCountsReceiver", "row %#v", row)

	return err
}

//

type UnitAddressRow struct{
	Unit		types.UnitT
	Address		types.AddressT
}

type UnitAddressesReceiver struct{
	Rows	[]UnitAddressRow
}

func (rcvr *UnitAddressesReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitAddressRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
		&lt.Address{ &row.Address },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitAddressesReceiver", "row %#v", row)

	return err
}

//

type UnitMCISequenceAddressRow struct{
	Unit		types.UnitT
	Main_chain_index  types.MCIndexT
	Good		bool
	Address		types.AddressT
	X_skip		bool
}

type UnitMCISequenceAddressesReceiver struct{
	Rows	[]UnitMCISequenceAddressRow
}

func (rcvr *UnitMCISequenceAddressesReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitMCISequenceAddressRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
		&lt.MCIndex{ &row.Main_chain_index },
		&row.Good,
		&lt.Address{ &row.Address },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitMCISequenceAddressesReceiver", "row %#v", row)

	return err
}

//

type UnitMCIRow struct{
	Unit		types.UnitT
	Main_chain_index  types.MCIndexT
}

type UnitMCIsReceiver struct{
	Rows	[]UnitMCIRow
}

func (rcvr *UnitMCIsReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitMCIRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
		&lt.MCIndex{ &row.Main_chain_index },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitMCIsReceiver", "row %#v", row)

	return err
}

//

type UnitLIMCIRow struct{
	Unit		types.UnitT
	Latest_included_mc_index  types.MCIndexT
}

type UnitLIMCIsReceiver struct{
	Rows	[]UnitLIMCIRow
}

func (rcvr *UnitLIMCIsReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitLIMCIRow{}

	params := []interface{}{
		&lt.MCIndex{ &row.Latest_included_mc_index },
		&lt.Unit{ &row.Unit },
	//	&lt.MCIndex{ &row.Latest_included_mc_index },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitLIMCIReceiver", "row %#v", row)

	return err
}

//

type UnitContentHashRow struct{
	Unit		types.UnitT
	Content_hash	types.ContentHashT
}

type UnitContentHashsReceiver struct{
	Rows	[]UnitContentHashRow
}

func (rcvr *UnitContentHashsReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitContentHashRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
		&lt.CHash{ &row.Content_hash },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitMCIReceiver", "row %#v", row)

	return err
}

//

// [tbd] will not work for all UPR uses - specialize params for each sql stmt

type UnitIsFreesReceiver struct{
	Rows	[]UnitPropsRow
}

func (rcvr *UnitIsFreesReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitPropsRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
		&row.Is_free,
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitIsFreesReceiver", "row %#v", row)

	return err
}

//

type UnitPropsRow struct{
	Unit		types.UnitT
	Is_on_main_chain  int
	Main_chain_index  types.MCIndexT
	Level		types.LevelT
	Is_free		int
}

type UnitPropsReceiver struct{
	Rows	[]UnitPropsRow
}

func (rcvr *UnitPropsReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitPropsRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
		&row.Is_on_main_chain,
		&lt.MCIndex{ &row.Main_chain_index },
		&lt.Level{ &row.Level },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitPropsReceiver", "row %#v", row)

	return err
}

//

type UnitMCPropsRow struct{
	Unit		types.UnitT
	Level		types.LevelT
	Latest_included_mc_index  types.MCIndexT
	Main_chain_index  types.MCIndexT
	Is_on_main_chain  int
}

type UnitMCPropsReceiver struct{
	Rows	[]UnitMCPropsRow
}

func (rcvr *UnitMCPropsReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitMCPropsRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
		&lt.Level{ &row.Level },
		&lt.MCIndex{ &row.Latest_included_mc_index },
		&lt.MCIndex{ &row.Main_chain_index },
		&row.Is_on_main_chain,
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitMCPropsReceiver", "row %#v", row)

	return err
}

//

type zzUnitPropsRow struct{
	Unit		types.UnitT
	Is_on_main_chain  int
	Main_chain_index  types.MCIndexT
	Level		types.LevelT
	Is_free		int
}

type RUPUnitPropsReceiver struct{
	Rows	[]UnitContentRow
}

func (rcvr *RUPUnitPropsReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitContentRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
		&lt.Level{ &row.Level },
		&lt.MCIndex{ &row.Latest_included_mc_index },
		&lt.MCIndex{ &row.Main_chain_index },
		&row.Is_on_main_chain,
		&row.Is_free,
		&row.Is_stable,
		&lt.Level{ &row.Witnessed_level },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("RUPUnitPropsReceiver", "row %#v ", row)

	return err
}

//

// [tbd] depends on order of columns in units.*

/**
        unit CHAR(44) NOT NULL PRIMARY KEY,
        creation_date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        version VARCHAR(3) NOT NULL DEFAULT '1.0',
        alt VARCHAR(3) NOT NULL DEFAULT '1',
        witness_list_unit CHAR(44) NULL,
        last_ball_unit CHAR(44) NULL,
        content_hash CHAR(44) NULL,
        headers_commission INT NOT NULL,
        payload_commission INT NOT NULL,
        is_free TINYINT NOT NULL DEFAULT 1,
        is_on_main_chain TINYINT NOT NULL DEFAULT 0,
        main_chain_index INT NULL,
        latest_included_mc_index INT NULL,
        level INT NULL,
        witnessed_level INT NULL,
        is_stable TINYINT NOT NULL DEFAULT 0,
        sequence TEXT CHECK (sequence IN('good','temp-bad','final-bad')) NOT NULL DEFAULT 'good',
        best_parent_unit CHAR(44) NULL,
 **/

/**
type UnitContentRow struct{
	Unit		types.UnitT
	types.PropsT
}
 **/
type UnitContentRow	= types.UnitPropsT

type UnitContentsReceiver struct{
	Rows	[]UnitContentRow
}

func (rcvr *UnitContentsReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitContentRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
//		&row.Creation_date,
		&lt.Time{ &row.Creation_date },
		&row.Version,
		&row.Alt,
		&lt.Unit{ &row.Witness_list_unit },
		&lt.Unit{ &row.Last_ball_unit },
		&lt.CHash{ &row.Content_hash },
		&row.Headers_commission,
		&row.Payload_commission,
		&row.Is_free,
		&row.Is_on_main_chain,
		&lt.MCIndex{ &row.Main_chain_index },
		&lt.MCIndex{ &row.Latest_included_mc_index },
		&lt.Level{ &row.Level },
		&lt.Level{ &row.Witnessed_level },
		&row.Is_stable,
		&row.Sequence,
		&lt.Unit{ &row.Best_parent_unit },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitContentsReceiver", "row %#v", row)

	return err
}

//

/**
type UnitContentBallRow struct{
	UnitContentRow
	Ball		types.BallT
}
 **/
type UnitContentBallRow	= types.UnitPropsBallT

type UnitContentBallsReceiver struct{
	Rows	[]UnitContentBallRow
}

func (rcvr *UnitContentBallsReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitContentBallRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
//		&row.Creation_date,
		&lt.Time{ &row.Creation_date },
		&row.Version,
		&row.Alt,
		&lt.Unit{ &row.Witness_list_unit },
		&lt.Unit{ &row.Last_ball_unit },
		&lt.CHash{ &row.Content_hash },
		&row.Headers_commission,
		&row.Payload_commission,
		&row.Is_free,
		&row.Is_on_main_chain,
		&lt.MCIndex{ &row.Main_chain_index },
		&lt.MCIndex{ &row.Latest_included_mc_index },
		&lt.Level{ &row.Level },
		&lt.Level{ &row.Witnessed_level },
		&row.Is_stable,
		&row.Sequence,
		&lt.Unit{ &row.Best_parent_unit },
		&lt.Ball{ &row.Ball },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitContentsBallReceiver", "row %#v", row)

	return err
}

//

type UnitBallRow struct{
	Unit		types.UnitT
	Ball		types.BallT
}

type UnitBallsReceiver struct{
	Rows	[]UnitBallRow
}

func (rcvr *UnitBallsReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitBallRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
		&lt.Ball{ &row.Ball },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitBallsReceiver", "row %#v", row)

	return err
}

//

type BallRow struct{
	Ball		types.BallT
}

type BallsReceiver struct{
	Rows	[]BallRow
}

func (rcvr *BallsReceiver) Scan(sqlRows *sql.Rows) error {
	row := BallRow{}

	params := []interface{}{
		&lt.Ball{ &row.Ball },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitBallsReceiver", "row %#v", row)

	return err
}

//

type AddressRow struct{ Address types.AddressT }

type AddressesReceiver struct{
	Rows	[]AddressRow
}

func (rcvr *AddressesReceiver) Scan(sqlRows *sql.Rows) error {
	row := AddressRow{}

	params := []interface{}{
		&lt.Address{ &row.Address },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("AddressesReceiver", "addr %#v", row)

	return err
}

//

type AddressDelayRow struct{
	Address		types.AddressT
	Delay		types.DelayT
}

type AddressDelaysReceiver struct{
	Rows	[]AddressDelayRow
}

func (rcvr *AddressDelaysReceiver) Scan(sqlRows *sql.Rows) error {
	row := AddressDelayRow{}

	params := []interface{}{
		&lt.Address{ &row.Address },
		&lt.Delay{ &row.Delay },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("AddressDelaysReceiver", "addr %#v", row)

	return err
}

//

type AddressDenominationAssetRow struct{
	Address		types.AddressT
	Denomination	types.DenominationT
	Asset		types.AssetT
}

type AddressDenominationAssetsReceiver struct{
	Rows	[]AddressDenominationAssetRow
}

func (rcvr *AddressDenominationAssetsReceiver) Scan(sqlRows *sql.Rows) error {
	row := AddressDenominationAssetRow{}

	params := []interface{}{
		&lt.Address{ &row.Address },
		&row.Denomination,
//		&row.Asset,
		&lt.Asset{ &row.Asset },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("AddressDenominationAssetsReceiver", "addr %#v", row)

	return err
}

//

type LevelRow struct{ Level types.LevelT }

type LevelsReceiver struct{
	Rows	[]LevelRow
}

func (rcvr *LevelsReceiver) Scan(sqlRows *sql.Rows) error {
	row := LevelRow{}

	params := []interface{}{
		&lt.Level{ &row.Level },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("LevelsReceiver", "level %#v", row)

	return err
}

//

type MaxLevelRow struct{ Max_level types.LevelT }

type MaxLevelsReceiver struct{
	Rows	[]MaxLevelRow
}

func (rcvr *MaxLevelsReceiver) Scan(sqlRows *sql.Rows) error {
	row := MaxLevelRow{}

	params := []interface{}{
		&lt.Level{ &row.Max_level },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("MaxLevelsReceiver", "level %#v", row)

	return err
}

//

type MaxAltLevelRow struct{ Max_alt_level types.LevelT }

type MaxAltLevelsReceiver struct{
	Rows	[]MaxAltLevelRow
}

func (rcvr *MaxAltLevelsReceiver) Scan(sqlRows *sql.Rows) error {
	row := MaxAltLevelRow{}

	params := []interface{}{
		&lt.Level{ &row.Max_alt_level },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("MaxAltLevelsReceiver", "level %#v", row)

	return err
}

//

type WitnessedLevelRow struct{ Witnessed_level types.LevelT }

type WitnessedLevelsReceiver struct{
	Rows	[]WitnessedLevelRow
}

func (rcvr *WitnessedLevelsReceiver) Scan(sqlRows *sql.Rows) error {
	row := WitnessedLevelRow{}

	params := []interface{}{
		&lt.Level{ &row.Witnessed_level },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("WitnesssedLevelsReceiver", "level %#v", row)

	return err
}

//

type MinMCWLRow struct{ Min_mc_wl types.LevelT }

type MinMCWLsReceiver struct{
	Rows	[]MinMCWLRow
}

func (rcvr *MinMCWLsReceiver) Scan(sqlRows *sql.Rows) error {
	row := MinMCWLRow{}

	params := []interface{}{
		&lt.Level{ &row.Min_mc_wl },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("MinMCWLsReceiver", "level %#v", row)

	return err
}

//

type MinRetrievableMCIRow struct{
	Min_retrievable_mci  types.MCIndexT
}

type MinRetrievableMCIsReceiver struct{
	Rows	[]MinRetrievableMCIRow
}

func (rcvr *MinRetrievableMCIsReceiver) Scan(sqlRows *sql.Rows) error {
	row := MinRetrievableMCIRow{}

	params := []interface{}{
		&lt.MCIndex{ &row.Min_retrievable_mci },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("MinRetrievableMCIsReceiver", "row %#v", row)

	return err
}

//

type MaxSpendableMCIRow struct{
	Max_spendable_mci  types.MCIndexT
}

type MaxSpendableMCIsReceiver struct{
	Rows	[]MaxSpendableMCIRow
}

func (rcvr *MaxSpendableMCIsReceiver) Scan(sqlRows *sql.Rows) error {
	row := MaxSpendableMCIRow{}

	params := []interface{}{
		&lt.MCIndex{ &row.Max_spendable_mci },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("MaxSpendableMCIsReceiver", "row %#v", row)

	return err
}

//

type MinMCIndexRow struct{
	Min_main_chain_index  types.MCIndexT
}

type MinMCIndexsReceiver struct{
	Rows	[]MinMCIndexRow
}

func (rcvr *MinMCIndexsReceiver) Scan(sqlRows *sql.Rows) error {
	row := MinMCIndexRow{}

	params := []interface{}{
		&lt.MCIndex{ &row.Min_main_chain_index },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("MinMCIndexsReceiver", "row %#v", row)

	return err
}

//

type MCIOnMCRow struct{
	Main_chain_index  types.MCIndexT
	Is_on_main_chain  int
}

type MCIOnMCsReceiver struct{
	Rows	[]MCIOnMCRow
}

func (rcvr *MCIOnMCsReceiver) Scan(sqlRows *sql.Rows) error {
	row := MCIOnMCRow{}

	params := []interface{}{
		&lt.MCIndex{ &row.Main_chain_index },
		&row.Is_on_main_chain,
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("MCIOnMCReceiver", "row %#v", row)

	return err
}

//

type PaymentInfoRow struct{
	Child_unit	types.UnitT
	Headers_commission  types.AmountT
	Next_mc_unit	types.UnitT
	Payer_unit	types.UnitT
}

type PaymentInfosReceiver struct{
	Rows	[]PaymentInfoRow
}

func (rcvr *PaymentInfosReceiver) Scan(sqlRows *sql.Rows) error {
	row := PaymentInfoRow{}

	params := []interface{}{
		&lt.Unit{ &row.Child_unit },
		&row.Headers_commission,
		&lt.Unit{ &row.Next_mc_unit },
		&lt.Unit{ &row.Payer_unit },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("PaymentInfoReceiver", "row %#v", row)

	return err
}

//

type UnitAuthorShareRow struct{
	Unit		types.UnitT
	Address		types.AddressT
	Earned_headers_commission_share types.ShareT
}

type UnitAuthorSharesReceiver struct{
	Rows	[]UnitAuthorShareRow
}

func (rcvr *UnitAuthorSharesReceiver) Scan(sqlRows *sql.Rows) error {
	row := UnitAuthorShareRow{}

	params := []interface{}{
		&lt.Unit{ &row.Unit },
		&lt.Address{ &row.Address },
		&row.Earned_headers_commission_share,
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("UnitAuthorShareReceiver", "row %#v", row)

	return err
}

//

type PropsReceiver struct{
	Rows	[]types.PropsT
}

func (rcvr *PropsReceiver) Scan(sqlRows *sql.Rows) error {
	row := types.PropsT{}

	params := []interface{}{
		&lt.Level{ &row.Level },
		&lt.Level{ &row.Witnessed_level },
		&lt.Unit{ &row.Best_parent_unit },
		&lt.Unit{ &row.Witness_list_unit },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("PropsReceiver", "props %#v", row)

	return err
}

//

type CountRow struct{ Count int }

type CountsReceiver struct{
	Rows	[]CountRow
}

func (rcvr *CountsReceiver) Scan(sqlRows *sql.Rows) error {
	row := CountRow{}

	params := []interface{}{
		&row.Count,
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("CountsReceiver", "count %#v", row)

	return err
}

//

type CountOnStableMCRow struct{
	Count		int
	Count_on_stable_mc  int
}

type CountOnStableMCsReceiver struct{
	Rows	[]CountOnStableMCRow
}

func (rcvr *CountOnStableMCsReceiver) Scan(sqlRows *sql.Rows) error {
	row := CountOnStableMCRow{}

	params := []interface{}{
		&row.Count,
		&row.Count_on_stable_mc,
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.Rows = append(rcvr.Rows, row)

	_log("CountOnStableMCsReceiver", "row %#v", row)

	return err
}

//

type DBRowT struct{
	Kv map[string] interface{}
}

type DBRowsReceiver struct{
	rows DBRowsT
}

func (rcvr *DBRowsReceiver) Scan(sqlRows *sql.Rows) error {
	row := DBRowT{}

	colNames, ecols := sqlRows.Columns()
	if ecols != nil {
		return ecols
	}

	colNamesCt := len(colNames)
	params := make([]interface{}, colNamesCt, colNamesCt)
	for k, _ := range params {
		var v interface{}
		params[k] = &v
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		//log.Printf("DBRowsReceiverT.Scan %#v %#v", err, params)
		return err
	}

	row.Kv = make(map[string] interface{}, colNamesCt)
	for k, colName := range colNames {
		row.Kv[colName] = params[k]
	}

	rcvr.rows = append(rcvr.rows, row)

	_log("DBRowsReceiver", "rows %#v", rcvr.rows)

	return err
}

//

/**
type AccountsReceiver struct{
	rows []Account
}

func (rcvr *AccountsReceiver) Scan(sqlRows *sql.Rows) error {
	acct := Account{}

	params := []interface{}{
		&pg.UUID{ &acct.ID },
		&pg.Address{ &acct.Address },
		&pg.Numeric{ &acct.Balance },
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		return err
	}

	rcvr.rows = append(rcvr.rows, acct)

	//log.Printf("AccountReceiver", "acct %#v", acct)

	return err
}

func (database *Database) GetAccounts(addrs []common.Address) ([]Account, error) {

	// [tbd] remove duplicate addrs

	acs := make([]Account, 0, len(addrs))
	var err error
	for _, addr := range addrs {
		acs_, err_ := database.GetAccount(addr)
		if err_ != nil {
			err = err_
			log.Printf("DB.GetAccounts: %s", err_.Error)
			continue
		}
		acs = append(acs, acs_...)
	}

	//log.Printf("GetAccounts acs %#v", acs)

	return acs, err
}

func (database *Database) GetAccount(addr common.Address) ([]Account, error) {

	params := []interface{}{
		pg.Address{ &addr },
	}

	sql := `
SELECT
	acs.id,
	address,
	balance
FROM "Accounts" acs
JOIN "AccountBalances" abs
	ON abs.id = acs.id
--WHERE acs.address IN ($1::varchar[])
WHERE acs.address IN ($1::varchar)
	`

	rcvr := AccountsReceiver{
		rows: []Account{},
	}

	err := database.Select(sql, params, &rcvr)

	//log.Printf("DB.GetAccount rcvr/rows %#v %#v", rcvr, rcvr.rows)

	return rcvr.rows, err
}

//

type Statistic struct{
	kv map[string] int64
}

type StatisticsReceiver struct{
	rows []Statistic
}

func (rcvr *StatisticsReceiver) Scan(sqlRows *sql.Rows) error {
	stat := Statistic{}

	colNames, ecols := sqlRows.Columns()
	if ecols != nil {
		return ecols
	}

	colNamesCt := len(colNames)
	params := make([]interface{}, colNamesCt, colNamesCt)
	for k, _ := range params {
		var v int64
		params[k] = &v
	}

	err := sqlRows.Scan(params...)
	if err != nil {
		log.Printf("StatisticsReceiver.Scan %#v %#v", err, params)
		return err
	}

	stat.kv = make(map[string] int64, colNamesCt)
	for k, colName := range colNames {
		stat.kv[colName] = *(params[k].(*int64))
	}

	rcvr.rows = append(rcvr.rows, stat)

	//log.Printf("StatisticsReceiver stat/rows %#v %#v", stat, rcvr.rows)

	return err
}

func (database *Database) UpdateAccountBalance(acct Account) ([]Statistic, error) {

	params := []interface{}{
		pg.UUID{ &acct.ID },
		pg.Numeric{ &acct.Balance },
	}

	sql := `
WITH updateBalances AS(
	UPDATE "AccountBalances" abs
	SET balance = $2
	WHERE abs.id = $1::uuid
	RETURNING *
)
	SELECT
		(SELECT count(*) FROM updateBalances) "updateBalances"
	`

	rcvr := StatisticsReceiver{
		rows: []Statistic{},
	}

	err := database.Select(sql, params, &rcvr)

	return rcvr.rows, err
}

//

func (database *Database) AddCreditTransaction(txn CreditTransaction) ([]Statistic, error) {

	// [fyi] insert credit txn

	// [tbd] abstract txid -> (txidHead, txidTail)

	txID := txn.ID
	txhead, _ := uuid.FromBytes(txID.Bytes()[0:16])
	txtail, _ := uuid.FromBytes(txID.Bytes()[16:32])

	params := []interface{}{
		pg.UUID{ &txhead },
		pg.UUID{ &txtail },
		pg.UUID{ &txn.Destination.ID },
		pg.Address{ &txn.Origination },
		pg.Numeric{ &txn.Amount },
	}

	sql := `
WITH insertTransactions AS (
	INSERT INTO "CreditTransactions"
		("transactionIDHead", "transactionIDTail",
		 "account", "originationAddress", "amount"
		)
	VALUES ($1::uuid, $2::uuid, $3::uuid, $4, $5)
	RETURNING *
)
	SELECT
		(SELECT count(*) FROM insertTransactions) "insertTransactions"
	`

	rcvr := StatisticsReceiver{
		rows: []Statistic{},
	}

	err := database.Select(sql, params, &rcvr)

	return rcvr.rows, err
}

//

func (database *Database) AddDebitTransaction(txn DebitTransaction) ([]Statistic, error) {

	// [fyi] insert debit txn

	txID := txn.ID
	txhead, _ := uuid.FromBytes(txID.Bytes()[0:16])
	txtail, _ := uuid.FromBytes(txID.Bytes()[16:32])

	params := []interface{}{
		pg.UUID{ &txhead },
		pg.UUID{ &txtail },
		pg.UUID{ &txn.Origination.ID },
		pg.Address{ &txn.Destination },
		pg.Numeric{ &txn.Amount },
	}

	sql := `
WITH insertTransactions AS (
	INSERT INTO "DebitTransactions"
		("transactionIDHead", "transactionIDTail",
		 "account", "destinationAddress", "amount"
		)
	VALUES ($1::uuid, $2::uuid, $3::uuid, $4, $5)
	RETURNING *
)
	SELECT
		(SELECT count(*) FROM insertTransactions) "insertTransactions"
	`

	rcvr := StatisticsReceiver{
		rows: []Statistic{},
	}

	err := database.Select(sql, params, &rcvr)

	return rcvr.rows, err
}
 **/
