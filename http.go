package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	parallel := flagSet.Int("parallel", 10, "A parallelization factor, value must be a number greater than 0. Default is 10")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		panic(fmt.Errorf("error parsing arguments, received arguments %v", os.Args[1:]))
	}

	if parallel != nil && *parallel <= 0 {
		panic(fmt.Errorf("A parallelization factor value must be a number greater than 0, but recived %d", *parallel))
	}

	if parallel == nil {
		panic(fmt.Errorf("error parsing arguments, received arguments %v", os.Args[1:]))
	}

	addresses := flagSet.Args()

	if len(addresses) == 0 {
		log.Printf("No addresses were provided. You need at least 1 address to send a request")
		return
	}

	limiter := make(chan byte, *parallel)
	wg := sync.WaitGroup{}
	for _, address := range addresses {
		limiter <- byte(1)
		wg.Add(1)
		go func(addr string) {
			defer func() {
				wg.Done()
				<- limiter
			}()
			originalAddress := addr
			if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
				addr = "http://" + addr
			}
			response, err := http.Get(addr)
			if err != nil {
				log.Printf("error while sending request to the address: %s, error was: %v", addr, err)
				return
			}

			responseBody, err := ioutil.ReadAll(response.Body)
			err = response.Body.Close()
			if err != nil {
				log.Printf("[WARN] error while closing reponse body from the address: %s, error was: %v", originalAddress, err)
			}

			md5Hash := md5.Sum(responseBody)

			fmt.Printf("%s %x\n", originalAddress, md5Hash)
		}(address)
	}

	wg.Wait()
}
