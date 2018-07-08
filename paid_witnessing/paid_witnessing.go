
// converted golang begin

package paid_witnessing

import(
	"time"
	"math"

 _core	"nodejs/core"
	"nodejs/console"

 .	"github.com/byteball/go-byteballcore/types"

)

import(
// _		"lodash"
//		"async"
		"github.com/byteball/go-byteballcore/storage"
		"github.com/byteball/go-byteballcore/graph"
		"github.com/byteball/go-byteballcore/db"
		"github.com/byteball/go-byteballcore/constants"
//		"github.com/byteball/go-byteballcore/conf"
		"github.com/byteball/go-byteballcore/mc_outputs"
		"github.com/byteball/go-byteballcore/profiler"
)

type(
	DBConnT		= db.DBConnT
	refDBConnT	= *DBConnT

	DBParamsT	= db.DBParamsT

)


//func calcWitnessEarnings(conn DBConnT, type typeT, from_main_chain_index MCIndexT, to_main_chain_index MCIndexT, address AddressT, callbacks callbacksT)  {
func CalcWitnessEarnings(conn refDBConnT, type_ string, from_main_chain_index MCIndexT, to_main_chain_index MCIndexT,
	address AddressT, callbacks mc_outputs.CalcEarningsCbT)  {
/**
	count_rows := /* await * /
	conn.query_sync("SELECT COUNT(*) AS count FROM units WHERE is_on_main_chain=1 AND is_stable=1 AND main_chain_index>=? AND main_chain_index<=?", DBParamsT{
		to_main_chain_index,
		to_main_chain_index + MCIndexT(constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING + 1),
	})
 **/
	rcvr := db.CountsReceiver{}
	conn.MustQuery("SELECT COUNT(*) AS count FROM units WHERE is_on_main_chain=1 AND is_stable=1 AND main_chain_index>=? AND main_chain_index<=?", DBParamsT{
		to_main_chain_index,
		to_main_chain_index + MCIndexT(constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING + 1),
	}, &rcvr)
	count_rows := rcvr.Rows
	// << flattened continuation for conn.query:15:1
//	if count_rows[0].count != constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING + 2 {
	if count_rows[0].Count != constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING + 2 {
//		callbacks.ifError("not enough stable MC units after to_main_chain_index")
		callbacks.IfError("not enough stable MC units after to_main_chain_index")
		return
	}
//	mc_outputs.calcEarnings(conn, type, from_main_chain_index, to_main_chain_index, address, callbacks)
	mc_outputs.CalcEarnings(conn, type_, from_main_chain_index, to_main_chain_index, address, callbacks)
	// >> flattened continuation for conn.query:15:1
}

/*
function readMaxWitnessSpendableMcIndex(conn, handleMaxSpendableMcIndex){
	conn.query("SELECT MAX(main_chain_index) AS max_mc_index FROM units WHERE is_on_main_chain=1 AND is_stable=1", function(rows){
		var max_mc_index = rows[0].max_mc_index;
		var max_spendable_mc_index = max_mc_index - constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING - 1;
		if (max_spendable_mc_index <= 0)
			return handleMaxSpendableMcIndex(max_spendable_mc_index);
		/ *
		function checkIfMajorityWitnessedByParentsAndAdjust(){
			readUnitOnMcIndex(conn, max_spendable_mc_index, function(unit){
				determineIfMajorityWitnessedByDescendants(conn, unit, arrParents, function(bWitnessed){
					if (!bWitnessed){
						max_spendable_mc_index--;
						checkIfMajorityWitnessedByParentsAndAdjust();
					}
					else
						handleMaxSpendableMcIndex(max_spendable_mc_index);
				});
			});
		}
		* /
		//arrParents ? checkIfMajorityWitnessedByParentsAndAdjust() : 
		handleMaxSpendableMcIndex(max_spendable_mc_index);
	});
}

 */

func readUnitOnMcIndex_sync(conn DBConnT, main_chain_index MCIndexT) UnitT {
/**
	rows := /* await * /
	conn.query_sync("SELECT unit FROM units WHERE is_on_main_chain=1 AND main_chain_index=?", DBParamsT{ main_chain_index })
 **/
	rcvr := db.UnitsReceiver{}
	conn.MustQuery("SELECT unit FROM units WHERE is_on_main_chain=1 AND main_chain_index=?", DBParamsT{ main_chain_index }, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:54:1
	if len(rows) != 1 {
//		_core.Throw("no units or more than one unit on MC index " + main_chain_index)
		_core.Throw("no units or more than one unit on MC index %d", main_chain_index)
	}
	// :: flattened return for handleUnit(rows[0].unit);
//	return rows[0].unit
	return rows[0].Unit
	// >> flattened continuation for conn.query:54:1
}

//func updatePaidWitnesses_sync(conn DBConnT)  {
func UpdatePaidWitnesses_sync(conn refDBConnT)  {
	console.Log("updating paid witnesses")
	profiler.Start()
	last_stable_mci := /* await */
//	storage.readLastStableMcIndex_sync(conn)
	storage.ReadLastStableMcIndex_sync(conn)
	// << flattened continuation for storage.readLastStableMcIndex:64:1
	profiler.Stop("mc-wc-readLastStableMCI")
//	max_spendable_mc_index := getMaxSpendableMciForLastBallMci(last_stable_mci)
	max_spendable_mc_index := GetMaxSpendableMciForLastBallMci(last_stable_mci)
	if max_spendable_mc_index > 0 {
		/* await */
		buildPaidWitnessesTillMainChainIndex_sync(conn, max_spendable_mc_index)
		// << flattened continuation for buildPaidWitnessesTillMainChainIndex:68:3
		// :: flattened return for cb();
		return 
		// >> flattened continuation for buildPaidWitnessesTillMainChainIndex:68:3
	} else {
		// :: flattened return for cb();
		return 
	}
	// >> flattened continuation for storage.readLastStableMcIndex:64:1
}

//func buildPaidWitnessesTillMainChainIndex_sync(conn DBConnT, to_main_chain_index MCIndexT)  {
func buildPaidWitnessesTillMainChainIndex_sync(conn refDBConnT, to_main_chain_index MCIndexT)  {
	var(
//		onIndexDone func (err errT) 
		onIndexDone func (err ErrorT) 
	)
	
	profiler.Start()
	// [tbd] refactor into db.getCross()?..
//	cross := (conf.storage == "sqlite" ? "CROSS": "")
	cross := "CROSS"
/**
	rows := /* await * /
	conn.query_sync("SELECT MIN(main_chain_index) AS min_main_chain_index FROM balls " + cross + " JOIN units USING(unit) WHERE count_paid_witnesses IS NULL")
 **/
	rcvr := db.MinMCIndexsReceiver{}
	conn.MustQuery("SELECT MIN(main_chain_index) AS min_main_chain_index FROM balls " + cross + " JOIN units USING(unit) WHERE count_paid_witnesses IS NULL", DBParamsT{}, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:78:1
	profiler.Stop("mc-wc-minMCI")
//	main_chain_index := rows[0].min_main_chain_index
	main_chain_index := rows[0].Min_main_chain_index
	if main_chain_index > to_main_chain_index {
		// :: flattened return for return cb();
		return 
	}
	
//	onIndexDone = func (err errT)  {
	onIndexDone = func (err ErrorT)  {
//		if err {
		if err != nil {
			// impossible
			_core.Throw("%s", err)
		} else {
			main_chain_index++
			if main_chain_index > to_main_chain_index {
				// :: flattened return for cb();
				return 
			} else {
				onIndexDone(/* await */
				buildPaidWitnessesForMainChainIndex_sync(conn, main_chain_index))
			}
		}
	}
	onIndexDone(/* await */
	buildPaidWitnessesForMainChainIndex_sync(conn, main_chain_index))
	// >> flattened continuation for conn.query:78:1
}

//func buildPaidWitnessesForMainChainIndex_sync(conn DBConnT, main_chain_index MCIndexT) ErrorT {
func buildPaidWitnessesForMainChainIndex_sync(conn refDBConnT, main_chain_index MCIndexT) ErrorT {
//	console.Log("updating paid witnesses mci " + main_chain_index)
	console.Log("updating paid witnesses mci %d", main_chain_index)
	profiler.Start()
/**
	rows := /* await * /
	conn.query_sync("SELECT COUNT(*) AS count, SUM(CASE WHEN is_stable=1 THEN 1 ELSE 0 END) AS count_on_stable_mc \n" +
		"		FROM units WHERE is_on_main_chain=1 AND main_chain_index>=? AND main_chain_index<=?", DBParamsT{
		main_chain_index,
		main_chain_index + MCIndexT(constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING + 1),
	})
 **/
	rcvr := db.CountOnStableMCsReceiver{}
	conn.MustQuery("SELECT COUNT(*) AS count, SUM(CASE WHEN is_stable=1 THEN 1 ELSE 0 END) AS count_on_stable_mc \n" +
		"		FROM units WHERE is_on_main_chain=1 AND main_chain_index>=? AND main_chain_index<=?", DBParamsT{
		main_chain_index,
		main_chain_index + MCIndexT(constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING + 1),
	}, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:106:1
	profiler.Stop("mc-wc-select-count")
//	count := rows[0].count
	count := rows[0].Count
//	count_on_stable_mc := rows[0].count_on_stable_mc
	count_on_stable_mc := rows[0].Count_on_stable_mc
	if count != constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING + 2 {
//		_core.Throw("main chain is not long enough yet for MC index " + main_chain_index)
		_core.Throw("main chain is not long enough yet for MC index %d", main_chain_index)
	}
	if count_on_stable_mc != count {
//		_core.Throw("not enough stable MC units yet after MC index " + main_chain_index + ": count_on_stable_mc=" + count_on_stable_mc + ", count=" + count)
		_core.Throw("not enough stable MC units yet after MC index %d: count_on_stable_mc=%d, count=%d", main_chain_index, count_on_stable_mc, count)
	}
	
	profiler.Start()
	arrWitnesses := /* await */
	readMcUnitWitnesses_sync(conn, main_chain_index)
	// << flattened continuation for readMcUnitWitnesses:121:3
/**
	/* await * /
	conn.query_sync("CREATE TEMPORARY TABLE paid_witness_events_tmp ( \n" +
		"					unit CHAR(44) NOT NULL, \n" +
		"					address CHAR(32) NOT NULL, \n" +
		"					delay TINYINT NULL)")
 **/
	conn.MustExec("CREATE TEMPORARY TABLE paid_witness_events_tmp ( \n" +
//[!!]	conn.MustExec("CREATE TEMPORARY TABLE IF NOT EXISTS paid_witness_events_tmp ( \n" +
		"					unit CHAR(44) NOT NULL, \n" +
		"					address CHAR(32) NOT NULL, \n" +
		"					delay TINYINT NULL)", DBParamsT{})
	// << flattened continuation for conn.query:122:4
/**
	rows := /* await * /
	conn.query_sync("SELECT * FROM units WHERE main_chain_index=?", DBParamsT{ main_chain_index })
 **/
	rcvr_1 := db.UnitContentsReceiver{}
	conn.MustQuery("SELECT * FROM units WHERE main_chain_index=?", DBParamsT{ main_chain_index }, &rcvr_1)
	rows_1 := rcvr_1.Rows
	// << flattened continuation for conn.query:126:6
	profiler.Stop("mc-wc-select-units")
//	et = 0
	et = time.Unix(0,0)
//	rt = 0
	rt = time.Duration(0)
	err := (func () ErrorT {
	  // :: inlined async.eachSeries:129:7
//	  for row := range rows {
	  for _, row := range rows_1 {
//	    _err := (func (row rowT) ErrorT {
	    _err := (func (row db.UnitContentRow) ErrorT {
	    	// the unit itself might be never majority witnessed by unit-designated witnesses (which might be far off), 
	    	// but its payload commission still belongs to and is spendable by the MC-unit-designated witnesses.
	    	//if (row.is_stable !== 1)
	    	//    throw "unit "+row.unit+" is not on stable MC yet";
	    	/* await */
	    	buildPaidWitnesses_sync(conn, row, arrWitnesses)
//	    	buildPaidWitnesses_sync(conn, row.PropsT, arrWitnesses)
	    	// << flattened continuation for buildPaidWitnesses:136:9
	    	// :: flattened return for cb2();
	    	// ** need 1 return(s) instead of 0
	    	return nil
	    	// >> flattened continuation for buildPaidWitnesses:136:9
	    })(row)
	    if _err != nil { return _err }
	  }
	  return nil
	})()
	// << flattened continuation for async.eachSeries:129:7
//	console.log(rt, et)
	console.Log("%v %v", rt, et)
//	if err {
	if err != nil {
		// impossible
//		_core.Throw(err)
		_core.Throw("%s", err)
	}
	//var t=Date.now();
	profiler.Start()
/**
	/* await * /
	conn.query_sync("INSERT INTO witnessing_outputs (main_chain_index, address, amount) \n" +
		"										SELECT main_chain_index, address, \n" +
		"											SUM(CASE WHEN sequence='good' THEN ROUND(1.0*payload_commission/count_paid_witnesses) ELSE 0 END) \n" +
		"										FROM balls \n" +
		"										JOIN units USING(unit) \n" +
		"										JOIN paid_witness_events_tmp USING(unit) \n" +
		"										WHERE main_chain_index=? \n" +
		"										GROUP BY address", DBParamsT{ main_chain_index })
 **/
	conn.MustExec("INSERT INTO witnessing_outputs (main_chain_index, address, amount) \n" +
		"SELECT main_chain_index, address, \n" +
		"	SUM(CASE WHEN sequence='good' THEN ROUND(1.0*payload_commission/count_paid_witnesses) ELSE 0 END) \n" +
		"FROM balls \n" +
		"JOIN units USING(unit) \n" +
		"JOIN paid_witness_events_tmp USING(unit) \n" +
		"WHERE main_chain_index=? \n" +
		"GROUP BY address", DBParamsT{ main_chain_index })
	// << flattened continuation for conn.query:144:9
	//console.log(Date.now()-t);
/**
	/* await * /
	conn.query_sync(conn.dropTemporaryTable("paid_witness_events_tmp"))
 **/
	// [tbd] "database table is locked"?..
	conn.MustExec(conn.DropTemporaryTable("paid_witness_events_tmp"), nil)
//!!	conn.MustExec("DELETE FROM paid_witness_events_tmp", nil)
	// << flattened continuation for conn.query:156:11
	profiler.Stop("mc-wc-aggregate-events")
	// :: flattened return for cb();
	// ** need 1 return(s) instead of 0
	return nil
	// >> flattened continuation for conn.query:156:11
	// >> flattened continuation for conn.query:144:9
	// >> flattened continuation for async.eachSeries:129:7
	// >> flattened continuation for conn.query:126:6
	// >> flattened continuation for conn.query:122:4
	// >> flattened continuation for readMcUnitWitnesses:121:3
	// >> flattened continuation for conn.query:106:1
}


//func readMcUnitWitnesses_sync(conn DBConnT, main_chain_index MCIndexT) AddressesT {
func readMcUnitWitnesses_sync(conn refDBConnT, main_chain_index MCIndexT) AddressesT {
/**
	rows := /* await * /
	conn.query_sync("SELECT witness_list_unit, unit FROM units WHERE main_chain_index=? AND is_on_main_chain=1", DBParamsT{ main_chain_index })
 **/
	rcvr := db.WitnessListUnitAndUnitsReceiver{}
	conn.MustQuery("SELECT witness_list_unit, unit FROM units WHERE main_chain_index=? AND is_on_main_chain=1", DBParamsT{ main_chain_index }, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:174:1
	if len(rows) != 1 {
//		_core.Throw("not 1 row on MC " + main_chain_index)
		_core.Throw("not 1 row on MC %d", main_chain_index)
	}
//	witness_list_unit := rows[0].unit
	witness_list_unit := rows[0].Unit
//	if rows[0].witness_list_unit {
	if ! rows[0].Witness_list_unit.IsNull() {
//		witness_list_unit = rows[0].witness_list_unit
		witness_list_unit = rows[0].Witness_list_unit
	}
	// :: flattened return for handleWitnesses(storage.readWitnessList(conn, witness_list_unit));
//	return /* await */
//	storage.readWitnessList_sync(conn, witness_list_unit)
	return storage.ReadWitnessList_sync(conn, witness_list_unit, false)
	// >> flattened continuation for conn.query:174:1
}

//et := {*init:null*}
var et TimeT
//rt := {*init:null*}
var rt time.Duration

type(
	DelayByAddressMapT = map[AddressT] DelayT
)

//func buildPaidWitnesses_sync(conn DBConnT, objUnitProps PropsT, arrWitnesses AddressesT)  {
func buildPaidWitnesses_sync(conn refDBConnT, objUnitProps db.UnitContentRow, arrWitnesses AddressesT)  {
	var(
//		updateCountPaidWitnesses func (count_paid_witnesses count_paid_witnessesT) 
		updateCountPaidWitnesses func (count_paid_witnesses int) 
	)
	
//	updateCountPaidWitnesses = func (count_paid_witnesses count_paid_witnessesT)  {
	updateCountPaidWitnesses = func (count_paid_witnesses int)  {
/**
		/* await * /
		conn.query_sync("UPDATE balls SET count_paid_witnesses=? WHERE unit=?", DBParamsT{
			count_paid_witnesses,
//			objUnitProps.unit,
			objUnitProps.Unit,
		})
 **/
		conn.MustExec("UPDATE balls SET count_paid_witnesses=? WHERE unit=?", DBParamsT{
			count_paid_witnesses,
//			objUnitProps.unit,
			objUnitProps.Unit,
		})
		// << flattened continuation for conn.query:187:2
		profiler.Stop("mc-wc-insert-events")
		// :: flattened return for onDone();
		return 
		// >> flattened continuation for conn.query:187:2
	}
	
//	unit := objUnitProps.unit
	unit := objUnitProps.Unit
//	to_main_chain_index := objUnitProps.main_chain_index + constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING
	to_main_chain_index := objUnitProps.Main_chain_index + MCIndexT(constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING)
	
//	t := Date.now()
	t := time.Now()
	arrUnits := /* await */
//	graph.readDescendantUnitsByAuthorsBeforeMcIndex_sync(conn, objUnitProps, arrWitnesses, to_main_chain_index)
	graph.ReadDescendantUnitsByAuthorsBeforeMcIndex_sync(conn, objUnitProps, arrWitnesses, to_main_chain_index)
	// << flattened continuation for graph.readDescendantUnitsByAuthorsBeforeMcIndex:197:1
//	rt += Date.now() - t
	rt += time.Now().Sub(t)
//	t = Date.now()
	t = time.Now()
/**
	strUnitsList := (len(arrUnits) == 0 ? "NULL": // .. not flattening for Array.map
	arrUnits.map(func (unit UnitT) {*returns*} {
		return conn.escape(unit)
	}).join(", "))
 **/
	//throw "no witnesses before mc "+to_main_chain_index+" for unit "+objUnitProps.unit;
	profiler.Start()
/**
	rows := /* await * /
	conn.query_sync(// we don't care if the unit is majority witnessed by the unit-designated witnesses
	// _left_ join forces use of indexes in units
	// can't get rid of filtering by address because units can be co-authored by witness with somebody else
	"SELECT address, MIN(main_chain_index-?) AS delay \n" +
		"			FROM units \n" +
		"			LEFT JOIN unit_authors USING(unit) \n" +
		"			WHERE unit IN(" + strUnitsList + ") AND address IN(?) AND +sequence='good' \n" +
		"			GROUP BY address", DBParamsT{
		objUnitProps.main_chain_index,
		arrWitnesses,
	})
 **/
	// [fyi] query is removed entirely as arrUnits returned by graph
	// [fyi] method just before already contain everything unit-related
	minDelayByAddress := make(DelayByAddressMapT)
	for _, objUnit := range arrUnits {
		if !objUnit.Good { continue }
		address := objUnit.Address
		mdba := DelayT(math.MaxInt64)
		if mdba_, _exists := minDelayByAddress[address]; _exists {
			mdba = mdba_
		}
		delay := DelayT(objUnit.Main_chain_index - objUnitProps.Main_chain_index)
		if delay < mdba {
			minDelayByAddress[address] = delay
		}
	}
	rows := make([]db.AddressDelayRow, 0, len(minDelayByAddress))
	for address, mdba := range minDelayByAddress {
		if arrWitnesses.IndexOf(address) == -1 { continue }
		rows = append(rows, db.AddressDelayRow{
			Address: address,
			Delay: mdba,
		})
	}
	// << flattened continuation for conn.query:203:2
//	et += Date.now() - t
	et.Add(time.Now().Sub(t))
	count_paid_witnesses := len(rows)
//	arrValues := {*init:null*}
	arrValues := PaidWitnessEventsT{}
	if count_paid_witnesses == 0 {
		// nobody witnessed, pay equally to all
		count_paid_witnesses = len(arrWitnesses)
/**
		arrValues = // .. not flattening for Array.map
		arrWitnesses.map(func (address AddressT) {*returns*} {
			return "(" + conn.escape(unit) + ", " + conn.escape(address) + ", NULL)"
		})
 **/
		for _, address := range arrWitnesses {
			arrValues = append(arrValues, PaidWitnessEventT{
				unit,
				address,
				DelayT_Null,
			})
		}
	} else {
/**
		arrValues = // .. not flattening for Array.map
		rows.map(func (row rowT) {*returns*} {
			return "(" + conn.escape(unit) + ", " + conn.escape(row.address) + ", " + row.delay + ")"
		})
 **/
		for _, row := range rows {
			arrValues = append(arrValues, PaidWitnessEventT{
				unit,
				row.Address,
				row.Delay,
			})
		}
	}
	profiler.Stop("mc-wc-select-events")
	profiler.Start()
	/* await */
/**
	conn.query_sync("INSERT INTO paid_witness_events_tmp (unit, address, delay) VALUES " + arrValues.join(", "))
 **/
	{{
	queryParams := DBParamsT{}
	valuesSql := queryParams.AddPaidWitnessEvents(arrValues)
	conn.MustExec("INSERT INTO paid_witness_events_tmp (unit, address, delay) VALUES " + valuesSql, queryParams)
	}}
	// << flattened continuation for conn.query:224:4
	updateCountPaidWitnesses(count_paid_witnesses)
	// >> flattened continuation for conn.query:224:4
	// >> flattened continuation for conn.query:203:2
	// >> flattened continuation for graph.readDescendantUnitsByAuthorsBeforeMcIndex:197:1
}

func GetMaxSpendableMciForLastBallMci(last_ball_mci MCIndexT) MCIndexT {
//	return last_ball_mci - 1 - constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING
	return last_ball_mci - MCIndexT(1 + constants.COUNT_MC_BALLS_FOR_PAID_WITNESSING)
}


//exports.updatePaidWitnesses = updatePaidWitnesses
//exports.calcWitnessEarnings = calcWitnessEarnings
//exports.getMaxSpendableMciForLastBallMci = getMaxSpendableMciForLastBallMci


// converted golang end

