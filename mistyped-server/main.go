package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var keyMap = map[string]string{
	"q": "aw",
	"w": "se",
	"e": "dr",
	"r": "ft",
	"t": "gz",
	"z": "hu",
	"u": "ji",
	"i": "ko",
	"o": "lp",
	"p": "öü",
	"ü": "ä",
	"a": "ys",
	"s": "xd",
	"d": "cf",
	"f": "vg",
	"g": "bh",
	"h": "nj",
	"j": "mk",
	"k": "l",
	"l": "ö",
	"ö": "ä",
	"ä": "",
	"y": "sx",
	"x": "dc",
	"c": "fv",
	"v": "gb",
	"b": "hn",
	"n": "jm",
	"m": "k",
}

var port = 8080

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ENDPOINT HOMEPAGE")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/checkurl/", handleCheckUrl)
	portString := ":" + strconv.Itoa(port)
	fmt.Println("mistyped-server up and running on port", portString)
	log.Fatal(http.ListenAndServe(portString, router))
}

func handleCheckUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	check(err)
	url := r.Form.Get("url")

	candidates := getCandidates(url)

	c1 := make(chan []string)
	c2 := make(chan []string)
	go checkUrlAvailability(c1, c2, candidates)
	availableCandidates := <-c1
	unavailableCandidates := <-c2

	result := make(map[string][]string)
	result["available"] = availableCandidates
	result["unavailable"] = unavailableCandidates
	fmt.Println(availableCandidates)
	fmt.Println(unavailableCandidates)

	json.NewEncoder(w).Encode(result)
}

func checkUrlAvailability(c1 chan []string, c2 chan []string, candidates []string) {
	fmt.Println("START GOROUTINE CheckUrlAvailability")
	availableCandidates := make([]string, 0)
	unavailableCandidates := make([]string, 0)

	for _, candidate := range candidates {
		client := http.Client{
			Timeout: 1 * time.Second,
		}
		res, err := client.Get(candidate)
		if err != nil && res != nil && (res.Status == strconv.Itoa(200) || res.Status == strconv.Itoa(204)) {
			availableCandidates = append(availableCandidates, candidate)
		} else {
			unavailableCandidates = append(unavailableCandidates, candidate)
		}
	}
	c1 <- availableCandidates
	c2 <- unavailableCandidates

	close(c1)
	close(c2)
}

func getCandidates(url string) []string {
	candidates := make([]string, 0)
	urlSplit := strings.Split(url, ".")
	hostName := urlSplit[1]
	for _, urlCharacter := range urlSplit[1] {
		for _, possibleCharacters := range keyMap[string(urlCharacter)] {
			for _, possibleCharacter := range string(possibleCharacters) {
				possibleHostName := replace(hostName, string(urlCharacter), string(possibleCharacter))
				urlSplit[1] = possibleHostName
				possibleUrl := getString(urlSplit)
				candidates = append(candidates, possibleUrl)
			}
		}
	}

	return candidates
}

func replace(replaceIn string, replaceThis string, replaceWith string) string {
	result := ""
	for _, c := range replaceIn {
		if string(c) == replaceThis {
			result += replaceWith
		} else {
			result += string(c)
		}
	}
	return result
}

func getString(stringArray []string) string {
	result := ""
	for i, s := range stringArray {
		result += s

		if i != len(stringArray)-1 {
			result += "."
		}
	}

	return result
}
