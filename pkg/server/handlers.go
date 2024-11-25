package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func validateUrl(url string) bool {
	validUrl := regexp.MustCompile(`^(http(s):\/\/.)[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`)
	isValid := validUrl.MatchString(url)
	return isValid
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func createUrl(w http.ResponseWriter, r *http.Request) {
	type reqJson struct {
		FullUrl string
	}
	enableCors(&w)
	var j reqJson
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&j)

	if err != nil {
		log.Fatal(err)
	}

	if !validateUrl(j.FullUrl) {
		w.Write([]byte(fmt.Sprintf("url %s not valid", j.FullUrl)))
		return
	}

	hash, err := dbConnection.CreateUrl(j.FullUrl, GetCurrentServer().host)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-type", "application/json")
	_, err = w.Write([]byte(fmt.Sprintf("{\"shortenedUrl\": \"%s:%s/%s\"}", GetCurrentServer().host, GetCurrentServer().port, hash)))

	if err != nil {
		log.Fatal(err)
	}
}

func getUrl(w http.ResponseWriter, r *http.Request) {

}

func getAllUrls(w http.ResponseWriter, r *http.Request) {

}

func redirectToOriginalUrl(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	path := r.URL.Path
	fullUrl, err := dbConnection.GetFullUrl(path[1:])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, fullUrl, http.StatusPermanentRedirect)

}
