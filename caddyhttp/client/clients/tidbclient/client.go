package tidbclient

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

type DBInfo struct {
	S3Source             string
	CaddySource          string
	DBMaxIdleConns       int
	DBMaxOpenConns       int
	DBConnMaxLifeSeconds int
}

type TidbClient struct {
	ClientS3       *sql.DB
	ClientBusiness *sql.DB
}

func NewTidbClient(db DBInfo) *TidbClient {
	cli := &TidbClient{}
	s3Conn, err := sql.Open("mysql", db.S3Source)
	if err != nil {
		os.Exit(1)
	}
	s3Conn.SetMaxIdleConns(db.DBMaxIdleConns)
	s3Conn.SetMaxOpenConns(db.DBMaxOpenConns)
	s3Conn.SetConnMaxLifetime(time.Duration(db.DBConnMaxLifeSeconds) * time.Second)
	cli.ClientS3 = s3Conn
	bsConn, err := sql.Open("mysql", db.CaddySource)
	if err != nil {
		os.Exit(1)
	}
	bsConn.SetMaxIdleConns(db.DBMaxIdleConns)
	bsConn.SetMaxOpenConns(db.DBMaxOpenConns)
	bsConn.SetConnMaxLifetime(time.Duration(db.DBConnMaxLifeSeconds) * time.Second)
	cli.ClientBusiness = bsConn
	return cli
}
