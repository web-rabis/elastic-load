package config

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

type DatabaseConfig struct {
	DSName  string `short:"n" long:"ds" env:"DATASTORE" description:"DataStore name (format: mongo/null)" required:"false" default:"postgres"`
	DSDB    string `short:"d" long:"ds-db" env:"DATASTORE_DB" description:"DataStore database name (format: acquiring)" required:"false" default:"nlrk"`
	DSURL   string `short:"u" long:"ds-url" env:"DATASTORE_URL" description:"DataStore URL (format: mongodb://localhost:27017)" required:"false" default:"postgres://postgres:postgres@localhost:5432/nlrk"`
	ESURL   string `long:"es-url" env:"ELASTICSEARCH_URL" description:"Elasticsearch URL" required:"false" default:"http://192.168.7.175:9200"`
	ESINDEX string `long:"es-index" env:"ELASTICSEARCH_INDEX" description:"Elasticsearch INDEX" required:"false" default:"books-index"`
}
type ServerConfig struct {
	ListenAddr string `short:"l" long:"listen" env:"LISTEN" description:"Listen Address (format: :8080|127.0.0.1:8080)" required:"false" default:":9191"`
	BasePath   string `long:"base-path" env:"BASE_PATH" description:"base path of the host" required:"false" default:"/reader"`
	FilesDir   string `long:"files-directory" env:"FILES_DIR" description:"Directory where all static files are located" required:"false" default:"/usr/share/reader"`
	CertFile   string `short:"c" long:"cert" env:"CERT_FILE" description:"Location of the SSL/TLS cert file" required:"false" default:""`
	KeyFile    string `short:"k" long:"key" env:"KEY_FILE" description:"Location of the SSL/TLS key file" required:"false" default:""`
}
type AuthConfig struct {
	JWTKey string `long:"jwt-key" env:"JWT_KEY" description:"JWT secret key" required:"false" default:"airbapay-secret"`
}

type GeneralConfig struct {
	Dbg       bool `long:"dbg" env:"DEBUG" description:"debug mode"`
	IsTesting bool `long:"testing" env:"APP_TESTING" description:"testing mode"`
}

func Parse(c interface{}) interface{} {
	p := flags.NewParser(c, flags.Default)

	if _, err := p.Parse(); err != nil {
		log.Println("[ERROR] Error while parsing config options:", err)
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	return c
}
