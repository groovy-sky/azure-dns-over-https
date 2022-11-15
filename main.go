package main

import (
	"fmt"
	"os"
)

func main() {
	var aztable AzureTable

	connStr, exist := os.LookupEnv("AzureWebJobsStorage")

	if !exist {
		fmt.Println("[ERR] Couldn't obtain connection string")
		return
	}
	tableName := "table3"

	err := table.aztable.init(connStr, tableName)

	if err != nil {
		panic(err)
	}
	aztable.setEntry("test.com")
	_, blocked := aztable.getEntry("aaa.test.com")

	fmt.Println(blocked)
}
