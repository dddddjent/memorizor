package util

func EncodePassword(password string) (string, error) {
	// TODO: add actual encode function
	return password, nil
}

func ComparePassword(store, given string) (bool, error) {
    // TODO
	if store == given {
		return true, nil
	}
	return false, nil
}
