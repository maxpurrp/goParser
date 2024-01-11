package vars

import (
	"crypto/tls"
	"net/http"
	"sync"
)

var tr *http.Transport = &http.Transport{
	TLSClientConfig: &tls.Config{
		MaxVersion: tls.VersionTLS13,
	},
}
var Client *http.Client = &http.Client{
	Transport: tr,
}
var Pages int = 10
var MainLink string = "https://platesmania.com/"
var WgMain sync.WaitGroup
var Wg sync.WaitGroup
var MaxCountGour int = 2
