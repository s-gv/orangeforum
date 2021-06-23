package models

import "github.com/golang/glog"

func InitializeModelsFromDB() {
	initializeBannedIpListModelFromDB()
}

func initializeBannedIpListModelFromDB() {
	ipAddressList, err := GetAllIpAddressesFromBannedIpsTable()
	if err != nil {
		glog.Errorf("Failed to load banned ip address from the DB")
		return
	}

	BannedIpAddresses = ipAddressList
}
