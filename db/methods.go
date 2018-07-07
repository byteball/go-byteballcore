package db

import(
	"database/sql"
	"context"
//	"fmt"
	"log"
	"strings"
//	"regexp"
	"time"

//	lt "db/wrappers"

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
}

//func (database *Database) MustSelect(sql DBSqlT, params DBParamsT, rcvr Scanner) {
func MustSelect(sql DBSqlT, params DBParamsT, rcvr Receiver) {
	database := Instance()
	t0 := time.Now()
	err := database.Select(sql, params, rcvr)
	TQuery += time.Now().Sub(t0).Nanoseconds()
	if err != nil {
		log.Printf("sql %s\n", sql)
		log.Printf("params %#v\n", params)
		log.Panicf("db.MustSelect: %#v", err)
	}
}

//

func (conn *DBConnT) BeginTransaction() {
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
//	conn.ExecContext = txn.ExecContext
//	conn.QueryContext = txn.QueryContext
}

func (conn *DBConnT) Commit() {
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
//	conn.ExecContext = conn.dbconn.ExecContext
//	conn.QueryContext = conn.dbconn.QueryContext
}

func (conn *DBConnT) Rollback() {
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
//	conn.ExecContext = conn.dbconn.ExecContext
//	conn.QueryContext = conn.dbconn.QueryContext
}

//

func (conn *DBConnT) ExecContext(ctx context.Context, sql_ string, params... interface{}) (sql.Result, error) {
	if conn.txn != nil {
		log_Printf("txn.Exec %s", strings.Split(sql_, "\n")[0])
		return conn.txn.ExecContext(ctx, sql_, params...)
	}
	log_Printf("conn.Exec %s", strings.Split(sql_, "\n")[0])
	return conn.dbconn.ExecContext(ctx, sql_, params...)
}

func (conn *DBConnT) QueryContext(ctx context.Context, sql_ string, params... interface{}) (*sql.Rows, error) {
	if conn.txn != nil {
		log_Printf("txn.Query %s", strings.Split(sql_, "\n")[0])
		return conn.txn.QueryContext(ctx, sql_, params...)
	}
	log_Printf("conn.Query %s", strings.Split(sql_, "\n")[0])
	return conn.dbconn.QueryContext(ctx, sql_, params...)
}

//

func (conn *DBConnT) MustQuery(sql DBSqlT, params DBParamsT, rcvr Receiver) {
	t0 := time.Now()
	err := conn.Query(sql, params, rcvr)
	TQuery += time.Now().Sub(t0).Nanoseconds()
	if err != nil {
		panic("conn.MustQuery: " + err.Error())
		//log.Fatalf("conn.MustQuery: %#v %#v", err, conn)
	}
}

func (conn *DBConnT) MustExec(sql DBSqlT, params DBParamsT) *DBQueryResultT {
	t0 := time.Now()
	err, res := conn.Exec(sql, params)
	TExec += time.Now().Sub(t0).Nanoseconds()
	if err != nil {
		log.Printf("sql %s\n", sql)
		log.Printf("params %#v\n", params)
		log.Panicf("conn.MustExec: %#v", err)
		//panic("conn.MustExec: " + err.Error())
		//log.Fatalf("conn.MustExec: %#v %#v", err, conn)
	}
	return res
}

//

