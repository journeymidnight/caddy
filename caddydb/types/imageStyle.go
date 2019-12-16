package types

type ImageStyle struct {
	Bucket    string
	StyleName string
	StyleCode string
}

func (i *ImageStyle) InsertStyle() (string, []interface{}) {
	sql := "insert into pipa(bucket_name, style_name, style_code) values(?,?,?)"
	args := []interface{}{i.Bucket, i.StyleName, i.StyleCode}
	return sql, args
}

func (i *ImageStyle) UpdateStyle() (string, []interface{}) {
	sql := "update pipa set style_code=? where bucket_name=? and style_name=?"
	args := []interface{}{i.StyleCode, i.Bucket, i.StyleName}
	return sql, args
}

func (i *ImageStyle) DeleteStyle() (string, []interface{}) {
	sql := "delete from pipa where bucket_name=? and style_name=?"
	args := []interface{}{i.Bucket, i.StyleName}
	return sql, args
}
