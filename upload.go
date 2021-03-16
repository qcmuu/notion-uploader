package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/google/uuid"
	"github.com/kjk/notionapi"
)

func upload(f string) string {

	client := &Client{notionapi.Client{AuthToken: *token}}
	page, err := client.DownloadPage(*pageid)
	if err != nil {
		log.Fatalf("DownloadPage() failed with %s\n", err)
	}

	root := page.BlockByID(page.ID)
	// PrintStruct(root)

	file, err := os.Open(f)
	info, err := os.Stat(f)

	bar := pb.Full.Start64(info.Size())
	bar.Set(pb.Bytes, true)
	reader := bar.NewProxyReader(file)

	fileID, fileURL, err := client.UploadFile(reader, info.Name(), info.Size())
	lastBlock := root.Content[len(root.Content)-1]

	userID := root.LastEditedByID
	spaceID := root.ParentID
	newBlockID := uuid.New().String()

	ops := []*Operation{
		buildOp(newBlockID, CommandSet, []string{}, map[string]interface{}{
			"type":    "file",
			"id":      newBlockID,
			"version": 1,
		}),
		buildOp(newBlockID, CommandUpdate, []string{}, map[string]interface{}{
			"parent_id":    root.ID,
			"parent_table": "block",
			"alive":        true,
		}),
		buildOp(root.ID, CommandListAfter, []string{"content"}, map[string]string{
			"id":    newBlockID,
			"after": lastBlock.ID,
		}),
		buildOp(newBlockID, CommandSet, []string{"created_by_id"}, userID),
		buildOp(newBlockID, CommandSet, []string{"created_by_table"}, "notion_user"),
		buildOp(newBlockID, CommandSet, []string{"created_time"}, time.Now().UnixNano()),
		buildOp(newBlockID, CommandSet, []string{"last_edited_time"}, time.Now().UnixNano()),
		buildOp(newBlockID, CommandSet, []string{"last_edited_by_id"}, userID),
		buildOp(newBlockID, CommandSet, []string{"last_edited_by_table"}, "notion_user"),
		buildOp(newBlockID, CommandUpdate, []string{"properties"}, map[string]interface{}{
			"source": [][]string{{fileURL}},
			"size":   [][]string{{ByteCountIEC(info.Size())}},
			"title":  [][]string{{info.Name()}},
		}),
		buildOp(newBlockID, CommandListAfter, []string{"file_ids"}, map[string]string{
			"id": fileID,
		}),
	}
	fmt.Printf("syncing blocks..")
	end := DotTicker()

	SubmitTransaction(ops, spaceID)
	*end <- struct{}{}
	return fmt.Sprintf("%s/%s?table=block&id=%s&name=%s&userId=%s&cache=v2", signedURLPrefix, url.QueryEscape(fileURL), newBlockID, info.Name(), userID)
}

// UploadFile Uploads a file to notion's asset hosting(aws s3)
func (c *Client) UploadFile(file io.Reader, name string, size int64) (fileID, fileURL string, err error) {
	ext := path.Ext(name)
	mt := mime.TypeByExtension(ext)
	if mt == "" {
		mt = "application/octet-stream"
	}
	// 1. getUploadFileURL
	uploadFileURLResp, err := c.getUploadFileURL(name, mt)
	if err != nil {
		err = fmt.Errorf("get upload file URL error: %s", err)
		return
	}

	// 2. Upload file to amazon - PUT
	httpClient := http.DefaultClient

	req, err := http.NewRequest(http.MethodPut, uploadFileURLResp.SignedPutURL, file)
	if err != nil {
		return
	}
	req.ContentLength = size
	req.TransferEncoding = []string{"identity"} // disable chunked (unsupported by aws)
	req.Header.Set("Content-Type", mt)
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var contents []byte
		contents, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			contents = []byte(fmt.Sprintf("Error from ReadAll: %s", err))
		}

		err = fmt.Errorf("http PUT '%s' failed with status %s: %s", req.URL, resp.Status, string(contents))
		return
	}

	return uploadFileURLResp.FileID, uploadFileURLResp.URL, nil
}

// getUploadFileURL executes a raw API call: POST /api/v3/getUploadFileUrl
func (c *Client) getUploadFileURL(name, contentType string) (*GetUploadFileUrlResponse, error) {
	const apiURL = "/api/v3/getUploadFileUrl"

	req := &getUploadFileUrlRequest{
		Bucket:      "secure",
		ContentType: contentType,
		Name:        name,
	}

	var rsp GetUploadFileUrlResponse
	var err error
	rsp.RawJSON, err = doNotionAPI(c, apiURL, req, &rsp)
	if err != nil {
		return nil, err
	}

	rsp.Parse()

	return &rsp, nil
}

func (r *GetUploadFileUrlResponse) Parse() {
	r.FileID = strings.Split(r.URL[len(s3URLPrefix):], "/")[0]
}
