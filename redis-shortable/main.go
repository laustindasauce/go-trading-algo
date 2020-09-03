package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/gomodule/redigo/redis"
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
	nonTradeable []string
}

var alpacaClient alpacaClientContainer

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
	alpacaClient.getShortable()
	fmt.Println("Shortable Assets:", len(alpacaClient.shortable))
	fmt.Println("non-Shortable Assets:", len(alpacaClient.nonShortable))
	fmt.Println("non-Tradeable Assets:", len(alpacaClient.nonTradeable))
	if len(alpacaClient.shortable) > len(alpacaClient.nonShortable) {
		setRedis(true)
	} else {
		setRedis(false)
	}
}

func (alp *alpacaClientContainer) setAssets() {
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
}

func (alp *alpacaClientContainer) getShortable() {
	fmt.Println("Nasdaq size:", len(alp.nasdaq))
	for i := 0; i < len(alp.nasdaq); i++ {
		if alp.nasdaq[i].Tradable {
			if alp.nasdaq[i].Shortable {
				alp.shortable = append(alp.shortable, alp.nasdaq[i].Symbol)
			} else {
				alp.nonShortable = append(alp.nonShortable, alp.nasdaq[i].Symbol)
			}
		} else {
			alp.nonTradeable = append(alp.nonTradeable, alp.nasdaq[i].Symbol)
		}
	}
}

func setRedis(shortable bool) {
	secret := goDotEnvVariable("REDIS")

	client, err := redis.Dial("tcp", "10.10.10.1:6379")
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Do("AUTH", secret)
	if err != nil {
		log.Fatal(err)
	}
	if shortable {
		client.Do("INCR", "goShort")
	}
	client.Do("INCR", "totalGoShort")
	shortableDays, _ := redis.Int(client.Do("GET", "goShort"))
	totalDays, _ := redis.Int(client.Do("GET", "totalGoShort"))
	getPercent(shortableDays, totalDays)
}

func getPercent(shortableDays, totalDays int) {
	pct := (shortableDays / totalDays) * 100
	fmt.Print(pct, "% of the time the Nasdaq has been more than 50% shortable over the past ", totalDays, " days\n")
}
