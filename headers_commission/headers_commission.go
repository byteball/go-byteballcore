
// converted golang begin

package headers_commission

import(
	"sort"

 _core	"nodejs/core"
	"nodejs/console"

 .	"github.com/byteball/go-byteballcore/types"

)

import(
//		"crypto"
		"github.com/byteball/go-byteballcore/object_hash"
//		"async"
		"github.com/byteball/go-byteballcore/db"
//		"github.com/byteball/go-byteballcore/conf"
)

type(
	DBConnT		= db.DBConnT
	refDBConnT	= *DBConnT

	DBParamsT	= db.DBParamsT

	ChildrenInfoT struct{
		Headers_commission  AmountT
		Children	ChildsT
	}
	ChildrenInfosT	= []ChildrenInfoT

	ChildT struct{
		db.PaymentInfoRow
		Hash		refSHA1HexT
	}
	ChildsT		= []ChildT

	refSHA1HexT	= *SHA1HexT

	ChildrenInfoByUnitMapT  = map[UnitT] refChildrenInfoT

	refChildrenInfoT  = *ChildrenInfoT

	refChildT  	= *ChildT
	refChildsT  	= *ChildsT

	WonAmountsByUnitMapT  = map[UnitT] refWonAmountsByPayerMapT
	WonAmountsByPayerMapT  = map[UnitT] AmountT

	refWonAmountsByPayerMapT = *WonAmountsByPayerMapT
)


var max_spendable_mci MCIndexT = MCIndexT_Null

//func calcHeadersCommissions(conn DBConnT, onDone onDoneT)  {
func CalcHeadersCommissions_sync(conn refDBConnT)  {
	// we don't require neither source nor recipient to be majority witnessed -- we don't want to return many times to the same MC index.
	console.Log("will calc h-comm")
//	if max_spendable_mci == nil {
	if max_spendable_mci.IsNull() {
		// first calc after restart only
		/* await */
		initMaxSpendableMci_sync(conn)
/**
		return
		// << flattened continuation for initMaxSpendableMci:14:9
		calcHeadersCommissions(conn, onDone)
		// >> flattened continuation for initMaxSpendableMci:14:9
 **/
	}
	// max_spendable_mci is old, it was last updated after previous calc
	since_mc_index := max_spendable_mci
	
	(func () ErrorT {
	  // :: inlined async.series:19:1
//	  for _f := range (AsyncFunctorsT{
	  for _, _f := range (AsyncFunctorsT{
	  	func () ErrorT {
/**
//	  		if conf.storage == "mysql" {
	  		if conf.Storage == "mysql" {
	  			best_child_sql := "SELECT unit \n" +
	  				"					FROM parenthoods \n" +
	  				"					JOIN units AS alt_child_units ON parenthoods.child_unit=alt_child_units.unit \n" +
	  				"					WHERE parent_unit=punits.unit AND alt_child_units.main_chain_index-punits.main_chain_index<=1 AND +alt_child_units.sequence='good' \n" +
	  				"					ORDER BY SHA1(CONCAT(alt_child_units.unit, next_mc_units.unit)) \n" +
	  				"					LIMIT 1"
	  			// headers commissions to single unit author
	  			/* await * /
	  			conn.query_sync("INSERT INTO headers_commission_contributions (unit, address, amount) \n" +
	  				"					SELECT punits.unit, address, punits.headers_commission AS hc \n" +
	  				"					FROM units AS chunits \n" +
	  				"					JOIN unit_authors USING(unit) \n" +
	  				"					JOIN parenthoods ON chunits.unit=parenthoods.child_unit \n" +
	  				"					JOIN units AS punits ON parenthoods.parent_unit=punits.unit \n" +
	  				"					JOIN units AS next_mc_units ON next_mc_units.is_on_main_chain=1 AND next_mc_units.main_chain_index=punits.main_chain_index+1 \n" +
	  				"					WHERE chunits.is_stable=1 \n" +
	  				"						AND +chunits.sequence='good' \n" +
	  				"						AND punits.main_chain_index>? \n" +
	  				"						AND chunits.main_chain_index-punits.main_chain_index<=1 \n" +
	  				"						AND +punits.sequence='good' \n" +
	  				"						AND punits.is_stable=1 \n" +
	  				"						AND next_mc_units.is_stable=1 \n" +
	  				"						AND chunits.unit=( " + best_child_sql + " ) \n" +
	  				"						AND (SELECT COUNT(*) FROM unit_authors WHERE unit=chunits.unit)=1 \n" +
	  				"						AND (SELECT COUNT(*) FROM earned_headers_commission_recipients WHERE unit=chunits.unit)=0 \n" +
	  				"					UNION ALL \n" +
	  				"					SELECT punits.unit, earned_headers_commission_recipients.address, \n" +
	  				"						ROUND(punits.headers_commission*earned_headers_commission_share/100.0) AS hc \n" +
	  				"					FROM units AS chunits \n" +
	  				"					JOIN earned_headers_commission_recipients USING(unit) \n" +
	  				"					JOIN parenthoods ON chunits.unit=parenthoods.child_unit \n" +
	  				"					JOIN units AS punits ON parenthoods.parent_unit=punits.unit \n" +
	  				"					JOIN units AS next_mc_units ON next_mc_units.is_on_main_chain=1 AND next_mc_units.main_chain_index=punits.main_chain_index+1 \n" +
	  				"					WHERE chunits.is_stable=1 \n" +
	  				"						AND +chunits.sequence='good' \n" +
	  				"						AND punits.main_chain_index>? \n" +
	  				"						AND chunits.main_chain_index-punits.main_chain_index<=1 \n" +
	  				"						AND +punits.sequence='good' \n" +
	  				"						AND punits.is_stable=1 \n" +
	  				"						AND next_mc_units.is_stable=1 \n" +
	  				"						AND chunits.unit=( " + best_child_sql + " )", DBParamsT{
	  				since_mc_index,
	  				since_mc_index,
	  			})
	  			// << flattened continuation for conn.query:29:4
	  			// :: flattened return for cb();
	  			// ** need 1 return(s) instead of 0
	  			return nil
	  			// >> flattened continuation for conn.query:29:4
	  		} else {
 **/
			{
/**
	  			rows := /* await * /
	  			conn.query_sync(// chunits is any child unit and contender for headers commission, punits is hc-payer unit
	  			"SELECT chunits.unit AS child_unit, punits.headers_commission, next_mc_units.unit AS next_mc_unit, punits.unit AS payer_unit \n" +
	  				"					FROM units AS chunits \n" +
	  				"					JOIN parenthoods ON chunits.unit=parenthoods.child_unit \n" +
	  				"					JOIN units AS punits ON parenthoods.parent_unit=punits.unit \n" +
	  				"					JOIN units AS next_mc_units ON next_mc_units.is_on_main_chain=1 AND next_mc_units.main_chain_index=punits.main_chain_index+1 \n" +
	  				"					WHERE chunits.is_stable=1 \n" +
	  				"						AND +chunits.sequence='good' \n" +
	  				"						AND punits.main_chain_index>? \n" +
	  				"						AND +punits.sequence='good' \n" +
	  				"						AND punits.is_stable=1 \n" +
	  				"						AND chunits.main_chain_index-punits.main_chain_index<=1 \n" +
	  				"						AND next_mc_units.is_stable=1", DBParamsT{ since_mc_index })
 **/
				rcvr := db.PaymentInfosReceiver{}
	  			conn.MustQuery(// chunits is any child unit and contender for headers commission, punits is hc-payer unit
	  			"SELECT chunits.unit AS child_unit, punits.headers_commission, next_mc_units.unit AS next_mc_unit, punits.unit AS payer_unit \n" +
	  				"FROM units AS chunits \n" +
	  				"JOIN parenthoods ON chunits.unit=parenthoods.child_unit \n" +
	  				"JOIN units AS punits ON parenthoods.parent_unit=punits.unit \n" +
	  				"JOIN units AS next_mc_units ON next_mc_units.is_on_main_chain=1 AND next_mc_units.main_chain_index=punits.main_chain_index+1 \n" +
	  				"WHERE chunits.is_stable=1 \n" +
	  				"	AND +chunits.sequence='good' \n" +
	  				"	AND punits.main_chain_index>? \n" +
	  				"	AND +punits.sequence='good' \n" +
	  				"	AND punits.is_stable=1 \n" +
	  				"	AND chunits.main_chain_index-punits.main_chain_index<=1 \n" +
	  				"	AND next_mc_units.is_stable=1", DBParamsT{ since_mc_index }, &rcvr)
				rows := rcvr.Rows
	  			// << flattened continuation for conn.query:68:4
//	  			assocChildrenInfos := [*ObjectExpression*]
	  			assocChildrenInfos := make(ChildrenInfoByUnitMapT)
	  			// .. not flattening for Array.forEach
//	  			for row, _ := range rows {
	  			for _, row := range rows {
//	  				payer_unit := row.payer_unit
	  				payer_unit := row.Payer_unit
//	  				child_unit := row.child_unit
//[uu]	  				child_unit := row.Child_unit
//	  				if ! assocChildrenInfos[payer_unit] {
	  				if _, _exists := assocChildrenInfos[payer_unit]; ! _exists {
//	  					assocChildrenInfos[payer_unit] = [*ObjectExpression*]
	  					assocChildrenInfos[payer_unit] = &ChildrenInfoT{
							Headers_commission: row.Headers_commission,
							Children: ChildsT{},
						}
	  				} else {
//	  					if assocChildrenInfos[payer_unit].headers_commission != row.headers_commission {
	  					if assocChildrenInfos[payer_unit].Headers_commission != row.Headers_commission {
	  						_core.Throw("different headers_commission")
	  					}
	  				}
//	  				row.headers_commission = nil
	  				row.Headers_commission = AmountT_Null
//	  				row.payer_unit = nil
	  				row.Payer_unit = UnitT_Null
//	  				assocChildrenInfos[payer_unit].children = append(assocChildrenInfos[payer_unit].children, row)
	  				assocChildrenInfos[payer_unit].Children = append(assocChildrenInfos[payer_unit].Children, ChildT{
						PaymentInfoRow: row,
					})
	  			}
//	  			assocWonAmounts := [*ObjectExpression*]
	  			assocWonAmounts := make(WonAmountsByUnitMapT)
	  			// amounts won, indexed by child unit who won the hc, and payer unit
	  			for payer_unit := range assocChildrenInfos {
//	  				headers_commission := assocChildrenInfos[payer_unit].headers_commission
	  				headers_commission := assocChildrenInfos[payer_unit].Headers_commission
//	  				winnerChildInfo := getWinnerInfo(assocChildrenInfos[payer_unit].children)
	  				winnerChildInfo := getWinnerInfo(assocChildrenInfos[payer_unit].Children)
//	  				child_unit := winnerChildInfo.child_unit
	  				child_unit := winnerChildInfo.Child_unit
//	  				if ! assocWonAmounts[child_unit] {
	  				if _, _exists := assocWonAmounts[child_unit]; ! _exists {
//	  					assocWonAmounts[child_unit] = [*ObjectExpression*]
	  					wasMap := make(WonAmountsByPayerMapT)
	  					assocWonAmounts[child_unit] = &wasMap
	  				}
//	  				assocWonAmounts[child_unit][payer_unit] = headers_commission
	  				(*assocWonAmounts[child_unit])[payer_unit] = headers_commission
	  			}
	  			//console.log(assocWonAmounts);
//	  			arrWinnerUnits := Object.keys(assocWonAmounts)
	  			arrWinnerUnits := make(UnitsT, 0, len(assocWonAmounts))
				for winnerUnit := range assocWonAmounts {
					arrWinnerUnits = append(arrWinnerUnits, winnerUnit)
				}
	  			if len(arrWinnerUnits) == 0 {
	  				// :: flattened return for return cb();
	  				// ** need 1 return(s) instead of 0
	  				return nil
	  			}

/**
	  			strWinnerUnitsList := arrWinnerUnits.map(db.escape).join(", ")
	  			profit_distribution_rows := /* await * /
	  			conn.query_sync("SELECT \n" +
	  				"								unit_authors.unit, \n" +
	  				"								unit_authors.address, \n" +
	  				"								100 AS earned_headers_commission_share \n" +
	  				"							FROM unit_authors \n" +
	  				"							LEFT JOIN earned_headers_commission_recipients USING(unit) \n" +
	  				"							WHERE unit_authors.unit IN(" + strWinnerUnitsList + ") AND earned_headers_commission_recipients.unit IS NULL \n" +
	  				"							UNION ALL \n" +
	  				"							SELECT \n" +
	  				"								unit, \n" +
	  				"								address, \n" +
	  				"								earned_headers_commission_share \n" +
	  				"							FROM earned_headers_commission_recipients \n" +
	  				"							WHERE unit IN(" + strWinnerUnitsList + ")")
 **/
				rcvr_1 := db.UnitAuthorSharesReceiver{}
				queryParams := DBParamsT{}
				winnerUnitsSql := queryParams.AddUnits(arrWinnerUnits)
				winnerUnitsSql = queryParams.AddUnits(arrWinnerUnits)
	  			conn.MustQuery("SELECT \n" +
	  				"	unit_authors.unit, \n" +
	  				"	unit_authors.address, \n" +
	  				"	100 AS earned_headers_commission_share \n" +
	  				"FROM unit_authors \n" +
	  				"LEFT JOIN earned_headers_commission_recipients USING(unit) \n" +
	  				"WHERE unit_authors.unit IN(" + winnerUnitsSql + ") AND earned_headers_commission_recipients.unit IS NULL \n" +
	  				"UNION ALL \n" +
	  				"SELECT \n" +
	  				"	unit, \n" +
	  				"	address, \n" +
	  				"	earned_headers_commission_share \n" +
	  				"FROM earned_headers_commission_recipients \n" +
	  				"WHERE unit IN(" + winnerUnitsSql + ")", queryParams, &rcvr_1)
				profit_distribution_rows := rcvr_1.Rows
	  			// << flattened continuation for conn.query:110:6
//	  			arrValues := {*ArrayExpression*}
	  			arrValues := ContributionsT{}
	  			// .. not flattening for Array.forEach
//	  			for row, _ := range profit_distribution_rows {
	  			for _, row := range profit_distribution_rows {
//	  				child_unit := row.unit
	  				child_unit := row.Unit
//	  				for payer_unit := range assocWonAmounts[child_unit] {
	  				for payer_unit := range *assocWonAmounts[child_unit] {
//	  					full_amount := assocWonAmounts[child_unit][payer_unit]
	  					full_amount := (*assocWonAmounts[child_unit])[payer_unit]
//	  					if ! full_amount {
	  					if full_amount.IsNull() {
//	  						_core.Throw("no amount for child unit " + child_unit + ", payer unit " + payer_unit)
	  						_core.Throw("no amount for child unit %s, payer unit %s", child_unit, payer_unit)
	  					}
	  					// note that we round _before_ summing up header commissions won from several parent units
//	  					amount := (row.earned_headers_commission_share == 100 ? full_amount: Math.round(full_amount * row.earned_headers_commission_share / 100.0))
	  					amount := full_amount
	  					if !(row.Earned_headers_commission_share == 100) {
							// [tbd] use math/big instead?...
							amount = AmountT(float64(full_amount) * float64(row.Earned_headers_commission_share) / 100.0 + 0.5)
						}
	  					// hc outputs will be indexed by mci of _payer_ unit
//	  					arrValues = append(arrValues, "('" + payer_unit + "', '" + row.address + "', " + amount + ")")
	  					arrValues = append(arrValues, ContributionT{
							Payer_unit: payer_unit,
							Address: row.Address,
							Amount: amount,
						})
	  				}
	  			}
/**
	  			/* await * /
	  			conn.query_sync("INSERT INTO headers_commission_contributions (unit, address, amount) VALUES " + arrValues.join(", "))
 **/
				queryParams = DBParamsT{}
				cbnsSql := queryParams.AddContributions(arrValues)
	  			conn.MustExec("INSERT INTO headers_commission_contributions (unit, address, amount) VALUES " + cbnsSql, queryParams)
	  			// << flattened continuation for conn.query:141:8
	  			// :: flattened return for cb();
	  			// ** need 1 return(s) instead of 0
	  			return nil
	  			// >> flattened continuation for conn.query:141:8
	  			// >> flattened continuation for conn.query:110:6
	  			// >> flattened continuation for conn.query:68:4
	  		}
	  	},
	  	func () ErrorT {
/**
	  		/* await * /
	  		conn.query_sync("INSERT INTO headers_commission_outputs (main_chain_index, address, amount) \n" +
	  			"				SELECT main_chain_index, address, SUM(amount) FROM headers_commission_contributions JOIN units USING(unit) \n" +
	  			"				WHERE main_chain_index>? \n" +
	  			"				GROUP BY main_chain_index, address", DBParamsT{ since_mc_index })
 **/
	  		conn.MustExec("INSERT INTO headers_commission_outputs (main_chain_index, address, amount) \n" +
	  			"SELECT main_chain_index, address, SUM(amount) FROM headers_commission_contributions JOIN units USING(unit) \n" +
	  			"WHERE main_chain_index>? \n" +
	  			"GROUP BY main_chain_index, address", DBParamsT{ since_mc_index })
	  		// << flattened continuation for conn.query:151:3
	  		// :: flattened return for cb();
	  		// ** need 1 return(s) instead of 0
	  		return nil
	  		// >> flattened continuation for conn.query:151:3
	  	},
	  	func () ErrorT {
/**
	  		rows := /* await * /
	  		conn.query_sync("SELECT MAX(main_chain_index) AS max_spendable_mci FROM headers_commission_outputs")
 **/
			rcvr := db.MaxSpendableMCIsReceiver{}
	  		conn.MustQuery("SELECT MAX(main_chain_index) AS max_spendable_mci FROM headers_commission_outputs", DBParamsT{}, &rcvr)
			rows := rcvr.Rows
	  		// << flattened continuation for conn.query:161:3
//	  		max_spendable_mci = rows[0].max_spendable_mci
	  		max_spendable_mci = rows[0].Max_spendable_mci
	  		// :: flattened return for cb();
	  		// ** need 1 return(s) instead of 0
	  		return nil
	  		// >> flattened continuation for conn.query:161:3
	  	},
	  }) {
	    if _err := _f() ; _err != nil { return _err }
	  }
	  return nil
	})()
}

func getWinnerInfo(arrChildren ChildsT) refChildT {
	if len(arrChildren) == 1 {
		return &arrChildren[0]
	}
	// .. not flattening for Array.forEach
//	for child, _ := range arrChildren {
	for _, child := range arrChildren {
//		child.hash = crypto.createHash("sha1").update(child.child_unit + child.next_mc_unit, "utf8").digest("hex")
		child.Hash = object_hash.SHA1Hex(string(child.Child_unit) + string(child.Next_mc_unit))
	}
	// .. not flattening for Array.sort
//	arrChildren.sort(func (a aT, b bT) {
//		return (a.hash < b.hash ? - 1: 1)
//	})
	sort.Slice(arrChildren, func (i, j int) bool {
		return *arrChildren[i].Hash < *arrChildren[j].Hash
	})
	return &arrChildren[0]
}

//func initMaxSpendableMci_sync(conn DBConnT)  {
func initMaxSpendableMci_sync(conn refDBConnT)  {
/**
	rows := /* await * /
	conn.query_sync("SELECT MAX(main_chain_index) AS max_spendable_mci FROM headers_commission_outputs")
 **/
	rcvr := db.MaxSpendableMCIsReceiver{}
	conn.MustQuery("SELECT MAX(main_chain_index) AS max_spendable_mci FROM headers_commission_outputs", nil, &rcvr)
	rows := rcvr.Rows
	// << flattened continuation for conn.query:181:1
//	max_spendable_mci = rows[0].max_spendable_mci || 0
	max_spendable_mci = MCIndexT(0)
	if 0 < len(rows) { max_spendable_mci = rows[0].Max_spendable_mci }
//	if onDone {
	if true {
		// :: flattened return for onDone();
		return 
	}
	// >> flattened continuation for conn.query:181:1
}

func getMaxSpendableMciForLastBallMci(last_ball_mci MCIndexT) MCIndexT {
	return last_ball_mci - 1
}

func Init() {
//	initMaxSpendableMci_sync()
}


//exports.calcHeadersCommissions = calcHeadersCommissions
//exports.getMaxSpendableMciForLastBallMci = getMaxSpendableMciForLastBallMci


// converted golang end

