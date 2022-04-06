package utils

import "../store"

func CheckEmptyFields(user *store.User) []string {
	var emptyFields []string
	if user.FirstName == "" {
		emptyFields = append(emptyFields, "first_name")
	}

	if user.LastName == "" {
		emptyFields = append(emptyFields, "last_name")
	}

	if user.Nickname == "" {
		emptyFields = append(emptyFields, "nickname")
	}

	if user.Password == "" {
		emptyFields = append(emptyFields, "password")
	}

	if user.Email == "" {
		emptyFields = append(emptyFields, "email")
	}

	if user.Country == "" {
		emptyFields = append(emptyFields, "country")
	}

	return emptyFields
}
