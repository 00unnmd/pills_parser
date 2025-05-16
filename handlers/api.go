package handlers

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func makeAPIRequest(serviceName string, method, url string, queryParams map[string]string, bodyData interface{}) ([]byte, error) {
	var body io.Reader

	if bodyData != nil {
		jsonData, err := json.Marshal(bodyData)
		if err != nil {
			return nil, fmt.Errorf("makeAPIRequest! error marshaling request body: %w", err)
		}

		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("makeAPIRequest! error creating request: %w", err)
	}

	if queryParams != nil {
		q := req.URL.Query()

		for key, value := range queryParams {
			q.Add(key, value)
		}

		req.URL.RawQuery = q.Encode()
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	if serviceName == "AR" {
		req.Header.Add("Accept", "application/json, text/plain, */*")
		req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
		req.Header.Add("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8,ru;q=0.7")
		req.Header.Add("Authorization", os.Getenv("AR_AUTH_TOKEN"))
		req.Header.Add("Connection", "keep-alive")
		req.Header.Add("Device-id", os.Getenv("AR_DEVICE_ID"))
		req.Header.Add("Host", os.Getenv("AR_HOST"))
		req.Header.Add("Origin", os.Getenv("AR_ORIGIN"))
		req.Header.Add("Referer", os.Getenv("AR_REFERER"))
		req.Header.Add("Sec-ch-ua", `"Chromium";v="136", "Google Chrome";v="136", "Not.A/Brand";v="99"`)
		req.Header.Add("Sec-ch-ua-mobile", "?0")
		req.Header.Add("Sec-ch-ua-platform", "Windows")
		req.Header.Add("Sec-Fetch-Mode", "cors")
		req.Header.Add("Sec-Fetch-Site", "same-site")
		req.Header.Add("User-session-id", os.Getenv("AR_USER_SESSION_ID"))
		req.Header.Add("X-active-exp", os.Getenv("AR_X_ACTIVE_EXP"))
		req.Header.Add("X-ciid-b", os.Getenv("AR_X_CIID_B"))
		req.Header.Add("X-ciid-h", os.Getenv("AR_X_CIID_H"))
		req.Header.Add("X-dreg-tst", os.Getenv("AR_X_DREG_TST"))
		req.Header.Add("ym-aru-visorc", os.Getenv("AR_YM_ARU_VISORC"))
	}

	if serviceName == "AR" && url == os.Getenv("AR_REQ_USER_CITY") {
		req.Header.Add("Content-length", "53")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("makeAPIRequest! error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("makeAPIRequest! unexpected status code: %d", resp.StatusCode)
	}

	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("makeAPIRequest! gzip decompression failed: %w", err)
		}

		defer gzReader.Close()
		reader = gzReader
	case "deflate":
		reader = flate.NewReader(resp.Body)
	default:
		reader = resp.Body
	}

	return io.ReadAll(reader)
}
