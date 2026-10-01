package logger

import "log/slog"

func LogErr(message string, err error) error {
	slog.Error(message, "error", err)
	return err
}
