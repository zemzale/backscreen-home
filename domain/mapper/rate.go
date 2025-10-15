package mapper

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/zemzale/backscreen-home/domain/entity"
)

type LVBankRSSRateFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Title   string   `xml:"title"`
	Item    []Item   `xml:"item"`
}

type Item struct {
	XMLName         xml.Name `xml:"item"`
	Title           string   `xml:"title"`
	Description     string   `xml:"description"`
	PublicationDate string   `xml:"pubDate"`
}

func RatesFromXML(reader io.Reader) ([]entity.Rate, error) {
	var feed LVBankRSSRateFeed
	if err := xml.NewDecoder(reader).Decode(&feed); err != nil {
		return nil, err
	}

	if len(feed.Channel.Item) == 0 {
		return nil, errors.New("no rates found")
	}

	var rates []entity.Rate

	for _, item := range feed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PublicationDate)
		if err != nil {
			return nil, errors.New("failed to parse publication date")
		}

		rateValues := strings.Split(strings.TrimSpace(item.Description), " ")
		if len(rateValues) < 2 {
			return nil, errors.New("invalid rate format")
		}

		for i := 0; i < len(rateValues); i += 2 {
			currencyCode := rateValues[i]
			value := rateValues[i+1]

			if len(currencyCode) != 3 {
				return nil, errors.New("invalid currency code")
			}

			if len(value) == 0 {
				return nil, errors.New("invalid value")
			}

			rates = append(rates, entity.Rate{
				PublishedAt: publishedAt,
				Code:        currencyCode,
				Value:       value,
			})
		}
	}

	return rates, nil
}
