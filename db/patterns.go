// [fyi] high-level abstractions for SQL client

package db

import(
	"database/sql"
//	"fmt"
	"log"
)

type Receiver interface{
	Scan(*sql.Rows) error
}

//

func (database *Database) Exec (sql string, params []interface{}) (error, *DBQueryResultT) {

	//log.Printf("db.Exec: sql %#v %#v", sql, params)

	qry, err := database.db.Prepare(sql)
	if err != nil {
		//log.Printf("Exec.Prepare: %s", err.Error())
		return err, nil
	}

	res, err := qry.Exec(params...)
	if err != nil {
		//log.Printf("Exec.Exec: %s", err.Error())
		return err, nil
	}

	ars, _ := res.RowsAffected()
	//log.Printf(">>> ars %#v\n", ars)
	res_ := DBQueryResultT{
		AffectedRows: ars,
	}
	return err, &res_
}

func (database *Database) Select(sql string, params []interface{}, rcvr Receiver) error {

	qry, err := database.db.Prepare(sql)
	if err != nil {
		//log.Printf("Select.Prepare: %s", err.Error())
		return err
	}

	rows, err := qry.Query(params...)
	if err != nil {
		//log.Printf("Select.Query: %s", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rcvr.Scan(rows)
		if err != nil {
			//log.Printf("Select.Scan: %s", err.Error())
			return err
		}
	}

	return nil
}

//

func (conn *DBConnT) Exec (prep refDBPreparedT, params []interface{}) (error, *DBQueryResultT) {

	//log.Printf("conn.Exec: sql %#v %#v", prep.Sql, params)

	res, err := conn.ExecPrepared(prep, params)
	if err != nil {
		log.Printf("conn.ExecPrepared: %s", err.Error())
		return err, nil
	}

	ars, _ := res.RowsAffected()
	//log.Printf(">>> ars %#v\n", ars)
	res_ := DBQueryResultT{
		AffectedRows: ars,
	}
	return err, &res_
}

func (conn *DBConnT) Query(prep refDBPreparedT, params []interface{}, rcvr Receiver) error {

	//log.Printf("conn.Query: sql %#v %#v", prep.Sql, params)

	rows, err := conn.QueryPrepared(prep, params)
	if err != nil {
		log.Printf("conn.QueryPrepared: %s", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rcvr.Scan(rows)
		if err != nil {
			log.Printf("conn.Query.Scan: %s", err.Error())
			return err
		}
	}

	return nil
}
