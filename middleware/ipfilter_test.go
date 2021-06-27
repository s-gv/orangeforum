package middleware

import (
	"testing"

	"github.com/s-gv/orangeforum/models"
)

func TestTrieBasedIpAddressSearch(t *testing.T) {

	ipList1 := []string{"192.168.1.0", "19.23.5.67"}
	ipList2 := []string{"10.23.56.234"}
	ipList3 := []string{"12.23.45.67", "123.34.67.8", "45.56.78.145"}
	domainId1 := 1
	domainId2 := 2
	domainId3 := 3
	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
		domainId2: ipList2,
		domainId3: ipList3,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddresTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Errorf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId3]
	if trieToBeSearched == nil {
		t.Errorf("Failed to create trie for the specified domain %d", domainId3)
	}

	addressToBeSearched := "45.56.78.145"
	found, err := trieToBeSearched.SearchIpv4AddressInTrie(addressToBeSearched)
	if err != nil {
		t.Errorf("Error while searching for ip %s in the trie", addressToBeSearched)
		t.Error(err)
	}

	if found == false {
		t.Errorf("Ipv4 addresss %s Not found", addressToBeSearched)
	}

	if found == true {
		t.Logf("Ipv4 address %s found", addressToBeSearched)
	}
}

func TestTrieBasedIpAddressSearchForTheFailureCase(t *testing.T) {

	ipList1 := []string{"192.168.1.0", "19.23.5.67"}
	ipList2 := []string{"10.23.56.234"}
	ipList3 := []string{"12.23.45.67", "123.34.67.8", "45.56.78.145"}
	domainId1 := 1
	domainId2 := 2
	domainId3 := 3
	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
		domainId2: ipList2,
		domainId3: ipList3,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddresTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Errorf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId2]
	if trieToBeSearched == nil {
		t.Errorf("Failed to create trie for the specified domain %d", domainId3)
	}

	addressToBeSearched := "45.56.78.145"
	found, err := trieToBeSearched.SearchIpv4AddressInTrie(addressToBeSearched)
	if err != nil {
		t.Errorf("Error while searching for ip %s in the trie", addressToBeSearched)
		t.Error(err)
	}

	if found == true {
		t.Errorf("Ipv4 address %s found, while it was not expected", addressToBeSearched)
	}

	if found == false {
		t.Logf("Ipv4 addresss %s Not found as expected", addressToBeSearched)
	}
}
