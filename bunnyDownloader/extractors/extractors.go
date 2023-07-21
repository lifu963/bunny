package extractors

import (
	"app/request"
	"app/utils"
	"github.com/pkg/errors"
	"net/url"
	"strings"
	"sync"
)

var lock sync.RWMutex
var extractorMap = make(map[string]Extractor)

// Register registers an Extractor.
func Register(domain string, e Extractor) {
	lock.Lock()
	extractorMap[domain] = e
	lock.Unlock()
}

func Extract(u string, option Options, requester *request.Requester) ([]*Data, error) {
	u = strings.TrimSpace(u)
	var domain string
	bilibiliShortLink := utils.MatchOneOf(u, `^(av|BV|ep)\w+`)
	if len(bilibiliShortLink) > 1 {
		bilibiliURL := map[string]string{
			"av": "https://www.bilibili.com/video/",
			"BV": "https://www.bilibili.com/video/",
			"ep": "https://www.bilibili.com/bangumi/play/",
		}
		domain = "bilibili"
		u = bilibiliURL[bilibiliShortLink[1]] + u
	} else {
		u, err := url.ParseRequestURI(u)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		domain = utils.Domain(u.Host)
	}
	extractor := extractorMap[domain]
	if extractor == nil {
		extractor = extractorMap[""]
	}
	videos, err := extractor.Extract(u, option, requester)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, v := range videos {
		v.FillUpStreamsData()
	}
	return videos, nil
}
