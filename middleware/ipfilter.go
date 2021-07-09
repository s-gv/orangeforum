package middleware

import (
	"net"
	"net/http"

	"github.com/golang/glog"
	"github.com/s-gv/orangeforum/models"
)

//TODO: Duplicate definition of domain id key, need to finalize a common place if needed
type contextKey string

const (
	ctxDomainID = contextKey("domain_id")
)

func IpFilter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainID := r.Context().Value(ctxDomainID).(int)

		// go-chi's middleware/realip.go already parses for RealIp and xForwardedFor, further sets RemoteAddr to xForwardedFor ip if available.
		ipAddress, _, splitHostPortError := net.SplitHostPort(r.RemoteAddr)
		if splitHostPortError != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		parsedIp := net.ParseIP(ipAddress)
		if parsedIp == nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		isIpBanned, ipFilterCheckError := checkIfIpAddressIsBanned(domainID, parsedIp.String())

		if ipFilterCheckError != nil {
			glog.Errorf("Error while searching for banned ip : %s", ipFilterCheckError.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if isIpBanned {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func checkIfIpAddressIsBanned(domainId int, ipAddress string) (bool, error) {
	ipv4AddressTrieRoot := models.BannedIpv4AddressTriesPerDomain[domainId]

	if ipv4AddressTrieRoot == nil {
		return false, nil
	}

	return ipv4AddressTrieRoot.SearchIpv4AddressInTrie(ipAddress)
}
