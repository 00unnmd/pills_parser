package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/00unnmd/pills_parser/internals/utils"
	"github.com/00unnmd/pills_parser/models"
)

func getBodyByte(searchQuery string, regionKey string) models.ZSRequestBody {
	body := models.ZSRequestBody{
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

		Variables: models.ZSRequestBodyVariables{
			Search: models.ZSSearch{
				AdvertisementKey: "search",
				Filters: models.ZSSearchFilters{
					Simple: [1]models.ZSSimple{
						{
							Id:  "query",
							Val: searchQuery,
						},
					},
				},
				Paginator: models.ZSSearchPaginator{
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

	return body
}

func GetZSPills(pillValue string, regionKey string, regionValue string) ([]models.ParsedItem, error) {
	body := getBodyByte(pillValue, regionKey)

	respBodyByte, err := makeAPIRequest(
		"ZS",
		"POST",
		os.Getenv("ZS_REQ_SEARCH"),
		nil,
		body,
	)
	if err != nil {
		return nil, err
	}

	var respBody models.ZSSearchBody
	if err := json.Unmarshal([]byte(respBodyByte), &respBody); err != nil {
		return nil, fmt.Errorf("GetZSPills error unmarshaling response: %w", err)
	}

	filteredProductItems := utils.FilterByProducer(respBody.Data.Products.Items, pillValue)
	parsedProductItems := utils.ParseRawData(filteredProductItems)

	for i := range parsedProductItems {
		parsedProductItems[i].Region = regionValue
	}

	return parsedProductItems, nil
}
