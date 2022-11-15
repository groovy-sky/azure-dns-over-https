package main

import (
	"fmt"
	"os"

	"github.com/groovy-sky/azuredoh/v2/pkg/table"
)

func main() {
	var aztable table.AzureTable

	connStr, exist := os.LookupEnv("AzureWebJobsStorage")

	if !exist {
		fmt.Println("[ERR] Couldn't obtain connection string")
		return
	}
	tableName := "table3"

	err := aztable.Init(connStr, tableName)

	if err != nil {
		panic(err)
	}
	aztable.SetEntry("test.com")
	_, blocked := aztable.GetEntry("aaa.test.com")

	fmt.Println(blocked)
}
