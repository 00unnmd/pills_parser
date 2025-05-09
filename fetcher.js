import fetch from 'node-fetch';

const regionsList = {
    moscowregion: "Москва и область",
    spb: "Санкт-Петербург и область",
    novosibirsk: "Новосибирск",
    ekb: "Екатеринбург",
    kazan: "Казань",
    nnovgorod: "Нижний Новгород",
    ufa: "Уфа",
    rostov: "Ростов-на-Дону",
};

const pillsList = {
    antistenMV: "Антистен",
    artrafik: "артрафик",
    cerpehol: "Церпехол",
    citipigam: "Цитипигам",
    detravenol: "Детравенол",
    doksiciklin: "Доксициклин",
    jekurohol: "Экурохол",
    jekziter: "Экзитер",
    jeljufor: "Элюфор",
    esslial: "эсслиал форте",
    estetatet: "эстет а тет",
    jezlor: "Эзлор",
    fazostabil: "Фазостабил",
    flebofa: "Флебофа",
    frejmitus: "Фреймитус",
    holikron: "Холикрон",
    izislip: "Изислип",
    kontrazud: "Контразуд",
    ksilometazolin: "Ксилометазолин",
    kruoksaban: "круоксабан",
    magnepasit: "Магнепасит",
    mebespalin: "Мебеспалин",
    meladapt: "Меладапт",
    memoritab: "Меморитаб",
    motilorus: "Мотилорус",
    mukocil: "Мукоцил",
    naftoderil: "Нафтодерил",
    nekspra: "Некспра",
    nikurilly: "Никуриллы",
    noocil: "Нооцил",
    orungamin: "Орунгамин",
    pronokognil: "Пронокогнил",
    safiston: "Сафистон",
    seltavir: "Сельтавир",
    skinoklir: "Скиноклир",
    sumatrolid: "Суматролид",
    tiloram: "Тилорам",
    tolizor: "Толизор",
    ["toraksol-soljushn-tablets"]: "Тораксол солюшн таблетс",
    ulblok: "Ульблок",
    vildegra: "Вилдегра",
    ["vismuta-trikalija-dicitrat"]: "Висмута трикалия дицитрат",
};

const producerNames = [`ООО "Озон"`, "Озон ООО/Озон Фарм ООО", "Озон ООО", "Озон Фарм ООО"];

const isOzonProducer = (pillItem) => {
    return producerNames.some(item => item === pillItem.producer);
}

const getBody = (pillLabel, regionKey) => {
    return {
        operationName: "SearchQuery",
        query: `
            query SearchQuery($search: ProductSearch!, $regionID: ID!, $advertisementType: AdvertisementType!, $query: String, $skipFeaturing: Boolean = false) {
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
            }
        `,
        variables: {
            search: {
                advertisementKey: "search",
                filters: {
                    simple: [
                        {
                            id: "query",
                            val: pillLabel,
                        }
                    ]
                },
                paginator: {
                    limit: 100,
                    offset: 0
                },
                regionID: regionKey,
                query: pillLabel,
            },
            regionID: regionKey,
            advertisementType: "featuring_search",
            query: pillLabel,
            skipFeaturing: false
        }
    }

};

const getReqFields = (pillsValue, regionKey) => {
    const body = getBody(pillsValue, regionKey);

    return {
        url: "https://zdravcity.ru/bff/query",
        options: {
            method: 'POST',
            headers: {
                "Content-type": "application/json",
            },
            body: JSON.stringify(body),
        }
    }
};

const getPillsForRegion = async(pillValue, regionKey, regionValue) => {
    const req = getReqFields(pillValue, regionKey);

    return fetch(req.url, req.options)
        .then(async (res) => {
            const data = await res.json();
            if (!data || !data.data || !data.data.products || !data.data.products.items) {
                return null;
            }

            const productsArray = data.data.products.items;
            const productsFilteredByProducer = productsArray.filter(pillItem => isOzonProducer(pillItem));
            const productsWithRegion = productsFilteredByProducer.map(item => {
                return {
                    region: regionValue,
                    ...item,
                }
            })

            return productsWithRegion;
        })
        .catch(err => {
            console.error("err: ", err)
            return err;
        })
}

const getPillsForAllRegions = async(pillValue) => {
    let allRegionPill = [];

    for (const [regionKey, regionValue] of Object.entries(regionsList)) {
        const pillAllRegion = await getPillsForRegion(pillValue, regionKey, regionValue);
        allRegionPill.push(...pillAllRegion);
    }

    return allRegionPill;
}

export const fetchAllPills = async() => {
    const allPills = [];

    console.log("Загрузка...");

    for(const [pillKey, pillValue] of Object.entries(pillsList)) {
        const parsedPill = await getPillsForAllRegions(pillValue);
        allPills.push(...parsedPill);
    }

    return allPills;
}


// const bodyRAW = {
//     operationName: "SearchQuery",
//     query: `
//         query SearchQuery($search: ProductSearch!, $regionID: ID!, $advertisementType: AdvertisementType!, $query: String, $skipFeaturing: Boolean = false) {
//             products(search: $search) { 
//                 items {
//                     ...ProductSummaryFragment
//                     reviewsCount
//                 }
//                 total
//                 description
//             }
//             featuring: advertisement(
//                 regionID: $regionID
//                 type: $advertisementType
//                 query: $query
//             ) @skip(if: $skipFeaturing) {
//                 title
//                 banner
//                 advertiser
//                 counters {
//                     ...AdvertisedCountersFragment
//                 }
//                 products {
//                     ...ProductCardRegularFragment
//                 }
//             }
//         }
//         fragment ProductSummaryFragment on ProductSummary {
//             id
//             lastPrice
//             alias
//             sku
//             availableForOrder
//             availableForBooking
//             hasOnlyBookingPrices
//             canPayOnPickup
//             brand {
//                 alias
//                 name
//             }
//             badges {
//                 description
//                 title
//                 type
//                 dateEnd
//                 src
//                 color
//             }
//             maxQuantity
//             expirationDate
//             deliveryDate
//             discount
//             url
//             image
//             isFavorite
//             price
//             priceOld
//             seoBasketText
//             seoPriceText
//             priceTypeID
//             name
//             pickupDate
//             rating
//             mnns {
//                 alias
//                 href
//                 title
//             }
//             producer
//             category {
//                 ...CategoryFragment
//             }
//             producerCountry
//             warning
//             counters {
//                 ...CountersFragment
//             }
//             patternBlock
//             bitrixID
//             patternBlockPosition
//             isAdv
//             reviewsDisabled
//             prices {
//                 dateExpired
//                 discount
//                 maxQuantity
//                 price
//                 priceOld
//                 priceTypeID
//                 dateExpiredOpened
//                 isExpirating
//                 hasPrefix
//             }
//             isBundle
//             bundleItemsSimple {
//                 id
//                 quantity
//                 shortName
//             }
//         }
//         fragment CategoryFragment on Category {
//             id
//             name
//             alias
//             url
//             path
//             level
//             img
//             background
//         }
//         fragment AdvertisedCountersFragment on AdvertisementCounters {
//             imps
//             track
//             creative
//             id
//             name
//             position
//         }
//         fragment CountersFragment on Counters {
//             yandex {
//                 ...YandexCountersFragment
//             }
//             advertisement {
//                 ...AdvertisedCountersFragment
//             }
//         }
//         fragment YandexCountersFragment on YandexCounters {
//             items {
//                 value
//                 utmConfig {
//                     key
//                     value
//                 }
//                 token
//             }
//         }
//         fragment ProductCardRegularFragment on ProductSummary {
//             id
//             lastPrice
//             alias
//             sku
//             availableForOrder
//             availableForBooking
//             hasOnlyBookingPrices
//             canPayOnPickup
//             maxQuantity
//             expirationDate
//             deliveryDate
//             discount
//             url
//             image
//             isFavorite
//             price
//             priceOld
//             seoBasketText
//             seoPriceText
//             priceTypeID
//             name
//             pickupDate
//             rating
//             producerCountry
//             warning
//             isAdv
//             isBundle
//             brand {
//                 alias
//                 name
//             }
//             badges {
//                 description
//                 title
//                 type
//                 dateEnd
//                 src
//                 color
//             }
//             category {
//                 id
//                 name
//                 path
//             }
//             counters {
//                 ...CountersFragment
//             }
//             prices {
//                 dateExpired
//                 discount
//                 maxQuantity
//                 price
//                 priceOld
//                 priceTypeID
//                 dateExpiredOpened
//                 isExpirating
//                 hasPrefix
//             }
//             bundleItemsSimple {
//                 id
//                 quantity
//                 shortName
//             }
//                 reviewsDisabled
//         }
//     `,
//     variables: {
//         search: {
//             advertisementKey: "search",
//             filters: {
//                 simple: [
//                     {
//                         id: "query",
//                         val: "эсслиал форте"
//                     }
//                 ]
//             },
//             paginator: {
//                 limit: 24,
//                 offset: 0
//             },
//             regionID: "rostov",
//             query: "эсслиал форте"
//         },
//         regionID: "rostov",
//         advertisementType: "featuring_search",
//         query: "эсслиал форте",
//         skipFeaturing: false
//     }
// };