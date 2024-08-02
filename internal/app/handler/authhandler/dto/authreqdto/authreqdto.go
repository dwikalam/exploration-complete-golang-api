package authreqdto

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/dwikalam/ecommerce-service/internal/app/type/wrappertype"
)

type RegisterUser struct {
	FullName string
	Email    string
	Password string
}

func (p *RegisterUser) Valid(ctx context.Context) wrappertype.ProblemsMap {
	var (
		problems wrappertype.ProblemsMap = make(map[string]string)

		nameProblem     <-chan error
		emailProblem    <-chan error
		passwordProblem <-chan error
	)

	if err := p.validatePayloadStructure(); err != nil {
		problems["payload"] = err.Error()

		return problems
	}

	nameProblem = p.validateFullName(ctx)
	emailProblem = p.validateEmail(ctx)
	passwordProblem = p.validatePassword(ctx)

	if err := <-nameProblem; err != nil {
		problems["name"] = err.Error()
	}

	if err := <-emailProblem; err != nil {
		problems["email"] = err.Error()
	}

	if err := <-passwordProblem; err != nil {
		problems["password"] = err.Error()
	}

	return problems
}

func (p *RegisterUser) validatePayloadStructure() error {
	if p.FullName == "" || p.Email == "" || p.Password == "" {
		return errors.New("payload structure not valid")
	}

	return nil
}

func (p *RegisterUser) validateFullName(ctx context.Context) <-chan error {
	var (
		ch = make(chan error)

		validate = func() error {
			const (
				nameRegexPattern string = `^[a-zA-Z\s'-]+$`
			)

			var (
				validName = strings.TrimSpace(p.FullName)
			)

			if len(validName) != len(p.FullName) {
				return errors.New("fullName contains leading and trailing white spaces")
			}

			if len(validName) < 3 {
				return errors.New("fullName length less than 3")
			}

			if len(validName) > 50 {
				return errors.New("fullName length greater than 50")
			}

			if ok := regexp.MustCompile(nameRegexPattern).MatchString(validName); !ok {
				return errors.New("fullName can only contain letter, spaces, hyphens, and apostrophes")
			}

			return nil
		}
	)

	go func() {
		defer close(ch)

		select {
		case <-ctx.Done():
			return
		case ch <- validate():
		}
	}()

	return ch
}

func (p *RegisterUser) validateEmail(ctx context.Context) <-chan error {
	var (
		ch = make(chan error)

		validate = func() error {
			const (
				emailRegexPattern string = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
			)

			if ok := regexp.MustCompile(emailRegexPattern).MatchString(p.Email); !ok {
				return errors.New("email structure not valid")
			}

			return nil
		}
	)

	go func() {
		defer close(ch)

		select {
		case <-ctx.Done():
			return
		case ch <- validate():
		}
	}()

	return ch
}

func (p *RegisterUser) validatePassword(ctx context.Context) <-chan error {
	var (
		ch = make(chan error)

		validate = func() error {
			const (
				minLength = 8
				maxLength = 50

				errmsg string = "password must contain minimal length of 8, maximal length of 50, and uppercase, lowercase, number, and special characters"
			)

			var (
				passwordLength int = len(p.Password)

				hasUpper   bool = false
				hasLower   bool = false
				hasNumber  bool = false
				hasSpecial bool = false
			)

			if passwordLength < minLength || passwordLength > maxLength {
				return errors.New(errmsg)
			}

			for _, char := range p.Password {
				switch {
				case unicode.IsUpper(char):
					hasUpper = true
				case unicode.IsLower(char):
					hasLower = true
				case unicode.IsNumber(char):
					hasNumber = true
				case unicode.IsPunct(char) || unicode.IsSymbol(char):
					hasSpecial = true
				}
			}

			if !(hasUpper && hasLower && hasNumber && hasSpecial) {
				return errors.New(errmsg)
			}

			return nil
		}
	)

	go func() {
		defer close(ch)

		select {
		case <-ctx.Done():
			return
		case ch <- validate():
		}
	}()

	return ch
}
