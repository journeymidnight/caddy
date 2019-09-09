package tidbclient

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

const MAX_OPEN_CONNS = 1024

type TidbClient struct {
	ClientS3       *sql.DB
	ClientBusiness *sql.DB
}

func NewTidbClient(s3Source string, businessSource string) *TidbClient {
	cli := &TidbClient{}
	s3Conn, err := sql.Open("mysql", s3Source)
	if err != nil {
		os.Exit(1)
	}
	s3Conn.SetMaxIdleConns(0)
	s3Conn.SetMaxOpenConns(MAX_OPEN_CONNS)
	cli.ClientS3 = s3Conn
	bsConn, err := sql.Open("mysql", businessSource)
	if err != nil {
		os.Exit(1)
	}
	bsConn.SetMaxIdleConns(0)
	bsConn.SetMaxOpenConns(MAX_OPEN_CONNS)
	cli.ClientBusiness = bsConn
	return cli
}
