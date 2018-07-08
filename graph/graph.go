
// converted golang begin

package graph

import(

 _core	"nodejs/core"
//	"nodejs/console"

 .	"github.com/byteball/go-byteballcore/types"

)

import(
// _		"lodash"
//		"async"
		"github.com/byteball/go-byteballcore/storage"
		"github.com/byteball/go-byteballcore/db"
		"github.com/byteball/go-byteballcore/profiler"
)

type(
	DBConnT		= db.DBConnT
	refDBConnT	= *DBConnT

	DBParamsT	= db.DBParamsT

)

/**
func DetermineIfIncludedOrEqual_sync(conn refDBConnT, unit UnitT, units UnitsT) bool {
	panic("[tbd] DetermineIfIncludedOrEqual_sync")
}
 **/

func CompareUnitsByProps_sync(conn refDBConnT, objUnitProps1 db.UnitContentRow, objUnitProps2 db.UnitContentRow) *int {
	panic("[tbd] CompareUnitsByProps_sync")
}

/**
func ReadDescendantUnitsByAuthorsBeforeMcIndex_sync(conn refDBConnT, objUnitProps db.UnitContentRow, arrWitnesses refAddressesT, mci MCIndexT) refUnitsT {
	panic("[tbd] ReadDescendantUnitsByAuthorsBeforeMcIndex_sync")
}
 **/


const _compareUnits = `
func compareUnits(conn DBConnT, unit1 unit1T, unit2 unit2T, handleResult handleResultT)  {
	if unit1 == unit2 {
		return handleResult(0)
	}
	rows := /* await */
	conn.query_sync("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain, is_free FROM units WHERE unit IN(?)", DBParamsT{ {*ArrayExpression*} })
	// << flattened continuation for conn.query:14:1
	if len(rows) != 2 {
		_core.Throw("not 2 rows")
	}
	objUnitProps1 := (rows[0].unit == unit1 ? rows[0]: rows[1])
	objUnitProps2 := (rows[0].unit == unit2 ? rows[0]: rows[1])
	compareUnitsByProps(conn, objUnitProps1, objUnitProps2, handleResult)
	// >> flattened continuation for conn.query:14:1
}

func compareUnitsByProps(conn DBConnT, objUnitProps1 objUnitProps1T, objUnitProps2 objUnitProps2T, handleResult handleResultT)  {
	var(
		goUp_1 func (arrStartUnits UnitsT) 
		goDown_1 func (arrStartUnits UnitsT) 
	)
	
	if objUnitProps1.unit == objUnitProps2.unit {
		return handleResult(0)
	}
	if objUnitProps1.level == objUnitProps2.level {
		return handleResult(nil)
	}
	if objUnitProps1.is_free == 1 && objUnitProps2.is_free == 1 {
		// free units
		return handleResult(nil)
	}
	
	// genesis
	if objUnitProps1.latest_included_mc_index == nil {
		return handleResult(- 1)
	}
	if objUnitProps2.latest_included_mc_index == nil {
		return handleResult(+ 1)
	}
	
	if objUnitProps1.latest_included_mc_index >= objUnitProps2.main_chain_index && objUnitProps2.main_chain_index != nil {
		return handleResult(+ 1)
	}
	if objUnitProps2.latest_included_mc_index >= objUnitProps1.main_chain_index && objUnitProps1.main_chain_index != nil {
		return handleResult(- 1)
	}
	
	if objUnitProps1.level <= objUnitProps2.level && objUnitProps1.latest_included_mc_index <= objUnitProps2.latest_included_mc_index && objUnitProps1.main_chain_index <= objUnitProps2.main_chain_index && objUnitProps1.main_chain_index != nil && objUnitProps2.main_chain_index != nil || objUnitProps1.main_chain_index == nil || objUnitProps2.main_chain_index == nil || objUnitProps1.level >= objUnitProps2.level && objUnitProps1.latest_included_mc_index >= objUnitProps2.latest_included_mc_index && objUnitProps1.main_chain_index >= objUnitProps2.main_chain_index && objUnitProps1.main_chain_index != nil && objUnitProps2.main_chain_index != nil || objUnitProps1.main_chain_index == nil || objUnitProps2.main_chain_index == nil {
	} else {
		return handleResult(nil)
	}
	
	objEarlierUnit := (objUnitProps1.level < objUnitProps2.level ? objUnitProps1: objUnitProps2)
	objLaterUnit := (objUnitProps1.level < objUnitProps2.level ? objUnitProps2: objUnitProps1)
	resultIfFound := (objUnitProps1.level < objUnitProps2.level ? - 1: 1)
	
	// can be negative if main_chain_index === null but that doesn't matter
	earlier_unit_delta := objEarlierUnit.main_chain_index - objEarlierUnit.latest_included_mc_index
	later_unit_delta := objLaterUnit.main_chain_index - objLaterUnit.latest_included_mc_index
	
	goUp_1 = func (arrStartUnits UnitsT)  {
		rows := /* await */
		conn.query_sync("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain \n" +
			"			FROM parenthoods JOIN units ON parent_unit=unit \n" +
			"			WHERE child_unit IN(?)", DBParamsT{ arrStartUnits })
		// << flattened continuation for conn.query:73:2
		arrNewStartUnits := UnitsT{}
		for i := 0; i < len(rows); i++ {
			objUnitProps := rows[i]
			if objUnitProps.unit == objEarlierUnit.unit {
				return handleResult(resultIfFound)
			}
			if objUnitProps.is_on_main_chain == 0 && objUnitProps.level > objEarlierUnit.level {
				arrNewStartUnits = append(arrNewStartUnits, objUnitProps.unit)
			}
		}
		(len(arrNewStartUnits) > 0 ? goUp_1(arrNewStartUnits): handleResult(nil))
		// >> flattened continuation for conn.query:73:2
	}
	
	goDown_1 = func (arrStartUnits UnitsT)  {
		rows := /* await */
		conn.query_sync("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain \n" +
			"			FROM parenthoods JOIN units ON child_unit=unit \n" +
			"			WHERE parent_unit IN(?)", DBParamsT{ arrStartUnits })
		// << flattened continuation for conn.query:93:2
		arrNewStartUnits := UnitsT{}
		for i := 0; i < len(rows); i++ {
			objUnitProps := rows[i]
			if objUnitProps.unit == objLaterUnit.unit {
				return handleResult(resultIfFound)
			}
			if objUnitProps.is_on_main_chain == 0 && objUnitProps.level < objLaterUnit.level {
				arrNewStartUnits = append(arrNewStartUnits, objUnitProps.unit)
			}
		}
		(len(arrNewStartUnits) > 0 ? goDown_1(arrNewStartUnits): handleResult(nil))
		// >> flattened continuation for conn.query:93:2
	}
	
	(later_unit_delta > earlier_unit_delta ? goUp_1(UnitsT{ objLaterUnit.unit }): goDown_1(UnitsT{ objEarlierUnit.unit }))
}
`

// determines if earlier_unit is included by at least one of arrLaterUnits 
//func determineIfIncluded_sync(conn DBConnT, earlier_unit UnitT, arrLaterUnits UnitsT) bool {
func determineIfIncluded_sync(conn refDBConnT, earlier_unit UnitT, arrLaterUnits UnitsT) bool {
	var(
		goUp_2 func (arrStartUnits UnitsT) bool
	)
	
//	if ! earlier_unit {
	if earlier_unit.IsNull() {
		_core.Throw("no earlier_unit")
	}
//	if storage.isGenesisUnit(earlier_unit) {
	if storage.IsGenesisUnit(earlier_unit) {
		// :: flattened return for return handleResult(true);
		return true
	}
//	( objEarlierUnitProps, arrLaterUnitProps ) := /* await */
	objEarlierUnitProps, arrLaterUnitProps := /* await */
//	storage.readPropsOfUnits_sync(conn, earlier_unit, arrLaterUnits)
	storage.ReadPropsOfUnits_sync(conn, earlier_unit, arrLaterUnits)
	// << flattened continuation for storage.readPropsOfUnits:122:1
//	if objEarlierUnitProps.is_free == 1 {
	if objEarlierUnitProps.Is_free == 1 {
		// :: flattened return for return handleResult(false);
		return false
	}

/**	
	max_later_limci := Math.max.apply(nil, // .. not flattening for Array.map
	arrLaterUnitProps.map(func (objLaterUnitProps PropsT) {*returns*} {
		return objLaterUnitProps.latest_included_mc_index
	}))
 **/
	max_later_limci := MCIndexT(-1)
	for _, objLaterUnitProps := range arrLaterUnitProps {
		limci := objLaterUnitProps.Latest_included_mc_index
		if max_later_limci < limci {
			max_later_limci = limci
		}
	}
	//console.log("max limci "+max_later_limci+", earlier mci "+objEarlierUnitProps.main_chain_index);
//	if objEarlierUnitProps.main_chain_index != nil && max_later_limci >= objEarlierUnitProps.main_chain_index {
	if ! objEarlierUnitProps.Main_chain_index.IsNull() && max_later_limci >= objEarlierUnitProps.Main_chain_index {
		// :: flattened return for return handleResult(true);
		return true
	}

/**	
	max_later_level := Math.max.apply(nil, // .. not flattening for Array.map
	arrLaterUnitProps.map(func (objLaterUnitProps PropsT) {*returns*} {
		return objLaterUnitProps.level
	}))
 **/
	max_later_level := LevelT(-1)
	for _, objLaterUnitProps := range arrLaterUnitProps {
		level := objLaterUnitProps.Level
		if max_later_level < level {
			max_later_level = level
		}
	}	
//	if max_later_level < objEarlierUnitProps.level {
	if max_later_level < objEarlierUnitProps.Level {
		// :: flattened return for return handleResult(false);
		return false
	}
	
	goUp_2 = func (arrStartUnits UnitsT) bool {
/**
		rows := /* await * /
		conn.query_sync("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain \n" +
			"				FROM parenthoods JOIN units ON parent_unit=unit \n" +
			"				WHERE child_unit IN(?)", DBParamsT{ arrStartUnits })
 **/
		rcvr := db.UnitPropsReceiver{}
		queryParams := DBParamsT{}
		susSql := queryParams.AddUnits(arrStartUnits)
		conn.MustQuery("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain \n" +
			"FROM parenthoods JOIN units ON parent_unit=unit \n" +
			"WHERE child_unit IN(" + susSql + ")", queryParams, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:139:3
		arrNewStartUnits := UnitsT{}
		for i := 0; i < len(rows); i++ {
			objUnitProps := rows[i]
//			if objUnitProps.unit == earlier_unit {
			if objUnitProps.Unit == earlier_unit {
				// :: flattened return for return handleResult(true);
				return true
			}
//			if objUnitProps.is_on_main_chain == 0 && objUnitProps.level > objEarlierUnitProps.level {
			if objUnitProps.Is_on_main_chain == 0 && objUnitProps.Level > objEarlierUnitProps.Level {
//				arrNewStartUnits = append(arrNewStartUnits, objUnitProps.unit)
				arrNewStartUnits = append(arrNewStartUnits, objUnitProps.Unit)
			}
		}
		if len(arrNewStartUnits) <= 0 {
			// :: flattened return for return handleResult(false);
			return false
		}
//		return goUp_2(_.uniq(arrNewStartUnits))
		return goUp_2(arrNewStartUnits.Uniq())
		// >> flattened continuation for conn.query:139:3
	}
	
	return goUp_2(arrLaterUnits)
	// >> flattened continuation for storage.readPropsOfUnits:122:1
}

//func determineIfIncludedOrEqual_sync(conn DBConnT, earlier_unit UnitT, arrLaterUnits UnitsT) bool {
func DetermineIfIncludedOrEqual_sync(conn refDBConnT, earlier_unit UnitT, arrLaterUnits UnitsT) bool {
//	if arrLaterUnits.indexOf(earlier_unit) >= 0 {
	if arrLaterUnits.IndexOf(earlier_unit) >= 0 {
		// :: flattened return for return handleResult(true);
		return true
	}
	// :: flattened return for handleResult(determineIfIncluded(conn, earlier_unit, arrLaterUnits));
	return determineIfIncluded_sync(conn, earlier_unit, arrLaterUnits)
}


// excludes earlier unit

type(
//	ReadDescendantUnitsByAuthorsBeforeMcIndexReturnT = UnitsT
	ReadDescendantUnitsByAuthorsBeforeMcIndexReturnT = []db.UnitMCISequenceAddressRow
)

//func readDescendantUnitsByAuthorsBeforeMcIndex_sync(conn DBConnT, objEarlierUnitProps UnitPropsT, arrAuthorAddresses AddressesT, to_main_chain_index MCIndexT) UnitsT {
func ReadDescendantUnitsByAuthorsBeforeMcIndex_sync(conn refDBConnT, objEarlierUnitProps UnitPropsT, arrAuthorAddresses AddressesT, to_main_chain_index MCIndexT) ReadDescendantUnitsByAuthorsBeforeMcIndexReturnT {
	var(
		goDown_2 func (arrStartUnits UnitsT) ReadDescendantUnitsByAuthorsBeforeMcIndexReturnT
	)
	
	arrUnits := ReadDescendantUnitsByAuthorsBeforeMcIndexReturnT{}
	
	goDown_2 = func (arrStartUnits UnitsT) ReadDescendantUnitsByAuthorsBeforeMcIndexReturnT {
		profiler.Start()
/**
		rows := /* await * /
		conn.query_sync("SELECT units.unit, unit_authors.address AS author_in_list \n" +
			"			FROM parenthoods \n" +
			"			JOIN units ON child_unit=units.unit \n" +
			"			LEFT JOIN unit_authors ON unit_authors.unit=units.unit AND address IN(?) \n" +
			"			WHERE parent_unit IN(?) AND latest_included_mc_index<? AND main_chain_index<=?", DBParamsT{
			arrAuthorAddresses,
			arrStartUnits,
//			objEarlierUnitProps.main_chain_index,
			objEarlierUnitProps.Main_chain_index,
			to_main_chain_index,
		})
 **/
		rcvr := db.UnitMCISequenceAddressesReceiver{}
		queryParams := DBParamsT{}
		aasSql := queryParams.AddAddresses(arrAuthorAddresses)
		susSql := queryParams.AddUnits(arrStartUnits)
		queryParams = append(queryParams,
//			objEarlierUnitProps.main_chain_index,
			objEarlierUnitProps.Main_chain_index,
			to_main_chain_index)
		conn.MustQuery("SELECT units.unit, \n" +
			"	main_chain_index AS mci, (+sequence='good') AS good, \n" +
			"	unit_authors.address AS author_in_list \n" +
			"FROM parenthoods \n" +
			"JOIN units ON child_unit=units.unit \n" +
			"LEFT JOIN unit_authors ON unit_authors.unit=units.unit AND address IN(" + aasSql + ") \n" +
			"WHERE parent_unit IN(" + susSql + ") AND latest_included_mc_index<? AND main_chain_index<=?",
			queryParams, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:180:2
		arrNewStartUnits := UnitsT{}
		for i := 0; i < len(rows); i++ {
			objUnitProps := rows[i]
//			arrNewStartUnits = append(arrNewStartUnits, objUnitProps.unit)
			arrNewStartUnits = append(arrNewStartUnits, objUnitProps.Unit)
//			if objUnitProps.author_in_list {
			if ! objUnitProps.Address.IsNull() {
//				arrUnits = append(arrUnits, objUnitProps.unit)
				arrUnits = append(arrUnits, objUnitProps)
			}
		}
		profiler.Stop("mc-wc-descendants-goDown")
		if len(arrNewStartUnits) <= 0 {
			// :: flattened return for return handleUnits(arrUnits);
			return arrUnits
		}
		return goDown_2(arrNewStartUnits)
		// >> flattened continuation for conn.query:180:2
	}
	
	profiler.Start()
/**
	rows := /* await * /
	conn.query_sync(// _left_ join forces use of indexes in units
//	"SELECT unit FROM units " + db.forceIndex("byMcIndex") + " LEFT JOIN unit_authors USING(unit) \n" +
	"SELECT unit FROM units " + conn.ForceIndex("byMcIndex") + " LEFT JOIN unit_authors USING(unit) \n" +
		"		WHERE latest_included_mc_index>=? AND main_chain_index>? AND main_chain_index<=? AND latest_included_mc_index<? AND address IN(?)", DBParamsT{
//		objEarlierUnitProps.main_chain_index,
		objEarlierUnitProps.Main_chain_index,
//		objEarlierUnitProps.main_chain_index,
		objEarlierUnitProps.Main_chain_index,
		to_main_chain_index,
		to_main_chain_index,
		arrAuthorAddresses,
	})
 **/
	rcvr := db.UnitMCISequenceAddressesReceiver{}
	conn.MustQuery("SELECT units.unit, \n" +
		"	main_chain_index AS mci, (+sequence='good') AS good, \n" +
		"	unit_authors.address \n" +
		"FROM units \n" +
		conn.ForceIndex("byMcIndex") + " \n" +
		"LEFT JOIN unit_authors USING(unit) \n" +
		"WHERE main_chain_index BETWEEN (?)+1 AND ? \n" +
		"	AND latest_included_mc_index BETWEEN ? AND (?)-1 \n", DBParamsT{
//			objEarlierUnitProps.main_chain_index,
			objEarlierUnitProps.Main_chain_index,
			to_main_chain_index,
//			objEarlierUnitProps.main_chain_index,
			objEarlierUnitProps.Main_chain_index,
			to_main_chain_index,
		}, &rcvr)
	rows := rcvr.Rows
	for _, row := range rows {
		if arrAuthorAddresses.IndexOf(row.Address) == -1 {
			row.X_skip = true
		}
	}
	// << flattened continuation for conn.query:204:1
/**
	arrUnits = // .. not flattening for Array.map
	rows.map(func (row rowT) {*returns*} {
		return row.unit
	})
 **/
	arrUnits = make(ReadDescendantUnitsByAuthorsBeforeMcIndexReturnT, 0, len(rows))
	for _, row := range rows {
		if row.X_skip { continue }
		arrUnits = append(arrUnits, row)
	}
	profiler.Stop("mc-wc-descendants-initial")
//	return goDown_2(UnitsT{ objEarlierUnitProps.unit })
	return goDown_2(UnitsT{ objEarlierUnitProps.Unit })
	// >> flattened continuation for conn.query:204:1
}



// excludes earlier unit
func readDescendantUnitsBeforeLandingOnMc_sync(conn DBConnT, objEarlierUnitProps UnitPropsT, arrLaterUnitProps UnitPropssT) (UnitsT, UnitsT) {
	var(
		goDown_3 func (arrStartUnits UnitsT) (UnitsT, UnitsT)
	)

/**	
	max_later_limci := Math.max.apply(nil, // .. not flattening for Array.map
	arrLaterUnitProps.map(func (objLaterUnitProps PropsT) {*returns*} {
		return objLaterUnitProps.latest_included_mc_index
	}))
 **/
	max_later_limci := MCIndexT(-1)
	for _, objLaterUnitProps := range arrLaterUnitProps {
		limci := objLaterUnitProps.Latest_included_mc_index
		if max_later_limci < limci {
			max_later_limci = limci
		}
	}
/**
	max_later_level := Math.max.apply(nil, // .. not flattening for Array.map
	arrLaterUnitProps.map(func (objLaterUnitProps PropsT) {*returns*} {
		return objLaterUnitProps.level
	}))
 **/
	max_later_level := LevelT(-1)
	for _, objLaterUnitProps := range arrLaterUnitProps {
		level := objLaterUnitProps.Level
		if max_later_level < level {
			max_later_level = level
		}
	}	
	arrLandedUnits := UnitsT{}
	// units that landed on MC before max_later_limci, they are already included in at least one of later units
	arrUnlandedUnits := UnitsT{}
	// direct shoots to later units, without touching the MC
	
	goDown_3 = func (arrStartUnits UnitsT) (UnitsT, UnitsT) {
/**
		rows := /* await * /
		conn.query_sync("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain \n" +
			"			FROM parenthoods JOIN units ON child_unit=unit \n" +
			"			WHERE parent_unit IN(?) AND latest_included_mc_index<? AND level<=?", DBParamsT{
			arrStartUnits,
//			objEarlierUnitProps.main_chain_index,
			objEarlierUnitProps.Main_chain_index,
			max_later_level,
		})
 **/
		rcvr := db.UnitMCPropsReceiver{}
		queryParams := DBParamsT{}
		susSql := queryParams.AddUnits(arrStartUnits)
		queryParams = append(queryParams,
//			objEarlierUnitProps.main_chain_index,
			objEarlierUnitProps.Main_chain_index,
			max_later_level)
		conn.MustQuery("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain \n" +
			"FROM parenthoods JOIN units ON child_unit=unit \n" +
			"WHERE parent_unit IN(" + susSql + ") AND latest_included_mc_index<? AND level<=?",
			queryParams, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:229:2
		arrNewStartUnits := UnitsT{}
		for i := 0; i < len(rows); i++ {
			objUnitProps := rows[i]
			//if (objUnitProps.latest_included_mc_index >= objEarlierUnitProps.main_chain_index)
			//    continue;
			//if (objUnitProps.level > max_later_level)
			//    continue;
//			arrNewStartUnits = append(arrNewStartUnits, objUnitProps.unit)
			arrNewStartUnits = append(arrNewStartUnits, objUnitProps.Unit)
//			if objUnitProps.main_chain_index != nil && objUnitProps.main_chain_index <= max_later_limci {
			if ! objUnitProps.Main_chain_index.IsNull() && objUnitProps.Main_chain_index <= max_later_limci {
				// exclude free balls!
//				arrLandedUnits = append(arrLandedUnits, objUnitProps.unit)
				arrLandedUnits = append(arrLandedUnits, objUnitProps.Unit)
			} else {
//				arrUnlandedUnits = append(arrUnlandedUnits, objUnitProps.unit)
				arrUnlandedUnits = append(arrUnlandedUnits, objUnitProps.Unit)
			}
		}
		if len(arrNewStartUnits) <= 0 {
			// :: flattened return for return handleUnits(arrLandedUnits, arrUnlandedUnits);
//			return meta.returnArguments(arrLandedUnits, arrUnlandedUnits)
			return arrLandedUnits, arrUnlandedUnits
		}
		return goDown_3(arrNewStartUnits)
		// >> flattened continuation for conn.query:229:2
	}
	
//	return goDown_3(UnitsT{ objEarlierUnitProps.unit })
	return goDown_3(UnitsT{ objEarlierUnitProps.Unit })
}

// includes later units
func readAscendantUnitsAfterTakingOffMc_sync(conn DBConnT, objEarlierUnitProps UnitPropsT, arrLaterUnitProps UnitPropssT) (UnitsT, UnitsT) {
	var(
		goUp_3 func (arrStartUnits UnitsT) (UnitsT, UnitsT)
	)
	
	arrLaterUnits := make(UnitsT, len(arrLaterUnitProps), len(arrLaterUnitProps))
//	for objLaterUnitProps, _k := range arrLaterUnitProps {
	for _k, objLaterUnitProps := range arrLaterUnitProps {
//		arrLaterUnits[_k] := objLaterUnitProps.unit
		arrLaterUnits[_k] = objLaterUnitProps.Unit
	}
/**
	max_later_limci := Math.max.apply(nil, // .. not flattening for Array.map
	arrLaterUnitProps.map(func (objLaterUnitProps PropsT) {*returns*} {
		return objLaterUnitProps.latest_included_mc_index
	}))
 **/
	max_later_limci := MCIndexT(-1)
	for _, objLaterUnitProps := range arrLaterUnitProps {
		limci := objLaterUnitProps.Latest_included_mc_index
		if max_later_limci < limci {
			max_later_limci = limci
		}
	}
	arrLandedUnits := UnitsT{}
	// units that took off MC after earlier unit's MCI, they already include the earlier unit
	arrUnlandedUnits := UnitsT{}
	// direct shoots from earlier units, without touching the MC
	
	// .. not flattening for Array.forEach
//	for objUnitProps, _ := range arrLaterUnitProps {
	for _, objUnitProps := range arrLaterUnitProps {
//		if objUnitProps.latest_included_mc_index >= objEarlierUnitProps.main_chain_index {
		if objUnitProps.Latest_included_mc_index >= objEarlierUnitProps.Main_chain_index {
//			arrLandedUnits = append(arrLandedUnits, objUnitProps.unit)
			arrLandedUnits = append(arrLandedUnits, objUnitProps.Unit)
		} else {
//			arrUnlandedUnits = append(arrUnlandedUnits, objUnitProps.unit)
			arrUnlandedUnits = append(arrUnlandedUnits, objUnitProps.Unit)
		}
	}
	
	goUp_3 = func (arrStartUnits UnitsT) (UnitsT, UnitsT) {
/**
		rows := /* await * /
		conn.query_sync("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain \n" +
			"			FROM parenthoods JOIN units ON parent_unit=unit \n" +
			"			WHERE child_unit IN(?) AND (main_chain_index>? OR main_chain_index IS NULL) AND level>=?", DBParamsT{
			arrStartUnits,
			max_later_limci,
//			objEarlierUnitProps.level,
			objEarlierUnitProps.Level,
		})
 **/
		rcvr := db.UnitMCPropsReceiver{}
		queryParams := DBParamsT{}
		susSql := queryParams.AddUnits(arrStartUnits)
		queryParams = append(queryParams,
			max_later_limci,
//			objEarlierUnitProps.level,
			objEarlierUnitProps.Level)
		conn.MustQuery("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain \n" +
			"FROM parenthoods JOIN units ON parent_unit=unit \n" +
			"WHERE child_unit IN(" + susSql + ") AND (main_chain_index>? OR main_chain_index IS NULL) AND level>=?",
			queryParams, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:272:2
		arrNewStartUnits := UnitsT{}
//		for i := 0; i < len(rows); i++ {
//			objUnitProps := rows[i]
		for _, objUnitProps := range rows {
			//if (objUnitProps.main_chain_index <= max_later_limci)
			//    continue;
			//if (objUnitProps.level < objEarlierUnitProps.level)
			//    continue;
//			arrNewStartUnits = append(arrNewStartUnits, objUnitProps.unit)
			arrNewStartUnits = append(arrNewStartUnits, objUnitProps.Unit)
//			if objUnitProps.latest_included_mc_index >= objEarlierUnitProps.main_chain_index {
			if objUnitProps.Latest_included_mc_index >= objEarlierUnitProps.Main_chain_index {
//				arrLandedUnits = append(arrLandedUnits, objUnitProps.unit)
				arrLandedUnits = append(arrLandedUnits, objUnitProps.Unit)
			} else {
//				arrUnlandedUnits = append(arrUnlandedUnits, objUnitProps.unit)
				arrUnlandedUnits = append(arrUnlandedUnits, objUnitProps.Unit)
			}
		}
		if len(arrNewStartUnits) <= 0 {
			// :: flattened return for return handleUnits(arrLandedUnits, arrUnlandedUnits);
//			return meta.returnArguments(arrLandedUnits, arrUnlandedUnits)
			return arrLandedUnits, arrUnlandedUnits
		}
		return goUp_3(arrNewStartUnits)
		// >> flattened continuation for conn.query:272:2
	}
	
	return goUp_3(arrLaterUnits)
}


//exports.compareUnitsByProps = compareUnitsByProps
//exports.compareUnits = compareUnits

//exports.determineIfIncluded = determineIfIncluded
//exports.determineIfIncludedOrEqual = determineIfIncludedOrEqual

//exports.readDescendantUnitsByAuthorsBeforeMcIndex = readDescendantUnitsByAuthorsBeforeMcIndex

// used only in majority_witnessing.js which is not used itself
//exports.readDescendantUnitsBeforeLandingOnMc = readDescendantUnitsBeforeLandingOnMc
//exports.readAscendantUnitsAfterTakingOffMc = readAscendantUnitsAfterTakingOffMc


// converted golang end

