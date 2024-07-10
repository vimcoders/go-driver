package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func (x *Handler) Dowload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cmd := exec.Command("go", "build", "-o", "e:/go-driver/bin/proxy.exe", "e:/go-driver/proxy/main.go") // 要执行的命令为 gm
	if _, err := cmd.Output(); err != nil {
		return
	}
	data, err := CompactServer()
	if err != nil {
		return
	}
	filename := fmt.Sprintf("\"server-%s.zip\"", time.Now().Format("2006年01月02日15时04分05秒"))
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Write(data)
}

func CompactServer() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	zw := zip.NewWriter(b)

	for _, name := range []string{"proxy.exe"} {
		f, err := zw.Create("server/" + name)
		if err != nil {
			return nil, err
		}
		data, err := os.ReadFile("e:/parkour/bin/" + name)
		if err != nil {
			return nil, err
		}
		if _, err := f.Write(data); err != nil {
			return nil, err
		}
	}
	zw.Close()
	return b.Bytes(), nil
}
