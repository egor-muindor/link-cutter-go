package models

import (
	"cutter-url-go/internal/vo"
)

type ShortLink struct {
	ShortLink vo.ShortURI `bson:"short_uri,omitempty"`
	FullURL   vo.FullURI  `bson:"full_uri,omitempty"`
}
