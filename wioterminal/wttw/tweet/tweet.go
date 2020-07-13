package tweet

import (
	"github.com/buger/jsonparser"
)

type Tweet2 struct {
	UserName      string
	ScreenName    string
	CreatedAt     string
	FullText      string
	FavoriteCount int64
	RetweetCount  int64
	IsRetweet     bool
	Entities      []string
}

func NewTweet2(b []byte) (*Tweet2, error) {
	ret := &Tweet2{}

	username, err := jsonparser.GetString(b, "UserName")
	if err != nil {
		return nil, err
	}
	ret.UserName = username

	screenName, err := jsonparser.GetString(b, "ScreenName")
	if err != nil {
		return nil, err
	}
	ret.ScreenName = screenName

	createdAt, err := jsonparser.GetString(b, "CreatedAt")
	if err != nil {
		return nil, err
	}
	ret.CreatedAt = createdAt

	fullText, err := jsonparser.GetString(b, "FullText")
	if err != nil {
		return nil, err
	}
	ret.FullText = fullText

	favoriteCount, err := jsonparser.GetInt(b, "FavoriteCount")
	if err != nil {
		return nil, err
	}
	ret.FavoriteCount = favoriteCount

	retweetCount, err := jsonparser.GetInt(b, "RetweetCount")
	if err != nil {
		return nil, err
	}
	ret.RetweetCount = retweetCount

	isRetweet, err := jsonparser.GetBoolean(b, "IsRetweet")
	if err != nil {
		return nil, err
	}
	ret.IsRetweet = isRetweet

	jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ret.Entities = append(ret.Entities, string(value))
	}, "Entities")

	return ret, nil
}
