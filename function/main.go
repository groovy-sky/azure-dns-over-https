package main

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"

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
	return err
}

func (t *azureTable) getEntry(dns string) (string, bool) {
	var result string
	exist := false

	h := sha512.New()
	h.Write([]byte(dns))

	hashed := hex.EncodeToString(h.Sum([]byte("com")))

	fmt.Println(hashed)

	fmt.Println(t.client.GetEntity(context.TODO(), hashed, "", nil))
	return result, exist
}

func main() {
	var aztable azureTable

	connStr, exists := os.LookupEnv("AzureWebJobsStorage")

	if !exists {
		fmt.Println("[ERR] Couldn't obtain connection string")
		return
	}
	tableName := "table1"

	err := aztable.init(connStr, tableName)

	if err != nil {
		panic(err)
	}

	aztable.getEntry("test.com")
}
