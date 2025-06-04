package calls

import (
	"encoding/json"
	"fmt"
	"github.com/00unnmd/pills_parser/internal/domain"
	"github.com/00unnmd/pills_parser/internal/transport"
	"github.com/00unnmd/pills_parser/pkg/utils"
	"os"
	"strconv"
	"strings"
	"time"
)

func getIdFromUrl(url string) (string, error) {
	parts := strings.Split(url, "-")
	if len(parts) == 0 {
		return "", fmt.Errorf("url length == 0")
	}

	id := parts[len(parts)-1]
	if id == "" {
		return "", fmt.Errorf("id is not found in url")
	}

	return id, nil
}

func getGroupId(resultItem domain.ARResultItem) (string, error) {
	if len(resultItem.ItemVariantsInfo) > 0 && resultItem.ItemVariantsInfo[0].Id != "" {
		return resultItem.ItemVariantsInfo[0].Id, nil
	} else {
		id, err := getIdFromUrl(resultItem.HumanableUrl)
		if err != nil {
			return "", fmt.Errorf("getIdFromUrl: %w", err)
		}

		return id, nil
	}
}

func getARGroupInfo(resultItem domain.ARResultItem) ([]domain.ARItemInfo, error) {
	groupId, err := getGroupId(resultItem)
	if err != nil {
		return nil, fmt.Errorf("getARGroupInfo error getting groupID: %w", err)
	}

	params := map[string]string{
		"itemGroupId": groupId,
	}

	respBodyByte, err := transport.MakeAPIRequest(
		"AR",
		"GET",
		os.Getenv("AR_REQ_GROUP_INFO"),
		params,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var respBody domain.ARGroupInfoBody
	if err := json.Unmarshal(respBodyByte, &respBody); err != nil {
		return nil, fmt.Errorf("getARGroupInfo error unmarshaling response: %w", err)
	}

	var groupItems []domain.ARItemInfo
	for _, groupItem := range respBody.GroupItems {
		groupItems = append(groupItems, groupItem.ItemInfos...)
	}

	return groupItems, nil
}

func ChangeARRegion(regionId string) (bool, error) {
	body := map[string]interface{}{
		"id":           regionId,
		"manualChange": true,
	}

	_, err := transport.MakeAPIRequest(
		"AR",
		"PUT",
		os.Getenv("AR_REQ_USER_CITY"),
		nil,
		body,
	)
	if err != nil {
		return false, fmt.Errorf("error changing region: %w", err)
	}

	return true, nil
}

func GetARPills(pillValue string, regionValue string, withFilter bool) ([]domain.ParsedItem, error) {
	params := map[string]string{
		"pageSize":      "100",
		"withprice":     "false",
		"withprofit":    "false",
		"withpromovits": "false",
		"phrase":        pillValue,
	}

	var rawResult []domain.ARResultItem
	processedCount := 0
	page := 0

	for {
		params["page"] = strconv.Itoa(page)

		respBodyByte, err := transport.MakeAPIRequest(
			"AR",
			"GET",
			os.Getenv("AR_REQ_SEARCH"),
			params,
			nil,
		)
		if err != nil {
			return nil, err
		}

		var respBody domain.ARSearchBody
		if err := json.Unmarshal(respBodyByte, &respBody); err != nil {
			return nil, fmt.Errorf("GetARPills error unmarshaling response: %w", err)
		}

		rawResult = append(rawResult, respBody.Result...)
		processedCount += respBody.CurrentCount

		if respBody.CurrentCount == 0 || processedCount >= respBody.TotalCount {
			break
		}

		page++
		time.Sleep(utils.RequestDelay)
	}

	filteredData := rawResult
	if withFilter == true {
		filteredData = utils.FilterByProducer(rawResult, pillValue)
	}

	if len(filteredData) == 0 {
		return nil, fmt.Errorf("не найдено препаратов удовлетворяющих запросу: len(filteredData) == 0")
	}

	var searchItems []domain.ARResultItem
	var groupItems []domain.ARResultItem

	for _, item := range filteredData {
		if item.ItemsCount < 2 {
			searchItems = append(searchItems, item)
		} else {
			groupItems = append(groupItems, item)
		}
	}

	result := utils.ParseRawData("aptekaru", regionValue, pillValue, searchItems)

	for _, item := range groupItems {
		time.Sleep(utils.RequestDelay)

		groupItemInfos, err := getARGroupInfo(item)
		if err != nil {
			pi := utils.CreatePIWithError(pillValue, regionValue, err, "aptekaru")
			result = append(result, pi...)
			break
		}

		parsedData := utils.ParseRawData("aptekaru", regionValue, pillValue, groupItemInfos)
		result = append(result, parsedData...)
	}

	return result, nil
}
