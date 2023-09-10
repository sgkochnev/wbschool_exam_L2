package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
Фасад — это структурный паттерн проектирования, который предоставляет
простой (но урезанный) интерфейс к сложной системе классов,
библиотеке или фреймворку.

Фасад позволяет скрыть сложную систему классов от клиентов,
делая ее более простой для использования.

Применение:
- Когда вам нужно представить простой или урезанный интерфейс
к сложной подсистеме.
- Когда вы хотите разложить подсистему на отдельные слои.

Плюсы:
- Изолирует клиентов от компонентов сложной подсистемы.

Минусы:
- Фасад рискует стать божественным объектом, привязанным ко всем классам программы.
*/

type Product struct {
	Name  string
	Price float64
}

type Shop struct {
	Name     string
	account  *Account
	Products map[string]Product
}

// фасад
func (s *Shop) Sell(user *User, product string) error {
	fmt.Println("Shop: Запрос баланса покупателя")
	if err := user.Card.CheckBalance(); err != nil {
		return err
	}
	fmt.Printf("Shop: Проверка может ли поукупатель %s купить товар %s\n", user.Name, product)

	if p, ok := s.Products[product]; ok {
		if p.Price > user.Balance() {
			return fmt.Errorf("Недостаточно средств на карте %s", user.Name)
		}

		err := s.account.Bank.TransferToAccount(user.Card.Number, s.account.Number, p.Price)
		if err != nil {
			return err
		}
		fmt.Printf("Shop: Покупка товара %s\n", product)
	} else {

	}
	return nil
}

type Bank struct {
	Name     string
	Cards    map[string]Card
	Accounts map[string]Account
}

func (b *Bank) CheckBalance(cardNumber string) error {
	fmt.Printf("Bank: Получение баланса карты %s\n", cardNumber)
	if v, ok := b.Cards[cardNumber]; ok {
		if v.Balance <= 0 {
			return fmt.Errorf("Недостаточно средств на карте %s", cardNumber)
		}
		fmt.Printf("Bank: Баланс карты %s: %.2f\n", cardNumber, v.Balance)
	}
	return nil
}

func (b *Bank) TransferToAccount(cardNumber string, accountNumber string, amount float64) error {
	fmt.Printf("Bank: Перевод с карты %s на счет %s: %.2f\n", cardNumber, accountNumber, amount)
	c, ok1 := b.Cards[cardNumber]
	a, ok2 := b.Accounts[accountNumber]
	if !ok1 || !ok2 {
		return fmt.Errorf("Неверный номер карты или счета")
	}

	c.Balance -= amount
	a.Balance += amount
	return nil
}

type Account struct {
	Number  string
	Balance float64
	Bank    *Bank
}

type Card struct {
	Number  string
	Balance float64
	Bank    *Bank
}

func (c *Card) CheckBalance() error {
	fmt.Println("Card: Запрос в банк для проверки баланса карты")
	return c.Bank.CheckBalance(c.Number)
}

type User struct {
	Name string
	Card *Card
}

func (u *User) Balance() float64 {
	return u.Card.Balance
}

var (
	bank1 = Bank{
		Name:     "Bank",
		Cards:    make(map[string]Card),
		Accounts: make(map[string]Account),
	}

	card1 = Card{
		Number:  "Card-1",
		Balance: 1200,
		Bank:    &bank1,
	}

	card2 = Card{
		Number:  "Card-2",
		Balance: 50,
		Bank:    &bank1,
	}

	account1 = Account{
		Number:  "Account-1",
		Balance: 0,
		Bank:    &bank1,
	}

	product1 = Product{
		Name:  "Product-1",
		Price: 100,
	}

	product2 = Product{
		Name:  "Product-2",
		Price: 200,
	}

	shop1 = Shop{
		Name:    "Shop",
		account: &account1,
		Products: map[string]Product{
			product1.Name: product1,
			product2.Name: product2,
		},
	}

	user1 = User{
		Name: "User-1",
		Card: &card1,
	}

	user2 = User{
		Name: "User-2",
		Card: &card2,
	}
)

func main() {

	bank1.Cards[card1.Number] = card1
	bank1.Cards[card2.Number] = card2
	bank1.Accounts[account1.Number] = account1

	err := shop1.Sell(&user1, product1.Name)
	if err != nil {
		fmt.Println(err)
	}
	err = shop1.Sell(&user2, product2.Name)
	if err != nil {
		fmt.Println(err)
	}
}
