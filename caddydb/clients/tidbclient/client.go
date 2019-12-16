package tidbclient

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

type DBInfo struct {
	DBMaxIdleConns       int
	DBMaxOpenConns       int
	DBConnMaxLifeSeconds int
}

type TidbClient struct {
	Client *sql.DB
}

func NewTidbClient(dbsource string, db DBInfo) *TidbClient {
	cli := &TidbClient{}
	conn, err := sql.Open("mysql", dbsource)
	if err != nil {
		os.Exit(1)
	}
	conn.SetMaxIdleConns(db.DBMaxIdleConns)
	conn.SetMaxOpenConns(db.DBMaxOpenConns)
	conn.SetConnMaxLifetime(time.Duration(db.DBConnMaxLifeSeconds) * time.Second)
	cli.Client = conn
	return cli
}
