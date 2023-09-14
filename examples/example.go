package main

import (
	"log"
	"os"

	"github.com/qerdcv/qreel/pkg/reelser"
)

func main() {
	r := reelser.New()

	u, err := r.GetVideoURL("https://www.instagram.com/reel/CvpqeLsJB4J/?igshid=MzRlODBiNWFlZA==")
	if err != nil {
		log.Fatalln("ERROR: get video url;", err.Error())
	}

	f, err := os.Create("output.mp4")
	if err != nil {
		log.Fatalln("ERROR: file create;", err.Error())
	}

	defer f.Close()

	if err = r.DownloadReel(u, f); err != nil {
		log.Fatalln("ERROR: download reel;", err.Error())
	}
}
