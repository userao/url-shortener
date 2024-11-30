package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func validateUrl(url string) error {
	validUrl := regexp.MustCompile(`^(http(s)?:\/\/.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`)
	isValid := validUrl.MatchString(url)

	if !isValid {
		return errors.New("provided url not valid")
	}

	return nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func createUrl(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		RawUrl string
	}

	enableCors(&w)
	w.Header().Set("Content-type", "application/json")

	var data reqBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	handleErr(err)

	err = validateUrl(data.RawUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"%s\"}", err)))
		return
	}

	parsedUrl, err := url.Parse(data.RawUrl)
	handleErr(err)

	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "https"
		parsedUrl.Host = parsedUrl.Path
		parsedUrl.Path = ""
	}

	if !strings.Contains(parsedUrl.Host, "www.") {
		parsedUrl.Host = "www." + parsedUrl.Host
	}

	hash, err := dbConnection.CreateUrl(parsedUrl.String())
	handleErr(err)

	_, err = w.Write([]byte(fmt.Sprintf("{\"urlHash\": \"%s\"}", hash)))
	handleErr(err)
}

func getUrl(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	path := r.URL.Path

	hash := strings.Split(path, "/")[2]

	fullUrl, err := dbConnection.GetFullUrl(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"%s\"}", err)))
		return
	}

	w.Write([]byte(fmt.Sprintf("{\"fullUrl\": \"%s\"}", fullUrl)))
}

func getAllUrls(w http.ResponseWriter, r *http.Request) {

}
