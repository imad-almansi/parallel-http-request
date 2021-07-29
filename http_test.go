package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMainSuccess(t *testing.T) {
	servers := []*httptest.Server{
		httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			_, _ = rw.Write([]byte("hello world"))
		})),
		httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			_, _ = rw.Write([]byte("hello world!"))
		})),
	}
	for _, test := range []struct {
		Name string
		Servers []*httptest.Server
		Args    []string
		Output  []string
		ExpectError bool
	}{
		{
			Name: "test default single address",
			Servers: servers[0:0],
			Args:   []string{"./test", servers[0].URL},
			Output: []string{
				fmt.Sprintf("%s %s\n", servers[0].URL, "5eb63bbbe01eeed093cb22bb8f5acdc3"),
			},
			ExpectError: false,
		},
		{
			Name: "test default multiple addresses",
			Args:   []string{"./test", servers[0].URL, servers[1].URL},
			Output: []string{
				fmt.Sprintf("%s %s\n", servers[0].URL, "5eb63bbbe01eeed093cb22bb8f5acdc3"),
				fmt.Sprintf("%s %s\n", servers[1].URL, "fc3ff98e8c6a0d3087d515c0473f8677"),
			},
			ExpectError: false,
		},
		{
			Name: "test invalid argument flag",
			Args:   []string{"./test", "-parallel", "0", servers[0].URL},
			Output: []string{
				fmt.Sprintf("%s %s\n", servers[0].URL, "5eb63bbbe01eeed093cb22bb8f5acdc3"),
			},
			ExpectError: true,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			os.Args = test.Args

			r, w, err := os.Pipe()
			if err != nil {
				log.Fatal(err)
			}

			origStdout := os.Stdout
			os.Stdout = w

			defer func() {
				rec := recover()

				// Restore
				os.Stdout = origStdout

				err = w.Close()
				if err != nil {
					log.Fatal(err)
				}

				actual, err := ioutil.ReadAll(r)
				if err != nil {
					log.Fatal(err)
				}

				if rec != nil {
					if test.ExpectError {
						return
					}
					t.Errorf("expected %s, got error %v", test.Output, rec)
					return
				}

				for _, expectedMessage := range test.Output {
					if !strings.Contains(string(actual), expectedMessage) {
						t.Errorf("expected %s, got %s", expectedMessage, actual)
					}
				}
			}()
			main()
		})
	}
}
