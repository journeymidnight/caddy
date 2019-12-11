package tidbclient

import (
	"database/sql"
	. "github.com/journeymidnight/yig-front-caddy/caddyerrors"
	. "github.com/journeymidnight/yig-front-caddy/caddyhttp/clients/types"
)

func (DB *TidbClient) InsertStyle(style ImageStyle) (err error) {
	var sqlTx *sql.Tx
	var tx interface{}
	tx, err = DB.ClientBusiness.Begin()
	if err != nil {
		return ErrSqlTransaction
	}
	defer func() {
		if err == nil {
			err = sqlTx.Commit()
		}
		if err != nil {
			sqlTx.Rollback()
		}
	}()
	sqlTx, _ = tx.(*sql.Tx)
	sqlL, args := style.InsertStyle()
	_, err = sqlTx.Exec(sqlL, args...)
	if err != nil {
		return ErrSqlInsert
	}
	return nil
}

func (DB *TidbClient) DelStyle(style ImageStyle) (err error) {
	var sqlTx *sql.Tx
	var tx interface{}
	tx, err = DB.ClientBusiness.Begin()
	if err != nil {
		return ErrSqlTransaction
	}
	defer func() {
		if err == nil {
			err = sqlTx.Commit()
		}
		if err != nil {
			sqlTx.Rollback()
		}
	}()
	sqlTx, _ = tx.(*sql.Tx)
	sql, args := style.DeleteStyle()
	_, err = sqlTx.Exec(sql, args...)
	if err != nil {
		return ErrSqlDelete
	}
	return nil
}

func (DB *TidbClient) UpdateStyle(style ImageStyle) (err error) {
	var sqlTx *sql.Tx
	var tx interface{}
	tx, err = DB.ClientBusiness.Begin()
	if err != nil {
		return ErrSqlTransaction
	}
	defer func() {
		if err == nil {
			err = sqlTx.Commit()
		}
		if err != nil {
			sqlTx.Rollback()
		}
	}()
	sqlTx, _ = tx.(*sql.Tx)
	sqlL, args := style.UpdateStyle()
	_, err = sqlTx.Exec(sqlL, args...)
	if err != nil {
		return ErrSqlUpdate
	}
	return nil
}

func (DB *TidbClient) GetStyles(bucket string) (styles []ImageStyle, err error) {
	sql := "select bucket_name,style_name,IFNULL(style_code,'') from pipa where bucket_name=?"
	args := []interface{}{bucket}
	rows, err := DB.ClientBusiness.Query(sql, args...)
	if err != nil {
		return styles, ErrInvalidSql
	}
	for rows.Next() {
		IImageStyle := ImageStyle{}
		err = rows.Scan(&IImageStyle.Bucket, &IImageStyle.StyleName, &IImageStyle.StyleCode)
		styles = append(styles, IImageStyle)
	}
	if err != nil {
		return styles, ErrNoSuchKey
	}
	defer rows.Close()
	return
}

func (DB *TidbClient) GetStyle(bucket string, styleName string) (style ImageStyle, err error) {
	sql := "select bucket_name,style_name,IFNULL(style_code,'') from pipa where bucket_name=? and style_name=?"
	args := []interface{}{bucket, styleName}
	rows, err := DB.ClientBusiness.Query(sql, args...)
	if err != nil {
		return style, ErrInvalidSql
	}
	for rows.Next() {
		err = rows.Scan(&style.Bucket, &style.StyleName, &style.StyleCode)
	}
	if err != nil {
		return style, ErrNoSuchKey
	}
	defer rows.Close()
	return
}
