package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/0w0mewo/line-sticker-downloader/linesticker"
	"github.com/0w0mewo/line-sticker-downloader/utils"
	log "github.com/sirupsen/logrus"
)

var dir string
var packs string

func init() {
	log.SetFormatter(
		&log.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		},
	)

	flag.StringVar(&dir, "path", ".", "where zipped sticker packages saved")
	flag.StringVar(&packs, "packs", "", "list of sticker packs (splited by ','), e.g.: 1234,5678,9000")
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

	tempDir, err := ioutil.TempDir(".", "sticker-*")
	if err != nil {
		log.Panicln(err)
	}
	defer os.RemoveAll(tempDir)

	// fetch list of sticker packages
	for _, packid := range packIds {
		wg.Add(1)
		go func(pack int) {
			defer wg.Done()
			stickerPack := linesticker.NewFetcher(ctx, http.DefaultClient)

			stickerPack.SetPackId(pack)
			stickerPack.SaveStickers(filepath.Join(tempDir, strconv.Itoa(pack)))
		}(packid)
	}

	wg.Wait()
	cancel()

	// zipping downloaded stickers
	z := utils.NewZip(filepath.Join(dir, filepath.Base(tempDir)))
	if err := z.Zip(tempDir); err != nil {
		log.Error(err)
		return
	}
	defer z.Close()

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
