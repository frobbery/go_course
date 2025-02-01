package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	//nolint:depguard
	"github.com/tidwall/gjson"
)

type User struct {
	ID int

	Name string

	Username string

	Email string

	Phone string

	Password string

	Address string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)

	result := make(DomainStat)

	re := regexp.MustCompile(domain + "$")

	for scanner.Scan() {
		email := gjson.Get(scanner.Text(), "Email").Str

		if re.MatchString(email) {
			result[strings.ToLower(strings.SplitN(email, "@", 2)[1])]++
		}
	}

	return result, nil
}
