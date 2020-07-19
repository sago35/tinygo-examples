package tweet

import (
	"github.com/buger/jsonparser"
)

type Tweet2 struct {
	Id            int64
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
	var err error
	ret := &Tweet2{}

	ret.Id, err = jsonparser.GetInt(b, "Id")
	if err != nil {
		return nil, err
	}

	ret.UserName, err = jsonparser.GetString(b, "UserName")
	if err != nil {
		return nil, err
	}

	ret.ScreenName, err = jsonparser.GetString(b, "ScreenName")
	if err != nil {
		return nil, err
	}

	ret.CreatedAt, err = jsonparser.GetString(b, "CreatedAt")
	if err != nil {
		return nil, err
	}

	ret.FullText, err = jsonparser.GetString(b, "FullText")
	if err != nil {
		return nil, err
	}

	ret.FavoriteCount, err = jsonparser.GetInt(b, "FavoriteCount")
	if err != nil {
		return nil, err
	}

	ret.RetweetCount, err = jsonparser.GetInt(b, "RetweetCount")
	if err != nil {
		return nil, err
	}

	ret.IsRetweet, err = jsonparser.GetBoolean(b, "IsRetweet")
	if err != nil {
		return nil, err
	}

	jsonparser.ArrayEach(b, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ret.Entities = append(ret.Entities, string(value))
	}, "Entities")

	return ret, nil
}
