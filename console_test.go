package console

import (
	"regexp"
	"testing"
)

func TestProcessLine(t *testing.T) {
	rx := regexp.MustCompile(LogRegexp)
	testCases := []struct {
		log  string
		want string
	}{
		{
			log:  "126.64.189.226 \"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)\" - [24/Feb/2016:21:41:00 +0100] \"GET /tags/open-source/list.html HTTP/1.1\" 200 343",
			want: "/tags/open-source",
		},
		{
			log:  "158.213.43.241 \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:23.0) Gecko/20100101 Firefox/23.0\" - [24/Feb/2016:21:41:00 +0100] \"GET /tags/python/header.html HTTP/1.1\" 302 1991",
			want: "/tags/python",
		},
		{
			log:  "70.102.185.212 \"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)\" - [24/Feb/2016:21:41:00 +0100] \"GET /articles/datapower-static-routes/item.html HTTP/1.1\" 200 1755",
			want: "/articles/datapower-static-routes",
		},
		{
			log:  "143.169.220.182 \"Mozilla/5.0 (iPhone; CPU iPhone OS 8_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B410 Safari/600.1.4\" - [24/Feb/2016:21:41:00 +0100] \"GET /articles/chess-board-in-objective-c/item.html HTTP/1.1\" 200 1847",
			want: "/articles/chess-board-in-objective-c",
		},
		{
			log:  "169.54.175.126 \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:23.0) Gecko/20100101 Firefox/23.0\" - [24/Feb/2016:21:41:00 +0100] \"GET /tags/datapower/header.html HTTP/1.1\" 200 1838",
			want: "/tags/datapower",
		},
		{
			log:  "75.12.129.248 \"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Win64; x64; Trident/6.0)\" - [24/Feb/2016:21:41:01 +0100] \"GET /articles/chess-board-in-objective-c/list.html HTTP/1.1\" 200 1646",
			want: "/articles/chess-board-in-objective-c",
		},
	}
	for _, c := range testCases {
		got := processLine(rx, c.log)
		if got != c.want {
			t.Errorf("Error processing %s\ngot: %s, want: %s", c.log, got, c.want)
		}
	}
}
