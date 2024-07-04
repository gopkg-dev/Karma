package validator

import (
	"net/mail"
	"strings"

	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"
)

func Register(v *validator.Validate, trans ut.Translator) {
	err := v.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
		return IsMobile(fl.Field().String())
	})
	if err != nil {
		panic(err)
	}
	err = v.RegisterTranslation("mobile", trans, registrationFunc("mobile", "必须是有效的手机号码", false), translateFunc)
	if err != nil {
		panic(err)
	}
}

// IsMobile is the validation function for validating if the current field's value is a valid e.164 formatted phone number.
func IsMobile(v string) bool {
	// Check the length of the phone number and the first digit
	if len(v) != 11 || v[0] != '1' {
		return false
	}
	// Check if all digits from the second to the last are numbers
	for _, digit := range v[1:] {
		if digit < '0' || digit > '9' {
			return false
		}
	}
	// Check if the first four digits match any of the specified valid prefixes
	validPrefixes := []string{
		"130", "131", "132", "133", "134", "135", "136", "137", "138", "139",
		"1400", "1410", "1440", "145", "146", "147", "148", "149",
		"150", "151", "152", "153", "154", "155", "156", "157", "158", "159",
		"162", "165", "166", "167",
		"170", "1703", "1704", "1705", "1706", "1707", "1708", "1709",
		"171", "172", "173", "1740", "1741", "1742", "1743", "1744", "1745", "1749",
		"175", "176", "177", "178",
		"180", "181", "182", "183", "184", "185", "186", "187", "188", "189",
		"190", "191", "192", "193", "195", "196", "197", "198", "199",
	}
	prefix := v[0:4]
	for _, validPrefix := range validPrefixes {
		if strings.HasPrefix(prefix, validPrefix) {
			return true
		}
	}
	return false
}

// IsEmail is the validation function for validating if the current field's value is a valid email address.
func IsEmail(addr string) bool {
	a, err := mail.ParseAddress(addr)
	if err != nil || strings.ContainsRune(addr, '<') {
		return false
	}

	addr = a.Address
	if len(addr) > 254 {
		return false
	}

	parts := strings.SplitN(addr, "@", 2)
	return len(parts[0]) <= 64 && isHostname(parts[1])
}

func isHostname(host string) bool {
	if len(host) > 253 {
		return false
	}

	s := strings.ToLower(strings.TrimSuffix(host, "."))
	// split hostname on '.' and validate each part
	for _, part := range strings.Split(s, ".") {
		// if part is empty, longer than 63 chars, or starts/ends with '-', it is invalid
		if l := len(part); l == 0 || l > 63 || part[0] == '-' || part[l-1] == '-' {
			return false
		}
		// for each character in part
		for _, ch := range part {
			// if the character is not a-z, 0-9, or '-', it is invalid
			if (ch < 'a' || ch > 'z') && (ch < '0' || ch > '9') && ch != '-' {
				return false
			}
		}
	}

	return true
}
