package services

import (
	"bufio"
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
	"github.com/gocolly/colly"
)

func loadCookieFromFile() string {
	file, err := os.Open("cookie.txt")

	if err != nil {
		fmt.Println("Error opening file:", err)

		return ""
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	return scanner.Text()
}

func saveCookieToFile(cookie string) {
	file, err := os.Create("cookie.txt")

	if err != nil {
		fmt.Println("Error creating file:", err)

		return
	}

	defer file.Close()

	file.WriteString(cookie)
}

func login(u string, p string, sid string) string {
	if cookie_from_file := loadCookieFromFile(); cookie_from_file != "" {
		return cookie_from_file
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
	// data.Set("autologin", "on")
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

	fmt.Println("Login response status:", resp.Status)

	defer resp.Body.Close()

	cookies := resp.Cookies()

	var cookie []string

	for _, c := range cookies {
		cookie = append(cookie, c.Name+"="+c.Value)
	}

	cookieStr := strings.Join(cookie, "; ")

	saveCookieToFile(cookieStr)

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
	wg := &sync.WaitGroup{}
	wg.Add(1)

	url := fmt.Sprintf("https://www.legendasdivx.pt/modules.php?name=Downloads&file=jz&d_op=search&op=_jz00&query=%s", imdbID)
	// cookie := "phpbb3_2z8zs_k=; PHPSESSID=tnoo8alhltjoektv2pu6plpgl2; phpbb3_2z8zs_u=66611; phpbb3_2z8zs_sid=1436ee4b7b342a8d954c2a6f269d6732; lang=portuguese; banner1=0"

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// Set "cookie" header
		r.Headers.Set("Cookie", cookie)
	})

	var subtitles []models.Subtitle = []models.Subtitle{}

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

		var re = regexp.MustCompile(`(?m)(^(\w{2,}[\.\s]){2,}.*-\w+$)`)

		desc := e.DOM.Find(".td_desc").Text()

		// Find all matched strings
		matches := re.FindAllString(desc, -1)

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

				// fmt.Printf("--- %s ---\n%v\n---------------", name, desc)
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

	c.OnScraped(func(*colly.Response) {
		wg.Done()
	})

	c.Visit(url)

	wg.Wait()

	return subtitles
}
