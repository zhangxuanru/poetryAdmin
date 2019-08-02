package validate

import "gopkg.in/go-playground/validator.v9"

func Validate() *validator.Validate {
	return G_Validate
}

//验证结构体
func ValidateStruct(str interface{}) (err error) {
	if G_Validate == nil {
		if err = InitValidate(); err != nil {
			return err
		}
	}
	if err = Validate().Struct(str); err != nil {
		return err
	}
	return nil
}

//错误转换为验证错误类型
func ErrToValidationErrors(err error) validator.ValidationErrors {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		return validationErrors
	}
	return nil
}
