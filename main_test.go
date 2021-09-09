package main

import (
	"testing"
)

func TestCreateProxy(t *testing.T) {
	_, err := createProxy("abc.test", "")
	if err == nil {
		t.Fatalf("Error should be returned because of empty content")
	}

	_, err = createProxy("abc.test", "1http://12ab:800a")
	if err == nil {
		t.Fatalf("Error should be returned because mallformed domain")
	}

	p, err := createProxy("abc.test", "8080")
	if err != nil {
		t.Fatalf("Error response for valid entry")
	}

	if p.Url.Host != "127.0.0.1:8080" {
		t.Fatalf("Localhost was not added")
	}

	if p.Url.Scheme != "http" {
		t.Fatalf("Scheme was added wrong")
	}

	p, err = createProxy("abc.test", "192.168.1.1:8080")
	if err != nil {
		t.Fatalf("Error response for valid entry")
	}

	if p.Url.Host != "192.168.1.1:8080" {
		t.Fatalf("Localhost was not added")
	}

	if p.Url.Scheme != "http" {
		t.Fatalf("Scheme was added wrong")
	}

	p, err = createProxy("abc.test", "8080")
	if err != nil {
		t.Fatalf("Error response for valid entry")
	}

	if p.Url.Host != "127.0.0.1:8080" {
		t.Fatalf("Localhost was not added")
	}

	if p.Url.Scheme != "http" {
		t.Fatalf("Scheme was added wrong")
	}

	p, err = createProxy("abc.test", "192.168.1.1:8080")
	if err != nil {
		t.Fatalf("Error response for valid entry")
	}

	if p.Url.Host != "192.168.1.1:8080" {
		t.Fatalf("Localhost was not added")
	}

	if p.Url.Scheme != "http" {
		t.Fatalf("Scheme was added wrong")
	}

	p, err = createProxy("abc.test", "https://google.com")
	if err != nil {
		t.Fatalf("Error response for valid entry")
	}

	if p.Url.Host != "google.com" {
		t.Fatalf("Localhost was not added")
	}

	if p.Url.Scheme != "https" {
		t.Fatalf("Scheme was added wrong")
	}
}
