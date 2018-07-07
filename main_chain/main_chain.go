
// converted golang begin

package main_chain

import(
	"fmt"
//	"strings"
	"errors"
	"sort"
	"math"

 _core	"nodejs/core"
 JSON	"nodejs/json"
	"nodejs/console"

 .	"github.com/byteball/go-byteballcore/types"

)

import(
// _		"lodash"
//		"async"
		"github.com/byteball/go-byteballcore/db"
		"github.com/byteball/go-byteballcore/constants"
		"github.com/byteball/go-byteballcore/storage"
		"github.com/byteball/go-byteballcore/graph"
 objectHash	"github.com/byteball/go-byteballcore/object_hash"
		"github.com/byteball/go-byteballcore/paid_witnessing"
		"github.com/byteball/go-byteballcore/headers_commission"
		"github.com/byteball/go-byteballcore/mutex"
 eventBus	"github.com/byteball/go-byteballcore/event_bus"
		"github.com/byteball/go-byteballcore/profiler"
		"github.com/byteball/go-byteballcore/breadcrumbs"
)

type(
	DBConnT		= db.DBConnT
	DBParamsT	= db.DBParamsT

	refDBConnT	= *DBConnT

	refPropsT	= *PropsT
	refXPropsT	= *XPropsT
	refJointT	= *JointT

	PropsByUnitMapT	= map[UnitT] refPropsT
	XPropsByUnitMapT = map[UnitT] refXPropsT
	MCIndexByUnitMapT  map[UnitT] MCIndexT

	XPropsT		= storage.XPropsT
)

func (map0 MCIndexByUnitMapT) Compare(map1 MCIndexByUnitMapT) int {
	d := len(map0) - len(map1)
	if d != 0 { return d }

	for unit0, mcindex0 := range map0 {
		mcindex1, _exists1 := map1[unit0]
		if ! _exists1 { return +1 }
		d := int(mcindex0 - mcindex1)
		if d != 0 { return d }
	}

	for unit1, _ := range map1 {
		_, _exists0 := map0[unit1]
		if ! _exists0 { return -1 }
	}

	return 0
}


// override when adding units which caused witnessed level to significantly retreat

//arrRetreatingUnits := UnitsT{
var arrRetreatingUnits UnitsT = UnitsT{
	"+5ntioHT58jcFb8oVc+Ff4UvO5UvYGRcrGfYIofGUW8=",
	"C/aPdM0sODPLC3NqJPWdZlqmV8B4xxf2N/+HSEi0sKU=",
	"sSev6hvQU86SZBemy9CW2lJIko2jZDoY55Lm3zf2QU4=",
	"19GglT3uZx1WmfWstLb3yIa85jTic+t01Kpe6s5gTTA=",
	"Hyi2XVdZ/5D3H/MhwDL/jRWHp3F/dQTmwemyUHW+Urg=",
	"xm0kFeKh6uqSXx6UUmc2ucgsNCU5h/e6wxSMWirhOTo=",
}


//func updateMainChain_sync(conn DBConnT, from_unit UnitT, last_added_unit UnitT)  {
func UpdateMainChain_sync(conn refDBConnT, from_unit UnitT, last_added_unit UnitT)  {
	var(
		findNextUpMainChainUnit_sync func (unit UnitT) UnitT
		goUpFromUnit func (unit UnitT) 
		checkNotRebuildingStableMainChainAndGoDown func (last_main_chain_index MCIndexT, last_main_chain_unit UnitT) 
		goDownAndUpdateMainChainIndex func (last_main_chain_index MCIndexT, last_main_chain_unit UnitT) 
		updateLatestIncludedMcIndex func (last_main_chain_index MCIndexT, bRebuiltMc bool) 
		readLastStableMcUnit_sync func () UnitT
		updateStableMcFlag func () 
		createListOfBestChildren_sync func (arrParentUnits UnitsT) UnitsT
		finish func () 
	)
	
	arrAllParents := UnitsT{}
	arrNewMcUnits := UnitsT{}
	
	// if unit === null, read free balls
	findNextUpMainChainUnit_sync = func (unit UnitT) UnitT {
		var(
//			handleProps func (props PropsT) UnitT
			handleProps func (props *db.BestParentUnitRow) UnitT
//			readLastUnitProps_sync func () PropsT
			readLastUnitProps_sync func () *db.BestParentUnitRow
		)
		
//		handleProps = func (props PropsT) UnitT {
		handleProps = func (props *db.BestParentUnitRow) UnitT {
//			if props.best_parent_unit == nil {
			if props.Best_parent_unit.IsNull() {
				_core.Throw("best parent is null")
			}
//			console.Log("unit " + unit + ", best parent " + props.best_parent_unit + ", wlevel " + props.witnessed_level)
			console.Log("unit %s, best parent %s, wlevel %d", unit, props.Best_parent_unit, props.Witnessed_level)
			// :: flattened return for handleUnit(props.best_parent_unit);
//			return props.best_parent_unit
			return props.Best_parent_unit
		}
//		readLastUnitProps_sync = func () PropsT {
		readLastUnitProps_sync = func () *db.BestParentUnitRow {
/**
			rows := /* await * /
			conn.query_sync("SELECT unit AS best_parent_unit, witnessed_level \n" +
				"				FROM units WHERE is_free=1 \n" +
				"				ORDER BY witnessed_level DESC, \n" +
				"					level-witnessed_level ASC, \n" +
				"					unit ASC \n" +
				"				LIMIT 5")
 **/
			//rcvr := db.Props{}
			rcvr := db.BestParentUnitsReceiver{}
			conn.MustQuery("SELECT unit AS best_parent_unit, witnessed_level \n" +
				"FROM units WHERE is_free=1 \n" +
				"ORDER BY witnessed_level DESC, \n" +
				"	level-witnessed_level ASC, \n" +
				"	unit ASC \n" +
				"LIMIT 5", DBParamsT{}, &rcvr)
			rows := rcvr.Rows
			// << flattened continuation for conn.query:42:3
			if len(rows) == 0 {
				_core.Throw("no free units?")
			}
			if len(rows) > 1 {
				arrParents := make(UnitsT, len(rows), len(rows))
//				for row, _k := range rows {
				for _k, row := range rows {
//					arrParents[_k] = row.best_parent_unit
					arrParents[_k] = row.Best_parent_unit
				}
				arrAllParents = arrParents
				for i := 0; i < len(arrRetreatingUnits); i++ {
//					n := arrParents.indexOf(arrRetreatingUnits[i])
					n := arrRetreatingUnits[i].IndexOf(arrParents)
					if n >= 0 {
						// :: flattened return for return handleLastUnitProps(rows[n]);
//						return rows[n]
						return &rows[n]
					}
				}
			}
			/*
								// override when adding +5ntioHT58jcFb8oVc+Ff4UvO5UvYGRcrGfYIofGUW8= which caused witnessed level to significantly retreat
								if (rows.length === 2 && (rows[1].best_parent_unit === '+5ntioHT58jcFb8oVc+Ff4UvO5UvYGRcrGfYIofGUW8=' || rows[1].best_parent_unit === 'C/aPdM0sODPLC3NqJPWdZlqmV8B4xxf2N/+HSEi0sKU=' || rows[1].best_parent_unit === 'sSev6hvQU86SZBemy9CW2lJIko2jZDoY55Lm3zf2QU4=') && (rows[0].best_parent_unit === '3XJT1iK8FpFeGjwWXd9+Yu7uJp7hM692Sfbb5zdqWCE=' || rows[0].best_parent_unit === 'TyY/CY8xLGvJhK6DaBumj2twaf4y4jPC6umigAsldIA=' || rows[0].best_parent_unit === 'VKX2Nsx2W1uQYT6YajMGHAntwNuSMpAAlxF7Y98tKj8='))
									return handleLastUnitProps(rows[1]);
								
			 */
			// :: flattened return for handleLastUnitProps(rows[0]);
//			return rows[0]
			return &rows[0]
			// >> flattened continuation for conn.query:42:3
		}

//		if unit {
		if ! unit.IsNull() {
			return handleProps(/* await */
//			storage.readStaticUnitProps_sync(conn, unit))
//			storage.ReadStaticUnitProps_sync(conn, unit))
			(func () *db.BestParentUnitRow {
				props := storage.ReadStaticUnitProps_sync(conn, unit)
				return &db.BestParentUnitRow{
					Best_parent_unit: props.Best_parent_unit,
					Witnessed_level: props.Witnessed_level,
				}
			})())
		} else {
			return handleProps(/* await */
			readLastUnitProps_sync())
		}
	}
	
	goUpFromUnit = func (unit UnitT)  {
//		if storage.isGenesisUnit(unit) {
		if storage.IsGenesisUnit(unit) {
			checkNotRebuildingStableMainChainAndGoDown(0, unit)
			return
		}
		
		profiler.Start()
		best_parent_unit := /* await */
		findNextUpMainChainUnit_sync(unit)
		// << flattened continuation for findNextUpMainChainUnit:81:2
		objBestParentUnitProps := /* await */
//		storage.readUnitProps_sync(conn, best_parent_unit)
		storage.ReadUnitProps_sync(conn, best_parent_unit)
		// << flattened continuation for storage.readUnitProps:82:3
//		objBestParentUnitProps2 := storage.assocUnstableUnits[best_parent_unit]
		objBestParentUnitProps2 := storage.AssocUnstableUnits[best_parent_unit]
//		if ! objBestParentUnitProps2 {
		if ! (objBestParentUnitProps2 != nil) {
//			if storage.isGenesisUnit(best_parent_unit) {
			if storage.IsGenesisUnit(best_parent_unit) {
//				objBestParentUnitProps2 = storage.assocStableUnits[best_parent_unit]
//				objBestParentUnitProps2 = storage.AssocStableUnits[best_parent_unit]
				objBestParentUnitProps2 = &XPropsT{
					PropsT: *storage.AssocStableUnits[best_parent_unit],
				}
			} else {
//				_core.Throw("unstable unit not found: " + best_parent_unit)
				_core.Throw("unstable unit not found: %s", best_parent_unit)
			}
		}
//		objBestParentUnitPropsForCheck := _.cloneDeep(objBestParentUnitProps2)
		objBestParentUnitPropsForCheck := *objBestParentUnitProps2
//		objBestParentUnitPropsForCheck.parent_units = nil
		objBestParentUnitPropsForCheck.Parent_units = nil
//		if ! _.isEqual(objBestParentUnitPropsForCheck, objBestParentUnitProps) {
		if objBestParentUnitPropsForCheck.PropsT != *objBestParentUnitProps {
//			throwError("different props, db: " + JSON.stringify(objBestParentUnitProps) + ", unstable: " + JSON.stringify(objBestParentUnitProps2))
			throwError(fmt.Sprintf("different props, db: %s, unstable: %s", JSON.Stringify(objBestParentUnitProps), JSON.Stringify(objBestParentUnitProps2)))
		}
//		if ! objBestParentUnitProps.is_on_main_chain {
		if ! (objBestParentUnitProps.Is_on_main_chain != 0) {
/**
			/* await * /
			conn.query_sync("UPDATE units SET is_on_main_chain=1, main_chain_index=NULL WHERE unit=?", DBParamsT{ best_parent_unit })
 **/
			conn.MustExec("UPDATE units SET is_on_main_chain=1, main_chain_index=NULL WHERE unit=?", DBParamsT{ best_parent_unit })
			// << flattened continuation for conn.query:95:5
//			objBestParentUnitProps2.is_on_main_chain = 1
			objBestParentUnitProps2.Is_on_main_chain = 1
//			objBestParentUnitProps2.main_chain_index = nil
			objBestParentUnitProps2.Main_chain_index = MCIndexT_Null
			arrNewMcUnits = append(arrNewMcUnits, best_parent_unit)
			profiler.Stop("mc-goUpFromUnit")
//			goUpFromUnit(best_parent_unit)
			goUpFromUnit(best_parent_unit)
			// >> flattened continuation for conn.query:95:5
		} else {
			profiler.Stop("mc-goUpFromUnit")
//			if unit == nil {
			if unit.IsNull() {
//				updateLatestIncludedMcIndex(objBestParentUnitProps.main_chain_index, false)
				updateLatestIncludedMcIndex(objBestParentUnitProps.Main_chain_index, false)
			} else {
//				checkNotRebuildingStableMainChainAndGoDown(objBestParentUnitProps.main_chain_index, best_parent_unit)
				checkNotRebuildingStableMainChainAndGoDown(objBestParentUnitProps.Main_chain_index, best_parent_unit)
			}
		}
		// >> flattened continuation for storage.readUnitProps:82:3
		// >> flattened continuation for findNextUpMainChainUnit:81:2
	}
	
	checkNotRebuildingStableMainChainAndGoDown = func (last_main_chain_index MCIndexT, last_main_chain_unit UnitT)  {
//		console.Log("checkNotRebuildingStableMainChainAndGoDown " + from_unit)
		console.Log("checkNotRebuildingStableMainChainAndGoDown %s", from_unit)
		profiler.Start()
/**
		rows := /* await * /
		conn.query_sync("SELECT unit FROM units WHERE is_on_main_chain=1 AND main_chain_index>? AND is_stable=1", DBParamsT{ last_main_chain_index })
 **/
		rcvr := db.UnitsReceiver{}
		conn.MustQuery("SELECT unit FROM units WHERE is_on_main_chain=1 AND main_chain_index>? AND is_stable=1", DBParamsT{ last_main_chain_index }, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:116:2
		profiler.Stop("mc-checkNotRebuilding")
		if len(rows) > 0 {
			units_ := make(UnitsT, len(rows), len(rows))
//			for row, _k := range rows {
			for _k, row := range rows {
//				units_[_k] := row.unit
				units_[_k] = row.Unit
			}
//			units := units_.join(", ")
			units := units_.Join(", ")
//			allps := arrAllParents.join(", ")
			allps := arrAllParents.Join(", ")
//			_core.Throw("removing stable units " + units + " from MC after adding " + last_added_unit + " with all parents " + allps)
			_core.Throw("removing stable units %s from MC after adding %s with all parents %s", units, last_added_unit, allps)
		}
		goDownAndUpdateMainChainIndex(last_main_chain_index, last_main_chain_unit)
		// >> flattened continuation for conn.query:116:2
	}
	
	goDownAndUpdateMainChainIndex = func (last_main_chain_index MCIndexT, last_main_chain_unit UnitT)  {
		profiler.Start()
/**
		/* await * /
		conn.query_sync(//"UPDATE units SET is_on_main_chain=0, main_chain_index=NULL WHERE is_on_main_chain=1 AND main_chain_index>?", 
 **/
		conn.MustExec(//"UPDATE units SET is_on_main_chain=0, main_chain_index=NULL WHERE is_on_main_chain=1 AND main_chain_index>?", 
		"UPDATE units SET is_on_main_chain=0, main_chain_index=NULL WHERE main_chain_index>?", DBParamsT{ last_main_chain_index })
		// << flattened continuation for conn.query:134:2
//		for unit := range storage.assocUnstableUnits {
		for unit := range storage.AssocUnstableUnits {
//			o := storage.assocUnstableUnits[unit]
			o := storage.AssocUnstableUnits[unit]
//			if o.Main_chain_index > last_main_chain_index {
			if o.Main_chain_index > last_main_chain_index {
				o.Is_on_main_chain = 0
				o.Main_chain_index = MCIndexT_Null
			}
		}
		main_chain_index := last_main_chain_index
//[uu]		main_chain_unit := last_main_chain_unit
/**
		rows := /* await * /
		conn.query_sync("SELECT unit FROM units WHERE is_on_main_chain=1 AND main_chain_index IS NULL ORDER BY level")
 **/
		rcvr := db.UnitsReceiver{}
		conn.MustQuery("SELECT unit FROM units WHERE is_on_main_chain=1 AND main_chain_index IS NULL ORDER BY level", DBParamsT{}, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:148:4
		if len(rows) == 0 {
			//if (last_main_chain_index > 0)
//			_core.Throw("no unindexed MC units after adding " + last_added_unit)
			_core.Throw("no unindexed MC units after adding %s", last_added_unit)
		}
		arrDbNewMcUnits := make(UnitsT, len(rows), len(rows))
//		for row, _k := range rows {
		for _k, row := range rows {
//			arrDbNewMcUnits[_k] := row.unit
			arrDbNewMcUnits[_k] = row.Unit
		}
//		arrNewMcUnits.reverse()
		for _l, _r := 0, len(arrNewMcUnits)-1; _l < _r; _l, _r = _l+1, _r-1 {
			arrNewMcUnits[_l], arrNewMcUnits[_r] = arrNewMcUnits[_r], arrNewMcUnits[_l]
		}
//		if ! _.isEqual(arrNewMcUnits, arrDbNewMcUnits) {
		if ! (arrNewMcUnits.Compare(arrDbNewMcUnits) == 0) {
//			throwError("different new MC units, arr: " + JSON.stringify(arrNewMcUnits) + ", db: " + JSON.stringify(arrDbNewMcUnits))
			throwError(fmt.Sprintf("different new MC units, arr: %s, db: %s", JSON.Stringify(arrNewMcUnits), JSON.Stringify(arrDbNewMcUnits)))
		}
		err := (func () ErrorT {
		  // :: inlined async.eachSeries:163:6
//		  for row := range rows {
		  for _, row := range rows {
//		    _err := (func (row rowT) ErrorT {
		    _err := (func (row struct{ Unit UnitT }) ErrorT {
		    	var(
		    		goUp_1 func (arrStartUnits UnitsT) 
		    		updateMc func () 
		    	)
		    	
		    	main_chain_index++
//		    	arrUnits := UnitsT{ row.unit }
		    	arrUnits := UnitsT{ row.Unit }
		    	
		    	goUp_1 = func (arrStartUnits UnitsT)  {
/**
		    		rows := /* await * /
		    		conn.query_sync("SELECT DISTINCT unit \n" +
		    			"										FROM parenthoods JOIN units ON parent_unit=unit \n" +
		    			"										WHERE child_unit IN(?) AND main_chain_index IS NULL", DBParamsT{ arrStartUnits })
 **/
				rcvr := db.UnitsReceiver{}
				queryParams := DBParamsT{}
				susSql := queryParams.AddUnits(arrStartUnits)
		    		conn.MustQuery("SELECT DISTINCT unit \n" +
		    			"FROM parenthoods JOIN units ON parent_unit=unit \n" +
		    			"WHERE child_unit IN("+ susSql +") AND main_chain_index IS NULL", queryParams, &rcvr)
				rows := rcvr.Rows
		    		// << flattened continuation for conn.query:170:9
		    		arrNewStartUnits := make(UnitsT, len(rows), len(rows))
//		    		for row, _k := range rows {
		    		for _k, row := range rows {
//		    			arrNewStartUnits[_k] := row.unit
		    			arrNewStartUnits[_k] = row.Unit
		    		}
		    		arrNewStartUnits2 := UnitsT{}
		    		// .. not flattening for Array.forEach
//		    		for start_unit, _ := range arrStartUnits {
		    		for _, start_unit := range arrStartUnits {
		    			// .. not flattening for Array.forEach
//		    			for parent_unit, _ := range storage.assocUnstableUnits[start_unit].parent_units {
		    			for _, parent_unit := range storage.AssocUnstableUnits[start_unit].Parent_units {
//		    				if storage.assocUnstableUnits[parent_unit] && storage.assocUnstableUnits[parent_unit].main_chain_index == nil && arrNewStartUnits2.indexOf(parent_unit) == - 1 {
		    				if unit, _exists := storage.AssocUnstableUnits[parent_unit] ; _exists && unit.Main_chain_index.IsNull() && parent_unit.IndexOf(arrNewStartUnits2) == - 1 {
		    					arrNewStartUnits2 = append(arrNewStartUnits2, parent_unit)
		    				}
		    			}
		    		}
//		    		if ! _.isEqual(arrNewStartUnits.sort(), arrNewStartUnits2.sort()) {
				sort.Slice(arrNewStartUnits, func (i, j int) bool {
					return arrNewStartUnits[i] < arrNewStartUnits[j]
				})
				sort.Slice(arrNewStartUnits2, func (i, j int) bool {
					return arrNewStartUnits2[i] < arrNewStartUnits2[j]
				})
		    		if ! (arrNewStartUnits.Compare(arrNewStartUnits2) == 0) {
//		    			throwError("different new start units, arr: " + JSON.stringify(arrNewStartUnits2) + ", db: " + JSON.stringify(arrNewStartUnits))
		    			throwError(fmt.Sprintf("different new start units, arr: %s, db: %s", JSON.Stringify(arrNewStartUnits2), JSON.Stringify(arrNewStartUnits)))
		    		}
		    		if len(arrNewStartUnits) == 0 {
		    			updateMc()
		    			return
		    		}
//		    		arrUnits = arrUnits.concat(arrNewStartUnits)
		    		arrUnits = append(arrUnits, arrNewStartUnits...)
		    		goUp_1(arrNewStartUnits)
		    		// >> flattened continuation for conn.query:170:9
		    	}
		    	
		    	updateMc = func ()  {
		    		// .. not flattening for Array.forEach
//		    		for unit, _ := range arrUnits {
		    		for _, unit := range arrUnits {
//		    			storage.assocUnstableUnits[unit].main_chain_index = main_chain_index
		    			storage.AssocUnstableUnits[unit].Main_chain_index = main_chain_index
		    		}
/**
		    		strUnitList := arrUnits.map(db.escape).join(", ")
		    		/* await * /
		    		conn.query_sync("UPDATE units SET main_chain_index=? WHERE unit IN(" + strUnitList + ")", DBParamsT{ main_chain_index })
 **/
				queryParams := DBParamsT{ main_chain_index }
				usSql := queryParams.AddUnits(arrUnits)
		    		conn.MustExec("UPDATE units SET main_chain_index=? WHERE unit IN(" + usSql + ")", queryParams)
		    		// << flattened continuation for conn.query:199:9
/**
		    		/* await * /
		    		conn.query_sync("UPDATE unit_authors SET _mci=? WHERE unit IN(" + strUnitList + ")", DBParamsT{ main_chain_index })
 **/
		    		conn.MustExec("UPDATE unit_authors SET _mci=? WHERE unit IN(" + usSql + ")", queryParams)
		    		// << flattened continuation for conn.query:200:10
		    		// :: flattened return for cb();
		    		// ** need 1 return(s) instead of 0
		    		return 
		    		// >> flattened continuation for conn.query:200:10
		    		// >> flattened continuation for conn.query:199:9
		    	}
		    	
		    	goUp_1(arrUnits)
			return nil
		    })(row)
		    if _err != nil { return _err }
		  }
		  return nil
		})()
		// << flattened continuation for async.eachSeries:163:6
		console.Log("goDownAndUpdateMainChainIndex done")
//		if err {
		if err != nil {
			_core.Throw("goDownAndUpdateMainChainIndex eachSeries failed")
		}
/**
		/* await * /
		conn.query_sync("UPDATE unit_authors SET _mci=NULL WHERE unit IN(SELECT unit FROM units WHERE main_chain_index IS NULL)")
 **/
		conn.MustExec("UPDATE unit_authors SET _mci=NULL WHERE unit IN(SELECT unit FROM units WHERE main_chain_index IS NULL)", DBParamsT{})
		// << flattened continuation for conn.query:213:8
		profiler.Stop("mc-goDown")
		updateLatestIncludedMcIndex(last_main_chain_index, true)
		// >> flattened continuation for conn.query:213:8
		// >> flattened continuation for async.eachSeries:163:6
		// >> flattened continuation for conn.query:148:4
		// >> flattened continuation for conn.query:134:2
	}
	
	updateLatestIncludedMcIndex = func (last_main_chain_index MCIndexT, bRebuiltMc bool)  {
		var(
			checkAllLatestIncludedMcIndexesAreSet func () 
			propagateLIMCI func () 
//			loadUnitProps_sync func (unit UnitT) PropsT
			loadUnitProps_sync func (unit UnitT) refXPropsT
			calcLIMCIs_sync func () 
		)

//		assocChangedUnits := make(PropsByUnitMapT)
		assocChangedUnits := make(XPropsByUnitMapT)
		assocLimcisByUnit := make(MCIndexByUnitMapT)
		assocDbLimcisByUnit := make(MCIndexByUnitMapT)
		
		checkAllLatestIncludedMcIndexesAreSet = func ()  {
			profiler.Start()
//			if ! _.isEqual(assocDbLimcisByUnit, assocLimcisByUnit) {
			if ! (assocDbLimcisByUnit.Compare(assocLimcisByUnit) == 0) {
//				throwError("different  LIMCIs, mem: " + JSON.stringify(assocLimcisByUnit) + ", db: " + JSON.stringify(assocDbLimcisByUnit))
				throwError(fmt.Sprintf("different  LIMCIs, mem: %s, db: %s", JSON.Stringify(assocLimcisByUnit), JSON.Stringify(assocDbLimcisByUnit)))
			}
/**
			rows := /* await * /
			conn.query_sync("SELECT unit FROM units WHERE latest_included_mc_index IS NULL AND level!=0")
 **/
			rcvr := db.UnitsReceiver{}
			conn.MustQuery("SELECT unit FROM units WHERE latest_included_mc_index IS NULL AND level!=0", DBParamsT{}, &rcvr)
			rows := rcvr.Rows
			// << flattened continuation for conn.query:234:3
			if len(rows) > 0 {
//				_core.Throw(len(rows) + " units have latest_included_mc_index=NULL, e.g. unit " + rows[0].unit)
				_core.Throw("%d units have latest_included_mc_index=NULL, e.g. unit %s", len(rows), rows[0].Unit)
			}
			profiler.Stop("mc-limci-check")
			updateStableMcFlag()
			// >> flattened continuation for conn.query:234:3
		}
		
		propagateLIMCI = func ()  {
//			console.Log("propagateLIMCI " + last_main_chain_index)
			console.Log("propagateLIMCI %d", last_main_chain_index)
			profiler.Start()
/**
			rows := /* await * /
			conn.query_sync(/*
							"UPDATE units AS punits \n\
							JOIN parenthoods ON punits.unit=parent_unit \n\
							JOIN units AS chunits ON child_unit=chunits.unit \n\
							SET chunits.latest_included_mc_index=punits.latest_included_mc_index \n\
							WHERE (chunits.main_chain_index > ? OR chunits.main_chain_index IS NULL) \n\
								AND (chunits.latest_included_mc_index IS NULL OR chunits.latest_included_mc_index < punits.latest_included_mc_index)",
							[last_main_chain_index],
							function(result){
								(result.affectedRows > 0) ? propagateLIMCI() : checkAllLatestIncludedMcIndexesAreSet();
							}
							
			 * /
			"SELECT punits.latest_included_mc_index, chunits.unit \n" +
				"				FROM units AS punits \n" +
				"				JOIN parenthoods ON punits.unit=parent_unit \n" +
				"				JOIN units AS chunits ON child_unit=chunits.unit \n" +
				"				WHERE (chunits.main_chain_index > ? OR chunits.main_chain_index IS NULL) \n" +
				"					AND (chunits.latest_included_mc_index IS NULL OR chunits.latest_included_mc_index < punits.latest_included_mc_index)", DBParamsT{ last_main_chain_index })
 **/
			rcvr := db.UnitLIMCIsReceiver{}
			conn.MustQuery(/*
							"UPDATE units AS punits \n\
							JOIN parenthoods ON punits.unit=parent_unit \n\
							JOIN units AS chunits ON child_unit=chunits.unit \n\
							SET chunits.latest_included_mc_index=punits.latest_included_mc_index \n\
							WHERE (chunits.main_chain_index > ? OR chunits.main_chain_index IS NULL) \n\
								AND (chunits.latest_included_mc_index IS NULL OR chunits.latest_included_mc_index < punits.latest_included_mc_index)",
							[last_main_chain_index],
							function(result){
								(result.affectedRows > 0) ? propagateLIMCI() : checkAllLatestIncludedMcIndexesAreSet();
							}
							
			 */
			"SELECT punits.latest_included_mc_index, chunits.unit \n" +
				"FROM units AS punits \n" +
				"JOIN parenthoods ON punits.unit=parent_unit \n" +
				"JOIN units AS chunits ON child_unit=chunits.unit \n" +
				"WHERE (chunits.main_chain_index > ? OR chunits.main_chain_index IS NULL) \n" +
				"	AND (chunits.latest_included_mc_index IS NULL OR chunits.latest_included_mc_index < punits.latest_included_mc_index)", DBParamsT{ last_main_chain_index }, &rcvr)
			rows := rcvr.Rows
			// << flattened continuation for conn.query:246:3
			profiler.Stop("mc-limci-select-propagate")
			if len(rows) == 0 {
				checkAllLatestIncludedMcIndexesAreSet()
				return
			}
			profiler.Start()
			(func () ErrorT {
			  // :: inlined async.eachSeries:271:5
//			  for row := range rows {
			  for _, row := range rows {
//			    _err := (func (row rowT) ErrorT {
			    _err := (func (row db.UnitLIMCIRow) ErrorT {
//			    	assocDbLimcisByUnit[row.unit] = row.latest_included_mc_index
			    	assocDbLimcisByUnit[row.Unit] = row.Latest_included_mc_index
/**
			    	/* await * /
			    	conn.query_sync("UPDATE units SET latest_included_mc_index=? WHERE unit=?", DBParamsT{
			    		row.latest_included_mc_index,
			    		row.unit,
			    	})
 **/
			    	conn.MustExec("UPDATE units SET latest_included_mc_index=? WHERE unit=?", DBParamsT{
			    		row.Latest_included_mc_index,
			    		row.Unit,
			    	})
			    	// << flattened continuation for conn.query:275:7
			    	// :: flattened return for cb();
			    	// ** need 1 return(s) instead of 0
			    	return nil
			    	// >> flattened continuation for conn.query:275:7
			    })(row)
			    if _err != nil { return _err }
			  }
			  return nil
			})()
			// << flattened continuation for async.eachSeries:271:5
			profiler.Stop("mc-limci-update-propagate")
			propagateLIMCI()
			// >> flattened continuation for async.eachSeries:271:5
			// >> flattened continuation for conn.query:246:3
		}
		
//		loadUnitProps_sync = func (unit UnitT) PropsT {
		loadUnitProps_sync = func (unit UnitT) refXPropsT {
//			if storage.assocUnstableUnits[unit] {
			if _, _exists := storage.AssocUnstableUnits[unit]; _exists {
				// :: flattened return for return handleProps(storage.assocUnstableUnits[unit]);
//				return storage.assocUnstableUnits[unit]
				return storage.AssocUnstableUnits[unit]
			}
			// :: flattened return for handleProps(storage.readUnitProps(conn, unit));
//			return /* await */
//			storage.readUnitProps_sync(conn, unit)
//			return storage.ReadUnitProps_sync(conn, unit)
			return &XPropsT{
				PropsT: *storage.ReadUnitProps_sync(conn, unit),
			}
		}
		
		calcLIMCIs_sync = func ()  {
			arrFilledUnits := UnitsT{}
			(func () ErrorT {
			  // :: inlined async.eachOfSeries:294:3
//			  for props, unit := range assocChangedUnits {
			  for unit, props := range assocChangedUnits {
//			    _err := (func (props PropsT, unit UnitT) ErrorT {
			    _err := (func (props refXPropsT, unit UnitT) ErrorT {
//			    	max_limci := - 1
			    	max_limci := MCIndexT(-1)
			    	err := (func () ErrorT {
			    	  // :: inlined async.eachSeries:298:5
//			    	  for parent_unit := range props.parent_units {
			    	  for _, parent_unit := range props.Parent_units {
			    	    _err := (func (parent_unit UnitT) ErrorT {
			    	    	parent_props := /* await */
			    	    	loadUnitProps_sync(parent_unit)
			    	    	// << flattened continuation for loadUnitProps:301:7
//			    	    	if parent_props.is_on_main_chain {
			    	    	if parent_props.Is_on_main_chain != 0 {
//			    	    		props.latest_included_mc_index = parent_props.main_chain_index
			    	    		props.Latest_included_mc_index = parent_props.Main_chain_index
//			    	    		assocLimcisByUnit[unit] = props.latest_included_mc_index
			    	    		assocLimcisByUnit[unit] = props.Latest_included_mc_index
			    	    		arrFilledUnits = append(arrFilledUnits, unit)
			    	    		// :: flattened return for return cb2('done');
//			    	    		return "done"
			    	    		return errors.New("done")
			    	    	}
//			    	    	if parent_props.latest_included_mc_index == nil {
			    	    	if parent_props.Latest_included_mc_index.IsNull() {
			    	    		// :: flattened return for return cb2('parent limci not known yet');
			    	    		return errors.New("parent limci not known yet")
			    	    	}
//			    	    	if parent_props.latest_included_mc_index > max_limci {
			    	    	if parent_props.Latest_included_mc_index > max_limci {
//			    	    		max_limci = parent_props.latest_included_mc_index
			    	    		max_limci = parent_props.Latest_included_mc_index
			    	    	}
			    	    	// :: flattened return for cb2();
			    	    	// ** need 1 return(s) instead of 0
			    	    	return nil
			    	    	// >> flattened continuation for loadUnitProps:301:7
			    	    })(parent_unit)
			    	    if _err != nil { return _err }
			    	  }
			    	  return nil
			    	})()
			    	// << flattened continuation for async.eachSeries:298:5
//			    	if err {
			    	if err != nil {
			    		// :: flattened return for return cb();
			    		// ** need 1 return(s) instead of 0
			    		return nil
			    	}
			    	if max_limci < 0 {
//			    		_core.Throw("max limci < 0 for unit " + unit)
			    		_core.Throw("max limci < 0 for unit %s", unit)
			    	}
//			    	props.latest_included_mc_index = max_limci
			    	props.Latest_included_mc_index = max_limci
//			    	assocLimcisByUnit[unit] = props.latest_included_mc_index
			    	assocLimcisByUnit[unit] = props.Latest_included_mc_index
			    	arrFilledUnits = append(arrFilledUnits, unit)
			    	// :: flattened return for cb();
			    	// ** need 1 return(s) instead of 0
			    	return nil
			    	// >> flattened continuation for async.eachSeries:298:5
			    })(props, unit)
			    if _err != nil { return _err }
			  }
			  return nil
			})()
			// << flattened continuation for async.forEachOfSeries:294:3
			// .. not flattening for Array.forEach
//			for unit, _ := range arrFilledUnits {
			for _, unit := range arrFilledUnits {
				delete(assocChangedUnits, unit)
			}
//			if len(Object.keys(assocChangedUnits)) > 0 {
			if len(assocChangedUnits) > 0 {
				/* await */
				calcLIMCIs_sync()
				// << flattened continuation for calcLIMCIs:332:6
				// :: flattened return for onUpdated();
				return 
				// >> flattened continuation for calcLIMCIs:332:6
			} else {
				// :: flattened return for onUpdated();
				return 
			}
			// >> flattened continuation for async.forEachOfSeries:294:3
		}
		
//		console.Log("updateLatestIncludedMcIndex " + last_main_chain_index)
		console.Log("updateLatestIncludedMcIndex %d", last_main_chain_index)
		profiler.Start()
/**
		assocChangedUnits := [*ObjectExpression*]
		assocLimcisByUnit := [*ObjectExpression*]
		assocDbLimcisByUnit := [*ObjectExpression*]
 **/

//		for unit := range storage.assocUnstableUnits {
		for unit := range storage.AssocUnstableUnits {
//			o := storage.assocUnstableUnits[unit]
			o := storage.AssocUnstableUnits[unit]
//			if o.main_chain_index > last_main_chain_index || o.main_chain_index == nil {
			if o.Main_chain_index.IsNull() || o.Main_chain_index > last_main_chain_index {
//				o.latest_included_mc_index = nil
				o.Latest_included_mc_index = MCIndexT_Null
				assocChangedUnits[unit] = o
			}
		}
		/* await */
		calcLIMCIs_sync()
		// << flattened continuation for calcLIMCIs:351:2
/**
		res := /* await * /
		conn.query_sync("UPDATE units SET latest_included_mc_index=NULL WHERE main_chain_index>? OR main_chain_index IS NULL", DBParamsT{ last_main_chain_index })
 **/
		res := conn.MustExec("UPDATE units SET latest_included_mc_index=NULL WHERE main_chain_index>? OR main_chain_index IS NULL", DBParamsT{ last_main_chain_index })
		// << flattened continuation for conn.query:352:3
//		console.Log("update LIMCI=NULL done, matched rows: " + res.affectedRows)
		console.Log("update LIMCI=NULL done, matched rows: %d", res.AffectedRows)
		profiler.Stop("mc-limci-set-null")
		profiler.Start()
/**
		rows := /* await * /
		conn.query_sync(// if these units have other parents, they cannot include later MC units (otherwise, the parents would've been redundant).
		// the 2nd condition in WHERE is the same that was used 1 query ago to NULL limcis.
		
		// I had to rewrite this single query because sqlite doesn't support JOINs in UPDATEs
		/*
							"UPDATE units AS punits \n\
							JOIN parenthoods ON punits.unit=parent_unit \n\
							JOIN units AS chunits ON child_unit=chunits.unit \n\
							SET chunits.latest_included_mc_index=punits.main_chain_index \n\
							WHERE punits.is_on_main_chain=1 \n\
								AND (chunits.main_chain_index > ? OR chunits.main_chain_index IS NULL) \n\
								AND chunits.latest_included_mc_index IS NULL", 
							[last_main_chain_index],
							function(result){
								if (result.affectedRows === 0 && bRebuiltMc)
									throw "no latest_included_mc_index updated";
								propagateLIMCI();
							}
							
		 * /
		"SELECT chunits.unit, punits.main_chain_index \n" +
			"					FROM units AS punits \n" +
			"					JOIN parenthoods ON punits.unit=parent_unit \n" +
			"					JOIN units AS chunits ON child_unit=chunits.unit \n" +
			"					WHERE punits.is_on_main_chain=1 \n" +
			"						AND (chunits.main_chain_index > ? OR chunits.main_chain_index IS NULL) \n" +
			"						AND chunits.latest_included_mc_index IS NULL", DBParamsT{ last_main_chain_index })
 **/
		rcvr := db.UnitMCIsReceiver{}
		conn.MustQuery(// if these units have other parents, they cannot include later MC units (otherwise, the parents would've been redundant).
		// the 2nd condition in WHERE is the same that was used 1 query ago to NULL limcis.
		
		// I had to rewrite this single query because sqlite doesn't support JOINs in UPDATEs
		/*
							"UPDATE units AS punits \n\
							JOIN parenthoods ON punits.unit=parent_unit \n\
							JOIN units AS chunits ON child_unit=chunits.unit \n\
							SET chunits.latest_included_mc_index=punits.main_chain_index \n\
							WHERE punits.is_on_main_chain=1 \n\
								AND (chunits.main_chain_index > ? OR chunits.main_chain_index IS NULL) \n\
								AND chunits.latest_included_mc_index IS NULL", 
							[last_main_chain_index],
							function(result){
								if (result.affectedRows === 0 && bRebuiltMc)
									throw "no latest_included_mc_index updated";
								propagateLIMCI();
							}
							
		 */
		"SELECT chunits.unit, punits.main_chain_index \n" +
			"					FROM units AS punits \n" +
			"					JOIN parenthoods ON punits.unit=parent_unit \n" +
			"					JOIN units AS chunits ON child_unit=chunits.unit \n" +
			"					WHERE punits.is_on_main_chain=1 \n" +
			"						AND (chunits.main_chain_index > ? OR chunits.main_chain_index IS NULL) \n" +
			"						AND chunits.latest_included_mc_index IS NULL", DBParamsT{ last_main_chain_index }, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:356:4
//		console.Log(len(rows) + " rows")
		console.Log("%d rows", len(rows))
		profiler.Stop("mc-limci-select-initial")
		profiler.Start()
		if len(rows) == 0 && bRebuiltMc {
//			_core.Throw("no latest_included_mc_index updated, last_mci=" + last_main_chain_index + ", affected=" + res.affectedRows)
			_core.Throw("no latest_included_mc_index updated, last_mci=%d, affected=%d", last_main_chain_index, res.AffectedRows)
		}
		(func () ErrorT {
		  // :: inlined async.eachSeries:390:6
//		  for row := range rows {
		  for _, row := range rows {
//		    _err := (func (row rowT) ErrorT {
		    _err := (func (row db.UnitMCIRow) ErrorT {
//		    	console.Log(row.main_chain_index, row.unit)
		    	console.Log("%d %s", row.Main_chain_index, row.Unit)
//		    	assocDbLimcisByUnit[row.unit] = row.main_chain_index
		    	assocDbLimcisByUnit[row.Unit] = row.Main_chain_index
/**
		    	/* await * /
		    	conn.query_sync("UPDATE units SET latest_included_mc_index=? WHERE unit=?", DBParamsT{
		    		row.main_chain_index,
		    		row.unit,
		    	})
 **/
		    	conn.MustExec("UPDATE units SET latest_included_mc_index=? WHERE unit=?", DBParamsT{
		    		row.Main_chain_index,
		    		row.Unit,
		    	})
		    	// << flattened continuation for conn.query:395:8
		    	// :: flattened return for cb();
		    	// ** need 1 return(s) instead of 0
		    	return nil
		    	// >> flattened continuation for conn.query:395:8
		    })(row)
		    if _err != nil { return _err }
		  }
		  return nil
		})()
		// << flattened continuation for async.eachSeries:390:6
		profiler.Stop("mc-limci-update-initial")
		propagateLIMCI()
		// >> flattened continuation for async.eachSeries:390:6
		// >> flattened continuation for conn.query:356:4
		// >> flattened continuation for conn.query:352:3
		// >> flattened continuation for calcLIMCIs:351:2
	}
	
	readLastStableMcUnit_sync = func () UnitT {
/**
		rows := /* await * /
		conn.query_sync("SELECT unit FROM units WHERE is_on_main_chain=1 AND is_stable=1 ORDER BY main_chain_index DESC LIMIT 1")
 **/
		rcvr := db.UnitsReceiver{}
		conn.MustQuery("SELECT unit FROM units WHERE is_on_main_chain=1 AND is_stable=1 ORDER BY main_chain_index DESC LIMIT 1", DBParamsT{}, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:409:2
		if len(rows) == 0 {
			_core.Throw("no units on stable MC?")
		}
		// :: flattened return for handleLastStableMcUnit(rows[0].unit);
//		return rows[0].unit
		return rows[0].Unit
		// >> flattened continuation for conn.query:409:2
	}
	
	
	updateStableMcFlag = func ()  {
		var(
			advanceLastStableMcUnitAndTryNext func () 
		)
		
		console.Log("updateStableMcFlag")
		profiler.Start()
		last_stable_mc_unit := /* await */
		readLastStableMcUnit_sync()
		// << flattened continuation for readLastStableMcUnit:420:2
//		console.Log("last stable mc unit " + last_stable_mc_unit)
		console.Log("last stable mc unit %s", last_stable_mc_unit)
		arrWitnesses := /* await */
//		storage.readWitnesses_sync(conn, last_stable_mc_unit)
		storage.ReadWitnesses_sync(conn, last_stable_mc_unit)
		// << flattened continuation for storage.readWitnesses:422:3
/**
		rows := /* await * /
		conn.query_sync("SELECT unit, is_on_main_chain, main_chain_index, level FROM units WHERE best_parent_unit=?", DBParamsT{ last_stable_mc_unit })
 **/
		rcvr_0 := db.UnitPropsReceiver{}
		conn.MustQuery("SELECT unit, is_on_main_chain, main_chain_index, level FROM units WHERE best_parent_unit=?", DBParamsT{ last_stable_mc_unit }, &rcvr_0)
		rows := rcvr_0.Rows
		// << flattened continuation for conn.query:423:4
		if len(rows) == 0 {
			//if (isGenesisUnit(last_stable_mc_unit))
			//    return finish();
//			_core.Throw("no best children of last stable MC unit " + last_stable_mc_unit + "?")
			_core.Throw("no best children of last stable MC unit %s?", last_stable_mc_unit)
		}
//		arrMcRows := make(arrMcRowsT, 0, len(rows))
		arrMcRows := make([]db.UnitPropsRow, 0, len(rows))
//		for row, _ := range rows {
		for _, row := range rows {
//			if ! (row.is_on_main_chain == 1) { continue }
			if ! (row.Is_on_main_chain == 1) { continue }
			arrMcRows = append(arrMcRows, row)
		}
		// only one element
//		arrAltRows := make(arrAltRowsT, 0, len(rows))
		arrAltRows := make([]db.UnitPropsRow, 0, len(rows))
//		for row, _ := range rows {
		for _, row := range rows {
//			if ! (row.is_on_main_chain == 0) { continue }
			if ! (row.Is_on_main_chain == 0) { continue }
			arrAltRows = append(arrAltRows, row)
		}
		if len(arrMcRows) != 1 {
			_core.Throw("not a single MC child?")
		}
//		first_unstable_mc_unit := arrMcRows[0].unit
//[uu]		first_unstable_mc_unit := arrMcRows[0].Unit
//		first_unstable_mc_index := arrMcRows[0].main_chain_index
		first_unstable_mc_index := arrMcRows[0].Main_chain_index
//		first_unstable_mc_level := arrMcRows[0].level
		first_unstable_mc_level := arrMcRows[0].Level
		arrAltBranchRootUnits := make(UnitsT, len(arrAltRows), len(arrAltRows))
//		for row, _k := range arrAltRows {
		for _k, row := range arrAltRows {
//			arrAltBranchRootUnits[_k] := row.unit
			arrAltBranchRootUnits[_k] = row.Unit
		}
		
		advanceLastStableMcUnitAndTryNext = func ()  {
			profiler.Stop("mc-stableFlag")
			/* await */
			markMcIndexStable_sync(conn, first_unstable_mc_index)
			// << flattened continuation for markMcIndexStable:440:6
			updateStableMcFlag()
			// >> flattened continuation for markMcIndexStable:440:6
		}
/**
		wl_rows := /* await * /
		conn.query_sync("SELECT witnessed_level FROM units WHERE is_free=1 AND is_on_main_chain=1")
 **/
		rcvr_1 := db.WitnessedLevelsReceiver{}
		conn.MustQuery("SELECT witnessed_level FROM units WHERE is_free=1 AND is_on_main_chain=1", DBParamsT{}, &rcvr_1)
		wl_rows := rcvr_1.Rows
		// << flattened continuation for conn.query:443:5
		if len(wl_rows) != 1 {
			_core.Throw("not a single mc wl")
		}
		// this is the level when we colect 7 witnesses if walking up the MC from its end
//		mc_end_witnessed_level := wl_rows[0].witnessed_level
		mc_end_witnessed_level := wl_rows[0].Witnessed_level
/**
		min_wl_rows := /* await * /
		conn.query_sync(// among these 7 witnesses, find min wl
		"SELECT MIN(witnessed_level) AS min_mc_wl FROM units LEFT JOIN unit_authors USING(unit) \n" +
			"							WHERE is_on_main_chain=1 AND level>=? AND address IN(?)", // _left_ join enforces the best query plan in sqlite
		DBParamsT{
			mc_end_witnessed_level,
			arrWitnesses,
		})
 **/
		queryParams := DBParamsT{
			mc_end_witnessed_level,
		}
		wsSql := queryParams.AddAddresses(arrWitnesses)
		rcvr_2 := db.MinMCWLsReceiver{}
		conn.MustQuery(// among these 7 witnesses, find min wl
			"SELECT MIN(witnessed_level) AS min_mc_wl FROM units LEFT JOIN unit_authors USING(unit) \n" +
			"WHERE is_on_main_chain=1 AND level>=? AND address IN(" + wsSql + ")", // _left_ join enforces the best query plan in sqlite
			queryParams, &rcvr_2)
		min_wl_rows := rcvr_2.Rows
		// << flattened continuation for conn.query:448:6
		if len(min_wl_rows) != 1 {
			_core.Throw("not a single min mc wl")
		}
//		min_mc_wl := min_wl_rows[0].min_mc_wl
		min_mc_wl := min_wl_rows[0].Min_mc_wl
		if len(arrAltBranchRootUnits) == 0 {
			// no alt branches
			if min_mc_wl >= first_unstable_mc_level {
				advanceLastStableMcUnitAndTryNext()
				return
			}
			finish()
			return
		}
		arrAltBestChildren := /* await */
		createListOfBestChildren_sync(arrAltBranchRootUnits)
		// << flattened continuation for createListOfBestChildren:477:8
/**
		max_alt_rows := /* await * /
		conn.query_sync("SELECT MAX(units.level) AS max_alt_level \n" +
			"										FROM units \n" +
			"										LEFT JOIN parenthoods ON units.unit=child_unit \n" +
			"										LEFT JOIN units AS punits ON parent_unit=punits.unit AND punits.witnessed_level >= units.witnessed_level \n" +
			"										WHERE units.unit IN(?) AND punits.unit IS NULL AND ( \n" +
			"											SELECT COUNT(*) \n" +
			"											FROM unit_witnesses \n" +
			"											WHERE unit_witnesses.unit IN(units.unit, units.witness_list_unit) AND unit_witnesses.address IN(?) \n" +
			"										)>=?", DBParamsT{
			arrAltBestChildren,
			arrWitnesses,
			constants.COUNT_WITNESSES - constants.MAX_WITNESS_LIST_MUTATIONS,
		})
 **/
		rcvr_3 := db.MaxAltLevelsReceiver{}
		{{
		queryParams := DBParamsT{}
		abcsSql := queryParams.AddUnits(arrAltBestChildren)
		wisSql := queryParams.AddAddresses(arrWitnesses)
		queryParams = append(queryParams,
			constants.COUNT_WITNESSES - constants.MAX_WITNESS_LIST_MUTATIONS)
		conn.MustQuery("SELECT MAX(units.level) AS max_alt_level \n" +
			"FROM units \n" +
			"LEFT JOIN parenthoods ON units.unit=child_unit \n" +
			"LEFT JOIN units AS punits ON parent_unit=punits.unit AND punits.witnessed_level >= units.witnessed_level \n" +
			"WHERE units.unit IN(" + abcsSql + ") AND punits.unit IS NULL AND ( \n" +
			"	SELECT COUNT(*) \n" +
			"	FROM unit_witnesses \n" +
			"	WHERE unit_witnesses.unit IN(units.unit, units.witness_list_unit) AND unit_witnesses.address IN(" + wisSql + ") \n" +
			")>=?", queryParams, &rcvr_3)
		}}
		max_alt_rows := rcvr_3.Rows
		// << flattened continuation for conn.query:481:9
		if len(max_alt_rows) != 1 {
			_core.Throw("not a single max alt level")
		}
//		max_alt_level := max_alt_rows[0].max_alt_level
		max_alt_level := max_alt_rows[0].Max_alt_level
		if min_mc_wl > max_alt_level {
			advanceLastStableMcUnitAndTryNext()
			return
		}
		finish()
		return
		// >> flattened continuation for conn.query:481:9
		// >> flattened continuation for createListOfBestChildren:477:8
		// >> flattened continuation for conn.query:448:6
		// >> flattened continuation for conn.query:443:5
		// >> flattened continuation for conn.query:423:4
		// >> flattened continuation for storage.readWitnesses:422:3
		// >> flattened continuation for readLastStableMcUnit:420:2
	}
	
	// also includes arrParentUnits
	createListOfBestChildren_sync = func (arrParentUnits UnitsT) UnitsT {
		var(
			goDownAndCollectBestChildren_1_sync func (arrStartUnits UnitsT) ErrorT
		)
		
		if len(arrParentUnits) == 0 {
			// :: flattened return for return handleBestChildrenList([]);
			return UnitsT{}
		}
//		arrBestChildren := arrParentUnits.slice()
		arrBestChildren := arrParentUnits[:]
		
		goDownAndCollectBestChildren_1_sync = func (arrStartUnits UnitsT) ErrorT {
/**
			rows := /* await * /
			conn.query_sync("SELECT unit, is_free FROM units WHERE best_parent_unit IN(?)", DBParamsT{ arrStartUnits })
 **/
			rcvr := db.UnitPropsReceiver{}
			queryParams := DBParamsT{}
			susSql := queryParams.AddUnits(arrStartUnits)
			conn.MustQuery("SELECT unit, is_free FROM units WHERE best_parent_unit IN(" + susSql + ")", queryParams, &rcvr)
			rows := rcvr.Rows
			// << flattened continuation for conn.query:517:3
			if len(rows) == 0 {
				// :: flattened return for return cb();
				// ** need 1 return(s) instead of 0
				return nil
			}
			// :: flattened return for cb(async.eachSeries(rows, function (row) {
			return (func () ErrorT {
			  // :: inlined async.eachSeries:521:4
//			  for row := range rows {
			  for _, row := range rows {
//			    _err := (func (row rowT) ErrorT {
			    _err := (func (row db.UnitPropsRow) ErrorT {
//			    	arrBestChildren = append(arrBestChildren, row.unit)
			    	arrBestChildren = append(arrBestChildren, row.Unit)
//			    	if row.is_free == 1 {
			    	if row.Is_free == 1 {
			    		// :: flattened return for cb2();
			    		// ** need 1 return(s) instead of 0
			    		return nil
			    	} else {
			    		/* await */
//			    		goDownAndCollectBestChildren_1_sync(UnitsT{ row.unit })
			    		goDownAndCollectBestChildren_1_sync(UnitsT{ row.Unit })
			    		// << flattened continuation for goDownAndCollectBestChildren_1:528:7
			    		// :: flattened return for cb2();
			    		// ** need 1 return(s) instead of 0
			    		return nil
			    		// >> flattened continuation for goDownAndCollectBestChildren_1:528:7
			    	}
			    })(row)
			    if _err != nil { return _err }
			  }
			  return nil
			})()
			// >> flattened continuation for conn.query:517:3
		}
		
		/* await */
		goDownAndCollectBestChildren_1_sync(arrParentUnits)
		// << flattened continuation for goDownAndCollectBestChildren_1:535:2
		// :: flattened return for handleBestChildrenList(arrBestChildren);
		return arrBestChildren
		// >> flattened continuation for goDownAndCollectBestChildren_1:535:2
	}
	
	
	
	finish = func ()  {
		profiler.Stop("mc-stableFlag")
		console.Log("done updating MC\n")
//		if onDone {
		if true {
			// :: flattened return for onDone();
			return 
		}
	}
	
	
	console.Log("will update MC")
	
	/*if (from_unit === null && arrRetreatingUnits.indexOf(last_added_unit) >= 0){
			conn.query("UPDATE units SET is_on_main_chain=1, main_chain_index=NULL WHERE unit=?", [last_added_unit], function(){
				goUpFromUnit(last_added_unit);
			});
		}
		else
	 */
	goUpFromUnit(from_unit)
}





/*
// blank line: []
// climbs up along best parent links up, returns list of units encountered with level >= min_level
function createListOfPrivateMcUnits(start_unit, min_level, handleList){
	var arrUnits = [];
// blank line: [	]
	function goUp_2(unit){
		conn.query(
			"SELECT best_parent_unit, level FROM units WHERE unit=?", [unit],
			function(rows){
				if (rows.length !== 1)
					throw "createListOfPrivateMcUnits: not 1 row";
				var row = rows[0];
				if (row.level < min_level) 
					return handleList(arrUnits);
				arrUnits.push(unit);
				goUp_2(row.best_parent_unit);
			}
		);
	}
// blank line: [	]
	goUp_2(start_unit);
}
// blank line: []

 */

//func determineIfStableInLaterUnits_sync(conn DBConnT, earlier_unit UnitT, arrLaterUnits UnitsT) bool {
func determineIfStableInLaterUnits_sync(conn refDBConnT, earlier_unit UnitT, arrLaterUnits UnitsT) bool {
	var(
//		findMinMcWitnessedLevel_sync func () 
		findMinMcWitnessedLevel_sync func () LevelT
		determineIfHasAltBranches_sync func () bool
		createListOfBestChildrenIncludedByLaterUnits_sync func (arrAltBranchRootUnits UnitsT) UnitsT
	)
	
//	if storage.isGenesisUnit(earlier_unit) {
	if storage.IsGenesisUnit(earlier_unit) {
		// :: flattened return for return handleResult(true);
		return true
	}
	// hack to workaround past validation error
//	if earlier_unit == "LGFzduLJNQNzEqJqUXdkXr58wDYx77V8WurDF3+GIws=" && arrLaterUnits.join(",") == "6O4t3j8kW0/Lo7n2nuS8ITDv2UbOhlL9fF1M6j/PrJ4=" {
	if earlier_unit == "LGFzduLJNQNzEqJqUXdkXr58wDYx77V8WurDF3+GIws=" && arrLaterUnits.Join(",") == "6O4t3j8kW0/Lo7n2nuS8ITDv2UbOhlL9fF1M6j/PrJ4=" {
		// :: flattened return for return handleResult(true);
		return true
	}
//	( objEarlierUnitProps, arrLaterUnitProps ) := /* await */
	objEarlierUnitProps, arrLaterUnitProps := /* await */
//	storage.readPropsOfUnits_sync(conn, earlier_unit, arrLaterUnits)
	storage.ReadPropsOfUnits_sync(conn, earlier_unit, arrLaterUnits)
	// << flattened continuation for storage.readPropsOfUnits:598:1
//	if objEarlierUnitProps.is_free == 1 {
	if objEarlierUnitProps.Is_free == 1 {
		// :: flattened return for return handleResult(false);
		return false
	}
//	arr_later_limcis := make(MCIndexesT, len(arrLaterUnitProps), len(arrLaterUnitProps))
	arr_later_limcis := make(MCIndexesT, len(arrLaterUnitProps), len(arrLaterUnitProps))
//	for objLaterUnitProps, _k := range arrLaterUnitProps {
	for _k, objLaterUnitProps := range arrLaterUnitProps {
//		arr_later_limcis[_k] := objLaterUnitProps.latest_included_mc_index
		arr_later_limcis[_k] = objLaterUnitProps.Latest_included_mc_index
	}
//	max_later_limci := Math.max.apply(nil, arr_later_limcis)
	max_later_limci := MCIndexT(-1)
	for _, later_limci := range arr_later_limcis {
		if max_later_limci < later_limci {
			max_later_limci = later_limci
		}
	}
//	( best_parent_unit, arrWitnesses ) := /* await */
	best_parent_unit, arrWitnesses := /* await */
	readBestParentAndItsWitnesses_sync(conn, earlier_unit)
	// << flattened continuation for readBestParentAndItsWitnesses:603:2
/**
	rows := /* await * /
	conn.query_sync("SELECT unit, is_on_main_chain, main_chain_index, level FROM units WHERE best_parent_unit=?", DBParamsT{ best_parent_unit })
 **/
//--	rows := *(func() *[]db.UnitPropsRow {

	rcvr := db.UnitPropsReceiver{}
	conn.MustQuery("SELECT unit, is_on_main_chain, main_chain_index, level FROM units WHERE best_parent_unit=?", DBParamsT{ best_parent_unit }, &rcvr)
	rows := rcvr.Rows
//--	return &rcvr.Rows

//--	}())
	// << flattened continuation for conn.query:604:3
	if len(rows) == 0 {
//		_core.Throw("no best children of " + best_parent_unit + "?")
		_core.Throw("no best children of %s?", best_parent_unit)
	}
	arrMcRows := make([]db.UnitPropsRow, 0, len(rows))
//	for row, _ := range rows {
	for _, row := range rows {
//		if ! (row.is_on_main_chain == 1) { continue }
		if ! (row.Is_on_main_chain == 1) { continue }
		arrMcRows = append(arrMcRows, row)
	}
	// only one element
	arrAltRows := make([]db.UnitPropsRow, 0, len(rows))
//	for row, _ := range rows {
	for _, row := range rows {
//		if ! (row.is_on_main_chain == 0) { continue }
		if ! (row.Is_on_main_chain == 0) { continue }
		arrAltRows = append(arrAltRows, row)
	}
	if len(arrMcRows) != 1 {
		_core.Throw("not a single MC child?")
	}
//	first_unstable_mc_unit := arrMcRows[0].unit
	first_unstable_mc_unit := arrMcRows[0].Unit
	if first_unstable_mc_unit != earlier_unit {
		_core.Throw("first unstable MC unit is not our input unit")
	}
//	first_unstable_mc_index := arrMcRows[0].main_chain_index
//[uu]	first_unstable_mc_index := arrMcRows[0].Main_chain_index
//	first_unstable_mc_level := arrMcRows[0].level
	first_unstable_mc_level := arrMcRows[0].Level
	arrAltBranchRootUnits := make(UnitsT, len(arrAltRows), len(arrAltRows))
//	for row, _k := range arrAltRows {
	for _k, row := range arrAltRows {
//		arrAltBranchRootUnits[_k] := row.unit
		arrAltBranchRootUnits[_k] = row.Unit
	}
	//console.Log("first_unstable_mc_index", first_unstable_mc_index);
	//console.Log("first_unstable_mc_level", first_unstable_mc_level);
	//console.Log("alt", arrAltBranchRootUnits);
	
//	findMinMcWitnessedLevel_sync = func ()  {
	findMinMcWitnessedLevel_sync = func () LevelT {
		var(
//			goUp_3 func (start_unit UnitT) 
			goUp_3 func (start_unit UnitT) LevelT
		)
		
//		min_mc_wl := Number.MAX_VALUE
		min_mc_wl := LevelT(math.MaxInt64)
		count := 0
		
//		goUp_3 = func (start_unit UnitT)  {
		goUp_3 = func (start_unit UnitT) LevelT {
/**
			rows := /* await * /
			conn.query_sync("SELECT best_parent_unit, witnessed_level, \n" +
				"								(SELECT COUNT(*) FROM unit_authors WHERE unit_authors.unit=units.unit AND address IN(?)) AS count \n" +
				"							FROM units WHERE unit=?", DBParamsT{
				arrWitnesses,
				start_unit,
			})
 **/
			rcvr := db.BestParentUnitCountsReceiver{}
			conn.MustQuery("SELECT best_parent_unit, witnessed_level, \n" +
				"(SELECT COUNT(*) FROM unit_authors WHERE unit_authors.unit=units.unit AND address IN(?)) AS count \n" +
				"FROM units WHERE unit=?", DBParamsT{
				arrWitnesses,
				start_unit,
			}, &rcvr)
			rows := rcvr.Rows
			// << flattened continuation for conn.query:626:6
			if len(rows) != 1 {
				_core.Throw("findMinMcWitnessedLevel: not 1 row")
			}
			row := rows[0]
//			if row.count > 0 && row.Witnessed_level < min_mc_wl {
			if row.Count > 0 && row.Witnessed_level < min_mc_wl {
//				min_mc_wl = row.witnessed_level
				min_mc_wl = row.Witnessed_level
			}
//			count += row.count
			count += row.Count
//			(count < constants.MAJORITY_OF_WITNESSES ? goUp_3(row.best_parent_unit): handleMinMcWl(min_mc_wl))
			if count >= constants.MAJORITY_OF_WITNESSES {
//				handleMinMcWl(min_mc_wl)
				return min_mc_wl
			}
//			return goUp_3(row.best_parent_unit)
			return goUp_3(row.Best_parent_unit)
			// >> flattened continuation for conn.query:626:6
		}
/**
		rows := /* await * /
		conn.query_sync("SELECT witnessed_level, best_parent_unit, \n" +
			"							(SELECT COUNT(*) FROM unit_authors WHERE unit_authors.unit=units.unit AND address IN(?)) AS count \n" +
			"						FROM units \n" +
			"						WHERE unit IN(?) \n" +
			"						ORDER BY witnessed_level DESC, \n" +
			"							level-witnessed_level ASC, \n" +
			"							unit ASC \n" +
			"						LIMIT 1", DBParamsT{
			arrWitnesses,
			arrLaterUnits,
		})
 **/
		rcvr := db.BestParentUnitCountsReceiver{}
		conn.MustQuery("SELECT best_parent_unit, witnessed_level, \n" +
			"	(SELECT COUNT(*) FROM unit_authors WHERE unit_authors.unit=units.unit AND address IN(?)) AS count \n" +
			"FROM units \n" +
			"WHERE unit IN(?) \n" +
			"ORDER BY witnessed_level DESC, \n" +
			"	level-witnessed_level ASC, \n" +
			"	unit ASC \n" +
			"LIMIT 1", DBParamsT{
			arrWitnesses,
			arrLaterUnits,
		}, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:642:5
		row := rows[0]
//		if row.count > 0 {
		if row.Count > 0 {
//			min_mc_wl = row.witnessed_level
			min_mc_wl = row.Witnessed_level
		}
//		count += row.count
		count += row.Count
//		goUp_3(row.best_parent_unit)
		return goUp_3(row.Best_parent_unit)
		// >> flattened continuation for conn.query:642:5
	}
	
	determineIfHasAltBranches_sync = func () bool {
		if len(arrAltBranchRootUnits) == 0 {
			// :: flattened return for return handleHasAltBranchesResult(false);
			return false
		}
		err := (func () ErrorT {
		  // :: inlined async.eachSeries:666:5
//		  for alt_root_unit := range arrAltBranchRootUnits {
		  for _, alt_root_unit := range arrAltBranchRootUnits {
		    _err := (func (alt_root_unit UnitT) ErrorT {
		    	bIncluded := /* await */
//		    	graph.determineIfIncludedOrEqual_sync(conn, alt_root_unit, arrLaterUnits)
		    	graph.DetermineIfIncludedOrEqual_sync(conn, alt_root_unit, arrLaterUnits)
		    	// << flattened continuation for graph.determineIfIncludedOrEqual:669:7
//		    	(bIncluded ? cb("included"): cb())
		    	if bIncluded { 
				return errors.New("included")
			}
			return nil
		    	// >> flattened continuation for graph.determineIfIncludedOrEqual:669:7
		    })(alt_root_unit)
		    if _err != nil { return _err }
		  }
		  return nil
		})()
		// << flattened continuation for async.eachSeries:666:5
		// :: flattened return for handleHasAltBranchesResult(err !== undefined);
		return err != nil
		// >> flattened continuation for async.eachSeries:666:5
	}
	
	// also includes arrAltBranchRootUnits
	createListOfBestChildrenIncludedByLaterUnits_sync = func (arrAltBranchRootUnits UnitsT) UnitsT {
		var(
			goDownAndCollectBestChildren_2_sync func (arrStartUnits UnitsT) ErrorT
			filterAltBranchRootUnits_sync func () 
		)
		
		if len(arrAltBranchRootUnits) == 0 {
			// :: flattened return for return handleBestChildrenList([]);
			return UnitsT{}
		}
		arrBestChildren := UnitsT{}
		
		goDownAndCollectBestChildren_2_sync = func (arrStartUnits UnitsT) ErrorT {
/**
			rows := /* await * /
			conn.query_sync("SELECT unit, is_free, main_chain_index FROM units WHERE best_parent_unit IN(?)", DBParamsT{ arrStartUnits })
 **/
			rcvr := db.UnitPropsReceiver{}
			conn.MustQuery("SELECT unit, is_free, main_chain_index FROM units WHERE best_parent_unit IN(?)", DBParamsT{ arrStartUnits }, &rcvr)
			rows := rcvr.Rows
			// << flattened continuation for conn.query:686:6
			if len(rows) == 0 {
				// :: flattened return for return cb();
				// ** need 1 return(s) instead of 0
				return nil
			}
			// :: flattened return for cb(async.eachSeries(rows, function (row) {
			return (func () ErrorT {
			  // :: inlined async.eachSeries:689:7
//			  for row := range rows {
			  for _, row := range rows {
//			    _err := (func (row rowT) ErrorT {
			    _err := (func (row db.UnitPropsRow) ErrorT {
			    	var(
			    		addUnit_1 func () ErrorT
			    	)
			    	
			    	addUnit_1 = func () ErrorT {
//			    		arrBestChildren = append(arrBestChildren, row.unit)
			    		arrBestChildren = append(arrBestChildren, row.Unit)
//			    		if row.is_free == 1 {
			    		if row.Is_free == 1 {
			    			// :: flattened return for cb2();
			    			// ** need 1 return(s) instead of 0
			    			return nil
			    		} else {
			    			/* await */
//			    			goDownAndCollectBestChildren_2_sync(UnitsT{ row.unit })
			    			goDownAndCollectBestChildren_2_sync(UnitsT{ row.Unit })
			    			// << flattened continuation for goDownAndCollectBestChildren_2:698:11
			    			// :: flattened return for cb2();
			    			// ** need 1 return(s) instead of 0
			    			return nil
			    			// >> flattened continuation for goDownAndCollectBestChildren_2:698:11
			    		}
			    	}
			    	
//			    	if row.main_chain_index != nil && row.main_chain_index <= max_later_limci {
			    	if ! row.Main_chain_index.IsNull() &&  row.Main_chain_index <= max_later_limci {
//			    		addUnit_1()
			    		return addUnit_1()
			    	} else {
			    		bIncluded := /* await */
//			    		graph.determineIfIncludedOrEqual_sync(conn, row.unit, arrLaterUnits)
			    		graph.DetermineIfIncludedOrEqual_sync(conn, row.Unit, arrLaterUnits)
			    		// << flattened continuation for graph.determineIfIncludedOrEqual:704:10
//			    		(bIncluded ? addUnit_1(): cb2())
			    		if ! bIncluded {
						//cb2()
						return nil
					}
					return addUnit_1()
			    		// >> flattened continuation for graph.determineIfIncludedOrEqual:704:10
			    	}
			    })(row)
			    if _err != nil { return _err }
			  }
			  return nil
			})()
			// >> flattened continuation for conn.query:686:6
		}
		
		// leaves only those roots that are included by later units
		filterAltBranchRootUnits_sync = func ()  {
			arrFilteredAltBranchRootUnits := UnitsT{}
/**
			rows := /* await * /
			conn.query_sync("SELECT unit, is_free, main_chain_index FROM units WHERE unit IN(?)", DBParamsT{ arrAltBranchRootUnits })
 **/
			rcvr := db.UnitPropsReceiver{}
			conn.MustQuery("SELECT unit, is_free, main_chain_index FROM units WHERE unit IN(?)", DBParamsT{ arrAltBranchRootUnits }, &rcvr)
			rows := rcvr.Rows
			// << flattened continuation for conn.query:716:6
			if len(rows) == 0 {
				_core.Throw("no alt branch root units?")
			}
			(func () ErrorT {
			  // :: inlined async.eachSeries:719:7
//			  for row := range rows {
			  for _, row := range rows {
//			    _err := (func (row rowT) ErrorT {
			    _err := (func (row db.UnitPropsRow) ErrorT {
			    	var(
			    		addUnit_2 func () ErrorT
			    	)
			    	
			    	addUnit_2 = func () ErrorT {
//			    		arrBestChildren = append(arrBestChildren, row.unit)
			    		arrBestChildren = append(arrBestChildren, row.Unit)
			    		//	if (row.is_free === 0) // seems no reason to exclude
//			    		arrFilteredAltBranchRootUnits = append(arrFilteredAltBranchRootUnits, row.unit)
			    		arrFilteredAltBranchRootUnits = append(arrFilteredAltBranchRootUnits, row.Unit)
			    		// :: flattened return for cb2();
			    		// ** need 1 return(s) instead of 0
			    		return nil
			    	}
			    	
//			    	if row.main_chain_index != nil && row.main_chain_index <= max_later_limci {
			    	if ! row.Main_chain_index.IsNull() && row.Main_chain_index <= max_later_limci {
//			    		addUnit_2()
			    		return addUnit_2()
			    	} else {
			    		bIncluded := /* await */
//			    		graph.determineIfIncludedOrEqual_sync(conn, row.unit, arrLaterUnits)
			    		graph.DetermineIfIncludedOrEqual_sync(conn, row.Unit, arrLaterUnits)
			    		// << flattened continuation for graph.determineIfIncludedOrEqual:733:10
//			    		(bIncluded ? addUnit_2(): cb2())
			    		if ! bIncluded {
						//cb2()
						return nil
					}
					return addUnit_2()
			    		// >> flattened continuation for graph.determineIfIncludedOrEqual:733:10
			    	}
			    })(row)
			    if _err != nil { return _err }
			  }
			  return nil
			})()
			// << flattened continuation for async.eachSeries:719:7
			//console.Log('filtered:', arrFilteredAltBranchRootUnits);
			/* await */
			goDownAndCollectBestChildren_2_sync(arrFilteredAltBranchRootUnits)
			// << flattened continuation for goDownAndCollectBestChildren_2:739:9
			// :: flattened return for cb();
			return 
			// >> flattened continuation for goDownAndCollectBestChildren_2:739:9
			// >> flattened continuation for async.eachSeries:719:7
			// >> flattened continuation for conn.query:716:6
		}
		
		/* await */
		filterAltBranchRootUnits_sync()
		// << flattened continuation for filterAltBranchRootUnits:745:5
		//console.Log('best children:', arrBestChildren);
		// :: flattened return for handleBestChildrenList(arrBestChildren);
		return arrBestChildren
		// >> flattened continuation for filterAltBranchRootUnits:745:5
	}
	min_mc_wl := /* await */
	findMinMcWitnessedLevel_sync()
	// << flattened continuation for findMinMcWitnessedLevel:751:4
	bHasAltBranches := /* await */
	determineIfHasAltBranches_sync()
	// << flattened continuation for determineIfHasAltBranches:753:5
	if ! bHasAltBranches {
		//console.Log("no alt");
		if min_mc_wl >= first_unstable_mc_level {
			// :: flattened return for return handleResult(true);
			return true
		}
		// :: flattened return for return handleResult(false);
		return false
	}
	arrAltBestChildren := /* await */
	createListOfBestChildrenIncludedByLaterUnits_sync(arrAltBranchRootUnits)
	// << flattened continuation for createListOfBestChildrenIncludedByLaterUnits:775:6
/**
	max_alt_rows := /* await * /
	conn.query_sync("SELECT MAX(units.level) AS max_alt_level \n" +
		"								FROM units \n" +
		"								LEFT JOIN parenthoods ON units.unit=child_unit \n" +
		"								LEFT JOIN units AS punits ON parent_unit=punits.unit AND punits.witnessed_level >= units.witnessed_level \n" +
		"								WHERE units.unit IN(?) AND punits.unit IS NULL AND ( \n" +
		"									SELECT COUNT(*) \n" +
		"									FROM unit_witnesses \n" +
		"									WHERE unit_witnesses.unit IN(units.unit, units.witness_list_unit) AND unit_witnesses.address IN(?) \n" +
		"								)>=?", DBParamsT{
		arrAltBestChildren,
		arrWitnesses,
		constants.COUNT_WITNESSES - constants.MAX_WITNESS_LIST_MUTATIONS,
	})
 **/
//--	max_alt_rows := *(func () *[]db.MaxAltLevelRow {

	rcvr_1 := db.MaxAltLevelsReceiver{}
	conn.MustQuery("SELECT MAX(units.level) AS max_alt_level \n" +
		"FROM units \n" +
		"LEFT JOIN parenthoods ON units.unit=child_unit \n" +
		"LEFT JOIN units AS punits ON parent_unit=punits.unit AND punits.witnessed_level >= units.witnessed_level \n" +
		"WHERE units.unit IN(?) AND punits.unit IS NULL AND ( \n" +
		"	SELECT COUNT(*) \n" +
		"	FROM unit_witnesses \n" +
		"	WHERE unit_witnesses.unit IN(units.unit, units.witness_list_unit) AND unit_witnesses.address IN(?) \n" +
		")>=?", DBParamsT{
		arrAltBestChildren,
		arrWitnesses,
		constants.COUNT_WITNESSES - constants.MAX_WITNESS_LIST_MUTATIONS,
	}, &rcvr_1)
	max_alt_rows := rcvr_1.Rows
//--	return &rcvr.Rows

//--	}())
	// << flattened continuation for conn.query:780:7
	if len(max_alt_rows) != 1 {
		_core.Throw("not a single max alt level")
	}
//	max_alt_level := max_alt_rows[0].max_alt_level
	max_alt_level := max_alt_rows[0].Max_alt_level
	// allow '=' since alt WL will *never* reach max_alt_level.
	// The comparison when moving the stability point above is still strict for compatibility
	// :: flattened return for handleResult(min_mc_wl >= max_alt_level);
	return min_mc_wl >= max_alt_level
	// >> flattened continuation for conn.query:780:7
	// >> flattened continuation for createListOfBestChildrenIncludedByLaterUnits:775:6
	// >> flattened continuation for determineIfHasAltBranches:753:5
	// >> flattened continuation for findMinMcWitnessedLevel:751:4
	// >> flattened continuation for conn.query:604:3
	// >> flattened continuation for readBestParentAndItsWitnesses:603:2
	// >> flattened continuation for storage.readPropsOfUnits:598:1
}

// It is assumed earlier_unit is not marked as stable yet
// If it appears to be stable, its MC index will be marked as stable, as well as all preceeding MC indexes
func determineIfStableInLaterUnitsAndUpdateStableMcFlag_sync(conn refDBConnT, earlier_unit UnitT, arrLaterUnits UnitsT, bStableInDb bool) bool {
	var(
		advanceLastStableMcUnitAndStepForward_sync func () 
	)
	
	bStable := /* await */
	determineIfStableInLaterUnits_sync(conn, earlier_unit, arrLaterUnits)
	// << flattened continuation for determineIfStableInLaterUnits:815:1
	console.Log("determineIfStableInLaterUnits", earlier_unit, arrLaterUnits, bStable)
	if ! bStable {
		// :: flattened return for return handleResult(bStable);
		return bStable
	}
	if bStable && bStableInDb {
		// :: flattened return for return handleResult(bStable);
		return bStable
	}
//	breadcrumbs.add("stable in parents, will wait for write lock")
	breadcrumbs.Add("stable in parents, will wait for write lock")
	unlock := /* await */
//	mutex.lock_sync({*ArrayExpression*})
	mutex.Lock_sync([]string{"write"})
	// << flattened continuation for mutex.lock:822:2
//	breadcrumbs.add("stable in parents, got write lock")
	breadcrumbs.Add("stable in parents, got write lock")
	last_stable_mci := /* await */
//	storage.readLastStableMcIndex_sync(conn)
	storage.ReadLastStableMcIndex_sync(conn)
	// << flattened continuation for storage.readLastStableMcIndex:824:3
	objEarlierUnitProps := /* await */
//	storage.readUnitProps_sync(conn, earlier_unit)
	storage.ReadUnitProps_sync(conn, earlier_unit)
	// << flattened continuation for storage.readUnitProps:825:4
//	new_last_stable_mci := objEarlierUnitProps.main_chain_index
	new_last_stable_mci := objEarlierUnitProps.Main_chain_index
	if new_last_stable_mci <= last_stable_mci {
		// fix: it could've been changed by parallel tasks - No, our SQL transaction doesn't see the changes
		_core.Throw("new last stable mci expected to be higher than existing")
	}
	mci := last_stable_mci
/**
	/* await * /
	advanceLastStableMcUnitAndStepForward_sync()
	// << flattened continuation for advanceLastStableMcUnitAndStepForward:830:5
	unlock()
	// :: flattened return for handleResult(bStable);
	return bStable
	// >> flattened continuation for advanceLastStableMcUnitAndStepForward:830:5
 **/	
	advanceLastStableMcUnitAndStepForward_sync = func ()  {
		mci++
		if mci <= new_last_stable_mci {
			/* await */
			markMcIndexStable_sync(conn, mci)
			// << flattened continuation for markMcIndexStable:838:7
			/* await */
			advanceLastStableMcUnitAndStepForward_sync()
			// << flattened continuation for advanceLastStableMcUnitAndStepForward:839:8
			// :: flattened return for onDone();
			return 
			// >> flattened continuation for advanceLastStableMcUnitAndStepForward:839:8
			// >> flattened continuation for markMcIndexStable:838:7
		} else {
			// :: flattened return for onDone();
			return 
		}
	}
	// >> flattened continuation for storage.readUnitProps:825:4
	// >> flattened continuation for storage.readLastStableMcIndex:824:3
	// >> flattened continuation for mutex.lock:822:2
	// >> flattened continuation for determineIfStableInLaterUnits:815:1

	/* await */
	advanceLastStableMcUnitAndStepForward_sync()
	// << flattened continuation for advanceLastStableMcUnitAndStepForward:830:5
	unlock()
	// :: flattened return for handleResult(bStable);
	return bStable
	// >> flattened continuation for advanceLastStableMcUnitAndStepForward:830:5
}




//func readBestParentAndItsWitnesses_sync(conn DBConnT, unit UnitT) (UnitT, UnitsT) {
func readBestParentAndItsWitnesses_sync(conn refDBConnT, unit UnitT) (UnitT, AddressesT) {
	props := /* await */
//	storage.readStaticUnitProps_sync(conn, unit)
	storage.ReadStaticUnitProps_sync(conn, unit)
	// << flattened continuation for storage.readStaticUnitProps:855:1
	arrWitnesses := /* await */
//	storage.readWitnesses_sync(conn, props.best_parent_unit)
	storage.ReadWitnesses_sync(conn, props.Best_parent_unit)
	// << flattened continuation for storage.readWitnesses:856:2
	// :: flattened return for handleBestParentAndItsWitnesses(props.best_parent_unit, arrWitnesses);
//	return meta.returnArguments(props.best_parent_unit, arrWitnesses)
	return props.Best_parent_unit, arrWitnesses
	// >> flattened continuation for storage.readWitnesses:856:2
	// >> flattened continuation for storage.readStaticUnitProps:855:1
}


//func markMcIndexStable_sync(conn DBConnT, mci MCIndexT)  {
func markMcIndexStable_sync(conn refDBConnT, mci MCIndexT)  {
	var(
		handleNonserialUnits func () 
		setContentHash_sync func (unit UnitT) 
//		findStableConflictingUnits_sync func (objUnitProps PropsT) UnitsT
		findStableConflictingUnits_sync func (objUnitProps db.UnitContentRow) UnitsT
		addBalls func () 
		updateRetrievable func () 
		calcCommissions func () 
	)
	
	profiler.Start()
	arrStabilizedUnits := UnitsT{}
//	for unit := range storage.assocUnstableUnits {
	for unit := range storage.AssocUnstableUnits {
//		o := storage.assocUnstableUnits[unit]
		o := storage.AssocUnstableUnits[unit]
//		if o.main_chain_index == mci && o.is_stable == 0 {
		if o.Main_chain_index == mci && o.Is_stable == 0 {
//			o.is_stable = 1
			o.Is_stable = 1
//			storage.assocStableUnits[unit] = o
			storage.AssocStableUnits[unit] = &o.PropsT
			arrStabilizedUnits = append(arrStabilizedUnits, unit)
		}
	}
	// .. not flattening for Array.forEach
//	for unit, _ := range arrStabilizedUnits {
	for _, unit := range arrStabilizedUnits {
//		delete(storage.assocUnstableUnits, unit)
		delete(storage.AssocUnstableUnits, unit)
	}
/**
	/* await * /
	conn.query_sync("UPDATE units SET is_stable=1 WHERE is_stable=0 AND main_chain_index=?", DBParamsT{ mci })
 **/
	conn.MustExec("UPDATE units SET is_stable=1 WHERE is_stable=0 AND main_chain_index=?", DBParamsT{ mci })
/**
	// [tbd] move after function definitions
	// << flattened continuation for conn.query:877:1
	// next op
	handleNonserialUnits()
	// >> flattened continuation for conn.query:877:1
 **/	
	
	handleNonserialUnits = func ()  {
/**
		rows := /* await * /
		conn.query_sync("SELECT * FROM units WHERE main_chain_index=? AND sequence!='good' ORDER BY unit", DBParamsT{ mci })
 **/
		rcvr := db.UnitContentsReceiver{}
		conn.MustQuery("SELECT * FROM units WHERE main_chain_index=? AND sequence!='good' ORDER BY unit", DBParamsT{ mci }, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:888:2
		(func () ErrorT {
		  // :: inlined async.eachSeries:891:4
//		  for row := range rows {
		  for _, row := range rows {
//		    _err := (func (row rowT) ErrorT {
		    _err := (func (row db.UnitContentRow) ErrorT {
//		    	if row.sequence == "final-bad" {
		    	if row.Sequence == "final-bad" {
//		    		if row.content_hash {
		    		if ! row.Content_hash.IsNull() {
		    			// :: flattened return for cb();
		    			// ** need 1 return(s) instead of 0
		    			return nil
		    		} else {
		    			/* await */
//		    			setContentHash_sync(row.unit)
		    			setContentHash_sync(row.Unit)
		    			// << flattened continuation for setContentHash:898:8
		    			// :: flattened return for cb();
		    			// ** need 1 return(s) instead of 0
		    			return nil
		    			// >> flattened continuation for setContentHash:898:8
		    		}
		    		return nil
		    	}
		    	// temp-bad
//		    	if row.content_hash {
		    	if ! row.Content_hash.IsNull() {
		    		_core.Throw("temp-bad and with content_hash?")
		    	}
		    	arrConflictingUnits := /* await */
		    	findStableConflictingUnits_sync(row)
		    	// << flattened continuation for findStableConflictingUnits:905:6
//		    	sequence := (len(arrConflictingUnits) > 0 ? "final-bad": "good")
		    	sequence := "good"
		    	if len(arrConflictingUnits) > 0 { sequence = "final-bad" }
//		    	console.Log("unit " + row.unit + " has competitors " + arrConflictingUnits + ", it becomes " + sequence)
		    	console.Log("unit %s has competitors %s , it becomes %s", row.Unit, arrConflictingUnits.Join(","), sequence)
/**
		    	/* await * /
		    	conn.query_sync("UPDATE units SET sequence=? WHERE unit=?", DBParamsT{
		    		sequence,
		    		row.unit,
		    	})
 **/
		    	conn.MustExec("UPDATE units SET sequence=? WHERE unit=?", DBParamsT{
		    		sequence,
		    		row.Unit,
		    	})
		    	// << flattened continuation for conn.query:908:7
		    	if sequence == "good" {
/**
		    		/* await * /
		    		conn.query_sync("UPDATE inputs SET is_unique=1 WHERE unit=?", DBParamsT{ row.unit })
 **/
		    		conn.MustExec("UPDATE inputs SET is_unique=1 WHERE unit=?", DBParamsT{ row.Unit })
		    		// << flattened continuation for conn.query:910:9
		    		// :: flattened return for cb();
		    		// ** need 1 return(s) instead of 0
		    		return nil
		    		// >> flattened continuation for conn.query:910:9
		    	} else {
		    		/* await */
//		    		setContentHash_sync(row.unit)
		    		setContentHash_sync(row.Unit)
		    		// << flattened continuation for setContentHash:912:9
		    		// :: flattened return for cb();
		    		// ** need 1 return(s) instead of 0
		    		return nil
		    		// >> flattened continuation for setContentHash:912:9
		    	}
		    	// >> flattened continuation for conn.query:908:7
		    	// >> flattened continuation for findStableConflictingUnits:905:6
		    })(row)
		    if _err != nil { return _err }
		  }
		  return nil
		})()
		// << flattened continuation for async.eachSeries:891:4
		//if (rows.length > 0)
		//    throw "stop";
		// next op
		addBalls()
		// >> flattened continuation for async.eachSeries:891:4
		// >> flattened continuation for conn.query:888:2
	}
	
	setContentHash_sync = func (unit UnitT)  {
//		storage.readJoint(conn, unit, [*ObjectExpression*])
		storage.ReadJoint(conn, unit, storage.ReadJointCbT{

			IfNotFound: func () {
//				throw Error("bad unit not found: "+unit);
				_core.Throw("bad unit not found: %s", unit)
			},

			IfFound: func (objJoint refJointT) {
//				var content_hash = objectHash.getUnitContentHash(objJoint.unit);
				content_hash := objectHash.GetUnitContentHash(objJoint.Unit)
//				conn.query("UPDATE units SET content_hash=? WHERE unit=?", [content_hash, unit], function(){
//					onSet();
//				});
				conn.MustExec("UPDATE units SET content_hash=? WHERE unit=?", DBParamsT{
					content_hash,
					unit,
				})
				//onSet()
			},

		})
	}
	
//	findStableConflictingUnits_sync = func (objUnitProps PropsT) UnitsT {
	findStableConflictingUnits_sync = func (objUnitProps db.UnitContentRow) UnitsT {
/**
		rows := /* await * /
		conn.query_sync("SELECT competitor_units.* \n" +
			"			FROM unit_authors AS this_unit_authors \n" +
			"			JOIN unit_authors AS competitor_unit_authors USING(address) \n" +
			"			JOIN units AS competitor_units ON competitor_unit_authors.unit=competitor_units.unit \n" +
			"			JOIN units AS this_unit ON this_unit_authors.unit=this_unit.unit \n" +
			"			WHERE this_unit_authors.unit=? AND competitor_units.is_stable=1 AND +competitor_units.sequence='good' \n" +
			"				-- if it were main_chain_index <= this_unit_limci, the competitor would've been included \n" +
			"				AND (competitor_units.main_chain_index > this_unit.latest_included_mc_index) \n" +
			"				AND (competitor_units.main_chain_index <= this_unit.main_chain_index)", // if on the same mci, the smallest unit wins becuse it got selected earlier and was assigned sequence=good
		DBParamsT{ objUnitProps.unit })
 **/
		rcvr := db.UnitContentsReceiver{}
		conn.MustQuery("SELECT competitor_units.* \n" +
			"FROM unit_authors AS this_unit_authors \n" +
			"JOIN unit_authors AS competitor_unit_authors USING(address) \n" +
			"JOIN units AS competitor_units ON competitor_unit_authors.unit=competitor_units.unit \n" +
			"JOIN units AS this_unit ON this_unit_authors.unit=this_unit.unit \n" +
			"WHERE this_unit_authors.unit=? AND competitor_units.is_stable=1 AND +competitor_units.sequence='good' \n" +
			"	-- if it were main_chain_index <= this_unit_limci, the competitor would've been included \n" +
			"	AND (competitor_units.main_chain_index > this_unit.latest_included_mc_index) \n" +
			"	AND (competitor_units.main_chain_index <= this_unit.main_chain_index)", // if on the same mci, the smallest unit wins becuse it got selected earlier and was assigned sequence=good
			DBParamsT{ objUnitProps.Unit }, &rcvr)
		rows := rcvr.Rows
		// << flattened continuation for conn.query:959:2
		arrConflictingUnits := UnitsT{}
		(func () ErrorT {
		  // :: inlined async.eachSeries:973:4
//		  for row := range rows {
		  for _, row := range rows {
//		    _err := (func (row rowT) ErrorT {
		    _err := (func (row db.UnitContentRow) ErrorT {
		    	result := /* await */
//		    	graph.compareUnitsByProps_sync(conn, row, objUnitProps)
		    	graph.CompareUnitsByProps_sync(conn, row, objUnitProps)
		    	// << flattened continuation for graph.compareUnitsByProps:976:6
		    	if result == nil {
//		    		arrConflictingUnits = append(arrConflictingUnits, row.unit)
		    		arrConflictingUnits = append(arrConflictingUnits, row.Unit)
		    	}
		    	// :: flattened return for cb();
		    	// ** need 1 return(s) instead of 0
		    	return nil
		    	// >> flattened continuation for graph.compareUnitsByProps:976:6
		    })(row)
		    if _err != nil { return _err }
		  }
		  return nil
		})()
		// << flattened continuation for async.eachSeries:973:4
		// :: flattened return for handleConflictingUnits(arrConflictingUnits);
		return arrConflictingUnits
		// >> flattened continuation for async.eachSeries:973:4
		// >> flattened continuation for conn.query:959:2
	}
	
	
	addBalls = func ()  {
/**
		unit_rows := /* await * /
		conn.query_sync("SELECT units.*, ball FROM units LEFT JOIN balls USING(unit) \n" +
			"			WHERE main_chain_index=? ORDER BY level", DBParamsT{ mci })
 **/
		rcvr := db.UnitContentBallsReceiver{}
		conn.MustQuery("SELECT units.*, ball FROM units LEFT JOIN balls USING(unit) \n" +
			"WHERE main_chain_index=? ORDER BY level", DBParamsT{ mci }, &rcvr)
		unit_rows := rcvr.Rows
		// << flattened continuation for conn.query:992:2
		if len(unit_rows) == 0 {
//			_core.Throw("no units on mci " + mci)
			_core.Throw("no units on mci %d", mci)
		}
		(func () ErrorT {
		  // :: inlined async.eachSeries:998:4
//		  for objUnitProps := range unit_rows {
		  for _, objUnitProps := range unit_rows {
//		    _err := (func (objUnitProps PropsT) ErrorT {
		    _err := (func (objUnitProps db.UnitContentBallRow) ErrorT {
		    	var(
		    		addBall func () 
		    	)
		    	
//		    	unit := objUnitProps.unit
		    	unit := objUnitProps.Unit
/**
		    	parent_ball_rows := /* await * /
		    	conn.query_sync("SELECT ball FROM parenthoods LEFT JOIN balls ON parent_unit=unit WHERE child_unit=? ORDER BY ball", DBParamsT{ unit })
 **/
			rcvr := db.BallsReceiver{}
		    	conn.MustQuery("SELECT ball FROM parenthoods LEFT JOIN balls ON parent_unit=unit WHERE child_unit=? ORDER BY ball", DBParamsT{ unit }, &rcvr)
			parent_ball_rows := rcvr.Rows
		    	// << flattened continuation for conn.query:1002:6
		    	bMissingParentBalls := false
//		    	for parent_ball_row, _ := range parent_ball_rows {
		    	for _, parent_ball_row := range parent_ball_rows {
//		    		if parent_ball_row.ball == nil { bMissingParentBalls = true; break }
		    		if parent_ball_row.Ball.IsNull() { bMissingParentBalls = true; break }
		    	}
		    	if bMissingParentBalls {
//		    		_core.Throw("some parent balls not found for unit " + unit)
		    		_core.Throw("some parent balls not found for unit %s", unit)
		    	}
		    	arrParentBalls := make(BallsT, len(parent_ball_rows), len(parent_ball_rows))
//		    	for parent_ball_row, _k := range parent_ball_rows {
		    	for _k, parent_ball_row := range parent_ball_rows {
//		    		arrParentBalls[_k] := parent_ball_row.ball
		    		arrParentBalls[_k] = parent_ball_row.Ball
		    	}
		    	arrSimilarMcis := getSimilarMcis(mci)
		    	arrSkiplistUnits := UnitsT{}
		    	arrSkiplistBalls := BallsT{}

/***
			// [tbd] move down after addball definition
//		    	if objUnitProps.is_on_main_chain == 1 && len(arrSimilarMcis) > 0 {
		    	if objUnitProps.Is_on_main_chain == 1 && len(arrSimilarMcis) > 0 {
				rcvr := db.UnitBallsReceiver{}
		    		conn.MustQuery("SELECT units.unit, ball FROM units LEFT JOIN balls USING(unit) \n" +
		    			"WHERE is_on_main_chain=1 AND main_chain_index IN(?)", DBParamsT{ arrSimilarMcis }, &rcvr)
				rows := rcvr.Rows
		    		// << flattened continuation for conn.query:1014:9
		    		// .. not flattening for Array.forEach
//		    		for row, _ := range rows {
		    		for _, row := range rows {
//		    			skiplist_unit := row.unit
		    			skiplist_unit := row.Unit
//		    			skiplist_ball := row.ball
		    			skiplist_ball := row.Ball
//		    			if ! skiplist_ball {
		    			if ! (! skiplist_ball.Undefined()) {
		    				_core.Throw("no skiplist ball")
		    			}
		    			arrSkiplistUnits = append(arrSkiplistUnits, skiplist_unit)
		    			arrSkiplistBalls = append(arrSkiplistBalls, skiplist_ball)
		    		}
		    		addBall()
		    		// >> flattened continuation for conn.query:1014:9
		    	} else {
		    		addBall()
		    	}
 ***/		    	
		    	addBall = func ()  {
//		    		ball := objectHash.getBallHash(unit, arrParentBalls, arrSkiplistBalls.sort(), objUnitProps.sequence == "final-bad")
				sort.Slice(arrSkiplistBalls, func (i, j int) bool {
					return arrSkiplistBalls[i] < arrSkiplistBalls[j]
				})
//		    		ball := objectHash.GetBallHash(unit, arrParentBalls, arrSkiplistBalls, objUnitProps.Sequence == "final-bad")
		    		ball_ := objectHash.GetBallHash(unit, arrParentBalls, arrSkiplistBalls, objUnitProps.Sequence == "final-bad")
				ball := BallT(*ball_)
//		    		console.Log("ball=" + ball)
		    		console.Log("ball=%s", ball)
//		    		if objUnitProps.ball {
		    		if ! objUnitProps.Ball.IsNull() {
		    			// already inserted
//		    			if objUnitProps.ball != ball {
		    			if objUnitProps.Ball != ball {
//		    				_core.Throw("stored and calculated ball hashes do not match, ball=" + ball + ", objUnitProps=" + JSON.stringify(objUnitProps))
		    				_core.Throw("stored and calculated ball hashes do not match, ball=%s, objUnitProps=%s", ball, JSON.Stringify(objUnitProps))
		    			}
		    			// :: flattened return for return cb();
		    			// ** need 1 return(s) instead of 0
		    			return
		    		}
/**
		    		/* await * /
		    		conn.query_sync("INSERT INTO balls (ball, unit) VALUES(?,?)", DBParamsT{
		    			ball,
		    			unit,
		    		})
 **/
		    		conn.MustExec("INSERT INTO balls (ball, unit) VALUES(?,?)", DBParamsT{
		    			ball,
		    			unit,
		    		})
		    		// << flattened continuation for conn.query:1042:9
/**
		    		/* await * /
		    		conn.query_sync("DELETE FROM hash_tree_balls WHERE ball=?", DBParamsT{ ball })
 **/
		    		conn.MustExec("DELETE FROM hash_tree_balls WHERE ball=?", DBParamsT{ ball })
		    		// << flattened continuation for conn.query:1043:10
		    		if len(arrSkiplistUnits) == 0 {
		    			// :: flattened return for return cb();
		    			// ** need 1 return(s) instead of 0
		    			return 
		    		}
/**
		    		slus := make(slusT, len(arrSkiplistUnits), len(arrSkiplistUnits))
		    		for skiplist_unit, _k := range arrSkiplistUnits {
		    			slus[_k] := "(" + conn.escape(unit) + ", " + conn.escape(skiplist_unit) + ")"
		    		}
		    		sluValues := slus.join(",")
		    		/* await * /
		    		conn.query_sync("INSERT INTO skiplist_units (unit, skiplist_unit) VALUES " + sluValues)
 **/
				queryParams := make(DBParamsT, 0, len(arrSkiplistUnits)*2)
				slusSql := queryParams.AddFromIterator(func (iterFn db.AddFromIteratorFunctorT) (int, string) {
					for _, skiplist_unit := range arrSkiplistUnits {
						iterFn(unit, skiplist_unit)
					}
					return len(arrSkiplistUnits), ",(?,?)"
				})
		    		conn.MustExec("INSERT INTO skiplist_units (unit, skiplist_unit) VALUES " + slusSql, queryParams)
		    		// << flattened continuation for conn.query:1051:11
		    		// :: flattened return for cb();
		    		// ** need 1 return(s) instead of 0
		    		return 
		    		// >> flattened continuation for conn.query:1051:11
		    		// >> flattened continuation for conn.query:1043:10
		    		// >> flattened continuation for conn.query:1042:9
		    	}
		    	// >> flattened continuation for conn.query:1002:6

			// [tbd] move down after addball definition
//		    	if objUnitProps.is_on_main_chain == 1 && len(arrSimilarMcis) > 0 {
		    	if objUnitProps.Is_on_main_chain == 1 && len(arrSimilarMcis) > 0 {
				rcvr := db.UnitBallsReceiver{}
				queryParams := DBParamsT{}
				smcisSql := queryParams.AddMCIndexes(arrSimilarMcis)
		    		conn.MustQuery("SELECT units.unit, ball FROM units LEFT JOIN balls USING(unit) \n" +
		    			"WHERE is_on_main_chain=1 AND main_chain_index IN(" + smcisSql + ")", queryParams, &rcvr)
				rows := rcvr.Rows
		    		// << flattened continuation for conn.query:1014:9
		    		// .. not flattening for Array.forEach
//		    		for row, _ := range rows {
		    		for _, row := range rows {
//		    			skiplist_unit := row.unit
		    			skiplist_unit := row.Unit
//		    			skiplist_ball := row.ball
		    			skiplist_ball := row.Ball
//		    			if ! skiplist_ball {
		    			if ! (! skiplist_ball.IsNull()) {
		    				_core.Throw("no skiplist ball")
		    			}
		    			arrSkiplistUnits = append(arrSkiplistUnits, skiplist_unit)
		    			arrSkiplistBalls = append(arrSkiplistBalls, skiplist_ball)
		    		}
		    		addBall()
		    		// >> flattened continuation for conn.query:1014:9
		    	} else {
		    		addBall()
		    	}

			return nil
		    })(objUnitProps)
		    if _err != nil { return _err }
		  }
		  return nil
		})()
		// << flattened continuation for async.eachSeries:998:4
		// next op
		updateRetrievable()
		// >> flattened continuation for async.eachSeries:998:4
		// >> flattened continuation for conn.query:992:2
	}
	
	updateRetrievable = func ()  {
//[uu]		min_retrievable_mci := /* await */
//		storage.updateMinRetrievableMciAfterStabilizingMci_sync(conn, mci)
		storage.UpdateMinRetrievableMciAfterStabilizingMci_sync(conn, mci)
		// << flattened continuation for storage.updateMinRetrievableMciAfterStabilizingMci:1071:2
		profiler.Stop("mc-mark-stable")
		calcCommissions()
		// >> flattened continuation for storage.updateMinRetrievableMciAfterStabilizingMci:1071:2
	}
	
	calcCommissions = func ()  {
		(func () ErrorT {
		  // :: inlined async.series:1078:2
//		  for _f := range (AsyncFunctorsT{
		  for _, _f := range (AsyncFunctorsT{
		  	func () ErrorT {
		  		profiler.Start()
		  		/* await */
//		  		headers_commission.calcHeadersCommissions_sync(conn)
		  		headers_commission.CalcHeadersCommissions_sync(conn)
		  		// << flattened continuation for headers_commission.calcHeadersCommissions:1081:4
		  		// :: flattened return for cb();
		  		// ** need 1 return(s) instead of 0
		  		return nil
		  		// >> flattened continuation for headers_commission.calcHeadersCommissions:1081:4
		  	},
		  	func () ErrorT {
		  		profiler.Stop("mc-headers-commissions")
		  		/* await */
//		  		paid_witnessing.updatePaidWitnesses_sync(conn)
		  		paid_witnessing.UpdatePaidWitnesses_sync(conn)
		  		// << flattened continuation for paid_witnessing.updatePaidWitnesses:1085:4
		  		// :: flattened return for cb();
		  		// ** need 1 return(s) instead of 0
		  		return nil
		  		// >> flattened continuation for paid_witnessing.updatePaidWitnesses:1085:4
		  	},
		  }) {
		    if _err := _f() ; _err != nil { return _err }
		  }
		  return nil
		})()
		// << flattened continuation for async.series:1078:2
		// .. not flattening for process.nextTick
///		process.nextTick(func () {*returns*} {
			// don't call it synchronously with event emitter
//			eventBus.emit("mci_became_stable", mci)
			eventBus.Emit("mci_became_stable", mci)
///		})
		// :: flattened return for onDone();
		return 
		// >> flattened continuation for async.series:1078:2
	}

	// << flattened continuation for conn.query:877:1
	// next op
	handleNonserialUnits()
	// >> flattened continuation for conn.query:877:1
}

// returns list of past MC indices for skiplist
func getSimilarMcis(mci MCIndexT) MCIndexesT {
	arrSimilarMcis := MCIndexesT{}
	divisor := 10
	for true {
		if int(mci) % divisor == 0 {
			arrSimilarMcis = append(arrSimilarMcis, MCIndexT(int(mci) - divisor))
			divisor *= 10
		} else {
			return arrSimilarMcis
		}
	}
	return nil
}

func throwError(msg string) {
//      if typeof window == "undefined" {
        if true {
                _core.Throw(msg)
        } else {
//              eventBus.emit("nonfatal_error", msg, [*NewExpression*])
                eventBus.Emit("nonfatal_error", msg, nil)
        }
}

//exports.updateMainChain = updateMainChain
//exports.determineIfStableInLaterUnitsAndUpdateStableMcFlag = determineIfStableInLaterUnitsAndUpdateStableMcFlag
//exports.determineIfStableInLaterUnits = determineIfStableInLaterUnits


// converted golang end

