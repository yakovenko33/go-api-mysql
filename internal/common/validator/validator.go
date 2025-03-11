package validator

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
)

type Validator struct {
	validatorInstance *validator.Validate
	translator        *ut.Translator
}

func NewValidator(lang string) *Validator {
	validatorInstance := validator.New()

	return &Validator{
		validatorInstance: validatorInstance,
		translator:        initTranslator(validatorInstance, lang),
	}
}

func initTranslator(validatorInstance *validator.Validate, lang string) *ut.Translator {
	enLocale := en.New()
	ruLocale := ru.New()

	uni := ut.New(enLocale, ruLocale)

	translator, _ := uni.GetTranslator(lang) // "en" "ru"
	if lang == "ru" {
		ru_translations.RegisterDefaultTranslations(validatorInstance, translator)
	} else {
		en_translations.RegisterDefaultTranslations(validatorInstance, translator)
	}

	return &translator
}

func (m *Validator) RegisterValidation(tag string, rule validator.Func, translate string) {
	m.validatorInstance.RegisterValidation(tag, rule)

	m.RegisterTranslation(tag, translate)
}

func (m *Validator) RegisterTranslation(tag string, translate string) {
	m.validatorInstance.RegisterTranslation(tag, *m.translator, func(ut ut.Translator) error {
		return ut.Add(tag, translate, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}

func (m *Validator) Validate(dto interface{}) map[string]string {
	err := m.validatorInstance.Struct(dto)
	errors := make(map[string]string)

	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range validationErrors {
				errors[e.Field()] = e.Translate(*m.translator)
			}
		}
	}
	return errors
}

// Custom validation function for email domain
// func emailDomainValidator(fl validator.FieldLevel) bool {
// 	email := fl.Field().String()
// 	if !strings.Contains(email, "@example.com") {
// 		return false
// 	}
// 	return true
// }
