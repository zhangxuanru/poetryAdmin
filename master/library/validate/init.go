package validate

import (
	"errors"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	trzh "gopkg.in/go-playground/validator.v9/translations/zh"
)

var G_Validate *validator.Validate
var G_Translator ut.Translator

func InitValidate() (err error) {
	var found bool
	G_Validate = validator.New()
	cn := zh.New()
	translator := ut.New(cn, cn)
	G_Translator, found = translator.GetTranslator("zh")
	if found == false {
		return errors.New("not locale zh")
	}
	err = trzh.RegisterDefaultTranslations(G_Validate, G_Translator)
	return err
}
