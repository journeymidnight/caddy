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
	ClientS3       *sql.DB
	ClientBusiness *sql.DB
}

func NewCustomDomainClient(s3Source string, caddySource string, db DBInfo) *TidbClient {
	cli := &TidbClient{}
	cli.ClientS3 = NewTidbClient(s3Source, db)
	cli.ClientBusiness = NewTidbClient(caddySource, db)
	return cli
}

func NewTidbClient(dbsource string, db DBInfo) *sql.DB {
	conn, err := sql.Open("mysql", dbsource)
	if err != nil {
		os.Exit(1)
	}
	conn.SetMaxIdleConns(db.DBMaxIdleConns)
	conn.SetMaxOpenConns(db.DBMaxOpenConns)
	conn.SetConnMaxLifetime(time.Duration(db.DBConnMaxLifeSeconds) * time.Second)
	return conn
}
