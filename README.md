# Go Expert Desafio Tec Rate Limiter 

## Instructions

### Run project

Execute `make docker-up`

#### Test with api.http

* Open the api.http file
* Run the request

### Test with go tests

* Run `make tests` 

### Configuration

* You can configure many token in .env file
* Add a name, for example: `PREMIUM_TOKEN_LIMIT=10`
* And create a new var for this new env var in rate_limiter.go, for example: 
`tokenLimit["premium"], _ = strconv.Atoi(os.Getenv("PREMIUM_TOKEN_LIMIT"))`
* Run the app again



