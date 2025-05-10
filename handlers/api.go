package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/00unnmd/pills_parser/internals/utils"
	"github.com/00unnmd/pills_parser/models"
)

func getBodyByte(searchQuery string, regionKey string) []byte {
	body := models.RequestBody{
		OperationName: "SearchQuery",
		Query: `query SearchQuery($search: ProductSearch!, $regionID: ID!, $advertisementType: AdvertisementType!, $query: String, $skipFeaturing: Boolean = false) {
                products(search: $search) { 
                    items {
                        ...ProductSummaryFragment
                        reviewsCount
                    }
                    total
                }
                featuring: advertisement(
                    regionID: $regionID
                    type: $advertisementType
                    query: $query
                ) @skip(if: $skipFeaturing) {
                    title
                    banner
                    advertiser
                    counters {
                        ...AdvertisedCountersFragment
                    }
                    products {
                        ...ProductCardRegularFragment
                    }
                }
            }
            fragment ProductSummaryFragment on ProductSummary {
                name
                price
                discount
                priceOld
                maxQuantity
                producer
                isBundle
                rating
            }
            fragment AdvertisedCountersFragment on AdvertisementCounters {
                imps
                track
                creative
                id
                name
                position
            }
            fragment CountersFragment on Counters {
                yandex {
                    ...YandexCountersFragment
                }
                advertisement {
                    ...AdvertisedCountersFragment
                }
            }
            fragment YandexCountersFragment on YandexCounters {
                items {
                    value
                    utmConfig {
                        key
                        value
                    }
                    token
                }
            }
            fragment ProductCardRegularFragment on ProductSummary {
                id
                lastPrice
                alias
                sku
                availableForOrder
                availableForBooking
                hasOnlyBookingPrices
                canPayOnPickup
                maxQuantity
                expirationDate
                deliveryDate
                discount
                url
                image
                isFavorite
                price
                priceOld
                seoBasketText
                seoPriceText
                priceTypeID
                name
                pickupDate
                rating
                producerCountry
                warning
                isAdv
                isBundle
                brand {
                    alias
                    name
                }
                badges {
                    description
                    title
                    type
                    dateEnd
                    src
                    color
                }
                category {
                    id
                    name
                    path
                }
                counters {
                    ...CountersFragment
                }
                prices {
                    dateExpired
                    discount
                    maxQuantity
                    price
                    priceOld
                    priceTypeID
                    dateExpiredOpened
                    isExpirating
                    hasPrefix
                }
                bundleItemsSimple {
                    id
                    quantity
                    shortName
                }
                    reviewsDisabled
            }`,

		Variables: models.RequestBodyVariables{
			Search: models.VariablesSearch{
				AdvertisementKey: "search",
				Filters: models.SearchFilters{
					Simple: [1]models.FiltersSimple{
						{"query", searchQuery},
					},
				},
				Paginator: models.SearchPaginator{
					Limit:  100,
					Offset: 0,
				},
				RegionID: regionKey,
				Query:    searchQuery,
			},
			RegionID:          regionKey,
			AdvertisementType: "featuring_search",
			Query:             searchQuery,
			SkipFeaturing:     false,
		},
	}

	bodyByte, err := json.Marshal(body)

	if err != nil {
		fmt.Println("Ошибка преобразования body в JSON:", err)
	}

	return bodyByte
}

func RequestZdravsitiData(pillValue string, regionKey string, regionValue string) []models.PillsItem {
	bodyByte := getBodyByte(pillValue, regionKey)

	pillsRequest, err := http.NewRequest(
		"POST",
		"https://zdravcity.ru/bff/query",
		bytes.NewBuffer(bodyByte),
	)
	pillsRequest.Header.Set("Content-Type", "application/json")
	pillsRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(pillsRequest)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBodyJSON, _ := io.ReadAll(resp.Body)
	respBody := models.ResponseBody{}
	json.Unmarshal([]byte(respBodyJSON), &respBody)
	productItems := respBody.Data.Products.Items
	filteredProductItems := utils.FilterByProducer(productItems, utils.ProducerNames)

	for i := range filteredProductItems {
		filteredProductItems[i].Region = regionValue
	}

	return filteredProductItems
}
