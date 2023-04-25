package application

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/qerr"
	"github.com/ilius/qt/core"
)

var (
	cacheDir   = config.GetCacheDir()
	audioCache = NewAudioCache()
)

func NewAudioCache() *AudioCache {
	dir := ""
	if cacheDir == "" {
		qerr.Error("cacheDir is empty")
	} else {
		dir = filepath.Join(cacheDir, "audio")
	}

	return &AudioCache{
		m:   map[string]*core.QUrl{},
		dir: dir,

		downloader: &http.Client{
			Timeout: conf.AudioDownloadTimeout,
		},
	}
}

type AudioCache struct {
	m     map[string]*core.QUrl
	mlock sync.RWMutex
	dir   string

	downloader *http.Client
}

func (c *AudioCache) ReloadConfig() {
	c.downloader.Timeout = conf.AudioDownloadTimeout
}

func (c *AudioCache) download(urlStr string, fpath string) error {
	res, err := c.downloader.Get(urlStr)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(fpath), 0o755)
	if err != nil {
		return err
	}
	err = os.WriteFile(fpath, data, 0o644)
	if err != nil {
		return err
	}
	log.Printf("Downloaded into %#v", fpath)
	return nil
}

func (c *AudioCache) Get(urlStr string) (*core.QUrl, error) {
	c.mlock.RLock()
	qUrl := c.m[urlStr]
	c.mlock.RUnlock()
	if qUrl != nil {
		return qUrl, nil
	}
	qUrl = core.NewQUrl3(urlStr, core.QUrl__TolerantMode)
	host := qUrl.Host(core.QUrl__FullyEncoded)
	// also add possible port?
	if c.dir == "" {
		return nil, fmt.Errorf("audio cache dir is empty")
	}
	fpath := filepath.Join(
		c.dir,
		host,
		qUrl.Path(core.QUrl__FullyEncoded),
	)
	qUrl.SetScheme("file")
	qUrl.SetHost("", core.QUrl__DecodedMode)
	qUrl.SetPath(fpath, core.QUrl__DecodedMode)
	if _, err := os.Stat(fpath); err != nil {
		if !os.IsNotExist(err) {
			qerr.Error(err)
		}
		err := c.download(urlStr, fpath)
		if err != nil {
			return nil, err
		}
	}
	c.mlock.Lock()
	c.m[urlStr] = qUrl
	c.mlock.Unlock()
	return qUrl, nil
}
