package jobs

import (
	"encoding/json"
	"fmt"
)

type iSocialPlatformFactory interface {
	MakeSpider(data []byte) iSpider
}

func GetSocialFactory(platform string) (iSocialPlatformFactory, error) {
	if platform == "dcard" {
		return &dcardFactory{}, nil
	}
	if platform == "ptt" {
		return &pttFactory{}, nil
	}

	return nil, fmt.Errorf("Wrong brand type passed")
}

type dcardFactory struct {
}

func (a *dcardFactory) MakeSpider(data []byte) iSpider {
	var dcardSpider DcardSpider
	json.Unmarshal(data, &dcardSpider)
	return &dcardSpider
}

type pttFactory struct {
}

func (a *pttFactory) MakeSpider(data []byte) iSpider {
	var pttSpider PttSpider
	json.Unmarshal(data, &pttSpider)
	return &pttSpider
}
