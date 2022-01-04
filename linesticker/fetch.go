package linesticker

import (
	"context"
	"fmt"
	"image/gif"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/0w0mewo/line-sticker-downloader/utils"
	"github.com/sirupsen/logrus"
)

// sticker package fetcher
type Fetcher struct {
	ctx    context.Context
	pack   *PackMeta
	client *http.Client
	logger *logrus.Logger
}

// new sticker package fetcher
func NewFetcher(ctx context.Context, client *http.Client) *Fetcher {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	ret := &Fetcher{
		client: client,
		logger: logger,
		ctx:    ctx,
	}

	return ret
}

func (wk *Fetcher) getmeta(packid int) {
	url := fmt.Sprintf(METAINFOURL, packid)
	wk.pack = &PackMeta{}

	err := utils.HttpGetJson(wk.ctx, wk.client, url, wk.pack)
	if err != nil {
		wk.logger.Panicln(err)

	}

}

// set sticker pack id and refetch package meta
func (wk *Fetcher) SetPackId(packid int) {
	wk.getmeta(packid)
}

// get stickers pack id
func (wk *Fetcher) GetPackId() int {
	return wk.pack.PackageId
}

// fetch stickers package and save it to saveToDir
func (wk *Fetcher) SaveStickers(saveToDir string) {
	packid := wk.pack.PackageId
	animated := wk.pack.HasGif

	stickerStorer := func(r io.Reader, s *Sticker) error {
		var fileExt string
		var folderName string

		if animated {
			fileExt = "gif"
			folderName = "animated"
		} else {
			fileExt = "gif"
			folderName = "not-animated"
		}

		stickerFolder := filepath.Join(saveToDir, folderName)
		path := filepath.Join(stickerFolder, strconv.Itoa(s.Id)+"."+fileExt)

		// make dir to store stickers
		err := os.MkdirAll(stickerFolder, 0o755)
		if err != nil {
			return err
		}

		// open file for writing sticker bytes
		fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			return err
		}
		defer fd.Close()

		// convert to gif, regardless it is animated or not
		if animated {
			_, err = io.Copy(fd, r)
			return err
		} else {

			pngImg, err := png.Decode(r)
			if err != nil {
				return err
			}

			return gif.Encode(fd, pngImg, nil)

		}

	}

	// fetch and save stickers pack
	for _, s := range wk.pack.Stickers {
		err := s.Fetch(wk.client, packid, animated, func(r io.Reader) error {
			wk.logger.Infof("downloading %d belongs to %d", s.Id, packid)
			return stickerStorer(r, s)
		})
		if err != nil {
			wk.logger.Errorln(err)
			continue
		}
	}

	wk.logger.Infoln("done on fetching: ", packid)
}
