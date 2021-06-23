package models

import (
	"database/sql"
	"log"
)

func GetIpAddressFromBannedIpsTable(ipAddressToBeQueried string) (string, error) {
	row := DB.QueryRow(`
								SELECT host(ip)
								FROM bannedips
								WHERE ip = $1`, ipAddressToBeQueried)

	var bannedIp string
	err := row.Scan(&bannedIp)

	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		log.Fatal(err)
		return "", err
	}

	return bannedIp, nil
}

func GetAllIpAddressesFromBannedIpsTable() ([]string, error) {
	rows, err := DB.Query(`
								SELECT host(ip)
								FROM bannedips`)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()
	ipAddressList := make([]string, 0)

	for rows.Next() {
		var ipAddress string
		scanRowErr := rows.Scan(&ipAddress)
		if scanRowErr != nil {
			log.Fatal(scanRowErr)
			return nil, scanRowErr
		}
		ipAddressList = append(ipAddressList, ipAddress)
	}

	return ipAddressList, nil
}
