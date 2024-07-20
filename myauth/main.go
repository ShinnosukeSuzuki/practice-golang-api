package main

import (
	"context"
	"fmt"

	"google.golang.org/api/idtoken"
)

func main() {
	googleClinetID := "[client_id]"
	idToken := `[id_token]`

	tokenValidator, err := idtoken.NewValidator(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	payload, err := tokenValidator.Validate(context.Background(), idToken, googleClinetID)
	if err != nil {
		fmt.Println("validate err:", err)
		return
	}

	fmt.Println(payload.Claims["name"])
}
