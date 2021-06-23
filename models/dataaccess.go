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

func GetAllIpAddressesFromBannedIpsTable() (map[int][]string, error) {
	rows, err := DB.Query(`
								SELECT domain_id, host(ip)
								FROM bannedips`)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()
	bannedIpAddressesBasedOnDomain := make(map[int][]string)

	for rows.Next() {
		var domainId int
		var ipAddress string
		scanRowErr := rows.Scan(&domainId, &ipAddress)
		if scanRowErr != nil {
			log.Fatal(scanRowErr)
			return nil, scanRowErr
		}

		bannedIpAddressesBasedOnDomain[domainId] = append(bannedIpAddressesBasedOnDomain[domainId], ipAddress)
	}

	return bannedIpAddressesBasedOnDomain, nil
}
