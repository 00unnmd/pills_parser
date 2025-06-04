package calls

import (
	"context"
	"fmt"
	"github.com/00unnmd/pills_parser/internal/domain"
	"github.com/00unnmd/pills_parser/pkg/utils"
	"math"
	"os"
	"strconv"
	"strings"

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

func getEAMNN(sel *goquery.Selection) string {
	mnnSel := sel.
		Find("span.listing-card__ingredient").
		ParentFiltered("p").
		Find("a")

	if mnnSel.Length() == 0 {
		return "mnn not found"
	}

	mnnName := mnnSel.Text()
	return mnnName
}

func getEADiscount(price float64, priceOld float64) (float64, int) {
	if priceOld == 0 {
		return 0, 0
	} else {
		discount := priceOld - price
		discountPercent := int(math.Round(discount / (priceOld / 100)))
		return discount, discountPercent
	}
}

func getEAPrices(sel *goquery.Selection) (float64, float64, int, error) {
	isInStock := sel.AttrOr("data-oldma-item-serp-is-in-stock", "0")
	if isInStock != "1" {
		return 0, 0, 0, nil
	}

	price := sel.AttrOr("data-oldma-item-serp-price", "0")
	priceOldSel := sel.Find("span.listing-card__price-old")
	priceOld := priceOldSel.AttrOr("data-old-price", "0")

	priceFl, errP := strconv.ParseFloat(strings.ReplaceAll(price, " ", ""), 64)
	priceOldFl, errPO := strconv.ParseFloat(strings.ReplaceAll(priceOld, " ", ""), 64)
	if errP != nil || errPO != nil {
		return 0, 0, 0, fmt.Errorf("getEAPrices err converting string to int: %w, %w", errP, errPO)
	}

	discount, discountPercent := getEADiscount(priceFl, priceOldFl)

	return priceFl, discount, discountPercent, nil
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

func getEARawItem(s *goquery.Selection) domain.EARawItem {
	name := s.AttrOr("data-oldma-item-serp-name", "N/A")
	errStr := ""

	mnn := getEAMNN(s)

	price, discount, discountPercent, err := getEAPrices(s)
	if err != nil {
		errStr = errStr + err.Error()
	}

	producer, err := getEAProducer(s)
	if err != nil {
		errStr = errStr + err.Error()
	}

	rating, reviewsCount, err := getEARatingReviews(s)
	if err != nil {
		errStr = errStr + err.Error()
	}

	return domain.EARawItem{
		Name:            name,
		Mnn:             mnn,
		Price:           price,
		Discount:        discount,
		DiscountPercent: discountPercent,
		Producer:        producer,
		Rating:          float64(rating),
		ReviewsCount:    reviewsCount,
		Error:           errStr,
	}
}

func parseEAHTMLData(html string, region string, pillValue string, withFilter bool) ([]domain.ParsedItem, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parseEAHTMLData new doc creation err: %w", err)
	}

	var rawData []domain.EARawItem
	doc.Find("article.listing-card.js-neon-item").Each(func(i int, s *goquery.Selection) {
		rawItem := getEARawItem(s)
		rawData = append(rawData, rawItem)
	})

	filteredData := rawData
	if withFilter == true {
		filteredData = utils.FilterByProducer(rawData, pillValue)
	}

	if len(filteredData) == 0 {
		return nil, fmt.Errorf(`не найдено препаратов удовлетворяющих запросу: len(filteredData) == 0`)
	}

	result := utils.ParseRawData("eapteka", region, pillValue, filteredData)
	return result, nil
}

func CreateEAContext() (context.Context, context.CancelFunc, error) {
	ctx := context.Background()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // TODO CHANGE TO TRUE
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("incognito", false),
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

func GetEAPills(ctx context.Context, pillValue string, regionKey string, regionValue string, withFilter bool) ([]domain.ParsedItem, error) {
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

	if strings.Contains(html, `<div class="sec-empty">`) {
		return nil, fmt.Errorf(`не найдено препаратов удовлетворяющих запросу: html.Contains("<div class="sec-empty">")`)
	}

	result, err := parseEAHTMLData(html, regionValue, pillValue, withFilter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
