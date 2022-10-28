package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

type azureTable struct {
	client *aztables.Client
}

func (t *azureTable) init(connStr, table string) error {
	var err error
	client, err := aztables.NewServiceClientFromConnectionString(connStr, nil)

	if err == nil {
		t.client = client.NewClient(table)
	}
	t.client.CreateTable(context.TODO(), nil)
	return err
}

func (t *azureTable) getEntry(dns string) (aztables.GetEntityResponse, bool) {
	var exist bool
	var result aztables.GetEntityResponse
	pKey, rKey, ok := parseDomain(dns)

	if ok {
		var err error
		result, err = t.client.GetEntity(context.TODO(), pKey, rKey, nil)

		if err == nil {
			exist = true
		}
	}

	return result, exist
}

func (t *azureTable) setEntry(dns string) {
	pKey, rKey, ok := parseDomain(dns)

	if ok {
		entity := aztables.Entity{
			PartitionKey: pKey,
			RowKey:       rKey,
		}

		newEntity, err := json.Marshal(entity)
		if err != nil {
			t.client.AddEntity(context.TODO(), newEntity, nil)
		}

	}

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

func main() {

	var aztable azureTable

	connStr, exist := os.LookupEnv("AzureWebJobsStorage")

	if !exist {
		fmt.Println("[ERR] Couldn't obtain connection string")
		return
	}
	tableName := "table3"

	err := aztable.init(connStr, tableName)

	if err != nil {
		panic(err)
	}
	aztable.setEntry("test.com")
	_, blocked := aztable.getEntry("aaa.test.com")

	fmt.Println(blocked)
}
