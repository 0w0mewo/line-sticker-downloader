package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/0w0mewo/line-sticker-downloader/linesticker"
	log "github.com/sirupsen/logrus"
)

var dir string
var packs string
var qqTrans bool

func init() {
	log.SetFormatter(
		&log.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		},
	)

	flag.StringVar(&dir, "path", ".", "where zipped sticker packages saved")
	flag.StringVar(&packs, "packs", "", "list of sticker packs (splited by ','), e.g.: 1234,5678,9000")
	flag.BoolVar(&qqTrans, "qqtrans", false, "whether support transparency of non-animated when import to qq")
	flag.Parse()

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	packIds, err := parsePackList(packs)
	if err != nil {
		log.Errorf("fail to parse stickers packages list: %v", err)
		return
	}

	// fetch list of sticker packages
	for _, packid := range packIds {
		wg.Add(1)
		go func(pack int) {
			defer wg.Done()

			// fetch only when the sticker package is not cached
			zipfilePath := filepath.Join(dir, strconv.Itoa(pack)+".zip")
			if isZipFileExist(zipfilePath) {
				log.Infof("%d is cached to %s", pack, dir)
				return
			}

			stickerPack := linesticker.NewFetcher(ctx, http.DefaultClient)
			stickerPack.SetPackId(pack)
			stickerPack.SaveStickers(zipfilePath, qqTrans)
		}(packid)
	}

	wg.Wait()
	cancel()

	log.Info("exit...")

}

func parsePackList(packs string) ([]int, error) {
	splited := strings.Split(strings.TrimSpace(packs), ",")
	packids := make([]int, 0)

	for _, p := range splited {
		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		packids = append(packids, id)
	}

	return packids, nil
}

func isZipFileExist(path string) bool {
	zipf, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	_, exist := os.Stat(zipf)

	return exist == nil
}
