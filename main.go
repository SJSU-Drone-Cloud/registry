//main.go
import (
	"os"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	err := gotdotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file!")
	}
	return os.Getenv(key)
}

func main() {
	a := App{}

	a.initialize(
		goDotEnvVariable("APP_DB_USERNAME"),
		goDotEnvVariable("APP_DB_PASSWORD")
		goDotEnvVariable("APP_DB_NAME")
	)
	a.run(":8010")
}