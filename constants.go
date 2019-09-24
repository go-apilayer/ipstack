package ipstack

type LangType string

const (
	EnglishUS        LangType = "en"
	German           LangType = "de"
	Spanish          LangType = "es"
	French           LangType = "fr"
	Japanese         LangType = "ja"
	Russian          LangType = "ru"
	Chinese          LangType = "zh"
	PortugueseBrazil LangType = "pt-br"
)

func (l LangType) String() string {
	return string(l)
}
