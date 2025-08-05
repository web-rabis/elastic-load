package elk

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

type Manager struct {
	elkClient *elasticsearch.Client
}

func NewElkManager() (*Manager, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("ошибка клиента: %s", err)
	}

}
