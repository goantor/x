package x

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

var (
	v *Validator
)

type IValidator interface {
	Valid(val interface{}) Errors
	Var(val interface{}, rule string) error
	AddRules(rules Rules)
	AddRule(name string, fn validator.Func) error
}

type ValidErrors []ValidError

type ValidError struct {
	Field   string `json:"field"`
	Content string `json:"content"`
}

func NewError(field string, content string) ValidError {
	return ValidError{Field: field, Content: content}
}

type Rules map[string]validator.Func

func NewValidator(entity *validator.Validate, rules Rules, tag string) *Validator {
	if v == nil {
		v = &Validator{
			entity:  entity,
			rules:   rules,
			tagName: tag,
		}
		v.initialize()
	}

	return v
}

type Validator struct {
	entity  *validator.Validate
	rules   Rules
	tagName string
}

func (v *Validator) Valid(val interface{}) (errs ValidErrors) {
	err := v.entity.Struct(val)
	if err == nil {
		return
	}

	return v.errors(err)
}

func (v *Validator) Var(value interface{}, rule string) error {
	return v.entity.Var(value, rule)
}

func (v *Validator) errors(err error) (errs ValidErrors) {
	switch err.(type) {
	case *validator.InvalidValidationError:
		errs = ValidErrors{
			ValidError{
				Field:   "all",
				Content: "参数验证错误",
			},
		}
		break
	case validator.ValidationErrors:
		validErrors := err.(validator.ValidationErrors)
		errs = make(ValidErrors, len(validErrors))
		for k, e := range validErrors {
			errs[k] = ValidError{
				Field:   e.StructField(),
				Content: e.Field(),
			}
		}
		break
	}

	return
}

// initialize 初始化数据
func (v *Validator) initialize() *Validator {
	v.bindRules()

	if v.tagName != "" {
		v.entity.RegisterTagNameFunc(validateTagFunc(v.tagName))
	}

	return v
}

// bindRules 绑定规则
func (v *Validator) bindRules() {
	v.AddRules(v.rules)
}

// AddRules 添加规则列表
func (v *Validator) AddRules(rules Rules) {
	if len(rules) == 0 {
		return
	}

	for name, fn := range rules {
		_ = v.AddRule(name, fn)
	}
}

// AddRule 添加规则
func (v *Validator) AddRule(name string, fn validator.Func) error {
	return v.entity.RegisterValidation(name, fn)
}

// validateTagFunc 获取tagName方法
func validateTagFunc(tag string) validator.TagNameFunc {
	return func(fl reflect.StructField) string {
		if label := fl.Tag.Get(tag); label != "" {
			return label
		}

		return fl.Name
	}
}
