package Emailhelper

import "strings"

func GenerateEmails(email string, role string) string {
	parsedEmail := strings.Split(email, "@")
	if len(parsedEmail) < 2 {
		return ""
	}

	var newEmail string

	if role == "school admin" {
		newEmail = parsedEmail[0] + "@school.edu"
	}
	return newEmail
}

func GenerateUsername(name string) string {

	return ""
}
