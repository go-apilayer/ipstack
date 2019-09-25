package ipstack

type LangType string

const (
	LangEnglishUS        LangType = "en"
	LangGerman           LangType = "de"
	LangSpanish          LangType = "es"
	LangFrench           LangType = "fr"
	LangJapanese         LangType = "ja"
	LangRussian          LangType = "ru"
	LangChinese          LangType = "zh"
	LangPortugueseBrazil LangType = "pt-br"
)

func (l LangType) String() string {
	return string(l)
}
