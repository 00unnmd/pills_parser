package handlers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/00unnmd/pills_parser/internals/utils"
	"github.com/00unnmd/pills_parser/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func getEAChangeRegionUrl(regionKey string) string {
	if regionKey == "msk" {
		return ""
	} else {
		return regionKey + "/"
	}
}

func getEASearchReqUrl(pillValue string, regionKey string) string {
	if regionKey == "msk" {
		return os.Getenv("EA_REQ_SEARCH") + pillValue
	} else {
		return os.Getenv("EA_REQ_BASE") + regionKey + "/search/?q=" + pillValue
	}
}

func getEADiscount(price int, priceOld int) int {
	if priceOld == 0 {
		return 0
	} else {
		return priceOld - price
	}
}

func getEAPrices(sel *goquery.Selection) (int, int, int, error) {
	isInStock := sel.AttrOr("data-oldma-item-serp-is-in-stock", "0")
	if isInStock != "1" {
		return 0, 0, 0, nil
	}

	price := sel.AttrOr("data-oldma-item-serp-price", "0")
	priceOldSel := sel.Find("span.listing-card__price-old")
	priceOld := priceOldSel.AttrOr("data-old-price", "0")

	priceInt, errP := strconv.Atoi(strings.ReplaceAll(price, " ", ""))
	priceOldInt, errPO := strconv.Atoi(strings.ReplaceAll(priceOld, " ", ""))
	if errP != nil || errPO != nil {
		return 0, 0, 0, fmt.Errorf("getEAPrices err converting string to int: %w, %w", errP, errPO)
	}

	discount := getEADiscount(priceInt, priceOldInt)
	return priceInt, priceOldInt, discount, nil
}

func getEAProducer(sel *goquery.Selection) (string, error) {
	producerSel := sel.
		Find("span.listing-card__manufacturer").
		ParentFiltered("p").
		Find("a")

	if producerSel.Length() == 0 {
		return "", fmt.Errorf("getEAProducer err: producer not found")
	}

	producerName := producerSel.Text()
	return producerName, nil
}

func getEARatingReviews(sel *goquery.Selection) (int, int, error) {
	containerSel := sel.Find("div.listing-card__rate")
	rating := containerSel.Find("meta[itemprop='ratingValue']").AttrOr("content", "0")
	reviews := containerSel.Find("meta[itemprop='reviewCount']").AttrOr("content", "0")

	ratingInt, errR := strconv.Atoi(rating)
	reviewsInt, errRV := strconv.Atoi(reviews)
	if errR != nil || errRV != nil {
		return 0, 0, fmt.Errorf("getEARatingReviews err converting string to int: %w, %w", errR, errRV)
	}

	return ratingInt, reviewsInt, nil
}

func parseEAHTMLData(html string, region string, pillValue string) ([]models.ParsedItem, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parseEAHTMLData new doc creation err: %w", err)
	}

	var result []models.ParsedItem

	doc.Find("article.listing-card.js-neon-item").Each(func(i int, s *goquery.Selection) {
		id := s.AttrOr("data-id", strconv.Itoa(i))
		name := s.AttrOr("data-oldma-item-serp-name", "N/A")

		price, priceOld, discount, err := getEAPrices(s)
		if err != nil {
			fmt.Println(err)
		}

		producer, err := getEAProducer(s)
		if err != nil {
			fmt.Println(err)
		}

		rating, reviewsCount, err := getEARatingReviews(s)
		if err != nil {
			fmt.Println(err)
		}

		result = append(result, models.ParsedItem{
			Id:           id,
			Region:       region,
			Name:         name,
			Price:        float64(price),
			Discount:     float64(discount),
			PriceOld:     float64(priceOld),
			MaxQuantity:  0,
			Producer:     producer,
			Rating:       float64(rating),
			ReviewsCount: reviewsCount,
		})
	})

	filteredResult := utils.FilterByProducer(result, pillValue)
	return filteredResult, nil
}

func CreateEAContext() (context.Context, context.CancelFunc, error) {
	ctx := context.Background()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // TODO CHANGE TO TRUE
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("incognito", true), // TODO SET TO FALSE MAYBE CAN RESOLVE PROBLEM WITH CAPTCHA
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)

	err := chromedp.Run(
		browserCtx,

		network.Enable(),
		chromedp.Navigate(os.Getenv("EA_REQ_SEARCH")+"some"),
		chromedp.Sleep(utils.RequestDelay),
	)
	if err != nil {
		allocCancel()
		return nil, nil, fmt.Errorf("CreateEAContext err: %w", err)
	}

	return browserCtx, func() {
		browserCancel()
		allocCancel()
	}, nil
}

func ChangeEARegion(ctx context.Context, regionKey string) (bool, error) {
	url := getEAChangeRegionUrl(regionKey)

	headers := network.Headers{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Encoding": "gzip, deflate, br, zstd",
		"Connection":      "keep-alive",
	}

	err := chromedp.Run(
		ctx,

		network.SetExtraHTTPHeaders(headers),
		chromedp.WaitVisible(`#single-spa-application\:\@front-office\/app-select-region`, chromedp.ByQuery),
		chromedp.Click(`#single-spa-application\:\@front-office\/app-select-region`, chromedp.ByQuery),

		chromedp.WaitVisible(`ul.Fc6yh6V li.yVFVbQm a[href="/`+url+`"]`, chromedp.ByQuery),
		chromedp.Click(`ul.Fc6yh6V li.yVFVbQm a[href="/`+url+`"]`, chromedp.ByQuery),

		chromedp.Click("div.JETIxY4 button.jwUdoOe", chromedp.ByQuery),
	)
	if err != nil {
		return false, fmt.Errorf("GetEAPills err changing region: %w", err)
	}

	return true, nil
}

func GetEAPills(ctx context.Context, pillValue string, regionKey string, regionValue string) ([]models.ParsedItem, error) {
	var html string
	reqUrl := getEASearchReqUrl(pillValue, regionKey)

	err := chromedp.Run(
		ctx,
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Navigate(reqUrl),

		chromedp.WaitVisible("div.listing.cc-group", chromedp.ByQuery),
		chromedp.OuterHTML("div.sec-categories__list.sec-search__list", &html),
	)
	if err != nil {
		return nil, fmt.Errorf("GetEAPills err get pills: %w", err)
	}

	return parseEAHTMLData(html, regionValue, pillValue)
}
