package web

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"office-convert/office"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const PORT = 10000
const SAVED_DIR = "files"

type ConvertRequestInfo struct {
	FileInUrl  string `json:"file_in_url"`
	SourceType string `json:"source_type"`
	TargetType string `json:"target_type"`
}

func logStackTrace(err ...interface{}) {
	log.Println(err)
	stack := string(debug.Stack())
	log.Println(stack)
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(503)
			fmt.Fprintln(w, r)
			logStackTrace(r)
		}
	}()
	if r.Method != "POST" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Method not support")
		return
	}

	var convertRequestInfo ConvertRequestInfo
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(reqBody, &convertRequestInfo)

	log.Println(convertRequestInfo)
	log.Println(convertRequestInfo.FileInUrl)

	downloadFile(convertRequestInfo.FileInUrl)

	fileOutAbsPath := getFileOutAbsPath(convertRequestInfo.FileInUrl, convertRequestInfo.TargetType)
	convert(convertRequestInfo)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	//文件过大的话考虑使用io.Copy进行流式拷贝
	outFileBytes, err := ioutil.ReadFile(fileOutAbsPath)
	if err != nil {
		panic(err)
	}
	w.Write(outFileBytes)

}

func convert(convertRequestInfo ConvertRequestInfo) {

	fileOutAbsPath := getFileOutAbsPath(convertRequestInfo.FileInUrl, convertRequestInfo.TargetType)
	switch convertRequestInfo.SourceType {
	case "doc", "docx":
		office.ConvertDoc2Pdf(getFileInAbsPath(convertRequestInfo.FileInUrl), fileOutAbsPath)
		break
	case "xls", "xlsx":
		office.ConvertXsl2Pdf(getFileInAbsPath(convertRequestInfo.FileInUrl), fileOutAbsPath)
		break
	case "ppt", "pptx":
		office.ConvertPPT2Pdf(getFileInAbsPath(convertRequestInfo.FileInUrl), fileOutAbsPath)
		break
	}
}

func getNameFromUrl(inputUrl string) string {
	u, err := url.Parse(inputUrl)
	if err != nil {
		panic(err)
	}
	return path.Base(u.Path)
}

func getCurrentWorkDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

func getFileInAbsPath(url string) string {
	fileName := getNameFromUrl(url)
	currentWorkDirectory := getCurrentWorkDirectory()
	absPath := filepath.Join(currentWorkDirectory, SAVED_DIR, fileName)
	return absPath
}

func getFileOutAbsPath(fileInUrl string, targetType string) string {
	return getFileInAbsPath(fileInUrl) + "." + targetType
}

func downloadFile(url string) {
	log.Println("Start download file url :", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fileInAbsPath := getFileInAbsPath(url)
	dir := filepath.Dir(fileInAbsPath)
	// log.Println("dir is " + dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println("dir is not exists")
		os.MkdirAll(dir, 0644)
	}
	out, err := os.Create(fileInAbsPath)
	log.Println("save file to " + fileInAbsPath)
	if err != nil {
		panic(err)
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	log.Println("Download file end url :", url)
}

func StartServer() {

	log.Println("start service ...")
	http.HandleFunc("/convert", convertHandler)
	http.ListenAndServe("127.0.0.1:"+strconv.Itoa(PORT), nil)
	// convertHandler中有recovery，不用下面的中间件了

	// r := http.NewServeMux()
	// r.HandleFunc("/convert", convertHandler)
	// http.ListenAndServe("127.0.0.1:"+strconv.Itoa(PORT),
	// 	RecoveryHandler(RecoveryLogger(log.StandardLogger()), PrintRecoveryStack(true))(r))
}
