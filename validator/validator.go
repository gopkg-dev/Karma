package validator

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/gopkg-dev/karma/errors"
)

var (
	Validator = validator.New()
	trans     ut.Translator

	ErrInvalidArgument = errors.BadRequest("INVALID_PARAMETER", "参数错误")
	ErrValidate        = errors.BadRequest("VALIDATE_ERROR", "参数校验错误")
)

func init() {
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("query")
		if name == "" {
			name = fld.Tag.Get("label")
		}
		return name
	})
	uni := ut.New(zh.New())
	trans, _ = uni.GetTranslator("zh")
	if err := zhTranslations.RegisterDefaultTranslations(Validator, trans); err != nil {
		panic(fmt.Sprintf("Error while registering english translations %+v", err))
	}
	Register(Validator, trans)
}

// Validate validates the input struct
func Validate(payload interface{}) error {
	err := Validator.Struct(payload)
	if err == nil {
		return nil
	}

	var invalidValidationError *validator.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		return ErrInvalidArgument.WithMessage(invalidValidationError.Error())
	}

	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		metadata := make(map[string]interface{}, 0)
		for _, fieldError := range errs {
			ns := fieldError.Field()
			metadata[ns] = fieldError.Translate(trans)
		}
		return ErrValidate.WithMetadata(metadata)
	}

	return nil
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}
		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		log.Printf("warning: error translating FieldError: %#v", fe)
		return fe.(error).Error()
	}
	return t
}
