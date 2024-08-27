package user

import "errors"

const (
	ThemeDark    = "dark"
	ThemeLight   = "light"
	ThemeDefault = ThemeDark
)

func IsTheme(s string) bool {
	switch s {
	case ThemeDark:
		return true
	case ThemeLight:
		return true
	default:
		return false
	}
}

func OtherTheme(theme string) (string, error) {
	if !IsTheme(theme) {
		return "", errors.New("that isn't a theme")
	}

	switch theme {
	case ThemeDark:
		return ThemeLight, nil
	case ThemeLight:
		return ThemeDark, nil
	default:
		return ThemeDefault, nil
	}
}
