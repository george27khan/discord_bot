package translator

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	opts option.ClientOption // opts global options for translator
	once sync.Once
)

// initOptions once init options for translator
func initOptions() {
	once.Do(func() {
		googleKey := os.Getenv("GOOGLE_KEY")
		opts = option.WithAPIKey(googleKey)
	})
}

// GetSupportedLang return string of support languages for translator
func GetSupportedLang(ctx context.Context) string {
	var res string
	initOptions()
	client, err := translate.NewClient(ctx, opts)
	if err != nil {
		log.Println("ERROR. Google translator init.", err)
		return ""
	}
	defer client.Close()
	if langs, err := client.SupportedLanguages(ctx, language.English); err != nil {
		return ""
	} else {
		for _, lang := range langs {
			res += lang.Name + ": " + lang.Tag.String() + "\n"
		}
	}
	return res
}

// getLangTagByName return lang tag by name
func getLangTagByName(ctx context.Context, langName string) (tag language.Tag) {
	initOptions()
	client, err := translate.NewClient(ctx, opts)
	if err != nil {
		log.Println("ERROR. Google translator init.", err)
		return language.Tag{}
	}
	defer client.Close()
	if langs, err := client.SupportedLanguages(ctx, language.English); err != nil {
		return language.Tag{}
	} else {
		for _, lang := range langs {
			if strings.ToLower(lang.Name) == strings.ToLower(langName) {
				tag = lang.Tag
			}
		}
	}
	return
}

// IsSupportedLang check lang name in google supported languages list
func IsSupportedLang(ctx context.Context, langName string) bool {
	initOptions()
	client, err := translate.NewClient(ctx, opts)
	if err != nil {
		log.Println("ERROR. Google translator init.", err)
		return false
	}
	defer client.Close()
	if langs, err := client.SupportedLanguages(ctx, language.English); err != nil {
		return false
	} else {
		for _, lang := range langs {
			if strings.ToLower(lang.Name) == strings.ToLower(langName) {
				return true
			}
		}
	}
	return false
}

// Translate translates text into the selected language
func Translate(ctx context.Context, text string, langTo string) (string, error) {
	var langTag language.Tag
	initOptions()
	trnsl, err := translate.NewClient(ctx, opts)
	if err != nil {
		log.Println("ERROR. Google translator init.", err)
		return "", fmt.Errorf("google translator init error")
	}
	defer trnsl.Close()

	if val := getLangTagByName(ctx, langTo); val.String() != "und" { //try to get tag by language long name
		langTag = val
	} else {
		langTags, _, err := language.ParseAcceptLanguage(langTo) //parse input language tag
		if err != nil {
			log.Println("ERROR. Google translator parse language.", err)
			return "", fmt.Errorf("target language error")
		}
		langTag = langTags[0]
	}

	if t, err := trnsl.Translate(ctx, []string{text}, langTag, &translate.Options{}); err != nil {
		log.Println("ERROR. Translate.", err)
		return "", fmt.Errorf("translate error")
	} else {
		return t[0].Text, nil
	}
}
