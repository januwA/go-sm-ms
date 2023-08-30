package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/januwA/go-ajahttp"
)

const BASIC_URL = "https://sm.ms/api/v2"

var smmsClint *ajahttp.AjaClient

func main() {
	// 定义一个命令行参数
	port := flag.String("port", "8080", "端口号")
	token := flag.String("token", "", "SM_MS_TOKEN")
	smmsUser := flag.String("user", "", "sm sm username")
	smmsPassword := flag.String("password", "", "sm sm password")

	// 解析命令行的参数
	flag.Parse()

	smmsClint = ajahttp.NewAjaClient()
	smmsClint.SetBaseURL(BASIC_URL)

	if *token == "" {
		SM_MS_TOKEN := os.Getenv("SM_MS_TOKEN")
		if SM_MS_TOKEN != "" {
			*token = SM_MS_TOKEN
		} else if *smmsUser != "" && *smmsPassword != "" {
			usertoken, err := smmsAuthLogin(*smmsUser, *smmsPassword)
			if err != nil {
				panic(err)
			}
			*token = usertoken
		}
	}

	if *token == "" {
		panic("需要 sm ms token才能运行")
	}

	smmsClint.Headers = map[string]string{
		"Authorization": *token,
	}

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		event := query.Get("e")

		resultJsonFunc := func(data any, err error) {
			if err != nil {
				http.Error(w, "解析body错误："+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		}

		smmsClintRespHandleFunc := func(resp *http.Response, err error) bool {

			fmt.Printf("%v\n", resp.Request.URL.String())

			if err != nil {
				http.Error(w, "请求错误："+err.Error(), http.StatusInternalServerError)
				return true
			}

			if resp.StatusCode != http.StatusOK {
				http.Error(w, "请求失败："+resp.Status, http.StatusInternalServerError)
				return true
			}

			return false
		}

		switch event {
		case "images":
			resp, err := smmsClint.Get("/upload_history?page=" + query.Get("page"))

			if smmsClintRespHandleFunc(resp, err) {
				return
			}

			var data map[string]any
			err = ajahttp.JsonResponse(resp, &data)
			resultJsonFunc(data, err)

		case "del":
			resp, err := smmsClint.Get("/delete/" + query.Get("hash"))

			if smmsClintRespHandleFunc(resp, err) {
				return
			}

			var data map[string]any
			err = ajahttp.JsonResponse(resp, &data)
			resultJsonFunc(data, err)

		case "upload":
			f, fHeader, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "上传错误："+err.Error(), http.StatusBadRequest)
				return
			}
			fd := ajahttp.NewFormData()
			fd.AppendFile("smfile", f, fHeader.Filename)

			resp, err := smmsClint.PostFormData("/upload", fd)
			if smmsClintRespHandleFunc(resp, err) {
				return
			}

			var data map[string]any
			err = ajahttp.JsonResponse(resp, &data)
			resultJsonFunc(data, err)

		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "请求错误: 未知的event --> %s", event)
		}
	})

	staticHandle := http.FileServer(http.Dir("./static"))
	http.Handle("/", staticHandle)

	fmt.Printf("服务器启动: http://localhost:%s\n", *port)
	http.ListenAndServe(":"+*port, nil)
}

func smmsAuthLogin(user, password string) (string, error) {
	fd := ajahttp.NewFormData()
	fd.Append("username", user)
	fd.Append("password", password)

	resp, err := smmsClint.PostFormData("/token", fd)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	var data map[string]any
	err = ajahttp.JsonResponse(resp, &data)

	if err != nil {
		return "", err
	}

	if !data["success"].(bool) {
		return "", errors.New(data["message"].(string))
	}

	return data["data"].(map[string]any)["token"].(string), nil
}
