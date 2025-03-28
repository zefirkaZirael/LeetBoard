package utils

import (
	"1337bo4rd/internal/domain"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Prints Help and Exit program with 0 code
func PrintHelp() {
	t := ` ./1337b04rd --help
hacker board

Usage:
  1337b04rd [--port <N>]  
  1337b04rd --help

Options:
  --help       Show this screen.
  --port N     Port number.
	`
	fmt.Println(t)
	os.Exit(0)
}

func CheckFlags() {
	flag.Parse()
	InvalidPort := errors.New("port flag range must be in 1025 and 65000")
	if *domain.Port < 1024 || *domain.Port > 65000 {
		slog.Error("Failed to start program", "CheckFlags err: ", InvalidPort)
		log.Fatal(InvalidPort)
	}
	if *domain.HelpFlag {
		PrintHelp()
	}
}

func ErrorPage(w http.ResponseWriter, message error, code int) error {
	var Error struct {
		Code    int
		Message string
	}
	Error.Code = code
	Error.Message = message.Error()

	temp, err := template.ParseFiles("web/templates/error.html")
	if err != nil {
		slog.Error("❌ Failed to send error message: " + err.Error())
		return err
	}
	w.WriteHeader(code)
	err = temp.Execute(w, Error)
	if err != nil {
		slog.Error("❌ Failed to send error message: " + err.Error())
		return err
	}
	return nil
}

func GenerateSessionID() string {
	// Create a slice with 16 length
	b := make([]byte, 16)

	// Reading random data to slice
	_, err := rand.Read(b)
	if err != nil {
		// If error occured, we return empty string
		return ""
	}

	// Encode slice to string in hexademical format
	return hex.EncodeToString(b)
}

func DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})
}

func ConvertTime(timestampStr string) time.Time {
	layout := "2006-01-02 15:04:05.999999 -0700"
	timestamp, err := time.Parse(layout, timestampStr)
	if err != nil {
		slog.Error("❌ Parsing Error:" + err.Error())
		return timestamp
	}
	return timestamp
}

func DetectType(file []byte) (string, error) {
	mime := http.DetectContentType(file)
	var extension string
	switch mime {
	case "image/jpeg":
		extension = ".jpeg"
	case "image/png":
		extension = ".png"
	case "image/bmp":
		extension = ".bmp"
	default:
		ErrUnsupport := errors.New("unsupported MIME type")
		return extension, ErrUnsupport
	}
	return extension, nil
}
