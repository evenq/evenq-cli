package config

import (
	"log"
	"testing"
)

func TestSetValue(t *testing.T) {
	SetValue("something", "hi")
}

func TestGetValue(t *testing.T) {
	log.Println(GetValue("something"))
}
