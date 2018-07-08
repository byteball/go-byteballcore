package db

import(
	"database/sql"
	"context"
//	"fmt"
	"log"
	"strings"
//	"regexp"
	"time"
	"sort"
	"math/bits"

//	lt "db/wrappers"

	"nodejs/console"

	"github.com/byteball/go-byteballcore/types"
)

type(
)

//

var(
	TExec		int64
	TQuery		int64
	TCommit		int64
	TRollback	int64
)

func TReset() {
	TExec = 0
	TQuery = 0
	TCommit = 0
	TRollback = 0
}

//

func log_Printf(fmt string, args... interface{}) {
	return
	log.Printf(fmt, args...)
}

//

func (params *DBParamsT) AddUnits(units types.UnitsT) DBSqlT {
	usSql := strings.Repeat(",?", len(units))[1:]
	for _, unit := range units {
		*params = append(*params, unit)
	}
	return usSql
}

// XPropsByUnitMapT

type(
	AddFromIteratorT	func (iterFn AddFromIteratorFunctorT) (int, string)
	AddFromIteratorFunctorT	func (params_... interface{})

	AddUnitIteratorT	func (iterFn AddUnitIteratorFunctorT) int
	AddUnitIteratorFunctorT	func (unit types.UnitT)
)

func (params *DBParamsT) AddFromIterator(iter AddFromIteratorT) DBSqlT {
	pslen, pSql := iter(func (params_... interface{}) {
		*params = append(*params, params_...)
	})
	//if pSql == nil { pSql = ",?" }
	psSql := strings.Repeat(pSql, pslen)[1:]
	return psSql
}

func (params *DBParamsT) AddUnitsFromIterator(iter AddUnitIteratorT) DBSqlT {
	uslen := iter(func (unit types.UnitT) {
		*params = append(*params, unit)
	})
	usSql := strings.Repeat(",?", uslen)[1:]
	return usSql
}

func (params *DBParamsT) AddAddresses(addresses types.AddressesT) DBSqlT {
	asSql := strings.Repeat(",?", len(addresses))[1:]
	for _, address := range addresses {
		*params = append(*params, address)
	}
	return asSql
}

func (params *DBParamsT) AddMCIndexes(mcindexes types.MCIndexesT) DBSqlT {
	mcisSql := strings.Repeat(",?", len(mcindexes))[1:]
	for _, mcindex := range mcindexes {
		*params = append(*params, mcindex)
	}
	return mcisSql
}

func (params *DBParamsT) AddContributions(cbns types.ContributionsT) DBSqlT {
	cbnsSql := strings.Repeat(",(?,?,?)", len(cbns))[1:]
	for _, cbn := range cbns {
		*params = append(*params, cbn.Payer_unit, cbn.Address, cbn.Amount)
	}
	return cbnsSql
}

func (params *DBParamsT) AddPaidWitnessEvents(pwes types.PaidWitnessEventsT) DBSqlT {
	pwesSql := strings.Repeat(",(?,?,?)", len(pwes))[1:]
	for _, pwe := range pwes {
		*params = append(*params, pwe.Unit, pwe.Address, pwe.Delay)
	}
	return pwesSql
}

//

func MustExec(sql DBSqlT, params DBParamsT) *DBQueryResultT {
/**
	database := Instance()
	t0 := time.Now()
	err, res := database.Exec(sql, params)
	TExec += time.Now().Sub(t0).Nanoseconds()
	if err != nil {
		log.Printf("sql %s\n", sql)
		log.Printf("params %#v\n", params)
		log.Panicf("db.MustExec: %#v", err)
	}
	return res
 **/
	conn := TakeConnectionFromPool_sync()
	defer conn.Release()
	return conn.MustExec(sql, params)
}

func MustQuery(sql DBSqlT, params DBParamsT, rcvr Receiver) {
/**
	database := Instance()
	t0 := time.Now()
	err := database.Query(sql, params, rcvr)
	TQuery += time.Now().Sub(t0).Nanoseconds()
	if err != nil {
		log.Printf("sql %s\n", sql)
		log.Printf("params %#v\n", params)
		log.Panicf("db.MustSelect: %#v", err)
	}
 **/
	conn := TakeConnectionFromPool_sync()
	defer conn.Release()
	conn.MustQuery(sql, params, rcvr)
}

//

func (conn *DBConnT) BeginTransaction() {

	// [fyi] emulated transactions over this connection

	conn.MustExec("BEGIN", nil)
	return

	// [fyi] using native sql transactions

	dbconn := conn.dbconn
	ctx := context.Background()

//	isol := sql.LevelWriteCommitted
	isol := sql.LevelLinearizable
	txn, err := dbconn.BeginTx(ctx, &sql.TxOptions{
		Isolation: isol,
	})
	if err != nil {
		log.Panicf("conn.BeginTransaction: %#v", err)
	}

	// [tbd] rollback existing txn

	log_Printf("conn.Begin: isol %d", isol)

	conn.txn = txn
}

func (conn *DBConnT) Commit() {

	// [fyi] emulated transactions over this connection
{{
	t0 := time.Now()
	conn.MustExec("COMMIT", nil)
	TCommit += time.Now().Sub(t0).Nanoseconds()
	return
}}
	// [fyi] using native sql transactions

	if conn.txn == nil {
		log.Printf("conn.Commit: txn == nil")
		return
	}

	t0 := time.Now()
	err := conn.txn.Commit()
	TCommit += time.Now().Sub(t0).Nanoseconds()
	if err != nil {
		//log.Printf("conn.Commit: " + err.Error())
		//return
		log.Panicf("conn.Commit: %#v", err)
	}

	log_Printf("conn.Commit: ok")

	conn.txn = nil
}

func (conn *DBConnT) Rollback() {

	// [fyi] emulated transactions over this connection
{{
	t0 := time.Now()
	conn.MustExec("ROLLBACK", nil)
	TRollback += time.Now().Sub(t0).Nanoseconds()
	return
}}
	// [fyi] using native sql transactions

	if conn.txn == nil {
		log.Printf("conn.Rollback: txn == nil")
		return
	}

	t0 := time.Now()
	err := conn.txn.Rollback()
	TRollback += time.Now().Sub(t0).Nanoseconds()
	if err != nil {
		log.Panicf("conn.Rollback: %#v", err)
	}

	log_Printf("conn.Rollback: ok")

	conn.txn = nil
}

//

func (conn *DBConnT) ExecPrepared(prep refDBPreparedT, params []interface{}) (sql.Result, error) {
	ctx := conn.ctx
	sql_ := prep.Sql
	if conn.txn != nil {
		return conn.txn.ExecContext(ctx, sql_, params...)
	}
	// [tbd] fix this
	log_Printf("conn.Exec %s", strings.Split(sql_, "\n")[0])
	if prep.coStmt != nil {
		prep.CoStmtExecCount++
		return prep.coStmt.ExecContext(ctx, params...)
	}
	return conn.dbconn.ExecContext(ctx, sql_, params...)
}

func (conn *DBConnT) QueryPrepared(prep refDBPreparedT, params []interface{}) (*sql.Rows, error) {
	ctx := conn.ctx
	sql_ := prep.Sql
	if conn.txn != nil {
		return conn.txn.QueryContext(ctx, sql_, params...)
	}
	// [tbd] fix this
	log_Printf("conn.Query %s", strings.Split(sql_, "\n")[0])
	if prep.coStmt != nil {
		prep.CoStmtExecCount++
		return prep.coStmt.QueryContext(ctx, params...)
	}
	return conn.dbconn.QueryContext(ctx, sql_, params...)
}

//

type(
	DBPreparedByTagT = map[DBPreparedTagT] refDBPreparedT

	DBPreparedTagT	uint16

	DBPreparedT struct{
		Sql		string
		Tag		string
		Hash		uint16
		LengthTag	int
		LengthSql	int
		LengthParams	int
		Count		int
		TotalTime	int64

		coStmt		*sql.Stmt

		CoStmtCount	int
		CoStmtTime	int64
		CoStmtExecCount	int
	}
	refDBPreparedT	= *DBPreparedT

)

func (prep *DBPreparedT) Sample (t int64) {
	prep.Count++
	prep.TotalTime += t

	if ! (t < 1 * 1000000 && prep.LengthParams < 10) {
//::		console.Log("!!  %3d  %8.3f  %#v", prep.LengthParams, float64(t)/1.0e6, printable(prep.Tag))
	}
}

func printable(tag string) string {
	tag = strings.Replace(tag, "\n", "\\n", -1)
	tag = strings.Replace(tag, "\t", "\\t", -1)
	return tag
}

func (conn *DBConnT) prepareSql(sql string, lenParams int) refDBPreparedT {
	lenTag := len(sql)
	if 48 < lenTag { lenTag = 48 }

	tag_ := sql[0:lenTag]

	// [fyi] random micro-hash
//	hash := uint16((len(sql) << 6) + lenTag)
	hash := uint16((len(sql) << 6) + lenParams)
	for _, c := range tag_ {
//	for _, c := range sql {
		hash = bits.RotateLeft16(hash ^ uint16(c), 5)
	}
	tag := DBPreparedTagT(hash)

	prep := conn.preparedByTag[tag]
	if prep == nil {
		prep = &DBPreparedT{
			Sql: sql,
			Tag: tag_,
			Hash: hash,
			LengthTag: lenTag,
			LengthSql: len(sql),
			LengthParams: lenParams,
		}
		conn.preparedByTag[tag] = prep
		var err error
		t0 := time.Now()
		prep.coStmt, err = conn.dbconn.PrepareContext(conn.ctx, sql)
		prep.CoStmtTime += time.Now().Sub(t0).Nanoseconds()
		prep.CoStmtCount++
		if err != nil { log.Panicf("conn.Prepare: %s", err.Error()) }
	} else {
		if !(prep.LengthSql == len(sql) && prep.LengthParams == lenParams) {
			log.Panicf("prepareSql: hash collision %#v %#v", prep.Sql, sql)
		}
	}

	return prep
}

func (conn *DBConnT) ShowPrepared() {
	preps := make([]refDBPreparedT, 0, len(conn.preparedByTag)*3)
	for _, prep := range conn.preparedByTag {
		preps = append(preps, prep)
	}
	sort.Slice(preps, func(i, j int) bool {
		if preps[i].Tag != preps[j].Tag {
			return preps[i].Tag < preps[j].Tag
		}
		//return preps[i].LengthParams < preps[j].LengthParams
		return preps[i].LengthSql < preps[j].LengthSql
	})
	for _, prep := range preps {
		count := prep.Count
		avgTime := float64(prep.TotalTime) / float64(count)
		avgCoStTime := float64(prep.CoStmtTime) / float64(prep.CoStmtCount)
		mark := ""
		if 0.5 < avgTime/1.0e6 { mark = "**" }
		console.Log("==  %6d  %2.2s  %6.3f  %6.3f %4d %4d  %2d %4d %3d  %04x  %s",
			count, mark, avgTime/1.0e6,
			avgCoStTime/1.0e6, prep.CoStmtCount, prep.CoStmtExecCount,
			prep.LengthTag, prep.LengthSql, prep.LengthParams,
			prep.Hash,
			printable(prep.Tag))
	}
}

func (conn *DBConnT) MustQuery(sql DBSqlT, params DBParamsT, rcvr Receiver) {
	prep := conn.prepareSql(sql, len(params))
	t0 := time.Now()

	err := conn.Query(prep, params, rcvr)

	tQ := time.Now().Sub(t0).Nanoseconds()
	TQuery += tQ
	prep.Sample(tQ)
	if err != nil {
		panic("conn.MustQuery: " + err.Error())
		//log.Fatalf("conn.MustQuery: %#v %#v", err, conn)
	}
}

func (conn *DBConnT) MustExec(sql DBSqlT, params DBParamsT) *DBQueryResultT {
	prep := conn.prepareSql(sql, len(params))
	t0 := time.Now()

	err, res := conn.Exec(prep, params)

	tX := time.Now().Sub(t0).Nanoseconds()
	TExec += tX
	prep.Sample(tX)
	if err != nil {
		log.Printf("sql %s\n", sql)
		log.Printf("params %#v\n", params)
		log.Panicf("conn.MustExec: %#v", err)
		//panic("conn.MustExec: " + err.Error())
		//log.Fatalf("conn.MustExec: %#v %#v", err, conn)
	}
	return res
}

