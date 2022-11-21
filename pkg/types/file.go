package types

import "time"

type File struct {
	Name         string    `json:"name" structs:"name"`
	MD5Hash      string    `json:"md5Hash" structs:"md5_hash"`
	Size         int64     `json:"size" structs:"size"`
	ModifiedTime time.Time `json:"modifiedTime" structs:"modified_time"`
}

func (f *File) Clone() *File {
	return &File{
		Name:         f.Name,
		MD5Hash:      f.MD5Hash,
		Size:         f.Size,
		ModifiedTime: f.ModifiedTime,
	}
}

type Bucket struct {
	UUID         string `json:"uid"`
	File         File   `json:"file"`
	Index        int64  `json:"index"`
	IsLastBucket bool   `json:"isLastBucket"`
	BucketSize   int64  `json:"bucketSize"`
	Data         []byte `json:"data"`
}

type FileResponse string

const (
	FileResponseCloned       FileResponse = "cloned"
	FileResponseUpToDate     FileResponse = "up_to_date"
	FileResponseNameExists   FileResponse = "name_exists"
	FileResponseNotAvailable FileResponse = "not_available"
)

const (
	OrderByASC = "asc"
	OrderByDsc = "dsc"
)

type FreqWordsRequest struct {
	Limit   uint   `json:"limit"`
	OrderBy string `json:"orderBy"`
}
