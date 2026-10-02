package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/mtkmkv/todo-app/internal/core/errors"
)

// Компилируем регулярное выражение один раз при старте программы для экономии ресурсов
var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

// User представляет основную доменную модель пользователя
type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

// NewUser создает экземпляр пользователя с заданными параметрами
func NewUser(id int, version int, fullname string, phonenumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullname,
		PhoneNumber: phonenumber,
	}
}

// NewUserUninitialized создает нового пользователя без присвоенного ID и версии (перед сохранением в БД)
func NewUserUninitialized(fullname string, phonenumber *string) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullname,
		phonenumber,
	)
}

// Validate проверяет корректность полей пользователя по бизнес-правилам
func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf("invalid `FullName` len: %d: %w", fullNameLength, core_errors.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneNumberLength := len([]rune(*u.PhoneNumber))
		if phoneNumberLength < 10 || phoneNumberLength > 15 {
			return fmt.Errorf("invalid `PhoneNumber` len: %d: %w", phoneNumberLength, core_errors.ErrInvalidArgument)
		}

		if !phoneRegex.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("invalid `PhoneNumber` format: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

// ApplyPatch применяет частичные изменения к пользователю с валидацией промежуточного состояния
func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set && patch.FullName.Value != nil {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}

// UserPatch описывает набор полей для частичного обновления пользователя
type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

// Validate проверяет допустимость значений в переданном патче
func (p *UserPatch) Validate() error {
	/*if !p.FullName.Set && !p.PhoneNumber.Set {
		return fmt.Errorf("patch is empty: at least one field must be provided: %w", core_errors.ErrInvalidArgument)
	}*/

	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("'FullName' can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}