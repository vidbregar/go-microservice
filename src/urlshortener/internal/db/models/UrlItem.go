package models

type UrlItem struct {
	Url      string `redis:"url" json:"url"`
	ExpireAt int64  `redis:"expireAt" json:"expireAt"`
}
