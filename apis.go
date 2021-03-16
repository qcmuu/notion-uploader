package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Command Types
const (
	CommandSet            = "set"
	CommandUpdate         = "update"
	CommandListAfter      = "listAfter"
	CommandListRemove     = "listRemove"
	notionHost            = "https://www.notion.so"
	signedURLPrefix       = "https://www.notion.so/signed"
	s3URLPrefix           = "https://s3-us-west-2.amazonaws.com/secure.notion-static.com/"
	notionSaveTranscation = "https://www.notion.so/api/v3/saveTransactions"
	userAgent             = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3483.0 Safari/537.36"
	acceptLang            = "en-US,en;q=0.9"
)

func SubmitTransaction(ops []*Operation, spaceID string) error {

	reqData := &submitTransactionRequest{
		RequestID: uuid.New().String(),
		Transaction: []Transaction{{
			ID:         uuid.New().String(),
			SpaceID:    spaceID,
			Operations: ops,
		}},
	}
	// PrintStruct(reqData)

	js, err := json.Marshal(reqData)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", notionSaveTranscation, bytes.NewBuffer(js))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", acceptLang)
	req.Header.Set("cookie", fmt.Sprintf("token_v2=%v", *token))
	var rsp *http.Response

	http.DefaultClient.Timeout = time.Second * 30
	rsp, err = http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	_ = rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return fmt.Errorf("http.Post returned non-200 status code of %d, returns: %s", rsp.StatusCode, d)
	}

	if !bytes.Equal(d, []byte("{}")) {
		return fmt.Errorf("unknown error: %s", d)
	}
	return nil
}

func doNotionAPI(c *Client, apiURL string, requestData interface{}, result interface{}) (map[string]interface{}, error) {
	var js []byte
	var err error
	if requestData != nil {
		js, err = json.Marshal(requestData)
		if err != nil {
			return nil, err
		}
	}
	uri := notionHost + apiURL
	body := bytes.NewBuffer(js)
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", acceptLang)
	if c.AuthToken != "" {
		req.Header.Set("cookie", fmt.Sprintf("token_v2=%v", c.AuthToken))
	}
	var rsp *http.Response

	rsp, err = http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer closeNoError(rsp.Body)

	if rsp.StatusCode != 200 {
		_, _ = ioutil.ReadAll(rsp.Body)
		return nil, fmt.Errorf("http.Post('%s') returned non-200 status code of %d", uri, rsp.StatusCode)
	}
	d, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(d, result)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(d, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func closeNoError(c io.Closer) {
	_ = c.Close()
}
