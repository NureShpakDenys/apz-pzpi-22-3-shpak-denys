package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
)

// User представляє собою структуру для зберігання контактної інформації користувача.
// Містить поля для різних каналів комунікації: Email, Phone та DeviceToken.
type User struct {
	// Email користувача для надсилання повідомлень електронною поштою
	Email string
	// Phone номер телефону користувача для SMS повідомлень
	Phone string
	// DeviceToken токен пристрою для push-повідомлень
	DeviceToken string
}

// EmailSender відповідає за надсилання електронних листів через SMTP.
type EmailSender struct {
	// smtpServer адреса SMTP сервера
	smtpServer string
	// smtpPort порт SMTP сервера
	smtpPort string
	// username ім'я користувача для автентифікації на SMTP сервері
	username string
	// password пароль для автентифікації на SMTP сервері
	password string
}

// Connect встановлює з'єднання з SMTP сервером з використанням TLS.
// Повертає клієнт SMTP або помилку, якщо з'єднання не вдалося встановити.
func (e *EmailSender) Connect() (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", e.smtpServer+":"+e.smtpPort, nil)
	if err != nil {
		return nil, err
	}
	client, err := smtp.NewClient(conn, e.smtpServer)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Authenticate виконує автентифікацію на SMTP сервері.
// Приймає клієнт SMTP і повертає помилку, якщо автентифікація не вдалася.
func (e *EmailSender) Authenticate(client *smtp.Client) error {
	auth := smtp.PlainAuth("", e.username, e.password, e.smtpServer)
	if err := client.Auth(auth); err != nil {
		return err
	}
	return nil
}

// Send надсилає електронний лист вказаному отримувачу.
// Приймає адресу отримувача, тему та тіло повідомлення.
// Повертає помилку, якщо відправлення не вдалося.
func (e *EmailSender) Send(to, subject, body string) error {
	client, err := e.Connect()
	if err != nil {
		return err
	}
	defer client.Quit()
	if err := e.Authenticate(client); err != nil {
		return err
	}
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	if err := client.Mail(e.username); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	_, err = writer.Write(msg)
	if err != nil {
		return err
	}
	return writer.Close()
}

// SMSSender відповідає за надсилання SMS повідомлень через API.
type SMSSender struct {
	// apiEndpoint URL API сервісу для відправки SMS
	apiEndpoint string
	// apiKey ключ для автентифікації в API сервісі
	apiKey string
}

// SendMessage надсилає SMS повідомлення на вказаний номер телефону.
// Приймає номер телефону та текст повідомлення.
// Повертає помилку, якщо відправлення не вдалося.
func (s *SMSSender) SendMessage(phoneNumber, message string) error {
	if phoneNumber == "" {
		return errors.New("phone number is not provided")
	}
	payload := map[string]string{
		"phone":   phoneNumber,
		"message": message,
		"apiKey":  s.apiKey,
	}
	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", s.apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("SMS API error: %s", string(body))
	}
	return nil
}

// PushNotifier відповідає за надсилання push-повідомлень через API.
type PushNotifier struct {
	// apiEndpoint URL API сервісу для відправки push-повідомлень
	apiEndpoint string
	// apiKey ключ для автентифікації в API сервісі
	apiKey string
}

// SendPush надсилає push-повідомлення на пристрій з вказаним токеном.
// Приймає токен пристрою, заголовок та текст повідомлення.
// Повертає помилку, якщо відправлення не вдалося.
func (p *PushNotifier) SendPush(deviceToken, title, message string) error {
	if deviceToken == "" {
		return errors.New("device token is not provided")
	}
	payload := map[string]string{
		"deviceToken": deviceToken,
		"title":       title,
		"message":     message,
		"apiKey":      p.apiKey,
	}
	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", p.apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Push API error: %s", string(body))
	}
	return nil
}

// NotificationFacade реалізує патерн Фасад для системи оповіщень,
// об'єднуючи різні канали комунікації в єдиний інтерфейс.
type NotificationFacade struct {
	// emailSender компонент для надсилання електронних листів
	emailSender *EmailSender
	// smsSender компонент для надсилання SMS
	smsSender *SMSSender
	// pushNotifier компонент для надсилання push-повідомлень
	pushNotifier *PushNotifier
}

// NewNotificationFacade створює та налаштовує новий екземпляр NotificationFacade
// з дефолтними параметрами для всіх типів оповіщень.
// Повертає готовий до використання екземпляр NotificationFacade.
func NewNotificationFacade() *NotificationFacade {
	return &NotificationFacade{
		emailSender:  &EmailSender{smtpServer: "smtp.example.com", smtpPort: "465", username: "user@example.com", password: "password"},
		smsSender:    &SMSSender{apiEndpoint: "https://api.example.com/sms", apiKey: "sms_api_key"},
		pushNotifier: &PushNotifier{apiEndpoint: "https://api.example.com/push", apiKey: "push_api_key"},
	}
}

// SendNotification надсилає повідомлення користувачу через всі доступні канали комунікації.
// Використовує наявні контактні дані користувача (email, телефон, токен пристрою).
// Повертає помилку, якщо виникла проблема з надсиланням хоча б через один канал.
func (n *NotificationFacade) SendNotification(user User, message string) error {
	if user.Email != "" {
		if err := n.emailSender.Send(user.Email, "Message", message); err != nil {
			return err
		}
	}
	if user.Phone != "" {
		if err := n.smsSender.SendMessage(user.Phone, message); err != nil {
			return err
		}
	}
	if user.DeviceToken != "" {
		if err := n.pushNotifier.SendPush(user.DeviceToken, "Message", message); err != nil {
			return err
		}
	}
	return nil
}

// main є точкою входу в програму, що демонструє використання системи оповіщень.
func main() {
	facade := NewNotificationFacade()
	user := User{Email: "test@example.com", Phone: "+380501234567", DeviceToken: "device_token_123"}
	err := facade.SendNotification(user, "This is test message!")
	if err != nil {
		fmt.Println("Error:", err)
	}
}
