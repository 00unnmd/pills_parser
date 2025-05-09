package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FiltersSimple struct {
	Id  string `json:"id"`
	Val string `json:"val"`
}

type SearchFilters struct {
	Simple [1]FiltersSimple `json:"simple"`
}

type SearchPaginator struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type VariablesSearch struct {
	AdvertisementKey string          `json:"advertisementKey"`
	Filters          SearchFilters   `json:"filters"`
	Paginator        SearchPaginator `json:"paginator"`
	RegionID         string          `json:"regionID"`
	Query            string          `json:"query"`
}

type RequestBodyVariables struct {
	Search            VariablesSearch `json:"search"`
	RegionID          string          `json:"regionID"`
	AdvertisementType string          `json:"advertisementType"`
	Query             string          `json:"query"`
	SkipFeaturing     bool            `json:"skipFeaturing"`
}

type RequestBody struct {
	OperationName string               `json:"operationName"`
	Query         string               `json:"query"`
	Variables     RequestBodyVariables `json:"variables"`
}

type ResponseBody struct {
	Data struct {
		Products struct {
			Items []struct {
				Name         string `json:"name"`
				Price        int    `json:"price"`
				Discount     int    `json:"discount"`
				PriceOld     int    `json:"priceOld"`
				MaxQuantity  int    `json:"maxQuantity"`
				Producer     string `json:"producer"`
				IsBundle     bool   `json:"isBundle"`
				Rating       int    `json:"rating"`
				ReviewsCount int    `json:"reviewsCount"`
			} `json:"items"`
			Total int `json:"total"`
		} `json:"products"`
		Featuring []struct {
			Title      string `json:"title"`
			Banner     string `json:"banner"`
			Advertiser string `json:"advertiser"`
			Counters   struct {
				Imps     interface{} `json:"imps"`
				Track    interface{} `json:"track"`
				Creative string      `json:"creative"`
				ID       string      `json:"id"`
				Name     string      `json:"name"`
				Position string      `json:"position"`
			} `json:"counters"`
			Products []interface{} `json:"products"`
		} `json:"featuring"`
	} `json:"data"`
}

func getBodyByte(searchQuery string, regionKey string) []byte {
	body := RequestBody{
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

		Variables: RequestBodyVariables{
			Search: VariablesSearch{
				AdvertisementKey: "search",
				Filters: SearchFilters{
					Simple: [1]FiltersSimple{
						{"query", searchQuery},
					},
				},
				Paginator: SearchPaginator{
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

func requestZdravsitiData(pillValue string, regionKey string) ResponseBody {
	bodyByte := getBodyByte(pillValue, regionKey)
	bodyJSON := bytes.NewBuffer(bodyByte)

	pillsRequest, err := http.NewRequest(
		"POST",
		"https://zdravcity.ru/bff/query",
		bodyJSON,
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
	respBody := ResponseBody{}
	json.Unmarshal([]byte(respBodyJSON), &respBody)
	// fmt.Println("respBodyJSON:", string(respBodyJSON))

	return respBody
}

func main() {
	responseBody := requestZdravsitiData("антистен", "moscowregion")

	fmt.Println("responseBody:", responseBody.Data.Products.Items)
}
