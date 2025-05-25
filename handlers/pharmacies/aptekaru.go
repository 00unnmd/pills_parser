package pharmacies

import (
	"encoding/json"
	"fmt"
	"github.com/00unnmd/pills_parser/models/pharmacies"
	"os"
	"time"

	"github.com/00unnmd/pills_parser/internals/utils"
	"github.com/00unnmd/pills_parser/models"
)

func getGroupId(resultItem pharmacies.ARResultItem) (string, error) {
	if resultItem.UniqueItemInfo.Id != "" {
		return resultItem.UniqueItemInfo.Id, nil
	}
	if len(resultItem.ItemVariantsInfo) > 0 && resultItem.ItemVariantsInfo[0].Id != "" {
		return resultItem.ItemVariantsInfo[0].Id, nil
	}

	return "", fmt.Errorf("getGroupId: no valid groupID found")
}

func getARGroupInfo(resultItem pharmacies.ARResultItem) ([]pharmacies.ARItemInfo, error) {
	groupId, err := getGroupId(resultItem)
	if err != nil {
		return nil, fmt.Errorf("getARGroupInfo error getting groudID: %w", err)
	}

	params := map[string]string{
		"itemGroupId": groupId,
	}

	respBodyByte, err := makeAPIRequest(
		"AR",
		"GET",
		os.Getenv("AR_REQ_GROUP_INFO"),
		params,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var respBody pharmacies.ARGroupInfoBody
	if err := json.Unmarshal([]byte(respBodyByte), &respBody); err != nil {
		return nil, fmt.Errorf("getARGroupInfo error unmarshaling response: %w", err)
	}

	var groupItems []pharmacies.ARItemInfo
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

	_, err := makeAPIRequest(
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

func GetARPills(pillValue string, regionValue string) ([]models.ParsedItem, error) {
	params := map[string]string{
		"page":          "0",
		"pageSize":      "25",
		"withprice":     "false",
		"withprofit":    "false",
		"withpromovits": "false",
		"phrase":        pillValue,
	}

	respBodyByte, err := makeAPIRequest(
		"AR",
		"GET",
		os.Getenv("AR_REQ_SEARCH"),
		params,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var respBody pharmacies.ARSearchBody
	if err := json.Unmarshal([]byte(respBodyByte), &respBody); err != nil {
		return nil, fmt.Errorf("GetARPills error unmarshaling response: %w", err)
	}

	filteredProductItems := utils.FilterByProducer(respBody.Result, pillValue)
	if len(filteredProductItems) == 0 {
		return nil, fmt.Errorf("не найдено препаратов удовлетворяющих запросу: len(filteredProductItems) == 0")
	}

	result := make([]models.ParsedItem, 0)

	for _, item := range filteredProductItems {
		time.Sleep(utils.RequestDelay)

		groupItems, err := getARGroupInfo(item)
		if err != nil {
			pi := utils.CreatePIWithError(pillValue, regionValue, err)
			result = append(result, pi...)
			break
		}

		parsedData := utils.ParseRawData(groupItems)
		for i := range parsedData {
			parsedData[i].Region = regionValue
		}
		result = append(result, parsedData...)
	}

	return result, nil
}
