package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"testing"
)

func TestSlog(t *testing.T) {
	slog.Info("Hello slog")
	slog.Info("Hello", "user", "test")

	logger := slog.Default()
	logger.Info("hello slog.Default()")

	textLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	textLogger.Info("hello", "type", "TextHandler")

	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	jsonLogger.Info("hello", "type", "JSONHandler")

	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"hello",
		slog.String("type", "Attrs"),
		slog.Int("num", 100),
	)

	file, err := os.OpenFile("test.log", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileJSONLogger := slog.New(slog.NewJSONHandler(file, nil))
	fileJSONLogger.Info("Hello")
	fileJSONLogger.Info("Hello", slog.String("foo", "bar"))
}

func TestSentinelAddr(t *testing.T) {
	addrs, err := net.LookupHost("google.com")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(addrs)
}
