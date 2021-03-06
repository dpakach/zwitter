package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/dpakach/zwitter/pkg/config"
	zlog "github.com/dpakach/zwitter/pkg/log"
)

type Hello struct {
	Log *zlog.ZwitLogger
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	rw.Write([]byte("Hello from the Web service"))
}

func NewHello(log *zlog.ZwitLogger) *Hello {
	return &Hello{log}
}

type Web struct {
	config    *config.WebServiceConfig
	indexPath string
	Log       *zlog.ZwitLogger
}

type saveFileOutput struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func (w *Web) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(w.config.AssetsPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(rw, r, filepath.Join(w.config.AssetsPath, w.indexPath))
		return
	} else if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(w.config.AssetsPath)).ServeHTTP(rw, r)
}

func NewWeb(config *config.WebServiceConfig, indexPath string, log *zlog.ZwitLogger) *Web {
	return &Web{config, indexPath, log}
}
