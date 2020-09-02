package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/my/repo/go/src/github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

type alpacaClientContainer struct {
	api          *alpaca.Client
	nasdaq       []alpaca.Asset
	shortable    []string
	nonShortable []string
	nonTradable  []string
}

var alpacaClient alpacaClientContainer

var wg sync.WaitGroup

func init() {
	os.Setenv(common.EnvApiKeyID, goDotEnvVariable("KEY"))
	os.Setenv(common.EnvApiSecretKey, goDotEnvVariable("SECRET"))

	// fmt.Printf("Running w/ credentials [%v %v]\n", common.Credentials().ID, common.Credentials().Secret)

	alpaca.SetBaseUrl(goDotEnvVariable("BASE_URL"))
	alpacaClient = alpacaClientContainer{
		alpaca.NewClient(common.Credentials()),
		[]alpaca.Asset{},
		[]string{},
		[]string{},
		[]string{},
	}
}

func runningtime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End:	", s, "took", endTime.Sub(startTime))
}

func main() {
	defer track(runningtime("main"))
	alpacaClient.setAssets()
	wg.Wait()
	fmt.Println("Nasdaq total assets:", len(alpacaClient.nasdaq))
	fmt.Println("Nasdaq shortable:", len(alpacaClient.shortable))
	fmt.Println("Nasdaq non-shortable:", len(alpacaClient.nonShortable))
	fmt.Println("Nasdaq non-tradable:", len(alpacaClient.nonShortable))
}

func (alp *alpacaClientContainer) setAssets() {
	defer wg.Done()
	// Get a list of all active assets.
	status := "active"
	assets, err := alp.api.ListAssets(&status)
	if err != nil {
		panic(err)
	}

	// Filter the assets down to just those on NASDAQ.
	for _, asset := range assets {
		if asset.Exchange == "NASDAQ" {
			alp.nasdaq = append(alp.nasdaq, asset)
		}
	}

	// Create sub slices so we can implement go-routines
	size := len(alp.nasdaq) / 20
	var j int
	for i := 0; i < len(alp.nasdaq); i += size {
		j += size
		if j > len(alp.nasdaq) {
			j = len(alp.nasdaq)
		}
		// do what do you want to with the sub-slice, here just printing the sub-slices
		wg.Add(1)
		go alp.getShortable(alp.nasdaq[i:j])
	}
}

func (alp *alpacaClientContainer) getShortable(subNasdaq []alpaca.Asset) {
	defer wg.Done()
	fmt.Println("subNasdaq size:", len(subNasdaq))
	for i := 0; i < len(subNasdaq); i++ {
		if subNasdaq[i].Tradable {
			if subNasdaq[i].Shortable {
				alp.shortable = append(alp.shortable, subNasdaq[i].Symbol)
			} else {
				alp.nonShortable = append(alp.nonShortable, subNasdaq[i].Symbol)
			}
		} else {
			alp.nonShortable = append(alp.nonShortable, subNasdaq[i].Symbol)
			alp.nonTradable = append(alp.nonTradable, subNasdaq[i].Symbol)
		}
	}
}
