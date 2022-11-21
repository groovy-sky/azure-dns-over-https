package table

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

type AzureTable struct {
	client *aztables.Client
}

func (t *AzureTable) Init(connStr, table string) error {
	var err error
	client, err := aztables.NewServiceClientFromConnectionString(connStr, nil)

	if err == nil {
		t.client = client.NewClient(table)
	}
	t.client.CreateTable(context.TODO(), nil)
	return err
}

func (t *AzureTable) GetEntry(dns string) (aztables.GetEntityResponse, bool) {
	var exist bool
	var result aztables.GetEntityResponse
	var err error
	pKey, rKey, ok := parseDomain(dns)

	if ok {
		rKeyLvl := strings.Count(rKey, ".")
		result, err = t.client.GetEntity(context.TODO(), pKey, rKey, nil)

		if err != nil {
			if rKeyLvl > 1 {
				rKey = strings.SplitAfter(rKey, ".")[1]
			} else {
				rKey = ""
			}
			result, err = t.client.GetEntity(context.TODO(), pKey, rKey, nil)
		}
		if err == nil {
			exist = true
		}

	}

	return result, exist
}

func (t *AzureTable) SetEntry(dns string) error {
	pKey, rKey, ok := parseDomain(dns)
	outerr := errors.New("couldn't parse DNS")
	if ok {
		entity := aztables.Entity{
			PartitionKey: pKey,
			RowKey:       rKey,
		}

		newEntity, err := json.Marshal(entity)
		if err == nil {
			_, outerr = t.client.AddEntity(context.TODO(), newEntity, nil)
		} else {
			outerr = err
		}

	}
	return outerr

}

func parseDomain(dns string) (string, string, bool) {
	dns = strings.ToLower(dns)
	var valid bool
	var domainName, subDomain string
	domains := strings.Split(dns, ".")
	domains_len := len(domains)
	switch {
	case domains_len >= 2:
		domainName = domains[domains_len-2] + "." + domains[domains_len-1]
		valid = true
		fallthrough
	case domains_len > 2:
		subDomain = dns[:len(dns)-len(domainName)]
	}
	return domainName, subDomain, valid
}
