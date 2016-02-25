package console

import (
	"container/ring"
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
			want: "/tags",
		},
		{
			log:  "158.213.43.241 \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:23.0) Gecko/20100101 Firefox/23.0\" - [24/Feb/2016:21:41:00 +0100] \"GET /tags/python/header.html HTTP/1.1\" 302 1991",
			want: "/tags",
		},
		{
			log:  "70.102.185.212 \"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)\" - [24/Feb/2016:21:41:00 +0100] \"GET /articles/datapower-static-routes/item.html HTTP/1.1\" 200 1755",
			want: "/articles",
		},
		{
			log:  "143.169.220.182 \"Mozilla/5.0 (iPhone; CPU iPhone OS 8_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B410 Safari/600.1.4\" - [24/Feb/2016:21:41:00 +0100] \"GET /articles/chess-board-in-objective-c/item.html HTTP/1.1\" 200 1847",
			want: "/articles",
		},
		{
			log:  "169.54.175.126 \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:23.0) Gecko/20100101 Firefox/23.0\" - [24/Feb/2016:21:41:00 +0100] \"GET /tags/datapower/header.html HTTP/1.1\" 200 1838",
			want: "/tags",
		},
		{
			log:  "75.12.129.248 \"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Win64; x64; Trident/6.0)\" - [24/Feb/2016:21:41:01 +0100] \"GET /articles/chess-board-in-objective-c/list.html HTTP/1.1\" 200 1646",
			want: "/articles",
		},
	}
	for _, c := range testCases {
		got := processLine(rx, c.log)
		if got != c.want {
			t.Errorf("Error processing %s\ngot: %s, want: %s", c.log, got, c.want)
		}
	}
}

func TestAverageSample(t *testing.T) {
	testCases := []struct {
		sample [5]int
		want   float64
	}{
		{
			sample: [5]int{2, 2, 2, 2, 2},
			want:   2,
		},
		{
			sample: [5]int{20, 20, 20, 20, 20},
			want:   20,
		},
		{
			sample: [5]int{1, 2, 3, 4, 5},
			want:   3,
		},
	}
	for _, c := range testCases {
		rs := ring.New(len(c.sample))
		for _, i := range c.sample {
			rs.Value = i
			rs = rs.Next()
		}
		got := averageSample(rs)
		if got != c.want {
			t.Errorf("Average error: got: %f, want: %f", got, c.want)
		}
	}
}

func TestCheckHighTraffic(t *testing.T) {
	testCases := []struct {
		sample    [5]int
		threshold int
		status    bool
		want      bool
	}{
		// under threshold and staying like that.
		{
			sample:    [5]int{2, 2, 2, 2, 2},
			threshold: 10,
			status:    false,
			want:      false,
		},
		// over threshold and getting under it.
		{
			sample:    [5]int{2, 2, 2, 2, 2},
			threshold: 10,
			status:    true,
			want:      false,
		},
		// under threshold and going over it.
		{
			sample:    [5]int{20, 20, 20, 20, 20},
			threshold: 10,
			status:    false,
			want:      true,
		},
		// over threshold and staying like that.
		{
			sample:    [5]int{20, 20, 20, 20, 20},
			threshold: 10,
			status:    true,
			want:      true,
		},
	}
	for _, c := range testCases {
		rs := ring.New(len(c.sample))
		for _, i := range c.sample {
			rs.Value = i
			rs = rs.Next()
		}
		got := checkHighTraffic(rs, c.threshold, c.status)
		if got != c.want {
			t.Errorf("High traffic check failed for %q; got %t, want %t (status: %t)", c.sample, got, c.want, c.status)
		}
	}
}
