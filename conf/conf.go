package conf

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config Struct
type Config struct {
	DBUrl     string
	AppSecret string
}

// GetConfig ...
func GetConfig(test bool) Config {

	if test {
		godotenv.Load(".test.env")
		log.Println("Load env for testing")
	} else {
		godotenv.Load()
		log.Println("Load default env")
	}

	db := os.Getenv("DB_URL")
	sc := os.Getenv("SECRET")

	if sc == "" {
		sc = "D3faultR4ndom53cretForJWT"
	}

	c := Config{
		DBUrl:     db,
		AppSecret: sc,
	}

	return c
}
