package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

const IPV4_ADDRESS_OCTETS = 4

func InitializeBannedIpsModelFromDB() {
	BannedIpsGroupedByDomain, err := GetAllIpAddressesFromBannedIpsTable()
	if err != nil {
		glog.Errorf("Failed to load banned ip address from the DB")
		return
	}

	BannedIpv4AddressTriesPerDomain, err = CreateIpv4BannedAddressTriePerDomain(BannedIpsGroupedByDomain)
	if err != nil {
		glog.Errorf("Failed to create trie model for banned ipv4 addresses : %s", err.Error())
		return
	}
}

func CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain map[int][]string) (map[int]*ipv4AddressTrie, error) {
	bannedIpsPerDomainMap := make(map[int]*ipv4AddressTrie)
	for domainId, ipAddressList := range bannedIpsByDomain {
		ipv4AddressTrieRootNode := createNewIpv4AddressTrie()
		err := addIpv4AddressesToTrieFromIpList(ipv4AddressTrieRootNode, ipAddressList)
		if err != nil {
			return nil, err
		}
		bannedIpsPerDomainMap[domainId] = ipv4AddressTrieRootNode
	}
	return bannedIpsPerDomainMap, nil
}

func createNewIpv4AddressTrie() *ipv4AddressTrie {
	return &ipv4AddressTrie{
		root: &ipv4AddressTrieNode{
			addresOctet:        0,
			children:           make(map[byte]*ipv4AddressTrieNode),
			octectIndex:        0,
			isLastAddressOctet: false,
		},
	}
}

func addIpv4AddressesToTrieFromIpList(root *ipv4AddressTrie, ipv4AddressList []string) error {
	for _, ipv4Address := range ipv4AddressList {
		err := root.insertIpv4AddressToTrie(ipv4Address)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *ipv4AddressTrie) insertIpv4AddressToTrie(ipv4Address string) error {
	cur := t.root
	addressIndex := 0

	addressOctets, parseError := parseIpv4Addresses(ipv4Address)
	if parseError != nil || addressOctets == nil {
		return parseError
	}

	for _, curAddressOctet := range addressOctets {
		if cur.children[curAddressOctet] == nil {
			addressIndex++
			cur.children[curAddressOctet] = &ipv4AddressTrieNode{
				addresOctet: curAddressOctet,
				children:    make(map[byte]*ipv4AddressTrieNode),
				octectIndex: addressIndex,
			}
		}

		cur = cur.children[curAddressOctet]
	}

	cur.isLastAddressOctet = true
	return nil
}

func (t *ipv4AddressTrie) SearchIpv4AddressInTrie(ipv4Address string) (bool, error) {
	node, err := t.traverseAllNodesInTheIpv4Address(ipv4Address)
	if err != nil {
		glog.Error(err.Error())
		return false, err
	}

	return (node != nil && node.isLastAddressOctet && node.octectIndex <= IPV4_ADDRESS_OCTETS), nil
}

func (t *ipv4AddressTrie) SearchIpv4AddressPrefixInTrie(ipv4AddressPrefix string) (bool, error) {
	node, err := t.traverseAllNodesInTheIpv4Address(ipv4AddressPrefix)
	if err != nil {
		glog.Error(err.Error())
		return false, err
	}
	return node != nil, nil
}

func (t *ipv4AddressTrie) traverseAllNodesInTheIpv4Address(ipv4Address string) (*ipv4AddressTrieNode, error) {
	cur := t.root
	addressOctets, parseError := parseIpv4Addresses(ipv4Address)
	if parseError != nil || addressOctets == nil {
		return nil, parseError
	}

	for _, curAddressOctet := range addressOctets {
		if cur.children[curAddressOctet] == nil {
			return nil, nil
		}

		cur = cur.children[curAddressOctet]
	}

	return cur, nil
}

func parseIpv4Addresses(ipv4Address string) ([]byte, error) {
	ipv4AddressOctets := make([]byte, 0)
	addressSegments := strings.Split(ipv4Address, ".")

	if len(addressSegments) > IPV4_ADDRESS_OCTETS {
		return nil, fmt.Errorf("error ipv4 address %s is invalid", ipv4Address)
	}

	for idx, segment := range addressSegments {
		if segment == "*" {
			wildCardValidationError := validateWildCardsInIpv4Address(addressSegments[idx:])
			if wildCardValidationError != nil {
				return nil, wildCardValidationError
			}

			// wild card validation successful no further entries are possible
			//Also, handle ipv4 address that has only wildcards(i.e there are no prior valid entries)
			if len(ipv4AddressOctets) == 0 {
				return nil, fmt.Errorf("error ipv4 address %s is invalid", ipv4Address)
			}

			return ipv4AddressOctets, nil
		}

		addressOctet, err := strconv.Atoi(segment)
		if err != nil {
			return nil, err
		}

		if addressOctet < 0 || addressOctet > 255 {
			return nil, fmt.Errorf("error ipv4 address octet %d exceeds expected range ", addressOctet)
		}

		ipv4AddressOctets = append(ipv4AddressOctets, byte(addressOctet))
	}
	return ipv4AddressOctets, nil
}

func validateWildCardsInIpv4Address(addressSegmentsToBeValidated []string) error {
	//check if the rest of address segments do not contain characters other than "*"
	for _, segment := range addressSegmentsToBeValidated {
		if segment != "*" {
			return fmt.Errorf("error ipv4 address octet %s is invalid ", segment)
		}
	}
	return nil
}
