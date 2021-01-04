package main

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type airtableClient struct {
	baseID string
	apiKey string
}

// record models fields on each Airtable record for
// our annotations table.
type record struct {
	ID     string     `json:"id"` // Airtable record ID
	Fields annotation `json:"fields"`
}

type annotation struct {
	FirstName    string   `json:"First"`
	LastName     string   `json:"Last"`
	Associations []string `json:"Associations"`
	Role         string   `json:"Role"`
}

// forEachRecord will apply the function rf to each record returned
// from the Airtable Master List of records
func (a *airtableClient) forEachRecord(rf func(r record) error) error {
	client := &http.Client{}
	offset := ""
	for {
		url := "https://api.airtable.com/v0/" + a.baseID + "/Master%20List?pageSize=100"
		if offset != "" {
			url += "&offset=" + offset
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return errors.Wrap(err, "creating request")
		}
		req.Header.Add("Authorization", "Bearer "+a.apiKey)
		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrap(err, "getting master list of records")
		}

		var body struct {
			Records []record `json:"records"`
			Offset  string   `json:"offset"`
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			return errors.Wrap(err, "unmarshaling response")
		}
		for _, r := range body.Records {
			err := rf(r)
			if err != nil {
				return errors.Wrap(err, "processing record")
			}
		}
		// Check offset to fetch next page of results
		if body.Offset != "" {
			offset = body.Offset
		} else {
			break
		}
	}
	return nil
}
