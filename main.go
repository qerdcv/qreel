package main

import (
	"flag"
	"log"
	"os"

	"github.com/qerdcv/qreel/pkg/reelser"
)

func main() {
	var (
		output string
		url    string
	)

	flag.StringVar(&url, "url", "", "Reels url")
	flag.StringVar(&output, "o", "output.mp4", "Output file")

	flag.Parse()

	if url == "" {
		log.Fatalln("ERROR: url not provided")
	}

	if output == "" {
		log.Fatalln("ERROR: output not provided")
	}

	f, err := os.Create(output)
	if err != nil {
		log.Fatalln("ERROR: create file", err.Error())
	}

	r := reelser.New()

	log.Println("Getting video source url...")
	u, err := r.GetVideoURL(url)
	if err != nil {
		log.Fatalln("ERROR: get video url;", err)
	}

	log.Println("Downloading reel...")
	if err = r.DownloadReel(u, f); err != nil {
		log.Fatalln("ERROR: download reel;", err)
	}
}
