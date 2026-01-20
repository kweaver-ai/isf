package enums

type AccountType string

func (a AccountType) String() string {
	return string(a)
}

const (
	AccountTypeUser AccountType = "user"
	AccountTypeApp  AccountType = "app"
)
