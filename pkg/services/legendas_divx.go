package services

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dronept/go-stremio-legendasdivx/pkg/models"
	"github.com/dronept/go-stremio-legendasdivx/pkg/utils"
	"github.com/gocolly/colly"
)

var cookiesMap map[string]string = map[string]string{}

func login(u string, p string, sid string) string {
	fmt.Println("Loggin to LegendasDivx with user:", u)

	if cookie_from_map := cookiesMap[u]; cookie_from_map != "" {
		return cookie_from_map
	}

	urlPath := "https://www.legendasdivx.pt/forum/ucp.php?mode=login"

	client := &http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	data := url.Values{}

	data.Set("username", u)
	data.Set("password", p)
	data.Set("autologin", "on")
	data.Set("redirect", "./ucp.php?mode=login")
	data.Set("sid", sid)
	data.Set("login", "Ligue-se")

	encodedData := data.Encode()

	req, err := http.NewRequest("POST", urlPath, strings.NewReader(encodedData))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encodedData)))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	cookies := resp.Cookies()

	var cookie []string

	for _, c := range cookies {
		cookie = append(cookie, c.Name+"="+c.Value)
	}

	cookieStr := strings.Join(cookie, "; ")

	cookiesMap[u] = cookieStr

	return cookieStr
}

func Login(u string, p string) string {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	url := "https://www.legendasdivx.pt/forum/ucp.php?mode=login"

	c := colly.NewCollector()

	var cookie string

	c.OnHTML("input[name=sid]", func(e *colly.HTMLElement) {
		if e.Attr("value") != "" {
			sid := e.Attr("value")

			cookie = login(u, p, sid)

			wg.Done()
		}
	})

	c.Visit(url)

	wg.Wait()

	return cookie
}

func FetchSubtitles(imdbID string, cookie string) []models.Subtitle {
	var count int

	wg := &sync.WaitGroup{}
	wg.Add(1)

	count++
	fmt.Printf("[%d]\n", count)

	url := fmt.Sprintf("https://www.legendasdivx.pt/modules.php?name=Downloads&file=jz&d_op=search&op=_jz00&imdbid=%s", strings.Replace(imdbID, "tt", "", 1))

	subtitles := []models.Subtitle{}
	pageVisited := map[string]bool{}
	pageRe := regexp.MustCompile(`page=(\d+)`)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// Set "cookie" header
		fmt.Printf("Requesting: %s\nCookie: %s\n", r.URL, cookie)
		r.Headers.Set("Cookie", cookie)
	})

	c.OnHTML(".sub_box", func(e *colly.HTMLElement) {
		langImageSrc := e.DOM.Find("tr").Eq(0).Find("img").Eq(0).AttrOr("src", "")

		var language string

		if strings.Contains(langImageSrc, "portugal") {
			language = "por"
		}

		if strings.Contains(langImageSrc, "brazil") {
			language = "pob"
		}

		if strings.Contains(langImageSrc, "fInglaterra") {
			language = "eng"
		}

		if strings.Contains(langImageSrc, "fEspanha") {
			language = "spa"
		}

		var subsRe = regexp.MustCompile(`(?m)(^(\w{2,}[\.\s]){2,}.*-[\[|\w.\]]+)`)

		desc := e.DOM.Find(".td_desc").Text()

		fmt.Printf("----------------------------------------\nDescription:\n\n%s\n\n----------------------------------------\n", desc)

		// Find all matched strings
		matches := subsRe.FindAllString(desc, -1)

		lidRe := regexp.MustCompile(`lid=(\d+)`)

		if len(matches) > 0 {
			// Print all matches
			for _, name := range matches {
				href := e.DOM.Find("a.sub_download").Eq(0).AttrOr("href", "")
				lid := lidRe.FindStringSubmatch(href)[1]

				subtitles = append(subtitles, models.Subtitle{
					DownloadUrl: lid,
					Language:    language,
					Name:        name,
				})
			}

		} else {
			href := e.DOM.Find("a.sub_download").Eq(0).AttrOr("href", "")
			lid := lidRe.FindStringSubmatch(href)[1]

			subtitles = append(subtitles, models.Subtitle{
				DownloadUrl: lid,
				Language:    language,
			})
		}
	})

	c.OnHTML(".pager_bar a", func(e *colly.HTMLElement) {
		if e.Text == "Seguinte" {
			page := pageRe.FindStringSubmatch(e.Attr("href"))[1]

			if _, ok := pageVisited[page]; !ok {
				wg.Add(1)
				count++

				pageVisited[page] = true

				e.Request.Visit(e.Attr("href"))
			}
		}
	})

	c.OnScraped(func(r *colly.Response) {
		count--
		fmt.Printf("[%d] Scrapped! %s\n", count, r.Request.URL.String())
		wg.Done()
	})

	c.Visit(url)

	wg.Wait()

	return subtitles
}

func Download(lid, cookie string) []string {
	fmt.Println("Downloading subtitle with lid:", lid)

	lidPath := "tmp/downloads/" + lid

	// If lid directory exists, return
	if _, err := os.Stat(lidPath); !os.IsNotExist(err) {
		fmt.Println("Subtitle already downloaded")

		files, _ := utils.ListFiles(lidPath)

		return files
	}

	// https://www.legendasdivx.pt/modules.php?name=Downloads&d_op=getit&lid=451481

	url := fmt.Sprintf("https://www.legendasdivx.pt/modules.php?name=Downloads&d_op=getit&lid=%s", lid)

	fmt.Println("Downloading subtitle from:", url)

	c := colly.NewCollector()
	wg := &sync.WaitGroup{}
	wg.Add(1)

	c.OnRequest(func(r *colly.Request) {
		// Set "cookie" header
		r.Headers.Set("Cookie", cookie)
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		// Get filename from headers Content-Disposition
		filename := r.Headers.Get("Content-Disposition")

		// Get filename from regex filename="(.+)"
		re := regexp.MustCompile(`filename="(.+)"`)
		filename = re.FindStringSubmatch(filename)[1]

		extensionSplit := strings.Split(filename, ".")
		extension := extensionSplit[len(extensionSplit)-1]

		fmt.Printf("Filename: %s\nExtension:%s\n", filename, extension)

		// Create directory under tmp/downloads/{lid}
		os.MkdirAll(lidPath, os.ModePerm)

		// Save file
		err := r.Save(lidPath + "/subtitle." + extension)

		if err != nil {
			log.Fatal(err)
		}

		utils.Extract(lid)

		wg.Done()
	})

	c.Visit(url)

	wg.Wait()

	files, _ := utils.ListFiles(lidPath)

	return files
}
