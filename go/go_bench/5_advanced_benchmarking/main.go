package main

import (
	"regexp"
)

// type EmailValidities struct {
// 	Email    string
// 	Validity bool
// }

// func IsValidEmail(emails []string) []EmailValidities {
// 	var emailValidities []EmailValidities

// 	for _, email := range emails {
// 		var re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// 		emailValidities = append(emailValidities, EmailValidities{
// 			Email:    email,
// 			Validity: re.MatchString(email),
// 		})
// 	}

// 	return emailValidities
// }

type EmailValidities struct {
	Email    string
	Validity bool
}

var re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsValidEmail(emails []string) []EmailValidities {
	emailValidities := make([]EmailValidities, len(emails))

	for i, email := range emails {
		emailValidities[i] = EmailValidities{
			Email:    email,
			Validity: re.MatchString(email),
		}
	}

	return emailValidities
}

// BenchmarkIsValidEmail-8            29851            200334 ns/op          608007 B/op       3308 allocs/op
// BenchmarkIsValidEmail-8          2196948              2741 ns/op             361 B/op          4 allocs/op
// BenchmarkIsValidEmail-8          2346427              2548 ns/op             176 B/op          1 allocs/op
