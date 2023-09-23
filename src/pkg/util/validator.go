package util

import "regexp"

func IsValidDomainName(domain string) bool {
	re := regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]\.)+[a-zA-Z]{2,6}$`)
	return re.MatchString(domain)
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsValidPassword(password string) bool {
	re := regexp.MustCompile(`^(.*[0-9])(.*[a-z])(.*[A-Z])(.*[#?!@$%*\-])[0-9A-Za-z#?!@$%*\-]{8,16}$`)
	return re.MatchString(password)
}
