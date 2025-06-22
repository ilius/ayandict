package application

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilius/ayandict/v3/pkg/config"
	qt "github.com/mappu/miqt/qt6"
)

var (
	cacheDir   = config.GetCacheDir()
	audioCache = NewAudioCache()
)

func NewAudioCache() *AudioCache {
	dir := ""
	if cacheDir == "" {
		slog.Error("cacheDir is empty")
	} else {
		dir = filepath.Join(cacheDir, "audio")
	}

	return &AudioCache{
		m:   map[string]*qt.QUrl{},
		dir: dir,

		downloader: &http.Client{
			Timeout: conf.AudioDownloadTimeout,
		},
	}
}

type AudioCache struct {
	m     map[string]*qt.QUrl
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
	slog.Debug("Downloaded audio", "fpath", fpath)
	return nil
}

func (c *AudioCache) Get(urlStr string) (*qt.QUrl, error) {
	c.mlock.RLock()
	qUrl := c.m[urlStr]
	c.mlock.RUnlock()
	if qUrl != nil {
		return qUrl, nil
	}
	qUrl = qt.NewQUrl4(urlStr, qt.QUrl__TolerantMode)
	host := qUrl.Host()
	// also add possible port?
	if c.dir == "" {
		return nil, fmt.Errorf("audio cache dir is empty")
	}
	fpath := filepath.Join(
		c.dir,
		host,
		qUrl.Path(),
	)
	qUrl.SetScheme("file")
	qUrl.SetHost("")
	qUrl.SetPath2(fpath, qt.QUrl__DecodedMode)
	if _, err := os.Stat(fpath); err != nil {
		if !os.IsNotExist(err) {
			slog.Error("error in Stat: "+err.Error(), "fpath", fpath)
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
