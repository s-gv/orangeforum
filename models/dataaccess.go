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
