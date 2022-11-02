package zinc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	host = "localhost"
)

type (
	Shit struct {
		Id      string `json:"_id"`
		Content string
		From    string
		To      string
		Subject string
	}
	Ssource struct {
		Source Shit `json:"_source"`
	}
	SQty struct {
		Value int
	}
	Shits struct {
		Total SQty
		Hits  []Ssource
	}
	Chits struct {
		Hits Shits
	}
	respCreate struct {
		Count int `json:"record_count"`
	}
)

func CreateIndex() error {
	indexDef :=
		`{
		"name": "email",
		"storage_type": "disk",
		"shard_num": 1,
		"mappings": {
			"properties": {
				"from": {
					"type": "text",
					"index": false,
					"store": false
				},
				"to": {
					"type": "text",
					"index": false,
					"store": false
				},
				"subject": {
					"type": "text",
					"index": true,
					"store": false,
					"highlightable": true
				},
				"content": {
					"type": "text",
					"index": true,
					"store": false,
					"highlightable": true
				}
			}
		}
	}`
	req, err := http.NewRequest("POST", "http://"+host+":4080/api/index", strings.NewReader(indexDef))
	if err != nil {
		return err
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	return nil
}

func DeleteIndex() (err error) {
	req, err := http.NewRequest("DELETE", "http://"+host+":4080/api/index/email", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	return nil
}

func CreateData(data string) (int, error) {
	req, err := http.NewRequest("POST", "http://"+host+":4080/api/email/_multi", strings.NewReader(data))
	if err != nil {
		return 0, err
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	res := respCreate{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return 0, err
	}
	return res.Count, nil
}

func Query(term string, from, max int) (res Chits, err error) {
	query := `{
        "search_type": "match",
        "query":
        {
            "term": "` + term + `"
        },
        "from": ` + fmt.Sprint(from) + `,
        "max_results": ` + fmt.Sprint(max) + `,
		"highlight": {
			"pre_tags": ["<pre>"],
        	"post_tags": ["</pre>"],
			"fields": {
				"content": {}
			}
		},
        "_source": []
    }`
	req, err := http.NewRequest("POST", "http://"+host+":4080/api/email/_search", strings.NewReader(query))
	if err != nil {
		return
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	return
}
