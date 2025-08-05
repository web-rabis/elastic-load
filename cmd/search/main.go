package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

// Структура результата
type Hit struct {
	Source Book `json:"_source"`
}

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []Hit `json:"hits"`
	} `json:"hits"`
}

func main() {
	// Создание клиента
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Запрос поиска
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]string{
				"author": "Абай",
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Выполнение запроса
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("books_index"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// Парсинг результата
	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}

	var sr SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		log.Fatalf("Error parsing response: %s", err)
	}

	// Вывод результатов
	for _, hit := range sr.Hits.Hits {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n",
			hit.Source.ID, hit.Source.Title, hit.Source.Author)
	}
	println(sr.Hits.Total.Value)
}
