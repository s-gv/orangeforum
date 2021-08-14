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
		t.Fatalf("Failed to create trie model for banned ipv4 addresses %s", err.Error())
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId3]
	if trieToBeSearched == nil {
		t.Fatalf("Failed to create trie for the specified domain %d", domainId3)
	}

	addressToBeSearched := "45.56.78.145"
	found, err := trieToBeSearched.SearchIpv4AddressInTrie(addressToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressToBeSearched, err.Error())
	}

	if found == false {
		t.Fatalf("Ipv4 addresss %s Not found", addressToBeSearched)
	}

	t.Logf("Ipv4 address %s found", addressToBeSearched)
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
		t.Fatalf("Failed to create trie model for banned ipv4 addresses %s", err.Error())
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId2]
	if trieToBeSearched == nil {
		t.Fatalf("Failed to create trie for the specified domain %d", domainId2)
	}

	addressToBeSearched := "45.56.78.145"
	found, err := trieToBeSearched.SearchIpv4AddressInTrie(addressToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressToBeSearched, err.Error())
	}

	if found == true {
		t.Fatalf("Ipv4 address %s found, while it was not expected", addressToBeSearched)
	}

	t.Logf("Ipv4 addresss %s Not found as expected", addressToBeSearched)
}

func Test_Trie_Based_IpAddress_Search_For_EmptyIpAddress(t *testing.T) {

	ipList1 := []string{"192", ""}
	domainId1 := 1
	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err == nil {
		t.Fatalf("should not create trie model for banned ipv4 addresses: %s", err.Error())
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched != nil {
		t.Fatalf("should not create trie for the specified domain %d ", domainId1)
	}

	t.Logf("Successful in not creating trie model for the invalid ip address list %s", ipList1)
}

func Test_Trie_Based_IpAddress_Search_For_Invalid_Address_Octet(t *testing.T) {

	ipList1 := []string{"192.555.0.1"}
	domainId1 := 1
	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err == nil {
		t.Fatalf("Should not create trie model for banned ipv4 addresses : %s", err.Error())
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched != nil {
		t.Fatalf("Should not create trie for the specified domain %d ", domainId1)
	}

	t.Logf("Successful in not creating trie model for the invalid ip address list %s", ipList1)
}

func Tes_Trie_Based_IpAddress_Search_For_InvalidInput(t *testing.T) {

	ipList1 := []string{"axyz"}
	domainId1 := 1
	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err == nil {
		t.Fatalf("Should not create trie model for banned ipv4 addresses : %s", err.Error())
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched != nil {
		t.Fatalf("Should not create trie for the specified domain %d ", domainId1)
	}

	t.Logf("Successful in not creating trie model for the invalid ip address list %s", ipList1)
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
	if err == nil {
		t.Fatalf("Should not create trie model for banned ipv4 addresses : %s", err.Error())
	}

	invalidTrieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if invalidTrieToBeSearched != nil {
		t.Fatalf("Should not create trie for the specified domain %d ", domainId1)
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId2]
	if trieToBeSearched != nil {
		t.Fatalf("Should not create trie for the specified domain %d ", domainId2)
	}

	t.Logf("Successful in not creating trie model for the given invalid ip address lists %s, %s", ipList1, ipList2)
}

func Test_Trie_Based_IpAddress_Search_Using_Prefix(t *testing.T) {

	ipList1 := []string{"192.168.1.0", "19.23.5.67"}
	domainId1 := 1

	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Fatalf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Fatalf("Failed to create trie for the specified domain %d", domainId1)
	}

	addressPrefixToBeSearched := "192.168"
	found, err := trieToBeSearched.SearchIpv4AddressPrefixInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == false {
		t.Fatalf("Ipv4 addresss prefix %s Not found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 address prefix %s found", addressPrefixToBeSearched)
}

func Test_Trie_Based_IpAddress_Search_Using_Non_Existent_Prefix(t *testing.T) {

	ipList1 := []string{"192.168.1.0", "19.23.5.67"}
	domainId1 := 1

	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Fatalf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Fatalf("Failed to create trie for the specified domain %d", domainId1)
	}

	addressPrefixToBeSearched := "192.169"
	found, err := trieToBeSearched.SearchIpv4AddressPrefixInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == true {
		t.Fatalf("Ipv4 address prefix %s found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 addresss prefix %s Not found as expected", addressPrefixToBeSearched)
}

func Test_Trie_Based_IpAddress_Search_Using_WildCard(t *testing.T) {

	ipList1 := []string{"192.168.1.122", "19.23.*"}
	domainId1 := 1

	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Fatalf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Fatalf("Failed to create trie for the specified domain %d", domainId1)
	}

	addressPrefixToBeSearched := "19.23.*"
	found, err := trieToBeSearched.SearchIpv4AddressPrefixInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == false {
		t.Fatalf("Ipv4 address prefix %s not found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 addresss prefix %s found as expected", addressPrefixToBeSearched)
}

func Test_Trie_Based_IpAddress_Search_Using_WildCard_For_Non_Existent_Address(t *testing.T) {

	ipList1 := []string{"192.168.1.122", "19.23.*"}
	domainId1 := 1

	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Fatalf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Fatalf("Failed to create trie for the specified domain %d", domainId1)
	}

	addressPrefixToBeSearched := "19.24.*"
	found, err := trieToBeSearched.SearchIpv4AddressPrefixInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == true {
		t.Fatalf("Ipv4 address prefix %s found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 addresss prefix %s not found as expected", addressPrefixToBeSearched)
}

func Test_Trie_Based_IpAddress_Search_Using_Invalid_WildCard_Address(t *testing.T) {

	ipList1 := []string{"*", "19.23.168.12"}
	domainId1 := 1

	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err == nil {
		t.Fatalf("Should not create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched != nil {
		t.Fatalf("Should not create trie for the specified domain %d", domainId1)
	}

	t.Logf("Successful in not creating trie model for the invalid ip address list %s", ipList1)
}

func Test_Trie_Based_IpAddress_Search_Using_WildCard_With_Repetitions(t *testing.T) {

	ipList1 := []string{"192.*.*.*", "19.23.*.*"}
	domainId1 := 1

	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Fatalf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Fatalf("Failed to create trie for the specified domain %d", domainId1)
	}

	addressPrefixToBeSearched := "19.23.*"
	found, err := trieToBeSearched.SearchIpv4AddressPrefixInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == false {
		t.Fatalf("Ipv4 address prefix %s not found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 addresss prefix %s  found as expected", addressPrefixToBeSearched)

	addressPrefixToBeSearched = "192"
	found, err = trieToBeSearched.SearchIpv4AddressPrefixInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == false {
		t.Fatalf("Ipv4 address prefix %s not found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 addresss prefix %s  found as expected", addressPrefixToBeSearched)
}

func Test_Trie_Based_IpAddress_Search_Using_WildCard_Without_Using_Prefix_Search_Method(t *testing.T) {

	ipList1 := []string{"192.168.1.122", "19.23.*", "100.*", "192"}
	domainId1 := 1

	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Fatalf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Fatalf("Failed to create trie for the specified domain %d", domainId1)
	}

	addressPrefixToBeSearched := "19.23.*"
	found, err := trieToBeSearched.SearchIpv4AddressInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == false {
		t.Fatalf("Ipv4 address prefix %s not found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 addresss prefix %s found as expected", addressPrefixToBeSearched)

	addressPrefixToBeSearched = "100.*"
	found, err = trieToBeSearched.SearchIpv4AddressInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == false {
		t.Fatalf("Ipv4 address prefix %s not found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 addresss prefix %s found as expected", addressPrefixToBeSearched)

	addressPrefixToBeSearched = "192"
	found, err = trieToBeSearched.SearchIpv4AddressInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == false {
		t.Fatalf("Ipv4 address prefix %s not found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 addresss prefix %s found as expected", addressPrefixToBeSearched)

	addressPrefixToBeSearched = "193"
	found, err = trieToBeSearched.SearchIpv4AddressInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == true {
		t.Fatalf("Ipv4 address prefix %s found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 addresss prefix %s not found as expected", addressPrefixToBeSearched)
}

func Test_Trie_Based_IpAddress_Search_Using_WildCard_That_Partially_Matches(t *testing.T) {

	ipList1 := []string{"192.168.1.122", "19.23.*", "10.*"}
	domainId1 := 1

	bannedIpsByDomain := map[int][]string{
		domainId1: ipList1,
	}

	bannedIpv4AddressTriesPerDomain, err := models.CreateIpv4BannedAddressTriePerDomain(bannedIpsByDomain)
	if err != nil {
		t.Fatalf("Failed to create trie model for banned ipv4 addresses")
	}

	trieToBeSearched := bannedIpv4AddressTriesPerDomain[domainId1]
	if trieToBeSearched == nil {
		t.Fatalf("Failed to create trie for the specified domain %d", domainId1)
	}

	addressPrefixToBeSearched := "192.168.*"
	found, err := trieToBeSearched.SearchIpv4AddressPrefixInTrie(addressPrefixToBeSearched)
	if err != nil {
		t.Fatalf("Error while searching for ip %s in the trie : %s", addressPrefixToBeSearched, err.Error())
	}

	if found == false {
		t.Fatalf("Ipv4 address prefix %s not found", addressPrefixToBeSearched)
	}

	t.Logf("Ipv4 addresss prefix %s found as expected", addressPrefixToBeSearched)
}
