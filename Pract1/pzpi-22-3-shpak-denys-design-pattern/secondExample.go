package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// PaymentProvider визначає інтерфейс для всіх платіжних провайдерів.
// Кожен провайдер повинен реалізувати методи Connect та Process.
type PaymentProvider interface {
	// Connect встановлює з'єднання з платіжним шлюзом.
	// Повертає помилку, якщо з'єднання не вдалося встановити.
	Connect() error

	// Process обробляє платіжну транзакцію з вказаними деталями та сумою.
	// Параметри:
	//   - details: карта з деталями платежу, специфічними для конкретного провайдера
	//   - amount: сума платежу
	// Повертає:
	//   - ідентифікатор транзакції у випадку успіху
	//   - помилку, якщо обробка не вдалася
	Process(details map[string]string, amount float64) (string, error)
}

// Payment є базовою структурою для всіх конкретних реалізацій платіжних провайдерів.
// Містить спільну функціональність для взаємодії з API платіжних шлюзів.
type Payment struct {
	// apiEndpoint базовий URL API платіжного шлюзу
	apiEndpoint string
	// apiKey ключ для автентифікації в API платіжного шлюзу
	apiKey string
}

// Connect реалізує метод інтерфейсу PaymentProvider.
// У базовій реалізації просто перевіряє наявність налаштувань.
// Конкретні провайдери можуть перевизначити цей метод для власної логіки з'єднання.
func (p *Payment) Connect() error {
	if p.apiEndpoint == "" || p.apiKey == "" {
		return errors.New("missing API endpoint or API key")
	}
	return nil
}

// request виконує HTTP-запит до платіжного шлюзу.
// Внутрішній метод, що використовується усіма платіжними провайдерами.
// Параметри:
//   - method: HTTP метод запиту (GET, POST, тощо)
//   - endpoint: шлях до API-ендпоінту відносно базового URL
//   - payload: дані, що будуть відправлені у форматі JSON
//
// Повертає:
//   - ідентифікатор транзакції у випадку успіху
//   - помилку, якщо запит не вдався
func (p *Payment) request(method, endpoint string, payload map[string]interface{}) (string, error) {
	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest(method, p.apiEndpoint+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.New("payment request failed")
	}
	return "transaction-id", nil
}

// StripePayment реалізує платіжний провайдер для Stripe.
// Вбудовує базову структуру Payment для використання спільної функціональності.
type StripePayment struct{ Payment }

// Process обробляє платіж через Stripe.
// Очікує token у details для ідентифікації платіжної картки.
// Повертає ідентифікатор транзакції або помилку.
func (s *StripePayment) Process(details map[string]string, amount float64) (string, error) {
	if token, ok := details["token"]; !ok || token == "" {
		return "", errors.New("missing token")
	}
	return s.request("POST", "/charge", map[string]interface{}{"token": details["token"], "amount": amount})
}

// PayPalPayment реалізує платіжний провайдер для PayPal.
// Вбудовує базову структуру Payment для використання спільної функціональності.
type PayPalPayment struct{ Payment }

// Process обробляє платіж через PayPal.
// Очікує account_id у details для ідентифікації рахунку PayPal.
// Повертає ідентифікатор транзакції або помилку.
func (p *PayPalPayment) Process(details map[string]string, amount float64) (string, error) {
	if accountID, ok := details["account_id"]; !ok || accountID == "" {
		return "", errors.New("missing account ID")
	}
	return p.request("POST", "/pay", map[string]interface{}{"accountID": details["account_id"], "amount": amount})
}

// ApplePayPayment реалізує платіжний провайдер для Apple Pay.
// Вбудовує базову структуру Payment для використання спільної функціональності.
type ApplePayPayment struct{ Payment }

// Process обробляє платіж через Apple Pay.
// Очікує device_data у details для ідентифікації пристрою та платіжної інформації.
// Повертає ідентифікатор транзакції або помилку.
func (a *ApplePayPayment) Process(details map[string]string, amount float64) (string, error) {
	if deviceData, ok := details["device_data"]; !ok || deviceData == "" {
		return "", errors.New("missing device data")
	}
	return a.request("POST", "/pay", map[string]interface{}{"deviceData": details["device_data"], "amount": amount})
}

// PaymentFacade реалізує патерн Фасад для системи обробки платежів.
// Спрощує взаємодію з різними платіжними провайдерами через єдиний інтерфейс.
type PaymentFacade struct {
	// providers карта доступних платіжних провайдерів за їх ідентифікаторами
	providers map[string]PaymentProvider
}

// NewPaymentFacade створює та налаштовує новий екземпляр PaymentFacade
// з передналаштованими платіжними провайдерами.
// Повертає готовий до використання екземпляр PaymentFacade.
func NewPaymentFacade() *PaymentFacade {
	return &PaymentFacade{
		providers: map[string]PaymentProvider{
			"stripe":   &StripePayment{Payment{"https://stripe.com/api", "stripe-key"}},
			"paypal":   &PayPalPayment{Payment{"https://paypal.com/api", "paypal-key"}},
			"applepay": &ApplePayPayment{Payment{"https://applepay.com/api", "applepay-key"}},
		},
	}
}

// ProcessPayment обробляє платіж через вказаний тип платіжного провайдера.
// Є основним методом фасаду для обробки платежів.
// Параметри:
//   - paymentType: ідентифікатор платіжного провайдера ("stripe", "paypal", "applepay")
//   - details: деталі платежу, специфічні для обраного провайдера
//   - amount: сума платежу
//
// Повертає:
//   - ідентифікатор транзакції у випадку успіху
//   - помилку, якщо обробка не вдалася або провайдер не підтримується
func (p *PaymentFacade) ProcessPayment(paymentType string, details map[string]string, amount float64) (string, error) {
	provider, exists := p.providers[paymentType]
	if !exists {
		return "", fmt.Errorf("unsupported payment type: %s", paymentType)
	}
	if err := provider.Connect(); err != nil {
		return "", err
	}
	return provider.Process(details, amount)
}

// AddProvider додає новий платіжний провайдер до фасаду.
// Дозволяє розширювати систему новими платіжними провайдерами під час виконання.
// Параметри:
//   - name: унікальний ідентифікатор провайдера
//   - provider: реалізація інтерфейсу PaymentProvider
//
// Повертає:
//   - помилку, якщо провайдер з таким іменем вже існує
func (p *PaymentFacade) AddProvider(name string, provider PaymentProvider) error {
	if _, exists := p.providers[name]; exists {
		return fmt.Errorf("provider %s already exists", name)
	}
	p.providers[name] = provider
	return nil
}
