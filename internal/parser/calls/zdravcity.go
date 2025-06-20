package calls

import (
	"encoding/json"
	"fmt"
	"github.com/00unnmd/pills_parser/internal/domain"
	"github.com/00unnmd/pills_parser/internal/transport"
	"github.com/00unnmd/pills_parser/pkg/utils"
	"os"
)

func getBodyByte(searchQuery string, regionKey string) domain.ZSRequestBody {
	body := domain.ZSRequestBody{
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
                producer
                rating
				mnns {
					title
				}
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

		Variables: domain.ZSRequestBodyVariables{
			Search: domain.ZSSearch{
				AdvertisementKey: "search",
				Filters: domain.ZSSearchFilters{
					Simple: [1]domain.ZSSimple{
						{
							Id:  "query",
							Val: searchQuery,
						},
					},
				},
				Paginator: domain.ZSSearchPaginator{
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

func GetZSPills(pillValue string, regionKey string, regionValue string, withFilter bool) ([]domain.ParsedItem, error) {
	body := getBodyByte(pillValue, regionKey)

	respBodyByte, err := transport.MakeAPIRequest(
		"ZS",
		"POST",
		os.Getenv("ZS_REQ_SEARCH"),
		nil,
		body,
	)
	if err != nil {
		return nil, err
	}

	var respBody domain.ZSSearchBody
	if err := json.Unmarshal(respBodyByte, &respBody); err != nil {
		return nil, fmt.Errorf("GetZSPills error unmarshaling response: %w", err)
	}

	filteredData := respBody.Data.Products.Items
	if withFilter == true {
		filteredData = utils.FilterByProducer(respBody.Data.Products.Items, pillValue)
	}

	if len(filteredData) == 0 {
		return nil, fmt.Errorf("не найдено препаратов удовлетворяющих запросу: len(filteredData) == 0")
	}

	result := utils.ParseRawData("zdravcity", regionValue, pillValue, filteredData)
	return result, nil
}
