package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jvzantvoort/vimtmpl"
)

func Ask(question string, input string) string {
	reader := bufio.NewReader(os.Stdin)
	if len(input) > 0 {
		question += "[" + input + "]"
	}

	fmt.Printf("%s: ", question)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic("at the disco")
	}
	text = strings.TrimSuffix(text, "\n")
	if len(input) > 0 && len(text) < 1 {
		return input
	}
	return text

}
func main() {
	uc := vimtmpl.UserConfig{}
	mymap := uc.Load()

	user := Ask("Account name", mymap["user"])
	uc.Set("user", user)

	username := Ask("Full user name", mymap["username"])
	uc.Set("username", username)

	mailaddress := Ask("Email", mymap["mailaddress"])
	uc.Set("mailaddress", mailaddress)

	company := Ask("Company", mymap["company"])
	uc.Set("company", company)

	copyright := Ask("Copyright", mymap["copyright"])
	uc.Set("copyright", copyright)

	license := Ask("License", mymap["license"])
	uc.Set("license", license)
}

// vim: noexpandtab filetype=go
