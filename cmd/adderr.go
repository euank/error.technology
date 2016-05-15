package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/euank/api.error.technology/errortech"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter a language: ")
	s.Scan()
	lang := s.Text()

	for {
		fmt.Println("Please enter an error for " + lang + ". It may be multiline. Terminate with the string 'EOE'.")
		errlines := ""
		for s.Scan() {
			if s.Text() == "EOE" {
				break
			}
			errlines += s.Text() + "\n"
		}
		fmt.Println("Enter the short error (one line): ")
		s.Scan()
		shorterr := s.Text()

		fmt.Println("Enter some tags (one line, comma delimited): ")
		s.Scan()
		tags := s.Text()

		inerr := errortech.Error{
			Language: lang,
			Tags:     strings.Split(tags, ","),
			Full:     errlines,
			Short:    shorterr,
		}
		errid := fmt.Sprintf("id-%v", rand.Int())

		data, err := json.Marshal(inerr)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(filepath.Join("base_errors", errid), data, 0777)
		fmt.Println("Wrote error: %v", errid)
	}
}
