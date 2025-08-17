package middlewares

import "rest/goAPI/utils"

func IsAuth(tokenString string) (bool, *utils.Payload) {
	payload, err := utils.ValidateToken(tokenString)
	if err != nil {
		return false, nil
	}

	return true, payload
}
