package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type obj = map[string]any

func main() {
	var larkPushKey = os.Getenv("LARK_PUSH_KEY")
	if larkPushKey == "" {
		fmt.Println("LARK_PUSH_KEY is empty")
		return
	}
	var qCloudCosDomain = os.Getenv("QCLOUD_COS_DOMAIN")
	if qCloudCosDomain == "" {
		fmt.Println("QCLOUD_COS_DOMAIN is empty")
		return
	}
	var ciCommitRefName = os.Getenv("CI_COMMIT_REF_NAME")
	if ciCommitRefName == "" {
		fmt.Println("CI_COMMIT_REF_NAME is empty")
		return
	}
	var url = fmt.Sprintf("https://%s/client/%s/pro.zip", qCloudCosDomain, ciCommitRefName)
	f, e := os.Open("out/zip/pro.zip")
	if e != nil {
		panic(e)
	}
	var cSign = `🔖`
	if os.Getenv("CI_COMMIT_TAG") == "" {
		cSign = `✔`
	}
	defer func() { _ = f.Close() }()
	var sha = sha256.New()
	_, e = io.Copy(sha, f)
	if e != nil {
		panic(e)
	}
	var hash = hex.EncodeToString(sha.Sum(nil))
	var t map[string]any
	_ = json.Unmarshal([]byte(template), &t)
	t["elements"].([]any)[0].(obj)["content"] = fmt.Sprintf("**📦Bundle：** %s\n**SHA-256: **%s", url, hash)
	t["header"].(obj)["title"].(obj)["content"] = fmt.Sprintf("%sClient - %s", cSign, ciCommitRefName)
	t = map[string]any{
		"msg_type": "interactive",
		"card":     t,
	}
	postData, e := json.Marshal(t)
	if e != nil {
		panic(e)
	}
	req, e := http.NewRequest("POST", "https://open.feishu.cn/open-apis/bot/v2/hook/"+larkPushKey, bytes.NewReader(postData))
	if e != nil {
		panic(e)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		panic(e)
	}
	fmt.Printf("status: %d - %s\n", resp.StatusCode, resp.Status)
	if resp.StatusCode > 299 {
		panic("failed to push message")
	}
}

var template = `
{
    "config": {
        "wide_screen_mode": true
    },
    "elements": [
        {
            "tag": "markdown",
            "content": ""
        }
    ],
    "header": {
        "template": "green",
        "title": {
            "content": "",
            "tag": "plain_text"
        }
    }
}
`
