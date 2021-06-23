package models

import (
	"fmt"

	"github.com/golang/glog"
)

func InitializeModelsFromDB() {
	initializeBannedIpListModelFromDB()
}

func initializeBannedIpListModelFromDB() {
	ipAddressList, err := GetAllIpAddressesFromBannedIpsTable()
	if err != nil {
		glog.Errorf("Failed to load banned ip address from the DB")
		return
	}

	BannedIpsGroupedByDomain = ipAddressList

	fmt.Println(BannedIpsGroupedByDomain)
}
