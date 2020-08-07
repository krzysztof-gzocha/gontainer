package example2

type Wallet struct {
	Value uint
}

func NewWallet(value uint) *Wallet {
	return &Wallet{Value: value}
}

type Person struct {
	Name   string
	Wallet *Wallet
}

func NewPerson(name string, wallet *Wallet) *Person {
	return &Person{Name: name, Wallet: wallet}
}
