package core

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/manifoldco/promptui"
)

// NumberPrompt asks a numerical questions using the label.
func NumberPrompt(ctx context.Context, label string, defaultValue string) string {
	logger := ctx.Value("logger").(*slog.Logger)
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:       label,
		Validate:    validate,
		Default:     defaultValue,
		AllowEdit:   true,
		HideEntered: false,
	}

	result, err := prompt.Run()
	if err != nil {
		logger.Error("Prompt failed", "err", err)
		return ""
	}

	logger.Debug("You selected a value", "value", result)
	return result
}

// StringPrompt asks a open string questions using the label.
func StringPrompt(ctx context.Context, label string, defaultValue string) string {
	logger := ctx.Value("logger").(*slog.Logger)
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:       label,
		Validate:    validate,
		Default:     defaultValue,
		AllowEdit:   true,
		HideEntered: false,
	}

	result, err := prompt.Run()
	if err != nil {
		logger.Error("Prompt failed", "err", err)
		return ""
	}

	logger.Debug("You selected a value", "value", result)
	return result
}

// PasswordPrompt asks a password questions using the label.
func PasswordPrompt(ctx context.Context, label string, defaultValue string) string {
	logger := ctx.Value("logger").(*slog.Logger)
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("password must have more than 6 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:       label,
		Validate:    validate,
		Mask:        '*',
		AllowEdit:   true,
		HideEntered: true,
	}

	result, err := prompt.Run()
	if err != nil {
		logger.Error("Prompt failed", "err", err)
		return ""
	}

	logger.Debug("You selected a value", "value", result)
	return result
}

// BoolPrompt asks a boolean questions using the label.
func BoolPrompt(ctx context.Context, label string, defaultValue string, isOpen bool) string {
	logger := ctx.Value("logger").(*slog.Logger)
	validate := func(input string) error {
		matched, err := regexp.MatchString("(?i)^(true|false|yes|no|y|n)", input)
		if err != nil {
			return errors.New("invalid value for bool")
		}
		if !matched {
			return errors.New("invalid value for bool")
		}
		return nil
	}

	if isOpen {
		prompt := promptui.Prompt{
			Label:       label,
			Default:     "FALSE",
			Validate:    validate,
			AllowEdit:   true,
			HideEntered: false,
		}

		result, err := prompt.Run()
		if err != nil {
			logger.Error("Prompt failed", "err", err)
			return "FALSE"
		}

		matched, err := regexp.MatchString("(?i)^(true|yes|y)", result)
		if err != nil {
			panic("error regex")
		}
		if matched {
			return "TRUE"
		} else {
			return "FALSE"
		}

	}

	// Select type for bool
	prompt := promptui.Select{
		Label: label,
		Items: []string{"TRUE", "FALSE"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		logger.Error("Prompt failed", "err", err)
		return ""
	}

	logger.Debug("You selected a value", "value", result)
	return result
}

// SelectPrompt asks a single select questions using the label.
func SingleSelectPrompt(ctx context.Context, label string, items []string) string {
	logger := ctx.Value("logger").(*slog.Logger)
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		logger.Error("Prompt failed", "err", err)
		return ""
	}

	logger.Debug("You selected a value", "value", result)
	return result
}
