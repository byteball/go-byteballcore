
// converted golang begin

package storage

import(
	"sort"
//	"strings"

 _core	"nodejs/core"
 JSON	"nodejs/json"
	"nodejs/console"

 .	"github.com/byteball/go-byteballcore/types"
//	"github.com/byteball/go-byteballcore/constants"

//	"github.com/byteball/go-byteballcore/db"
//	"github.com/byteball/go-byteballcore/archiving"
)

import(
//		"async"
// _		"lodash"
		"github.com/byteball/go-byteballcore/db"
//		"github.com/byteball/go-byteballcore/conf"
// objectHash	"github.com/byteball/go-byteballcore/object_hash"
		"github.com/byteball/go-byteballcore/constants"
//		"github.com/byteball/go-byteballcore/mutex"
		"github.com/byteball/go-byteballcore/archiving"
//		"github.com/byteball/go-byteballcore/profiler"
)

type(
	DBConnT		= db.DBConnT
	refDBConnT	= *DBConnT

	DBParamsT	= db.DBParamsT

	refPropsT	= *PropsT
	refJointT	= *JointT

	BoolByUnitMapT	= map[UnitT] bool

	PropsByUnitMapT = map[UnitT] refPropsT
	XPropsByUnitMapT = map[UnitT] refXPropsT

	AddressesByUnitMapT = map[UnitT] AddressesT

	XPropsT struct{
		Unit		UnitT
		PropsT
		Parent_units	UnitsT
	}
	refXPropsT = *XPropsT

)

//MAX_INT32 := Math.pow(2, 31) - 1

//genesis_ball := objectHash.getBallHash(constants.GENESIS_UNIT)

//MAX_ITEMS_IN_CACHE := 300
//assocKnownUnits := [*ObjectExpression*]
var AssocKnownUnits BoolByUnitMapT = make(BoolByUnitMapT)
//assocCachedUnits := [*ObjectExpression*]
var AssocCachedUnits PropsByUnitMapT = make(PropsByUnitMapT)
//assocCachedUnitAuthors := [*ObjectExpression*]
var AssocCachedUnitAuthors AddressesByUnitMapT = make(AddressesByUnitMapT)
//assocCachedUnitWitnesses := [*ObjectExpression*]
var AssocCachedUnitWitnesses AddressesByUnitMapT = make(AddressesByUnitMapT)
//assocCachedAssetInfos := [*ObjectExpression*]

//assocUnstableUnits := [*ObjectExpression*]
var AssocUnstableUnits XPropsByUnitMapT = make(XPropsByUnitMapT)
//assocStableUnits := [*ObjectExpression*]
var AssocStableUnits PropsByUnitMapT = make(PropsByUnitMapT)

//min_retrievable_mci := nil
var min_retrievable_mci MCIndexT = MCIndexT_Null
// [fyi] moved to .Init()
//initializeMinRetrievableMci()


const _readJoint = `
func readJoint(conn DBConnT, unit UnitT, callbacks callbacksT)  {
	if ! conf.bSaveJointJson {
		readJointDirectly(conn, unit, callbacks)
		return
	}
	rows := /* await */
	conn.query_sync("SELECT json FROM joints WHERE unit=?", DBParamsT{ unit })
	// << flattened continuation for conn.query:34:1
	if len(rows) == 0 {
		readJointDirectly(conn, unit, callbacks)
		return
	}
	callbacks.ifFound(JSON.parse(rows[0].json))
	// >> flattened continuation for conn.query:34:1
}
`
type(
	ReadJointCbT struct{
		IfFound func (objJoint refJointT)
		IfNotFound func ()
	}
)
func ReadJoint(conn refDBConnT, unit UnitT, cb ReadJointCbT) {
	panic("[tbd] ReadJoint")
}


const _readJointDirectly = `
func readJointDirectly(conn DBConnT, unit UnitT, callbacks callbacksT, bRetrying bool)  {
	console.log("reading unit " + unit)
	if min_retrievable_mci == nil {
		console.log("min_retrievable_mci not known yet")
		setTimeout(func () {*returns*} {
			readJointDirectly(conn, unit, callbacks)
		}, 1000)
		return 
	}
	unit_rows := /* await */
	conn.query_sync("SELECT units.unit, version, alt, witness_list_unit, last_ball_unit, balls.ball AS last_ball, is_stable, \n" +
		"			content_hash, headers_commission, payload_commission, main_chain_index, " + conn.getUnixTimestamp("units.creation_date") + " AS timestamp \n" +
		"		FROM units LEFT JOIN balls ON last_ball_unit=balls.unit WHERE units.unit=?", DBParamsT{ unit })
	// << flattened continuation for conn.query:51:1
	if len(unit_rows) == 0 {
		//profiler.stop('read');
		callbacks.ifNotFound()
		return
	}
	objUnit := unit_rows[0]
	objJoint := [*ObjectExpression*]
	main_chain_index := objUnit.main_chain_index
	//delete objUnit.main_chain_index;
	objUnit.timestamp = parseInt(objUnit.timestamp)
	bFinalBad := ! ! objUnit.content_hash
	bStable := objUnit.is_stable
	objUnit.is_stable = nil
	
	objectHash.cleanNulls(objUnit)
	bVoided := objUnit.content_hash && main_chain_index < min_retrievable_mci
	bRetrievable := main_chain_index >= min_retrievable_mci || main_chain_index == nil
	
	if ! conf.bLight && ! objUnit.last_ball && ! isGenesisUnit(unit) {
		_core.Throw("no last ball in unit " + JSON.stringify(objUnit))
	}
	
	// unit hash verification below will fail if:
	// 1. the unit was received already voided, i.e. its messages are stripped and content_hash is set
	// 2. the unit is still retrievable (e.g. we are syncing)
	// In this case, bVoided=false hence content_hash will be deleted but the messages are missing
	if bVoided {
		//delete objUnit.last_ball;
		//delete objUnit.last_ball_unit;
		objUnit.headers_commission = nil
		objUnit.payload_commission = nil
	} else {
		objUnit.content_hash = nil
	}
	
	(func () ErrorT {
	  // :: inlined async.series:90:3
	  for _f := range AsyncFunctorsT{
	  	func () ErrorT {
	  		rows := /* await */
	  		conn.query_sync("SELECT parent_unit \n" +
	  			"						FROM parenthoods \n" +
	  			"						WHERE child_unit=? \n" +
	  			"						ORDER BY parent_unit", DBParamsT{ unit })
	  		// << flattened continuation for conn.query:92:5
	  		if len(rows) == 0 {
	  			// :: flattened return for return callback();
	  			// ** need 1 return(s) instead of 0
	  			return 
	  		}
	  		objUnit.parent_units = // .. not flattening for Array.map
	  		rows.map(func (row rowT) {*returns*} {
	  			return row.parent_unit
	  		})
	  		// :: flattened return for callback();
	  		// ** need 1 return(s) instead of 0
	  		return 
	  		// >> flattened continuation for conn.query:92:5
	  	},
	  	func () ErrorT {
	  		// ball
	  		if bRetrievable && ! isGenesisUnit(unit) {
	  			// :: flattened return for return callback();
	  			// ** need 1 return(s) instead of 0
	  			return 
	  		}
	  		rows := /* await */
	  		conn.query_sync("SELECT ball FROM balls WHERE unit=?", DBParamsT{ unit })
	  		// << flattened continuation for conn.query:111:5
	  		if len(rows) == 0 {
	  			// :: flattened return for return callback();
	  			// ** need 1 return(s) instead of 0
	  			return 
	  		}
	  		objJoint.ball = rows[0].ball
	  		// :: flattened return for callback();
	  		// ** need 1 return(s) instead of 0
	  		return 
	  		// >> flattened continuation for conn.query:111:5
	  	},
	  	func () ErrorT {
	  		// skiplist
	  		if bRetrievable {
	  			// :: flattened return for return callback();
	  			// ** need 1 return(s) instead of 0
	  			return 
	  		}
	  		rows := /* await */
	  		conn.query_sync("SELECT skiplist_unit FROM skiplist_units WHERE unit=? ORDER BY skiplist_unit", DBParamsT{ unit })
	  		// << flattened continuation for conn.query:121:5
	  		if len(rows) == 0 {
	  			// :: flattened return for return callback();
	  			// ** need 1 return(s) instead of 0
	  			return 
	  		}
	  		objJoint.skiplist_units = // .. not flattening for Array.map
	  		rows.map(func (row rowT) {*returns*} {
	  			return row.skiplist_unit
	  		})
	  		// :: flattened return for callback();
	  		// ** need 1 return(s) instead of 0
	  		return 
	  		// >> flattened continuation for conn.query:121:5
	  	},
	  	func () ErrorT {
	  		rows := /* await */
	  		conn.query_sync("SELECT address FROM unit_witnesses WHERE unit=? ORDER BY address", DBParamsT{ unit })
	  		// << flattened continuation for conn.query:129:5
	  		if len(rows) > 0 {
	  			objUnit.witnesses = // .. not flattening for Array.map
	  			rows.map(func (row rowT) {*returns*} {
	  				return row.address
	  			})
	  		}
	  		// :: flattened return for callback();
	  		// ** need 1 return(s) instead of 0
	  		return 
	  		// >> flattened continuation for conn.query:129:5
	  	},
	  	func () ErrorT {
	  		// earned_headers_commission_recipients
	  		if bVoided {
	  			// :: flattened return for return callback();
	  			// ** need 1 return(s) instead of 0
	  			return 
	  		}
	  		rows := /* await */
	  		conn.query_sync("SELECT address, earned_headers_commission_share FROM earned_headers_commission_recipients " +
	  			"						WHERE unit=? ORDER BY address", DBParamsT{ unit })
	  		// << flattened continuation for conn.query:138:5
	  		if len(rows) > 0 {
	  			objUnit.earned_headers_commission_recipients = rows
	  		}
	  		// :: flattened return for callback();
	  		// ** need 1 return(s) instead of 0
	  		return 
	  		// >> flattened continuation for conn.query:138:5
	  	},
	  	func () ErrorT {
	  		rows := /* await */
	  		conn.query_sync("SELECT address, definition_chash FROM unit_authors WHERE unit=? ORDER BY address", DBParamsT{ unit })
	  		// << flattened continuation for conn.query:149:5
	  		objUnit.authors = {*ArrayExpression*}
	  		(func () ErrorT {
	  		  // :: inlined async.eachSeries:151:6
	  		  for row := range rows {
	  		    _err := (func (row rowT) ErrorT {
	  		    	var(
	  		    		onAuthorDone func () 
	  		    	)
	  		    	
	  		    	author := [*ObjectExpression*]
	  		    	
	  		    	onAuthorDone = func ()  {
	  		    		objUnit.authors = append(objUnit.authors, author)
	  		    		// :: flattened return for cb();
	  		    		// ** need 1 return(s) instead of 0
	  		    		return 
	  		    	}
	  		    	
	  		    	if bVoided {
	  		    		onAuthorDone()
	  		    		return
	  		    	}
	  		    	author.authentifiers = [*ObjectExpression*]
	  		    	sig_rows := /* await */
	  		    	conn.query_sync("SELECT path, authentifier FROM authentifiers WHERE unit=? AND address=?", DBParamsT{
	  		    		unit,
	  		    		author.address,
	  		    	})
	  		    	// << flattened continuation for conn.query:164:8
	  		    	for i := 0; i < len(sig_rows); i++ {
	  		    		author.authentifiers[sig_rows[i].path] = sig_rows[i].authentifier
	  		    	}
	  		    	
	  		    	// if definition_chash is defined:
	  		    	if row.definition_chash {
	  		    		readDefinition(conn, row.definition_chash, [*ObjectExpression*])
	  		    	} else {
	  		    		onAuthorDone()
	  		    	}
	  		    	// >> flattened continuation for conn.query:164:8
	  		    })(row)
	  		    if _err != nil { return _err }
	  		  }
	  		  return nil
	  		})()
	  		// << flattened continuation for async.eachSeries:151:6
	  		// :: flattened return for callback();
	  		// ** need 1 return(s) instead of 0
	  		return 
	  		// >> flattened continuation for async.eachSeries:151:6
	  		// >> flattened continuation for conn.query:149:5
	  	},
	  	func () ErrorT {
	  		// messages
	  		if bVoided {
	  			// :: flattened return for return callback();
	  			// ** need 1 return(s) instead of 0
	  			return 
	  		}
	  		rows := /* await */
	  		conn.query_sync("SELECT app, payload_hash, payload_location, payload, payload_uri, payload_uri_hash, message_index \n" +
	  			"						FROM messages WHERE unit=? ORDER BY message_index", DBParamsT{ unit })
	  		// << flattened continuation for conn.query:197:5
	  		if len(rows) == 0 {
	  			if conf.bLight {
	  				_core.Throw([*NewExpression*])
	  			}
	  			// :: flattened return for return callback();
	  			// ** need 1 return(s) instead of 0
	  			return 
	  		}
	  		objUnit.messages = {*ArrayExpression*}
	  		// :: flattened return for callback(async.eachSeries(rows, function (row) {
	  		return (func () ErrorT {
	  		  // :: inlined async.eachSeries:207:7
	  		  for row := range rows {
	  		    _err := (func (row rowT) ErrorT {
	  		    	var(
	  		    		addSpendProofs func () 
	  		    	)
	  		    	
	  		    	objMessage := row
	  		    	message_index := row.message_index
	  		    	objMessage.message_index = nil
	  		    	objectHash.cleanNulls(objMessage)
	  		    	objUnit.messages = append(objUnit.messages, objMessage)
	  		    	
	  		    	addSpendProofs = func ()  {
	  		    		proof_rows := /* await */
	  		    		conn.query_sync("SELECT spend_proof, address FROM spend_proofs WHERE unit=? AND message_index=? ORDER BY spend_proof_index", DBParamsT{
	  		    			unit,
	  		    			message_index,
	  		    		})
	  		    		// << flattened continuation for conn.query:217:10
	  		    		if len(proof_rows) == 0 {
	  		    			// :: flattened return for return cb();
	  		    			// ** need 1 return(s) instead of 0
	  		    			return 
	  		    		}
	  		    		objMessage.spend_proofs = {*ArrayExpression*}
	  		    		for i := 0; i < len(proof_rows); i++ {
	  		    			objSpendProof := proof_rows[i]
	  		    			if len(objUnit.authors) == 1 {
	  		    				// single-authored
	  		    				objSpendProof.address = nil
	  		    			}
	  		    			objMessage.spend_proofs = append(objMessage.spend_proofs, objSpendProof)
	  		    		}
	  		    		// :: flattened return for cb();
	  		    		// ** need 1 return(s) instead of 0
	  		    		return 
	  		    		// >> flattened continuation for conn.query:217:10
	  		    	}
	  		    	
	  		    	if objMessage.payload_location != "inline" {
	  		    		addSpendProofs()
	  		    		return
	  		    	}
	  		    	[*SwitchStatement*]
	  		    })(row)
	  		    if _err != nil { return _err }
	  		  }
	  		  return nil
	  		})()
	  		// >> flattened continuation for conn.query:197:5
	  	},
	  } {
	    if _err := _f() ; _err != nil { return _err }
	  }
	  return nil
	})()
	// << flattened continuation for async.series:90:3
	//profiler.stop('read');
	// verify unit hash. Might fail if the unit was archived while reading, in this case retry
	// light wallets don't have last_ball, don't verify their hashes
	if ! conf.bLight && ! isCorrectHash(objUnit, unit) {
		if bRetrying {
			_core.Throw("unit hash verification failed, unit: " + unit + ", objUnit: " + JSON.stringify(objUnit))
		}
		console.log("unit hash verification failed, will retry")
		return setTimeout(func () {*returns*} {
			readJointDirectly(conn, unit, callbacks, true)
		}, 60 * 1000)
	}
	if ! conf.bSaveJointJson || ! bStable || bFinalBad && bRetrievable || bRetrievable {
		callbacks.ifFound(objJoint)
		return
	}
	/* await */
	conn.query_sync("INSERT " + db.getIgnore() + " INTO joints (unit, json) VALUES (?,?)", DBParamsT{
		unit,
		JSON.stringify(objJoint),
	})
	// << flattened continuation for conn.query:494:4
	callbacks.ifFound(objJoint)
	// >> flattened continuation for conn.query:494:4
	// >> flattened continuation for async.series:90:3
	// >> flattened continuation for conn.query:51:1
}
`


const _isCorrectHash = `
func isCorrectHash(objUnit UnitT, unit UnitT)  {
	[*TryStatement*]
}
`


const _readJointWithBall = `
// add .ball even if it is not retrievable
func readJointWithBall(conn DBConnT, unit UnitT, handleJoint handleJointT)  {
	readJoint(conn, unit, [*ObjectExpression*])
}
`



//func readWitnessList_sync(conn DBConnT, unit UnitT, bAllowEmptyList bool) AddressesT {
func ReadWitnessList_sync(conn refDBConnT, unit UnitT, bAllowEmptyList bool) AddressesT {
//	arrWitnesses := assocCachedUnitWitnesses[unit]
	arrWitnesses := AssocCachedUnitWitnesses[unit]
//	if arrWitnesses {
	if arrWitnesses != nil {
		// :: flattened return for return handleWitnessList(arrWitnesses);
		return arrWitnesses
	}
	rcvr := db.AddressesReceiver{}
/**
	rows := /* await * /
	conn.query_sync("SELECT address FROM unit_witnesses WHERE unit=? ORDER BY address", DBParamsT{ unit })
 **/
	conn.MustQuery("SELECT address FROM unit_witnesses WHERE unit=? ORDER BY address", DBParamsT{ unit }, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:537:1
	if ! bAllowEmptyList && len(rows) == 0 {
//		_core.Throw("witness list of unit " + unit + " not found")
		_core.Throw("witness list of unit %s not found", unit)
	}
	if len(rows) > 0 && len(rows) != constants.COUNT_WITNESSES {
//		_core.Throw("wrong number of witnesses in unit " + unit)
		_core.Throw("wrong number of witnesses in unit %s", unit)
	}
/**
	// [tbd] map is not converted for AssignStatement
	arrWitnesses = // .. not flattening for Array.map
	rows.map(func (row rowT) {*returns*} {
		return row.address
	})
 **/
	arrWitnesses = make(AddressesT, 0, len(rows))
	for _, row := range rows {
		arrWitnesses = append(arrWitnesses, row.Address)
	}
	if len(rows) > 0 {
//		assocCachedUnitWitnesses[unit] = arrWitnesses
		AssocCachedUnitWitnesses[unit] = arrWitnesses
	}
	// :: flattened return for handleWitnessList(arrWitnesses);
	return arrWitnesses
	// >> flattened continuation for conn.query:537:1
}

//func readWitnesses_sync(conn DBConnT, unit UnitT) AddressesT {
func ReadWitnesses_sync(conn refDBConnT, unit UnitT) AddressesT {
//	arrWitnesses := assocCachedUnitWitnesses[unit]
	arrWitnesses := AssocCachedUnitWitnesses[unit]
//	if arrWitnesses {
	if arrWitnesses != nil {
		// :: flattened return for return handleWitnessList(arrWitnesses);
		return arrWitnesses
	}
/**
	rows := /* await * /
	conn.query_sync("SELECT witness_list_unit FROM units WHERE unit=?", DBParamsT{ unit })
 **/
	rcvr := db.WitnessListUnitsReceiver{}
	conn.MustQuery("SELECT witness_list_unit FROM units WHERE unit=?", DBParamsT{ unit }, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:553:1
	if len(rows) == 0 {
//		_core.Throw("unit " + unit + " not found")
		_core.Throw("unit %s not found", unit)
	}
//	witness_list_unit := rows[0].witness_list_unit
	witness_list_unit := rows[0].Witness_list_unit
	witness_list_unit_ := witness_list_unit
//	if witness_list_unit_ == nil {
	if witness_list_unit_.IsNull() {
		witness_list_unit_ = unit
	}
	// [tbd] flattening/inlining: rename variables or add scope when necessary
	{{
	arrWitnesses := /* await */
//	readWitnessList_sync(conn, witness_list_unit_)
	ReadWitnessList_sync(conn, witness_list_unit_, false)
	// << flattened continuation for readWitnessList:559:2
//	assocCachedUnitWitnesses[unit] = arrWitnesses
	AssocCachedUnitWitnesses[unit] = arrWitnesses
	// :: flattened return for handleWitnessList(arrWitnesses);
	return arrWitnesses
	}}
	// >> flattened continuation for readWitnessList:559:2
	// >> flattened continuation for conn.query:553:1
}

const _determineIfWitnessAddressDefinitionsHaveReferences_sync = `
func determineIfWitnessAddressDefinitionsHaveReferences_sync(conn DBConnT, arrWitnesses AddressesT) bool {
	rows := /* await */
	conn.query_sync("SELECT 1 FROM address_definition_changes JOIN definitions USING(definition_chash) \n" +
		"		WHERE address IN(?) AND has_references=1 \n" +
		"		UNION \n" +
		"		SELECT 1 FROM definitions WHERE definition_chash IN(?) AND has_references=1 \n" +
		"		LIMIT 1", DBParamsT{
		arrWitnesses,
		arrWitnesses,
	})
	// << flattened continuation for conn.query:567:1
	// :: flattened return for handleResult(rows.length > 0);
	return len(rows) > 0
	// >> flattened continuation for conn.query:567:1
}

func determineWitnessedLevelAndBestParent_sync(conn DBConnT, arrParentUnits UnitsT, arrWitnesses AddressesT) (int, UnitT) {
	var(
		addWitnessesAndGoUp func (start_unit UnitT) (int, UnitT)
	)
	
	arrCollectedWitnesses := AddressesT{}
	my_best_parent_unit := {*init:null*}
	
	addWitnessesAndGoUp = func (start_unit UnitT) (int, UnitT) {
		props := /* await */
		readStaticUnitProps_sync(conn, start_unit)
		// << flattened continuation for readStaticUnitProps:585:2
		best_parent_unit := props.best_parent_unit
		level := props.level
		if level == nil {
			_core.Throw("null level in updateWitnessedLevel")
		}
		if level == 0 {
			// genesis
			// :: flattened return for return handleWitnessedLevelAndBestParent(0, my_best_parent_unit);
			return meta.returnArguments(0, my_best_parent_unit)
		}
		arrAuthors := /* await */
		readUnitAuthors_sync(conn, start_unit)
		// << flattened continuation for readUnitAuthors:592:3
		for i := 0; i < len(arrAuthors); i++ {
			address := arrAuthors[i]
			if arrWitnesses.indexOf(address) != - 1 && arrCollectedWitnesses.indexOf(address) == - 1 {
				arrCollectedWitnesses = append(arrCollectedWitnesses, address)
			}
		}
		if len(arrCollectedWitnesses) >= constants.MAJORITY_OF_WITNESSES {
			// :: flattened return for return handleWitnessedLevelAndBestParent(level, my_best_parent_unit);
			return meta.returnArguments(level, my_best_parent_unit)
		}
		return addWitnessesAndGoUp(best_parent_unit)
		// >> flattened continuation for readUnitAuthors:592:3
		// >> flattened continuation for readStaticUnitProps:585:2
	}
	best_parent_unit := /* await */
	determineBestParent_sync(conn, [*ObjectExpression*], arrWitnesses)
	// << flattened continuation for determineBestParent:605:1
	if ! best_parent_unit {
		_core.Throw("no best parent of " + arrParentUnits.join(", ") + ", witnesses " + arrWitnesses.join(", "))
	}
	my_best_parent_unit = best_parent_unit
	return addWitnessesAndGoUp(best_parent_unit)
	// >> flattened continuation for determineBestParent:605:1
}
`


/*
function readWitnessesOnMcUnit(conn, main_chain_index, handleWitnesses){
	conn.query( // we read witnesses from MC unit (users can cheat with side-chains)
		"SELECT address \n\
		FROM units \n\
		JOIN unit_witnesses ON(units.unit=unit_witnesses.unit OR units.witness_list_unit=unit_witnesses.unit) \n\
		WHERE main_chain_index=? AND is_on_main_chain=1", 
		[main_chain_index],
		function(witness_rows){
			if (witness_rows.length === 0)
				throw "no witness list on MC unit "+main_chain_index;
			if (witness_rows.length !== constants.COUNT_WITNESSES)
				throw "wrong number of witnesses on MC unit "+main_chain_index;
			var arrWitnesses = witness_rows.map(function(witness_row){ return witness_row.address; });
			handleWitnesses(arrWitnesses);
		}
	);
}
 */


const _readDefinitionByAddress = `
// max_mci must be stable
func readDefinitionByAddress(conn DBConnT, address AddressT, max_mci MCIndexT, callbacks callbacksT)  {
	if max_mci == nil {
		max_mci = MAX_INT32
	}
	rows := /* await */
	conn.query_sync("SELECT definition_chash FROM address_definition_changes CROSS JOIN units USING(unit) \n" +
		"		WHERE address=? AND is_stable=1 AND sequence='good' AND main_chain_index<=? ORDER BY level DESC LIMIT 1", DBParamsT{
		address,
		max_mci,
	})
	// << flattened continuation for conn.query:639:1
	definition_chash := address
	if len(rows) > 0 {
		definition_chash = rows[0].definition_chash
	}
	readDefinitionAtMci(conn, definition_chash, max_mci, callbacks)
	// >> flattened continuation for conn.query:639:1
}

// max_mci must be stable
func readDefinitionAtMci(conn DBConnT, definition_chash definition_chashT, max_mci MCIndexT, callbacks callbacksT)  {
	sql := "SELECT definition FROM definitions CROSS JOIN unit_authors USING(definition_chash) CROSS JOIN units USING(unit) \n" +
		"		WHERE definition_chash=? AND is_stable=1 AND sequence='good' AND main_chain_index<=?"
	params := {*ArrayExpression*}
	rows := /* await */
	conn.query_sync(sql, params)
	// << flattened continuation for conn.query:656:1
	if len(rows) == 0 {
		callbacks.ifDefinitionNotFound(definition_chash)
		return
	}
	callbacks.ifFound(JSON.parse(rows[0].definition))
	// >> flattened continuation for conn.query:656:1
}

func readDefinition(conn DBConnT, definition_chash definition_chashT, callbacks callbacksT)  {
	rows := /* await */
	conn.query_sync("SELECT definition FROM definitions WHERE definition_chash=?", DBParamsT{ definition_chash })
	// << flattened continuation for conn.query:664:1
	if len(rows) == 0 {
		callbacks.ifDefinitionNotFound(definition_chash)
		return
	}
	callbacks.ifFound(JSON.parse(rows[0].definition))
	// >> flattened continuation for conn.query:664:1
}

func readFreeJoints_sync(ifFoundFreeBall BallT) ErrorT {
	rows := /* await */
	db.query_sync("SELECT units.unit FROM units LEFT JOIN archived_joints USING(unit) WHERE is_free=1 AND archived_joints.unit IS NULL")
	// << flattened continuation for db.query:672:1
	// :: flattened return for onDone(async.each(rows, function (row) {
	return (func () ErrorT {
	  // :: inlined async.each:673:2 !! [tbd] finish this
	  for row := range rows {
	    _err := (func (row rowT) ErrorT {
	    	readJoint(db, row.unit, [*ObjectExpression*])
	    })(row)
	    if _err != nil { return _err }
	  }
	  return nil
	})()
	// >> flattened continuation for db.query:672:1
}
`

//func isGenesisUnit(unit UnitT) bool {
func IsGenesisUnit(unit UnitT) bool {
	return unit == constants.GENESIS_UNIT
}

const _isGenesisBall = `
func isGenesisBall(ball BallT) bool {
	return ball == genesis_ball
}
`






//func readUnitProps_sync(conn DBConnT, unit UnitT) PropsT {
func ReadUnitProps_sync(conn refDBConnT, unit UnitT) refPropsT {
//	if assocStableUnits[unit] {
	if AssocStableUnits[unit] != nil {
		// :: flattened return for return handleProps(assocStableUnits[unit]);
//		return assocStableUnits[unit]
		return AssocStableUnits[unit]
	}
/**
	rows := /* await * /
	conn.query_sync("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain, is_free, is_stable, witnessed_level FROM units WHERE unit=?", DBParamsT{ unit })
 **/
	rcvr := db.RUPUnitPropsReceiver{}
	conn.MustQuery("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain, is_free, is_stable, witnessed_level FROM units WHERE unit=?", DBParamsT{ unit }, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:703:1
	if len(rows) != 1 {
		_core.Throw("not 1 row")
	}
	props := rows[0]
//	if props.is_stable {
	if props.Is_stable != 0 {
//		assocStableUnits[unit] = props
		AssocStableUnits[unit] = &props.PropsT
	} else {
//		props2 := _.cloneDeep(assocUnstableUnits[unit])
		// [fyi] not cloning
		props2 := AssocUnstableUnits[unit]
//		if ! props2 {
		if props2 == nil {
//		if false {
//			_core.Throw("no unstable props of " + unit)
			_core.Throw("no unstable props of %s", unit)
		}
//		props2.parent_units = nil
		// [fyi] props2 is not cloned - do not modify original
		//props2.Parent_units = nil
//		if ! _.isEqual(props, props2) {
		if ! (props.PropsT == props2.PropsT) {
//			_core.Throw("different props of " + unit + ", mem: " + JSON.stringify(props2) + ", db: " + JSON.stringify(props))
			_core.Throw("different props of %s, mem: %s, db: %s", unit, JSON.Stringify(props2.PropsT), JSON.Stringify(props.PropsT))
		}
	}
	// :: flattened return for handleProps(props);
	return &props.PropsT
	// >> flattened continuation for conn.query:703:1
}

const _readPropsOfUnits_sync = `
func readPropsOfUnits_sync(conn DBConnT, earlier_unit UnitT, arrLaterUnits UnitsT) (PropsT, PropssT) {
	bEarlierInLaterUnits := arrLaterUnits.indexOf(earlier_unit) != - 1
	rows := /* await */
	conn.query_sync("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain, is_free FROM units WHERE unit IN(?, ?)", DBParamsT{
		earlier_unit,
		arrLaterUnits,
	})
	// << flattened continuation for conn.query:727:1
	k := 0
	if ! bEarlierInLaterUnits {
		k = 1
	}
	if len(rows) != len(arrLaterUnits) + k {
		_core.Throw("wrong number of rows for earlier " + earlier_unit + ", later " + arrLaterUnits)
	}
	objEarlierUnitProps := {*init:null*}
	arrLaterUnitProps := PropssT{}
	for i := 0; i < len(rows); i++ {
		if rows[i].unit == earlier_unit {
			objEarlierUnitProps = rows[i]
		} else {
			arrLaterUnitProps = append(arrLaterUnitProps, rows[i])
		}
	}
	if bEarlierInLaterUnits {
		arrLaterUnitProps = append(arrLaterUnitProps, objEarlierUnitProps)
	}
	// :: flattened return for handleProps(objEarlierUnitProps, arrLaterUnitProps);
	return meta.returnArguments(objEarlierUnitProps, arrLaterUnitProps)
	// >> flattened continuation for conn.query:727:1
}
`
//func ReadPropsOfUnits_sync(conn refDBConnT, unit UnitT, units UnitsT) (refPropsT, refPropssT) {
func ReadPropsOfUnits_sync(conn refDBConnT, unit UnitT, units UnitsT) (refPropsT, PropssT) {
	panic("[tbd] ReadPropsOfUnits_sync")
}






//func readLastStableMcUnitProps_sync(conn DBConnT) PropsT {
func readLastStableMcUnitProps_sync(conn refDBConnT) refPropsT {
/**
	rows := /* await * /
	conn.query_sync("SELECT units.*, ball FROM units LEFT JOIN balls USING(unit) WHERE is_on_main_chain=1 AND is_stable=1 ORDER BY main_chain_index DESC LIMIT 1")
 **/
	rcvr := db.UnitContentBallsReceiver{}
	conn.MustQuery("SELECT units.*, ball FROM units LEFT JOIN balls USING(unit) WHERE is_on_main_chain=1 AND is_stable=1 ORDER BY main_chain_index DESC LIMIT 1", DBParamsT{}, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:755:1
	if len(rows) == 0 {
		// :: flattened return for return handleLastStableMcUnitProps(null);
		return nil
	}
	// empty database
	//throw "readLastStableMcUnitProps: no units on stable MC?";
//	if ! rows[0].ball {
	if rows[0].Ball.IsNull() {
//		_core.Throw("no ball for last stable unit " + rows[0].unit)
		_core.Throw("no ball for last stable unit %s", rows[0].Unit)
	}
	// :: flattened return for handleLastStableMcUnitProps(rows[0]);
//	return rows[0]
	return &rows[0].PropsT
	// >> flattened continuation for conn.query:755:1
}

//func readLastStableMcIndex_sync(conn DBConnT) MCIndexT {
func ReadLastStableMcIndex_sync(conn refDBConnT) MCIndexT {
	objLastStableMcUnitProps := /* await */
	readLastStableMcUnitProps_sync(conn)
	// << flattened continuation for readLastStableMcUnitProps:769:1
	if objLastStableMcUnitProps != nil {
		// :: flattened return for handleLastStableMcIndex(objLastStableMcUnitProps.main_chain_index);
//		return objLastStableMcUnitProps.main_chain_index
		return objLastStableMcUnitProps.Main_chain_index
	} else {
		// :: flattened return for handleLastStableMcIndex(0);
		return 0
	}
	// >> flattened continuation for readLastStableMcUnitProps:769:1
}


const _readLastMainChainIndex = `
func readLastMainChainIndex(handleLastMcIndex MCIndexT)  {
	rows := /* await */
	db.query_sync("SELECT MAX(main_chain_index) AS last_mc_index FROM units")
	// << flattened continuation for db.query:780:1
	last_mc_index := rows[0].last_mc_index
	if last_mc_index == nil {
		// empty database
		last_mc_index = 0
	}
	handleLastMcIndex(last_mc_index)
	// >> flattened continuation for db.query:780:1
}
`


//func findLastBallMciOfMci_sync(conn DBConnT, mci MCIndexT)  {
func findLastBallMciOfMci_sync(conn refDBConnT, mci MCIndexT) MCIndexT {
	if mci == 0 {
		_core.Throw("findLastBallMciOfMci called with mci=0")
	}
/**
	rows := /* await * /
	conn.query_sync("SELECT lb_units.main_chain_index, lb_units.is_on_main_chain \n" +
		"		FROM units JOIN units AS lb_units ON units.last_ball_unit=lb_units.unit \n" +
		"		WHERE units.is_on_main_chain=1 AND units.main_chain_index=?", DBParamsT{ mci })
 **/
	rcvr := db.MCIOnMCsReceiver{}
	conn.MustQuery("SELECT lb_units.main_chain_index, lb_units.is_on_main_chain \n" +
		"FROM units JOIN units AS lb_units ON units.last_ball_unit=lb_units.unit \n" +
		"WHERE units.is_on_main_chain=1 AND units.main_chain_index=?", DBParamsT{ mci }, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:792:1
	if len(rows) != 1 {
//		_core.Throw("last ball's mci count " + len(rows) + " !== 1, mci = " + mci)
		_core.Throw("last ball's mci count %d !== 1, mci = %d", len(rows), mci)
	}
//	if rows[0].is_on_main_chain != 1 {
	if rows[0].Is_on_main_chain != 1 {
		_core.Throw("lb is not on mc?")
	}
	// :: flattened return for handleLastBallMci(rows[0].main_chain_index);
	// ** need 0 return(s) instead of 1
//	return rows[0].main_chain_index
	return rows[0].Main_chain_index
	// >> flattened continuation for conn.query:792:1
}

//func getMinRetrievableMci()  {
func getMinRetrievableMci() MCIndexT {
	return min_retrievable_mci
}

//func updateMinRetrievableMciAfterStabilizingMci_sync(conn DBConnT, last_stable_mci MCIndexT) MCIndexT {
func UpdateMinRetrievableMciAfterStabilizingMci_sync(conn refDBConnT, last_stable_mci MCIndexT) MCIndexT {
//	console.log("updateMinRetrievableMciAfterStabilizingMci " + last_stable_mci)
	console.Log("updateMinRetrievableMciAfterStabilizingMci %d", last_stable_mci)
	last_ball_mci := /* await */
	findLastBallMciOfMci_sync(conn, last_stable_mci)
	// << flattened continuation for findLastBallMciOfMci:813:1
	if last_ball_mci <= min_retrievable_mci {
		// nothing new
		// :: flattened return for return handleMinRetrievableMci(min_retrievable_mci);
		return min_retrievable_mci
	}
	prev_min_retrievable_mci := min_retrievable_mci
	min_retrievable_mci = last_ball_mci
/**
	unit_rows := /* await * /
	conn.query_sync(// 'JOIN messages' filters units that are not stripped yet
	"SELECT DISTINCT unit, content_hash FROM units " + db.forceIndex("byMcIndex") + " CROSS JOIN messages USING(unit) \n" +
		"			WHERE main_chain_index<=? AND main_chain_index>=? AND sequence='final-bad'", DBParamsT{
		min_retrievable_mci,
		prev_min_retrievable_mci,
	})
 **/
	rcvr := db.UnitContentHashsReceiver{}
	conn.MustQuery(// 'JOIN messages' filters units that are not stripped yet
		"SELECT DISTINCT unit, content_hash FROM units " + conn.ForceIndex("byMcIndex") + " CROSS JOIN messages USING(unit) \n" +
		"WHERE main_chain_index<=? AND main_chain_index>=? AND sequence='final-bad'", DBParamsT{
		min_retrievable_mci,
		prev_min_retrievable_mci,
	}, &rcvr)
	unit_rows := rcvr.Rows
	// << flattened continuation for conn.query:820:2
	arrQueries := AsyncFunctorsT{}
	(func () ErrorT {
	  // :: inlined async.eachSeries:827:4
//	  for unit_row := range unit_rows {
	  for _, unit_row := range unit_rows {
//	    _err := (func (unit_row unit_rowT) ErrorT {
	    _err := (func (unit_row db.UnitContentHashRow) ErrorT {
//	    	unit := unit_row.unit
	    	unit := unit_row.Unit
//	    	if ! unit_row.content_hash {
	    	if unit_row.Content_hash.IsNull() {
//	    		_core.Throw("no content hash in bad unit " + unit)
	    		_core.Throw("no content hash in bad unit %s", unit)
	    	}
//		readJoint(conn, unit, [*ObjectExpression*])
	    	ReadJoint(conn, unit, ReadJointCbT{

			IfNotFound: func () {
		    		_core.Throw("bad unit not found: %s", unit)
			},

			IfFound: func (objJoint refJointT) {
				archiving.GenerateQueriesToArchiveJoint_sync(conn, objJoint, "voided", &arrQueries);
			},

		})
		return nil
	    })(unit_row)
	    if _err != nil { return _err }
	  }
	  return nil
	})()
	// << flattened continuation for async.eachSeries:827:4
	if len(arrQueries) == 0 {
		// :: flattened return for return handleMinRetrievableMci(min_retrievable_mci);
		return min_retrievable_mci
	}
	(func () ErrorT {
	  // :: inlined async.series:845:6
//	  for _f := range arrQueries {
	  for _, _f := range arrQueries {
	    if _err := _f() ; _err != nil { return _err }
	  }
	  return nil
	})()
	// << flattened continuation for async.series:845:6
	// .. not flattening for Array.forEach
//	for unit_row, _ := range unit_rows {
	for _, unit_row := range unit_rows {
		forgetUnit(unit_row.Unit)
	}
	// :: flattened return for handleMinRetrievableMci(min_retrievable_mci);
	return min_retrievable_mci
	// >> flattened continuation for async.series:845:6
	// >> flattened continuation for async.eachSeries:827:4
	// >> flattened continuation for conn.query:820:2
	// >> flattened continuation for findLastBallMciOfMci:813:1
}

func initializeMinRetrievableMci()  {
/**
	rows := /* await * /
	db.query_sync("SELECT MAX(lb_units.main_chain_index) AS min_retrievable_mci \n" +
		"		FROM units JOIN units AS lb_units ON units.last_ball_unit=lb_units.unit \n" +
		"		WHERE units.is_on_main_chain=1 AND units.is_stable=1")
 **/
	rcvr := db.MinRetrievableMCIsReceiver{}
	db.MustQuery("SELECT MAX(lb_units.main_chain_index) AS min_retrievable_mci \n" +
		"FROM units JOIN units AS lb_units ON units.last_ball_unit=lb_units.unit \n" +
		"WHERE units.is_on_main_chain=1 AND units.is_stable=1", DBParamsT{}, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for db.query:859:1
	if len(rows) != 1 {
		_core.Throw("MAX() no rows?")
	}
//	min_retrievable_mci = rows[0].min_retrievable_mci
	min_retrievable_mci = rows[0].Min_retrievable_mci
//	if min_retrievable_mci == nil {
	if min_retrievable_mci.IsNull() {
		min_retrievable_mci = 0
	}
	// >> flattened continuation for db.query:859:1
}


const _archiveJointAndDescendantsIfExists = `
func archiveJointAndDescendantsIfExists(from_unit UnitT)  {
	console.log("will archive if exists from unit " + from_unit)
	rows := /* await */
	db.query_sync("SELECT 1 FROM units WHERE unit=?", {*ArrayExpression*})
	// << flattened continuation for db.query:876:1
	if len(rows) > 0 {
		archiveJointAndDescendants(from_unit)
	}
	// >> flattened continuation for db.query:876:1
}

func archiveJointAndDescendants(from_unit UnitT)  {
	/* await */
	db.executeInTransaction_sync(func (conn DBConnT, cb cbT) {*returns*} {
		var(
			addChildren func (arrParentUnits UnitsT) 
			archive func () 
		)
		
		arrUnits := UnitsT{ from_unit }
		
		addChildren = func (arrParentUnits UnitsT)  {
			rows := /* await */
			conn.query_sync("SELECT DISTINCT child_unit FROM parenthoods WHERE parent_unit IN(?)", DBParamsT{ arrParentUnits })
			// << flattened continuation for conn.query:887:3
			if len(rows) == 0 {
				archive()
				return
			}
			arrChildUnits := make(UnitsT, len(rows), len(rows))
			for row, _k := range rows {
				arrChildUnits[_k] := row.child_unit
			}
			arrUnits = arrUnits.concat(arrChildUnits)
			addChildren(arrChildUnits)
			// >> flattened continuation for conn.query:887:3
		}
		
		archive = func ()  {
			arrUnits = _.uniq(arrUnits)
			// does not affect the order
			arrUnits.reverse()
			console.log("will archive", arrUnits)
			arrQueries := AsyncFunctorsT{}
			(func () ErrorT {
			  // :: inlined async.eachSeries:901:3
			  for unit := range arrUnits {
			    _err := (func (unit UnitT) ErrorT {
			    	readJoint(conn, unit, [*ObjectExpression*])
			    })(unit)
			    if _err != nil { return _err }
			  }
			  return nil
			})()
			// << flattened continuation for async.eachSeries:901:3
			conn.addQuery(arrQueries, "DELETE FROM known_bad_joints")
			console.log("will execute " + len(arrQueries) + " queries to archive")
			(func () ErrorT {
			  // :: inlined async.series:916:5
			  for _f := range arrQueries {
			    if _err := _f() ; _err != nil { return _err }
			  }
			  return nil
			})()
			// << flattened continuation for async.series:916:5
			for v, _ := range arrUnits {
				forgetUnit(v)// !! remove this
			}
			cb()
			// >> flattened continuation for async.series:916:5
			// >> flattened continuation for async.eachSeries:901:3
		}
		
		console.log("will archive from unit " + from_unit)
		addChildren(arrUnits)
	})
	// << flattened continuation for db.executeInTransaction:883:1
	console.log("done archiving from unit " + from_unit)
	// >> flattened continuation for db.executeInTransaction:883:1
}
`


//_______________________________________________________________________________________________
// Assets

const _readAssetInfo_sync = `
func readAssetInfo_sync(conn DBConnT, asset AssetT) AssetT {
	objAsset := assocCachedAssetInfos[asset]
	if objAsset {
		// :: flattened return for return handleAssetInfo(objAsset);
		return objAsset
	}
	rows := /* await */
	conn.query_sync("SELECT assets.*, main_chain_index, sequence, is_stable, address AS definer_address, unit AS asset \n" +
		"		FROM assets JOIN units USING(unit) JOIN unit_authors USING(unit) WHERE unit=?", DBParamsT{ asset })
	// << flattened continuation for conn.query:940:1
	if len(rows) > 1 {
		_core.Throw("more than one asset?")
	}
	if len(rows) == 0 {
		// :: flattened return for return handleAssetInfo(null);
		return nil
	}
	objAsset := rows[0]
	if objAsset.issue_condition {
		objAsset.issue_condition = JSON.parse(objAsset.issue_condition)
	}
	if objAsset.transfer_condition {
		objAsset.transfer_condition = JSON.parse(objAsset.transfer_condition)
	}
	if objAsset.is_stable {
		// cache only if stable
		assocCachedAssetInfos[asset] = objAsset
	}
	// :: flattened return for handleAssetInfo(objAsset);
	return objAsset
	// >> flattened continuation for conn.query:940:1
}

func readAsset_sync(conn DBConnT, asset AssetT, last_ball_mci MCIndexT) (ErrorT, AssetT) {
	if last_ball_mci == nil {
		if conf.bLight {
			last_ball_mci = MAX_INT32
		} else {
			last_stable_mci := /* await */
			readLastStableMcIndex_sync(conn)
			// << flattened continuation for readLastStableMcIndex:966:10
			// :: flattened return for handleAsset(readAsset(conn, asset, last_stable_mci));
			// ** need 2 return(s) instead of 1
			return /* await */
			readAsset_sync(conn, asset, last_stable_mci)
			// >> flattened continuation for readLastStableMcIndex:966:10
		}
	}
	objAsset := /* await */
	readAssetInfo_sync(conn, asset)
	// << flattened continuation for readAssetInfo:970:1
	if ! objAsset {
		// :: flattened return for return handleAsset('asset ' + asset + ' not found');
		// ** need 2 return(s) instead of 1
		return "asset " + asset + " not found"
	}
	if objAsset.main_chain_index > last_ball_mci {
		// :: flattened return for return handleAsset('asset definition must be before last ball');
		// ** need 2 return(s) instead of 1
		return "asset definition must be before last ball"
	}
	if objAsset.sequence != "good" {
		// :: flattened return for return handleAsset('asset definition is not serial');
		// ** need 2 return(s) instead of 1
		return "asset definition is not serial"
	}
	if ! objAsset.spender_attested {
		// :: flattened return for return handleAsset(null, objAsset);
		return meta.returnArguments(nil, objAsset)
	}
	latest_rows := /* await */
	conn.query_sync("SELECT MAX(level) AS max_level FROM asset_attestors CROSS JOIN units USING(unit) \n" +
		"			WHERE asset=? AND main_chain_index<=? AND is_stable=1 AND sequence='good'", DBParamsT{
		asset,
		last_ball_mci,
	})
	// << flattened continuation for conn.query:981:2
	max_level := latest_rows[0].max_level
	if ! max_level {
		_core.Throw("no max level of asset attestors")
	}
	att_rows := /* await */
	conn.query_sync("SELECT attestor_address FROM asset_attestors CROSS JOIN units USING(unit) \n" +
		"					WHERE asset=? AND level=? AND main_chain_index<=? AND is_stable=1 AND sequence='good'", DBParamsT{
		asset,
		max_level,
		last_ball_mci,
	})
	// << flattened continuation for conn.query:991:4
	if len(att_rows) == 0 {
		_core.Throw("no attestors?")
	}
	objAsset.arrAttestorAddresses = // .. not flattening for Array.map
	att_rows.map(func (att_row att_rowT) {*returns*} {
		return att_row.attestor_address
	})
	// :: flattened return for handleAsset(null, objAsset);
	return meta.returnArguments(nil, objAsset)
	// >> flattened continuation for conn.query:991:4
	// >> flattened continuation for conn.query:981:2
	// >> flattened continuation for readAssetInfo:970:1
}

// filter only those addresses that are attested (doesn't work for light clients)
func filterAttestedAddresses_sync(conn DBConnT, objAsset AssetT, last_ball_mci MCIndexT, arrAddresses AddressesT) AddressesT {
	addr_rows := /* await */
	conn.query_sync("SELECT DISTINCT address FROM attestations CROSS JOIN units USING(unit) \n" +
		"		WHERE attestor_address IN(?) AND address IN(?) AND main_chain_index<=? AND is_stable=1 AND sequence='good' \n" +
		"			AND main_chain_index>IFNULL( \n" +
		"				(SELECT main_chain_index FROM address_definition_changes JOIN units USING(unit) \n" +
		"				WHERE address_definition_changes.address=attestations.address ORDER BY main_chain_index DESC LIMIT 1), \n" +
		"			0)", DBParamsT{
		objAsset.arrAttestorAddresses,
		arrAddresses,
		last_ball_mci,
	})
	// << flattened continuation for conn.query:1009:1
	arrAttestedAddresses := make(AddressesT, len(addr_rows), len(addr_rows))
	for addr_row, _k := range addr_rows {
		arrAttestedAddresses[_k] := addr_row.address
	}
	// :: flattened return for handleAttestedAddresses(arrAttestedAddresses);
	return arrAttestedAddresses
	// >> flattened continuation for conn.query:1009:1
}

// note that light clients cannot check attestations
func loadAssetWithListOfAttestedAuthors_sync(conn DBConnT, asset AssetT, last_ball_mci MCIndexT, arrAuthorAddresses AddressesT) (ErrorT, AssetT) {
	( err, objAsset ) := /* await */
	readAsset_sync(conn, asset, last_ball_mci)
	// << flattened continuation for readAsset:1026:1
	if err {
		// :: flattened return for return handleAsset(err);
		// ** need 2 return(s) instead of 1
		return err
	}
	if ! objAsset.spender_attested {
		// :: flattened return for return handleAsset(null, objAsset);
		return meta.returnArguments(nil, objAsset)
	}
	arrAttestedAddresses := /* await */
	filterAttestedAddresses_sync(conn, objAsset, last_ball_mci, arrAuthorAddresses)
	// << flattened continuation for filterAttestedAddresses:1031:2
	objAsset.arrAttestedAddresses = arrAttestedAddresses
	// :: flattened return for handleAsset(null, objAsset);
	return meta.returnArguments(nil, objAsset)
	// >> flattened continuation for filterAttestedAddresses:1031:2
	// >> flattened continuation for readAsset:1026:1
}

func findWitnessListUnit_sync(conn DBConnT, arrWitnesses AddressesT, last_ball_mci MCIndexT) UnitT {
	rows := /* await */
	conn.query_sync("SELECT witness_list_hashes.witness_list_unit \n" +
		"		FROM witness_list_hashes CROSS JOIN units ON witness_list_hashes.witness_list_unit=unit \n" +
		"		WHERE witness_list_hash=? AND sequence='good' AND is_stable=1 AND main_chain_index<=?", DBParamsT{
		objectHash.getBase64Hash(arrWitnesses),
		last_ball_mci,
	})
	// << flattened continuation for conn.query:1039:1
	if len(rows) != 0 {
		// :: flattened return for handleWitnessListUnit(rows[0].witness_list_unit);
		return rows[0].witness_list_unit
	} else {
		// :: flattened return for handleWitnessListUnit(0);
		return 0
	}
	// >> flattened continuation for conn.query:1039:1
}

func sliceAndExecuteQuery_sync(query queryT, params paramsT, largeParam largeParamT)  {
	if typeof largeParam != "object" || len(largeParam) == 0 {
		// :: flattened return for return callback([]);
		// ** need 0 return(s) instead of 1
		return {*ArrayExpression*}
	}
	CHUNK_SIZE := 200
	length := len(largeParam)
	arrParams := {*ArrayExpression*}
	newParams := {*init:null*}
	largeParamPosition := params.indexOf(largeParam)
	
	for offset := 0; offset < length; offset += CHUNK_SIZE {
		newParams = params.slice(0)
		newParams[largeParamPosition] = largeParam.slice(offset, offset + CHUNK_SIZE)
		arrParams = append(arrParams, newParams)
	}
	
	result := {*ArrayExpression*}
	(func () ErrorT {
	  // :: inlined async.eachSeries:1069:1
	  for params := range arrParams {
	    _err := (func (params paramsT) ErrorT {
	    	rows := /* await */
	    	db.query_sync(query, params)
	    	// << flattened continuation for db.query:1070:2
	    	result = result.concat(rows)
	    	// :: flattened return for cb();
	    	// ** need 1 return(s) instead of 0
	    	return 
	    	// >> flattened continuation for db.query:1070:2
	    })(params)
	    if _err != nil { return _err }
	  }
	  return nil
	})()
	// << flattened continuation for async.eachSeries:1069:1
	// :: flattened return for callback(result);
	// ** need 0 return(s) instead of 1
	return result
	// >> flattened continuation for async.eachSeries:1069:1
}

func filterNewOrUnstableUnits_sync(arrUnits UnitsT) UnitsT {
	rows := /* await */
	sliceAndExecuteQuery_sync("SELECT unit FROM units WHERE unit IN(?) AND is_stable=1", {*ArrayExpression*}, arrUnits)
	// << flattened continuation for sliceAndExecuteQuery:1080:1
	arrKnownStableUnits := make(UnitsT, len(rows), len(rows))
	for row, _k := range rows {
		arrKnownStableUnits[_k] := row.unit
	}
	arrNewOrUnstableUnits := _.difference(arrUnits, arrKnownStableUnits)
	// :: flattened return for handleFilteredUnits(arrNewOrUnstableUnits);
	return arrNewOrUnstableUnits
	// >> flattened continuation for sliceAndExecuteQuery:1080:1
}

// for unit that is not saved to the db yet
func determineBestParent_sync(conn DBConnT, objUnit UnitT, arrWitnesses AddressesT) UnitT {
	rows := /* await */
	conn.query_sync("SELECT unit \n" +
		"		FROM units AS parent_units \n" +
		"		WHERE unit IN(?) \n" +
		"			AND (witness_list_unit=? OR ( \n" +
		"				SELECT COUNT(*) \n" +
		"				FROM unit_witnesses AS parent_witnesses \n" +
		"				WHERE parent_witnesses.unit IN(parent_units.unit, parent_units.witness_list_unit) AND address IN(?) \n" +
		"			)>=?) \n" +
		"		ORDER BY witnessed_level DESC, \n" +
		"			level-witnessed_level ASC, \n" +
		"			unit ASC \n" +
		"		LIMIT 1", DBParamsT{
		objUnit.parent_units,
		objUnit.witness_list_unit,
		arrWitnesses,
		constants.COUNT_WITNESSES - constants.MAX_WITNESS_LIST_MUTATIONS,
	})
	// << flattened continuation for conn.query:1090:1
	if len(rows) != 1 {
		// :: flattened return for return handleBestParent(null);
		return nil
	}
	best_parent_unit := rows[0].unit
	// :: flattened return for handleBestParent(best_parent_unit);
	return best_parent_unit
	// >> flattened continuation for conn.query:1090:1
}

func determineIfHasWitnessListMutationsAlongMc_sync(conn DBConnT, objUnit UnitT, last_ball_unit UnitT, arrWitnesses AddressesT) ErrorT {
	if ! objUnit.parent_units {
		// genesis
		// :: flattened return for return handleResult();
		// ** need 1 return(s) instead of 0
		return 
	}
	( bHasBestParent, arrMcUnits ) := /* await */
	buildListOfMcUnitsWithPotentiallyDifferentWitnesslists_sync(conn, objUnit, last_ball_unit, arrWitnesses)
	// << flattened continuation for buildListOfMcUnitsWithPotentiallyDifferentWitnesslists:1117:1
	if ! bHasBestParent {
		// :: flattened return for return handleResult('no compatible best parent');
		return "no compatible best parent"
	}
	console.log("###### MC units ", arrMcUnits)
	if len(arrMcUnits) == 0 {
		// :: flattened return for return handleResult();
		// ** need 1 return(s) instead of 0
		return 
	}
	rows := /* await */
	conn.query_sync("SELECT units.unit, COUNT(*) AS count_matching_witnesses \n" +
		"			FROM units CROSS JOIN unit_witnesses ON (units.unit=unit_witnesses.unit OR units.witness_list_unit=unit_witnesses.unit) AND address IN(?) \n" +
		"			WHERE units.unit IN(?) \n" +
		"			GROUP BY units.unit \n" +
		"			HAVING count_matching_witnesses<?", DBParamsT{
		arrWitnesses,
		arrMcUnits,
		constants.COUNT_WITNESSES - constants.MAX_WITNESS_LIST_MUTATIONS,
	})
	// << flattened continuation for conn.query:1123:2
	console.log(rows)
	if len(rows) > 0 {
		// :: flattened return for return handleResult('too many (' + (constants.COUNT_WITNESSES - rows[0].count_matching_witnesses) + ') witness list mutations relative to MC unit ' + rows[0].unit);
		return "too many (" + constants.COUNT_WITNESSES - rows[0].count_matching_witnesses + ") witness list mutations relative to MC unit " + rows[0].unit
	}
	// :: flattened return for handleResult();
	// ** need 1 return(s) instead of 0
	return 
	// >> flattened continuation for conn.query:1123:2
	// >> flattened continuation for buildListOfMcUnitsWithPotentiallyDifferentWitnesslists:1117:1
}

// the MC for this function is the MC built from this unit, not our current MC
func buildListOfMcUnitsWithPotentiallyDifferentWitnesslists_sync(conn DBConnT, objUnit UnitT, last_ball_unit UnitT, arrWitnesses AddressesT) (bool, UnitsT) {
	var(
		addAndGoUp_sync func (unit UnitT) (bool, UnitsT)
	)
	
	addAndGoUp_sync = func (unit UnitT) (bool, UnitsT) {
		props := /* await */
		readStaticUnitProps_sync(conn, unit)
		// << flattened continuation for readStaticUnitProps:1144:2
		// the parent has the same witness list and the parent has already passed the MC compatibility test
		if objUnit.witness_list_unit && objUnit.witness_list_unit == props.witness_list_unit {
			// :: flattened return for return handleList(true, arrMcUnits);
			return meta.returnArguments(true, arrMcUnits)
		} else {
			arrMcUnits = append(arrMcUnits, unit)
		}
		if unit == last_ball_unit {
			// :: flattened return for return handleList(true, arrMcUnits);
			return meta.returnArguments(true, arrMcUnits)
		}
		if ! props.best_parent_unit {
			_core.Throw("no best parent of unit " + unit + "?")
		}
		// :: flattened return for handleList(addAndGoUp(props.best_parent_unit));
		// ** need 2 return(s) instead of 1
		return /* await */
		addAndGoUp_sync(props.best_parent_unit)
		// >> flattened continuation for readStaticUnitProps:1144:2
	}
	
	arrMcUnits := UnitsT{}
	best_parent_unit := /* await */
	determineBestParent_sync(conn, objUnit, arrWitnesses)
	// << flattened continuation for determineBestParent:1159:1
	if ! best_parent_unit {
		// :: flattened return for return handleList(false);
		// ** need 2 return(s) instead of 1
		return false
	}
	// :: flattened return for handleList(addAndGoUp(best_parent_unit));
	// ** need 2 return(s) instead of 1
	return /* await */
	addAndGoUp_sync(best_parent_unit)
	// >> flattened continuation for determineBestParent:1159:1
}
`


//func readStaticUnitProps_sync(conn DBConnT, unit UnitT) refPropsT {
func ReadStaticUnitProps_sync(conn refDBConnT, unit UnitT) refPropsT {
//	props := assocCachedUnits[unit]
	props := AssocCachedUnits[unit]
//	if props {
	if props != nil {
		// :: flattened return for return handleProps(props);
		return props
	}
/**
	rows := /* await * /
	conn.query_sync("SELECT level, witnessed_level, best_parent_unit, witness_list_unit FROM units WHERE unit=?", DBParamsT{ unit })
 **/
	rcvr := db.PropsReceiver{}
	conn.MustQuery("SELECT level, witnessed_level, best_parent_unit, witness_list_unit FROM units WHERE unit=?", DBParamsT{ unit }, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:1171:1
	if len(rows) != 1 {
		_core.Throw("not 1 unit")
	}
	props = &rows[0]
//	assocCachedUnits[unit] = props
	AssocCachedUnits[unit] = props
	// :: flattened return for handleProps(props);
	return props
	// >> flattened continuation for conn.query:1171:1
}

//func readUnitAuthors_sync(conn DBConnT, unit UnitT)  {
func ReadUnitAuthors_sync(conn refDBConnT, unit UnitT) AddressesT {
//	arrAuthors := assocCachedUnitAuthors[unit]
	arrAuthors := AssocCachedUnitAuthors[unit]
//	if arrAuthors {
	if arrAuthors != nil {
		// :: flattened return for return handleAuthors(arrAuthors);
		// ** need 0 return(s) instead of 1
		return arrAuthors
	}
/**
	rows := /* await * /
	conn.query_sync("SELECT address FROM unit_authors WHERE unit=?", DBParamsT{ unit })
 **/
	rcvr := db.AddressesReceiver{}
	conn.MustQuery("SELECT address FROM unit_authors WHERE unit=?", DBParamsT{ unit }, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:1184:1
	if len(rows) == 0 {
		_core.Throw("no authors")
	}
/**
	arrAuthors2 := // .. not flattening for Array.map
	rows.map(func (row rowT) {*returns*} {
		return row.address
	}).sort()
 **/
	arrAuthors2 := make(AddressesT, 0, len(rows))
	for _, row := range rows {
		arrAuthors2 = append(arrAuthors2, row.Address)
	}
	sort.Slice(arrAuthors2, func (i, j int) bool {
		return arrAuthors2[i] < arrAuthors2[j]
	})
	//	if (arrAuthors && arrAuthors.join('-') !== arrAuthors2.join('-'))
	//		throw Error('cache is corrupt');
//	assocCachedUnitAuthors[unit] = arrAuthors2
	AssocCachedUnitAuthors[unit] = arrAuthors2
	// :: flattened return for handleAuthors(arrAuthors2);
	// ** need 0 return(s) instead of 1
	return arrAuthors2
	// >> flattened continuation for conn.query:1184:1
}

func isKnownUnit(unit UnitT) bool {
//	return assocCachedUnits[unit] != nil || assocKnownUnits[unit]
	return AssocCachedUnits[unit] != nil || AssocKnownUnits[unit]
}

func setUnitIsKnown(unit UnitT)  {
//	return assocKnownUnits[unit] = true
	AssocKnownUnits[unit] = true
}

func forgetUnit(unit UnitT)  {
	delete(AssocKnownUnits, unit)
	delete(AssocCachedUnits, unit)
	delete(AssocCachedUnitAuthors, unit)
	delete(AssocCachedUnitWitnesses, unit)
	delete(AssocUnstableUnits, unit)
	delete(AssocStableUnits, unit)
}

const _shrinkCache = `
func shrinkCache()  {
	if len(Object.keys(assocCachedAssetInfos)) > MAX_ITEMS_IN_CACHE {
		assocCachedAssetInfos = [*ObjectExpression*]
	}
	console.log(len(Object.keys(assocUnstableUnits)) + " unstable units")
	arrKnownUnits := Object.keys(assocKnownUnits)
	arrPropsUnits := Object.keys(assocCachedUnits)
	arrStableUnits := Object.keys(assocStableUnits)
	arrAuthorsUnits := Object.keys(assocCachedUnitAuthors)
	arrWitnessesUnits := Object.keys(assocCachedUnitWitnesses)
	if len(arrPropsUnits) < MAX_ITEMS_IN_CACHE && len(arrAuthorsUnits) < MAX_ITEMS_IN_CACHE && len(arrWitnessesUnits) < MAX_ITEMS_IN_CACHE && len(arrKnownUnits) < MAX_ITEMS_IN_CACHE && len(arrStableUnits) < MAX_ITEMS_IN_CACHE {
		return console.log("cache is small, will not shrink")
	}
	arrUnits := _.union(arrPropsUnits, arrAuthorsUnits, arrWitnessesUnits, arrKnownUnits, arrStableUnits)
	console.log("will shrink cache, total units: " + len(arrUnits))
	last_stable_mci := /* await */
	readLastStableMcIndex_sync(db)
	// << flattened continuation for readLastStableMcIndex:1225:1
	CHUNK_SIZE := 500
	// there is a limit on the number of query params
	for offset := 0; offset < len(arrUnits); offset += CHUNK_SIZE {
		rows := /* await */
		db.query_sync("SELECT unit FROM units WHERE unit IN(?) AND main_chain_index<? AND main_chain_index!=0", {*ArrayExpression*})
		// << flattened continuation for db.query:1229:3
		console.log("will remove " + len(rows) + " units from cache")
		// .. not flattening for Array.forEach
		for row, _ := range rows {
			delete(assocKnownUnits, row.unit)
			delete(assocCachedUnits, row.unit)
			delete(assocStableUnits, row.unit)
			delete(assocCachedUnitAuthors, row.unit)
			delete(assocCachedUnitWitnesses, row.unit)
		}
		// >> flattened continuation for db.query:1229:3
	}
	// >> flattened continuation for readLastStableMcIndex:1225:1
}
setInterval(shrinkCache, 300 * 1000)
`



func initUnstableUnits_sync()  {
/**
	rows := /* await * /
	db.query_sync("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain, is_free, is_stable, witnessed_level \n" +
		"		FROM units WHERE is_stable=0 ORDER BY +level")
 **/
	rcvr := db.RUPUnitPropsReceiver{}
	db.MustQuery("SELECT unit, level, latest_included_mc_index, main_chain_index, is_on_main_chain, is_free, is_stable, witnessed_level \n" +
		"		FROM units WHERE is_stable=0 ORDER BY +level", DBParamsT{}, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for db.query:1251:1
	//	assocUnstableUnits = {};
	// .. not flattening for Array.forEach
//	for row, _ := range rows {
	for _, row := range rows {
//		row.parent_units = [*ArrayExpression*]
		// [fyi] making copy instead of modding an original
		// [tbd] alternative: use XPropsT for the rows
		//row.Parent_units = UnitsT{}
//		assocUnstableUnits[row.unit] = row
		AssocUnstableUnits[row.Unit] = &XPropsT{
			// [tbd] use UnitPropsT
			Unit: row.Unit,
			PropsT: row.PropsT,
			Parent_units: UnitsT{},
		}
	}
//	console.log("initUnstableUnits 1 done")
	console.Log("initUnstableUnits 1 done")
//	if len(Object.keys(assocUnstableUnits)) == 0 {
	if len(AssocUnstableUnits) == 0 {
//		if onDone {
		if true {
			// :: flattened return for onDone();
			return 
		}
		return 
	}
/**
	prows := /* await * /
	db.query_sync("SELECT parent_unit, child_unit FROM parenthoods WHERE child_unit IN(" + Object.keys(assocUnstableUnits).map(db.escape) + ")")
 **/	
	queryParams := DBParamsT{}
	unitsSql := queryParams.AddUnitsFromIterator(func (iterFn db.AddUnitIteratorFunctorT) int {
		for unit, _ := range AssocUnstableUnits { iterFn(unit) }
		return len(AssocUnstableUnits)
	})
	rcvr_1 := db.ParentChildUnitsReceiver{}
	db.MustQuery("SELECT parent_unit, child_unit FROM parenthoods WHERE child_unit IN(" + unitsSql + ")", queryParams, &rcvr_1)
	prows := rcvr_1.Rows
	// << flattened continuation for db.query:1266:3
	// .. not flattening for Array.forEach
//	for prow, _ := range prows {
	for _, prow := range prows {
		// [tbd] optimize updates of the map element
//		assocUnstableUnits[prow.child_unit].parent_units = append(assocUnstableUnits[prow.child_unit].parent_units, prow.parent_unit)
		AssocUnstableUnits[prow.Child_unit].Parent_units = append(AssocUnstableUnits[prow.Child_unit].Parent_units, prow.Parent_unit)
	}
//	console.log("initUnstableUnits done")
	console.Log("initUnstableUnits done")
//	if onDone {
	if true {
		// :: flattened return for onDone();
		return 
	}
	// >> flattened continuation for db.query:1266:3
	// >> flattened continuation for db.query:1251:1
}

//func resetUnstableUnits_sync()  {
func ResetUnstableUnits_sync()  {
	// .. not flattening for Array.forEach
//	for unit, _ := range Object.keys(assocUnstableUnits) {
	for unit, _ := range AssocUnstableUnits {
//		delete(assocUnstableUnits, unit)
		delete(AssocUnstableUnits, unit)
	}
	/* await */
	initUnstableUnits_sync()
	// << flattened continuation for initUnstableUnits:1285:1
	// :: flattened return for onDone();
	return 
	// >> flattened continuation for initUnstableUnits:1285:1
}

//mutex.lock({*ArrayExpression*}, initUnstableUnits)

func Init() {
	initializeMinRetrievableMci()
	initUnstableUnits_sync()
}

//if ! conf.bLight {
//	archiveJointAndDescendantsIfExists("N6QadI9yg3zLxPMphfNGJcPfddW4yHPkoGMbbGZsWa0=")
//}


//exports.isGenesisUnit = isGenesisUnit
//exports.isGenesisBall = isGenesisBall

//exports.readWitnesses = readWitnesses
//exports.readWitnessList = readWitnessList
//exports.findWitnessListUnit = findWitnessListUnit
//exports.determineIfWitnessAddressDefinitionsHaveReferences = determineIfWitnessAddressDefinitionsHaveReferences

//exports.readUnitProps = readUnitProps
//exports.readPropsOfUnits = readPropsOfUnits

//exports.readJoint = readJoint
//exports.readJointWithBall = readJointWithBall
//exports.readFreeJoints = readFreeJoints

//exports.readDefinitionByAddress = readDefinitionByAddress
//exports.readDefinition = readDefinition

//exports.readLastMainChainIndex = readLastMainChainIndex

//exports.readLastStableMcUnitProps = readLastStableMcUnitProps
//exports.readLastStableMcIndex = readLastStableMcIndex


//exports.findLastBallMciOfMci = findLastBallMciOfMci
//exports.getMinRetrievableMci = getMinRetrievableMci
//exports.updateMinRetrievableMciAfterStabilizingMci = updateMinRetrievableMciAfterStabilizingMci

//exports.archiveJointAndDescendantsIfExists = archiveJointAndDescendantsIfExists

//exports.readAsset = readAsset
//exports.filterAttestedAddresses = filterAttestedAddresses
//exports.loadAssetWithListOfAttestedAuthors = loadAssetWithListOfAttestedAuthors

//exports.filterNewOrUnstableUnits = filterNewOrUnstableUnits

//exports.determineWitnessedLevelAndBestParent = determineWitnessedLevelAndBestParent
//exports.determineBestParent = determineBestParent
//exports.determineIfHasWitnessListMutationsAlongMc = determineIfHasWitnessListMutationsAlongMc

//exports.readStaticUnitProps = readStaticUnitProps
//exports.readUnitAuthors = readUnitAuthors

//exports.isKnownUnit = isKnownUnit
//exports.setUnitIsKnown = setUnitIsKnown
//exports.forgetUnit = forgetUnit

//exports.sliceAndExecuteQuery = sliceAndExecuteQuery

//exports.assocUnstableUnits = assocUnstableUnits
//exports.assocStableUnits = assocStableUnits
//exports.resetUnstableUnits = resetUnstableUnits


// converted golang end

