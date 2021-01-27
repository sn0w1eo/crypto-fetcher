package mysql

import (
	"database/sql"
	"fmt"
	"github.com/Sn0w1eo/crypto-fetcher/src/crypto"
	_ "github.com/go-sql-driver/mysql"
)

const DBName = "CryptoFetcher"
const TableTick = "Tick"

// MySQLConn struct implemented to satisfy Storage interface requirements
type MySQLConn struct {
	db *sql.DB
}

// Creates MySQLConn
func New() (*MySQLConn, error) {
	c := new(MySQLConn)
	return c, nil
}

// Opens connection based on DSN and runs init() function
func (mysqlConn *MySQLConn) Open(dsn string) (err error) {
	mysqlConn.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return mysqlConn.init()
}

// init used to invoke some required methods after successful connection to DB
func (mysqlConn *MySQLConn) init() error {
	tx, err := mysqlConn.db.Begin()
	if err != nil {
		return err
	}
	{
		stmt, err := tx.Prepare("CREATE DATABASE IF NOT EXISTS CryptoFecther;")
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				panic(rbErr)
			}
			return err
		}
		stmt.Close()
	}
	{
		scheme := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS CryptoFetcher.Ticks ( 
		%s integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
		%s BIGINT NOT NULL,
		%s VARCHAR(255) NOT NULL,
		%s DOUBLE NOT NULL,
		%s DOUBLE NOT NULL
		);`, "`id`", "`timestamp`", "`symbol`", "`bid`", "`ask`")
		stmt, err := tx.Prepare(scheme)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				panic(rbErr)
			}
			return err
		}
		stmt.Close()
	}
	return tx.Commit()
}

// Write crypto.Ticker to DB
func (mysqlConn *MySQLConn) WriteTick(ticker crypto.Ticker) error {
	timestamp := ticker.Timestamp()
	pair, _ := ticker.Pair()
	bestBid, _ := ticker.BestBid()
	bestAsk, err := ticker.BestAsk()

	symbol := pair.String('-')

	rows, err := mysqlConn.db.Query("INSERT INTO CryptoFetcher.Ticks (`timestamp`, `symbol`, `bid`, `ask`) VALUES(?, ?, ?, ?);", timestamp.Unix(), symbol, bestBid, bestAsk)
	if err != nil {
		return err
	}
	return rows.Close()
}

// Closes connection
func (mysqlConn *MySQLConn) Close() error {
	return mysqlConn.db.Close()
}
