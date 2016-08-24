package service

import (
	"fmt"

	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gen1us2k/log"
)

type CinematicaService struct {
	BaseService
	cinema *Cinematica

	logger log.Logger
}

func (fs *CinematicaService) Name() string {
	return "cinamatica_service"
}

func (fs *CinematicaService) Init(cinema *Cinematica) error {
	fs.cinema = cinema

	fs.logger = log.NewLogger(fs.Name())
	return nil
}
func (fs *CinematicaService) Run() error {
	fs.updatePremiere()
	for range time.Tick(24 * time.Hour) {
		fs.updatePremiere()
	}
	return nil
}

func (fs *CinematicaService) updatePremiere() {
	doc, err := goquery.NewDocument("http://cinematica.kg/premiere/")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".film-block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}
