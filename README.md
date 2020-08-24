# go trading algo
### main.go
* This is a simple algorithm that incorporates daily rebalancing to limit volatility and increase returns 
    * Set up so you can edit your holdings and weights however you like
* Environment Variable Setup
    * There are several ways to set and get environment variables with Go
    * I chose to use the GoDotEnv Package
        1. Install the GoDotEnv Package - $ go get github.com/joho/godotenv
        2. Create the .env file within the same directory as your algorithm
* Download Alpaca Packages
    * $ go get github.com/alpacahq/alpaca-trade-api-go/common
    * $ go get github.com/alpacahq/alpaca-trade-api-go/polygon
    * $ go get github.com/alpacahq/alpaca-trade-api-go/stream
    * $ go get github.com/alpacahq/alpaca-trade-api-go/alpaca
* Download Decimal Package
    * $ go get github.com/shopspring/decimal