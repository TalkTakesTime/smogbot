package smogbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	postContentFilter = regexp.MustCompile("<blockquote[^>]+?>(.+?)</blockquote>")
	pageNumCounter    = regexp.MustCompile("<span class=\"pageNavHeader\">.+?([0-9]+)</span>")
	replayRegex       = regexp.MustCompile("<a href=\".*?(replay.pokemonshowdown.com/[a-z0-9\\-]+?)\".*?>")

	playerRegex = regexp.MustCompile(`\|player\|(p1|p2)\|([^|]+)`)
	teamRegex   = regexp.MustCompile(`\|poke\|(p1|p2)\|([^,\n]+)`)
	winnerRegex = regexp.MustCompile(`\|win\|([^\n<]+)`)
)

type ReplayData struct {
	url     string
	players map[string]string
	teams   map[string][]string
	winner  string
}

func Start(url string) {
	url = strings.TrimRight(url, "/")
	html := getPage(url)
	n := getNumberOfPages(html)
	pages := getPageURLs(url, n)

	posts := filterPosts(html)
	for _, v := range pages {
		html = getPage(v)
		posts = append(posts, filterPosts(html)...)
	}

	var replays []string
	for _, v := range posts {
		replays = append(replays, extractReplays(v)...)
	}

	replays = getUnique(replays)

	for _, v := range replays {
		fmt.Println("_____________________________________________________________________")
		fmt.Println(v)
		fmt.Println()
		d := getReplayData(v)
		for k, v := range d.teams {
			fmt.Println(k + ":")
			fmt.Println(strings.Join(v, " / "))
		}
		fmt.Println("\nWon by: " + d.winner)
	}
}

func getUnique(slice []string) []string {
	unique := map[string]bool{}

	for _, v := range slice {
		unique[v] = true
	}

	u := []string{}
	for k := range unique {
		u = append(u, k)
	}

	return u
}

func getPage(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func getNumberOfPages(html string) int {
	numCounter := pageNumCounter.FindStringSubmatch(html)

	if numCounter == nil {
		return 1
	}

	n, err := strconv.Atoi(numCounter[1])
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func getPageURLs(baseURL string, nPages int) []string {
	urls := []string{}

	for i := 2; i <= nPages; i++ {
		urls = append(urls, baseURL+"/page-"+strconv.Itoa(i))
	}

	return urls
}

func filterPosts(data string) []string {
	data = strings.Replace(data, "\n", "", -1)
	posts := postContentFilter.FindAllStringSubmatch(data, -1)

	var postContents []string
	for _, v := range posts {
		postContents = append(postContents, strings.TrimSpace(v[1]))
	}

	return postContents
}

func extractReplays(data string) []string {
	replays := replayRegex.FindAllStringSubmatch(data, -1)

	var links []string
	for _, v := range replays {
		links = append(links, "http://"+v[1])
	}

	return links
}

func getReplayData(replay string) ReplayData {
	players := map[string]string{}
	teams := map[string][]string{}

	html := getPage(replay)

	p := playerRegex.FindAllStringSubmatch(html, -1)
	for _, v := range p {
		if _, ok := players[v[1]]; !ok {
			players[v[1]] = v[2]
			teams[v[2]] = []string{}
		}
	}

	t := teamRegex.FindAllStringSubmatch(html, -1)
	for _, v := range t {
		teams[players[v[1]]] = append(teams[players[v[1]]], v[2])
	}

	w := winnerRegex.FindStringSubmatch(html)
	winner := ""
	if w != nil {
		winner = w[1]
	}

	return ReplayData{
		replay,
		players,
		teams,
		winner,
	}
}
