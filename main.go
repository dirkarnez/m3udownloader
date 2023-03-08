package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

var (
	url string
)

func main() {
	flag.StringVar(&url, "url", "", "url")

	flag.Parse()
	if len(url) < 1 {
		log.Fatal("No --url is given")
	}

	_, m3uURL, err := getTitle(url)
	errExit(err)
	fmt.Print(m3uURL)
}

func errExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Transcoder struct {
}

func (w Transcoder) Write(b []byte) (n int, err error) {
	n = len(b)
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		fmt.Printf("%c", r)
		b = b[size:]
	}
	return n, err
}

func getTitle(urlstr string) (string, string, error) {
	login := os.Getenv("PROXY_LOGIN")
	password := os.Getenv("PROXY_PASSWORD")
	proxy_addr := os.Getenv("PROXY_ADDR") //127.0.0.1:1080
	schema := os.Getenv("PROXY_TYPE")     //socks5:// or http://

	allocOpts := chromedp.DefaultExecAllocatorOptions[:]

	var proxyURL string

	if len(login) > 0 {
		proxyURL = fmt.Sprintf("%s%s:%s@%s", schema, login, password, proxy_addr)
	} else {
		proxyURL = fmt.Sprintf("%s%s", schema, proxy_addr)
	}

	fmt.Println(proxyURL)

	allocOpts = append(allocOpts, chromedp.ProxyServer(proxyURL))
	actx, cancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(actx)
	defer cancel()

	var title string
	var url string

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			resp := ev.Response
			if strings.Contains(resp.URL, ".m3u8") && resp.MimeType == "application/vnd.apple.mpegurl" {
				url = resp.URL
			}
		}
	})

	req := `
(async () => new Promise((resolve, reject) => {
	var handle = NaN;

	(function animate() {
		if (!isNaN(handle)) {
			clearTimeout(handle);
		}

		if (document.title.length > 0 && !document.title.startsWith("http")) {
			resolve(document.title);
		} else {
			handle = setTimeout(animate, 1000);
		}
	}());
}));
`
	err := chromedp.Run(ctx,
		chromedp.Navigate(urlstr),
		//chromedp.Evaluate(`window.location.href`, &res),
		chromedp.Evaluate(req, nil, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
		chromedp.Title(&title),
	)

	return title, url, err
}

// //ffmpeg -i "%url%" "%filename%"
// func FFmpegDownload(url  string) error {
// 	command := exec.Command("ffmpeg", "-i", url, filename)
// 	command.Stdout = os.Stdout

// 	command.Stderr = os.Stderr

// 	err := command.Start()
// 	if err != nil {
// 		return err
// 	}
// 	go command.Wait()
// 	return nil
// }
