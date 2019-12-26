package types

type Imagestyles struct {
	ImageStyle []ImageStyle
}

type ImageStyle struct {
	Bucket    string
	StyleName string
	StyleCode string
}

func (i *ImageStyle) InsertStyle() (string, []interface{}) {
	sql := "insert into images_style(bucket_name, style_name, style_code) values(?,?,?)"
	args := []interface{}{i.Bucket, i.StyleName, i.StyleCode}
	return sql, args
}

func (i *ImageStyle) UpdateStyle() (string, []interface{}) {
	sql := "update images_style set style_code=? where bucket_name=? and style_name=?"
	args := []interface{}{i.StyleCode, i.Bucket, i.StyleName}
	return sql, args
}

func (i *ImageStyle) DeleteStyle() (string, []interface{}) {
	sql := "delete from images_style where bucket_name=? and style_name=?"
	args := []interface{}{i.Bucket, i.StyleName}
	return sql, args
}
