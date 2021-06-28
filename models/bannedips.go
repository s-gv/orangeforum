package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

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
			addresOctet: 0,
			children:    make(map[byte]*ipv4AddressTrieNode),
			octectIndex: 0,
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
	addressSegments := strings.Split(ipv4Address, ".")
	for _, segment := range addressSegments {
		curAddressOctet, err := getIpv4AddressOctet(segment)
		if err != nil {
			return err
		}

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
	return nil
}

func (t *ipv4AddressTrie) SearchIpv4AddressInTrie(ipv4Address string) (bool, error) {
	node, err := t.traverseAllNodesInTheIpv4Address(ipv4Address)
	if err != nil {
		return false, err
	}

	return (node != nil && node.octectIndex == 4), nil
}

func (t *ipv4AddressTrie) SearchIpv4AddressPrefixInTrie(ipv4AddressPrefix string) (bool, error) {
	node, err := t.traverseAllNodesInTheIpv4Address(ipv4AddressPrefix)
	if err != nil {
		return false, err
	}
	return node != nil, nil
}

func (t *ipv4AddressTrie) traverseAllNodesInTheIpv4Address(ipv4Address string) (*ipv4AddressTrieNode, error) {
	cur := t.root
	addressSegments := strings.Split(ipv4Address, ".")
	for _, segment := range addressSegments {
		curAddressOctet, err := getIpv4AddressOctet(segment)
		if err != nil {
			return nil, err
		}

		if cur.children[curAddressOctet] == nil {
			return nil, nil
		}
		cur = cur.children[curAddressOctet]
	}
	return cur, nil
}

func getIpv4AddressOctet(addressSegment string) (byte, error) {
	addressOctet, err := strconv.Atoi(addressSegment)
	if err != nil {
		return 0, err
	}
	if addressOctet < 0 || addressOctet > 255 {
		return 0, fmt.Errorf("error ipv4 address octet %d exceeds expected range ", addressOctet)
	}
	return byte(addressOctet), nil
}
