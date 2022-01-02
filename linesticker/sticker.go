package linesticker

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/0w0mewo/line-sticker-downloader/utils"
)

const PNGURL = "http://dl.stickershop.line.naver.jp/stickershop/v1/sticker/%d/iphone/sticker@2x.png"
const GIFURL = "https://sdl-stickershop.line.naver.jp/products/0/0/1/%d/iphone/animation/%d@2x.png"
const METAINFOURL = "http://dl.stickershop.line.naver.jp/products/0/0/1/%d/android/productInfo.meta"

type PackMeta struct {
	PackageId int        `json:"packageId"`
	Stickers  []*Sticker `json:"stickers"`
	HasGif    bool       `json:"hasAnimation"`
}

type Sticker struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	Id     int `json:"id"`
}

func (s *Sticker) Key() string {
	return fmt.Sprintf("%d", s.Id)
}

func (s *Sticker) Fetch(client *http.Client, packid int, isAnimated bool, fn func(r io.Reader) error) error {
	var stickerUrl string

	if isAnimated {
		stickerUrl = fmt.Sprintf(GIFURL, packid, s.Id)
	} else {
		stickerUrl = fmt.Sprintf(PNGURL, s.Id)
	}

	return utils.HttpGetWithProcessor(context.Background(), client, stickerUrl, fn)

}
