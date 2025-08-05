package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq" // <--- ВАЖНО: импорт драйвера
	"log"
	"net/http"
	"time"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
type Fields map[string]string
type Book1 struct {
	Id int `json:"id"`
	Fields
}

func main1() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=nlrk sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, author, title FROM ebook")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Author, &b.Title); err != nil {
			log.Fatal(err)
		}
		books = append(books, b)
	}
	println(len(books))
	chunks := chunkBooks(books, 1000) // батчи по 1000

	for i, chunk := range chunks {
		err := bulkUploadBooks(chunk)
		if err != nil {
			log.Printf("Ошибка в батче %d: %v", i, err)
			// можно retry здесь
		}
		log.Printf("Загружен батч %d из %d", i+1, len(chunks))
		time.Sleep(100 * time.Millisecond) // пауза между батчами (опционально)
	}

}
func chunkBooks(books []Book, size int) [][]Book {
	var chunks [][]Book
	for size < len(books) {
		books, chunks = books[size:], append(chunks, books[0:size:size])
	}
	chunks = append(chunks, books)
	return chunks
}

func bulkUploadBooks(books []Book) error {
	var buf bytes.Buffer

	for _, book := range books {
		// Meta строка
		meta := map[string]any{
			"index": map[string]any{
				"_index": "books_index",
				"_id":    book.ID, // Если не нужен фиксированный ID — можно не указывать
			},
		}
		metaJson, _ := json.Marshal(meta)
		buf.Write(metaJson)
		buf.WriteByte('\n')

		// Сам документ
		docJson, _ := json.Marshal(book)
		buf.Write(docJson)
		buf.WriteByte('\n')
	}

	// Отправка запроса
	resp, err := http.Post("http://localhost:9200/_bulk", "application/x-ndjson", &buf)
	if err != nil {
		return fmt.Errorf("bulk post error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bulk request failed: %s", resp.Status)
	}

	return nil
}
