package application

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/qerr"
	"github.com/therecipe/qt/core"
)

var cacheDir = config.GetCacheDir()
var downloader = http.Client{
	Timeout: 1000 * time.Millisecond,
}
var audioCache = NewAudioCache()

func NewAudioCache() *AudioCache {
	return &AudioCache{
		m: map[string]*core.QUrl{},
	}
}

type AudioCache struct {
	m     map[string]*core.QUrl
	mlock sync.RWMutex
}

func (c *AudioCache) download(urlStr string, fpath string) error {
	res, err := downloader.Get(urlStr)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(fpath), 0o755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fpath, data, 0o644)
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
	if cacheDir == "" {
		return nil, fmt.Errorf("cacheDir is empty")
	}
	fpath := filepath.Join(
		cacheDir,
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
