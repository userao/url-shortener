package server

import (
	"log"
	"net/http"
	"regexp"
)

func validateUrl(url string) bool {
	// validUrl := regexp.MustCompile(`(http(s)?:\/\/)?(www\.)?([a-zA-Z]+\.)+[a-zA-Z]{2,6}(\/[\w?=#-$.]+)+`)
	validUrl := regexp.MustCompile(`^(http(s):\/\/.)[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`)
	isValid := validUrl.MatchString(url)
	return isValid
}

func createUrl(w http.ResponseWriter, r *http.Request) {
	fullUrl := r.FormValue("fullUrl")

	if !validateUrl(fullUrl) {
		w.Write([]byte("url not valid"))
		return
	}

	shortenedUrl, err := dbConnection.CreateUrl(fullUrl, GetCurrentServer().host)

	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write([]byte(shortenedUrl))

	if err != nil {
		log.Fatal(err)
	}
}

func getUrl(w http.ResponseWriter, r *http.Request) {

}

func getAllUrls(w http.ResponseWriter, r *http.Request) {

}

func redirectToOriginalUrl(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fullUrl, err := dbConnection.GetFullUrl(path[1:])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, fullUrl, http.StatusPermanentRedirect)

}
