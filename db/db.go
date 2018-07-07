package db

import(
	"database/sql"
	// [fyi] register database driver
//	_ "github.com/lib/pq"
//	_ "github.com/mattn/go-sqlite3"
// _lt	"github.com/mattn/go-sqlite3"
 _lt	"github.com/mxk/go-sqlite/sqlite3"
	"context"

	"fmt"
	"log"
	//"bytes"
//	"math/big"
//	"crypto/rand"
	"regexp"
	"strings"

	"nodejs/console"

	"github.com/byteball/go-byteballcore/conf"
 .	"github.com/byteball/go-byteballcore/db/conf"
)

import(
// .	"github.com/byteball/go-byteballcore/types"
)

type(
	DBSqlT		= string

//	DBParamT	interface{}
	DBParamT	= interface{}
//	DBParamsT	= []DBParamT
	DBParamsT	[]DBParamT

	DBQueryT struct{
		Sql	DBSqlT
		Params	DBParamsT
		Cb	DBQueryCbT
	}
	DBQuerysT       = []DBQueryT
	refDBQuerysT	= *DBQuerysT

	DBConnT struct{
		dbconn	*sql.Conn
		ctx	context.Context
		txn	*sql.Tx
		querys	DBQuerysT
		preparedByTag  DBPreparedByTagT
	}
	refDBConnT	= *DBConnT

	DBQueryCbT	func (*DBQueryResultT)

	DBQueryResultT struct{
		AffectedRows int64
	}

//	DBRowT		interface{}
	DBRowsT		= []DBRowT
)

// ---

type(
	DBConnPoolT struct{
		items	chan refDBConnT
//		waiting	util.AtomicCounter
	}

	refDBConnPoolT	= *DBConnPoolT
)

func (pool *DBConnPoolT) Put(item refDBConnT) {
	pool.items <- item
}

func (pool *DBConnPoolT) Get() refDBConnT {
	item := <-pool.items
	return item
}

func mkDBConnPool(connMin int) refDBConnPoolT {
	pool := &DBConnPoolT{
		items: make(chan refDBConnT, connMin),
	}

        for k := 0; k < connMin; k++ {
                pool.Put(mkDBConn())
        }

	return pool
}

// ---

type Database struct{
	db		*sql.DB
	connPool	refDBConnPoolT
}

var(
	dbInstance *Database = nil
)

func Instance() *Database {
	return dbInstance
}

func Init() {
	dbInstance = &Database{}

	dbConf := conf.Instance().DB

	// [tbd] default conf

	dbConf.Filename = "byteball.sqlite"

	Instance().init(dbConf)
}

func (database *Database) init(dbConf Conf) {

//	ver, vn, sid := _lt.Version()
	ver, vn, sid := _lt.Version(), -1, "[* no source id *]"
	fmt.Printf("SQLlite ver. %s (%d), %s\n", ver, vn, sid)

/**
	dsn := fmt.Sprintf(""+
		"host=%s port=%d "+
		"user=%s password=%s dbname=%s"+
		" sslmode=disable",
		dbConf.Host, dbConf.Port,
		dbConf.User, dbConf.Password, dbConf.DBName)
 **/
	query := strings.Join([]string{
		"_busy_timeout=30000",
		"_journal=wal",
		"_synchronous=normal",
	}, "&")
	dsn := fmt.Sprintf("file:%s/%s?%s",
		conf.AppDirName(),
		dbConf.Filename,
		query)

	db, eopen := sql.Open("sqlite3", dsn)
	if eopen != nil {
		log.Printf("dsn %#v\n", dsn)
		log.Fatalf("sql.Open: %s", eopen.Error())
	}

	defer func() {
		if database.db == nil {
			db.Close()
		}
	}()

	// [fyi] serialize requests
	// [tbd] database.max_connections
	max_connections := 1
	db.SetMaxOpenConns(max_connections)

	if eping := db.Ping(); eping != nil {
		log.Printf("dsn %#v\n", dsn)
		log.Fatalf("db.ping: %s", eping.Error())
	}

	database.db = db

	database.connPool = mkDBConnPool(max_connections)

	for _, sql := range ([]string{
		"PRAGMA foreign_keys = 1",
		"PRAGMA busy_timeout=30000",
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA temp_store=MEMORY",
		"PRAGMA wal_autocheckpoint=32768",
	}) {
		console.Log("%s", sql)
		MustExec(sql, nil)
	}

	//log.Printf("db.init: %#v\n", db)
}


// ---

func mkDBConn() refDBConnT {
	ctx := context.Background()

	dbconn, err := Instance().db.Conn(ctx)
	if err != nil {
		log.Fatalf("db.Conn: %#v", err)
	}

	conn := DBConnT{
		dbconn: dbconn,
		ctx: ctx,
//		ExecContext: dbconn.ExecContext,
//		QueryContext: dbconn.QueryContext,
		preparedByTag: make(DBPreparedByTagT),
	}
	return &conn
}

// [fyi] must stick to single db connection
// [fyi] to avoid prepared statements' insanity

func TakeConnectionFromPool_sync() refDBConnT {
//	return mkDBConn()
	conn := Instance().connPool.Get()
	return conn
}

func (conn *DBConnT) Release() {
//	conn.dbconn.Close()
	Instance().connPool.Put(conn)
}

//

func (conn *DBConnT) ResetAddedQueries() {
	// [tbd] verify this
	conn.querys = make(DBQuerysT, 0, 99)
}

func (conn *DBConnT) AddQuery(sql DBSqlT, params DBParamsT) {
	conn.AddQueryCb(sql, params, nil)
}

func (conn *DBConnT) AddQueryCb(sql DBSqlT, params DBParamsT, cb DBQueryCbT) {
	query := DBQueryT{
		Sql: sql,
		Params: params,
		Cb: cb,
	}
	conn.querys = append(conn.querys, query)
}

var(
	sqlSelectRe = regexp.MustCompile("^\\s*SELECT")
)

func isSelectSql(sql DBSqlT) bool {
	return sqlSelectRe.MatchString(sql)
}

func (conn *DBConnT) ExecuteAddedQuerys() error {
	for _k, query := range conn.querys {
//		rows := conn.QueryCb_sync(query.Sql, query.Params, query.Cb)
//		rows = rows
		if !isSelectSql(query.Sql) {
//..			log.Printf("ExecuteAddedQueries[%d]: %s", _k, strings.Split(query.Sql,"\n")[0])
//			log.Printf("ExecuteAddedQueries[%d]: %s", _k, query.Sql)
//			log.Printf("ExecuteAddedQueries[%d]: %#v", _k, query.Params)
			_k = _k
			res := conn.MustExec(query.Sql, query.Params)
			if query.Cb != nil { query.Cb(res) }
		} else {
			panic("ExecuteAddedQuerys: SELECT")
		}
	}
	return nil
}

//

func (conn *DBConnT) GetFromUnixTime(ts string) string {
	return "datetime("+ts+", 'unixepoch')"
}

func (conn *DBConnT) GetIgnore() string {
	return "OR IGNORE "
}

func (conn *DBConnT) ForceIndex(index string) string {
	return "INDEXED BY " + index
}

func (conn *DBConnT) DropTemporaryTable(table string) string {
	return "DROP TABLE IF EXISTS " + table;
}

/***
	function addTime(interval){
		return "datetime('now', '"+interval+"')";
	}

	function getNow(){
		return "datetime('now')";
	}

	function getUnixTimestamp(date){
		return "strftime('%s', "+date+")";
	}

	function getFromUnixTime(ts){
		return "datetime("+ts+", 'unixepoch')";
	}

	function getRandom(){
		return "RANDOM()";
	}

	function forceIndex(index){
		return "INDEXED BY " + index;
	}

	function dropTemporaryTable(table) {
		return "DROP TABLE IF EXISTS " + table;
	}
 ***/
