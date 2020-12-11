package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/mux"
	"github.com/soheilhy/cmux"
	"github.com/zhenshaw/tesseract-ocr/orc"
)

var apiToken = ""

func HTTPServer(l net.Listener, token string) {
	apiToken = token
	r := mux.NewRouter()
	r = r.PathPrefix(os.Getenv("PREFIX")).Subrouter()
	r.HandleFunc("/recognize", recognize)

	s := &http.Server{
		Handler: r,
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

func recognize(w http.ResponseWriter, r *http.Request) {

	if r.URL.Query().Get("token") != apiToken {
		err := fmt.Errorf("rpc unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	rFile, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer rFile.Close()
	data, err := ioutil.ReadAll(rFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	res, err := orc.Recognize(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	logs.Debug("识别结果：", res)
	_, _ = w.Write([]byte(res))
}
