package tests

import (
	"testing"

	"github.com/s-gv/orangeforum/models"
)

func Test_Trie_Based_IpAddress_Search(t *testing.T) {

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

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
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

func Test_Trie_Based_IpAddress_Search_For_Non_Existent_Ip_Address(t *testing.T) {

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

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Errorf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId2]
	if trieToBeSearched == nil {
		t.Errorf("Failed to create trie for the specified domain %d", domainId2)
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

func Test_Trie_Based_IpAddress_Search_For_EmptyIpAddress(t *testing.T) {

	ipList1 := []string{"192", ""}
	domainId1 := 1
	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Log("Failed to create trie model for banned ipv4 addresses as expected")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Logf("Failed to create trie for the specified domain %d as expected", domainId1)
	}
}

func Tes_Trie_Based_IpAddress_Search_For_InvalidInput(t *testing.T) {

	ipList1 := []string{"axyz"}
	domainId1 := 1
	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Log("Failed to create trie model for banned ipv4 addresses as expected")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Logf("Failed to create trie for the specified domain %d as expected", domainId1)
	}
}

func Test_Trie_Based_IpAddress_Search_For_Valid_And_InvalidInputs(t *testing.T) {

	ipList1 := []string{"axyz"}
	ipList2 := []string{"192.168.2.3", "45.56.78.145"}
	domainId1 := 1
	domainId2 := 2
	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
		domainId2: ipList2,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Log("Failed to create trie model for banned ipv4 addresses as expected")
	}

	invalidTrieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if invalidTrieToBeSearched == nil {
		t.Logf("Failed to create trie for the specified domain %d as expected", domainId1)
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId2]
	if trieToBeSearched == nil {
		t.Logf("Failed to create trie for the specified domain %d as expected", domainId2)
	}
}

func Test_Trie_Based_IpAddress_Search_Using_Prefix(t *testing.T) {

	ipList1 := []string{"192.168.1.0", "19.23.5.67"}
	domainId1 := 1

	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Errorf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Errorf("Failed to create trie for the specified domain %d", domainId1)
	}

	addressPrefixToBeSearched := "192.168"
	found, err := trieToBeSearched.SearchIpv4AddressPrefixInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Errorf("Error while searching for ip %s in the trie", addressPrefixToBeSearched)
		t.Error(err)
	}

	if found == false {
		t.Errorf("Ipv4 addresss prefix %s Not found", addressPrefixToBeSearched)
	}

	if found == true {
		t.Logf("Ipv4 address prefix %s found", addressPrefixToBeSearched)
	}
}

func Test_Trie_Based_IpAddress_Search_Using_Non_Existent_Prefix(t *testing.T) {

	ipList1 := []string{"192.168.1.0", "19.23.5.67"}
	domainId1 := 1

	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Errorf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Errorf("Failed to create trie for the specified domain %d", domainId1)
	}

	addressPrefixToBeSearched := "192.169"
	found, err := trieToBeSearched.SearchIpv4AddressPrefixInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Errorf("Error while searching for ip %s in the trie", addressPrefixToBeSearched)
		t.Error(err)
	}

	if found == false {
		t.Logf("Ipv4 addresss prefix %s Not found as expected", addressPrefixToBeSearched)
	}

	if found == true {
		t.Errorf("Ipv4 address prefix %s found", addressPrefixToBeSearched)
	}
}
