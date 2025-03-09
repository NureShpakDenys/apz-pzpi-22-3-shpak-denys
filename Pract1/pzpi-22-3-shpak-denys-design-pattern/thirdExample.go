package main

import (
	"context"
	"database/sql"
	"fmt"

	// Драйвер для MySQL
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseProvider визначає загальний інтерфейс для всіх типів баз даних.
// Реалізації цього інтерфейсу можуть представляти різні технології баз даних.
type DatabaseProvider interface {
	// Connect встановлює з'єднання з базою даних.
	// Повертає помилку, якщо з'єднання не вдалося встановити.
	Connect() error

	// ExecuteQuery виконує запит до бази даних.
	// Приймає рядок запиту, специфічний для конкретного типу бази даних.
	// Повертає результат у вигляді рядка та помилку у випадку невдачі.
	ExecuteQuery(query string) (string, error)
}

// SQLDatabase представляє реалізацію провайдера для SQL-баз даних.
// Використовує стандартний пакет database/sql для взаємодії з SQL базами даних.
type SQLDatabase struct {
	// connectionString рядок підключення до SQL бази даних
	connectionString string
	// db об'єкт підключення до бази даних
	db *sql.DB
}

// Connect встановлює з'єднання з SQL базою даних.
// Використовує драйвер MySQL для підключення.
// Повертає помилку, якщо з'єднання не вдалося встановити або перевірити.
func (s *SQLDatabase) Connect() error {
	var err error
	s.db, err = sql.Open("mysql", s.connectionString)
	if err != nil {
		return fmt.Errorf("could not connect to SQL database: %v", err)
	}
	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("could not ping SQL database: %v", err)
	}
	return nil
}

// ExecuteQuery виконує SQL-запит до бази даних.
// Приймає SQL-запит у вигляді рядка.
// Повертає результат першого стовпця першого рядка у вигляді рядка.
// У випадку помилки виконання запиту або читання даних повертає помилку.
func (s *SQLDatabase) ExecuteQuery(query string) (string, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return "", fmt.Errorf("could not execute SQL query: %v", err)
	}
	defer rows.Close()
	var result string
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			return "", fmt.Errorf("could not scan row: %v", err)
		}
	}
	return result, nil
}

// NoSQLDatabase представляє реалізацію провайдера для NoSQL-баз даних.
// Використовує MongoDB в якості NoSQL рішення.
type NoSQLDatabase struct {
	// connectionURL URL для підключення до MongoDB
	connectionURL string
	// client клієнт MongoDB
	client *mongo.Client
}

// Connect встановлює з'єднання з MongoDB.
// Використовує офіційний драйвер MongoDB для Go.
// Повертає помилку, якщо з'єднання не вдалося встановити.
func (n *NoSQLDatabase) Connect() error {
	clientOptions := options.Client().ApplyURI(n.connectionURL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("could not connect to NoSQL database: %v", err)
	}
	n.client = client
	return nil
}

// ExecuteQuery виконує пошук документа в MongoDB за вказаним значенням.
// Для MongoDB, "запит" є просто значенням поля "user", за яким шукається документ.
// Повертає знайдений документ у вигляді рядка або помилку, якщо документ не знайдено.
func (n *NoSQLDatabase) ExecuteQuery(query string) (string, error) {
	collection := n.client.Database("test").Collection("users")
	var result bson.M
	err := collection.FindOne(context.Background(), bson.M{"user": query}).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("could not execute NoSQL query: %v", err)
	}
	return fmt.Sprintf("%v", result), nil
}

// DatabaseFacade реалізує патерн Фасад для роботи з різними базами даних.
// Надає єдиний спрощений інтерфейс для взаємодії з різними типами баз даних.
type DatabaseFacade struct {
	// providers карта доступних провайдерів баз даних за їх типами
	providers map[string]DatabaseProvider
}

// NewDatabaseFacade створює новий екземпляр DatabaseFacade з передналаштованими провайдерами.
// Ініціалізує SQL та NoSQL провайдери з дефолтними параметрами підключення.
// Повертає готовий до використання екземпляр DatabaseFacade.
func NewDatabaseFacade() *DatabaseFacade {
	return &DatabaseFacade{
		providers: map[string]DatabaseProvider{
			"sql":   &SQLDatabase{connectionString: "user:password@tcp(localhost:3306)/dbname"},
			"nosql": &NoSQLDatabase{connectionURL: "mongodb://localhost:27017"},
		},
	}
}

// ExecuteDatabaseOperation виконує операцію в базі даних вказаного типу.
// Є основним методом фасаду для роботи з базами даних.
// Параметри:
//   - dbType: тип бази даних ("sql" або "nosql")
//   - query: запит до бази даних, специфічний для обраного типу
//
// Повертає:
//   - результат запиту у вигляді рядка у випадку успіху
//   - помилку, якщо тип бази даних не підтримується або виникла помилка при виконанні запиту
func (d *DatabaseFacade) ExecuteDatabaseOperation(dbType string, query string) (string, error) {
	provider, exists := d.providers[dbType]
	if !exists {
		return "", fmt.Errorf("unsupported database type: %s", dbType)
	}
	if err := provider.Connect(); err != nil {
		return "", err
	}
	return provider.ExecuteQuery(query)
}

// AddProvider додає новий провайдер бази даних до фасаду.
// Дозволяє розширювати систему новими типами баз даних під час виконання.
// Параметри:
//   - name: унікальний ідентифікатор типу бази даних
//   - provider: реалізація інтерфейсу DatabaseProvider
//
// Повертає:
//   - помилку, якщо провайдер з таким іменем вже існує
func (d *DatabaseFacade) AddProvider(name string, provider DatabaseProvider) error {
	if _, exists := d.providers[name]; exists {
		return fmt.Errorf("database provider %s already exists", name)
	}
	d.providers[name] = provider
	return nil
}

// Close закриває всі відкриті з'єднання з базами даних.
// Рекомендується викликати цей метод перед завершенням роботи програми.
// Повертає помилку, якщо виникла проблема при закритті будь-якого з з'єднань.
func (d *DatabaseFacade) Close() error {
	for name, provider := range d.providers {
		switch p := provider.(type) {
		case *SQLDatabase:
			if p.db != nil {
				if err := p.db.Close(); err != nil {
					return fmt.Errorf("error closing SQL database %s: %v", name, err)
				}
			}
		case *NoSQLDatabase:
			if p.client != nil {
				if err := p.client.Disconnect(context.Background()); err != nil {
					return fmt.Errorf("error closing NoSQL database %s: %v", name, err)
				}
			}
		}
	}
	return nil
}

// main є точкою входу в програму, що демонструє використання фасаду баз даних.
func main() {
	facade := NewDatabaseFacade()
	defer facade.Close() // Закриваємо з'єднання після завершення

	sqlResult, err := facade.ExecuteDatabaseOperation("sql", "SELECT name FROM users LIMIT 1")
	if err != nil {
		fmt.Println("Error executing SQL query:", err)
	} else {
		fmt.Println("SQL result:", sqlResult)
	}

	nosqlResult, err := facade.ExecuteDatabaseOperation("nosql", "john")
	if err != nil {
		fmt.Println("Error executing NoSQL query:", err)
	} else {
		fmt.Println("NoSQL result:", nosqlResult)
	}
}
