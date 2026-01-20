package langcmp

//go:generate mockgen -source=./interface.go -destination ./langcmpmock/lang.go -package langcmpmock
type LangInterface interface {
	SetSysDefLang(lang string)
	GetSysDefaultLang() Lang
}
