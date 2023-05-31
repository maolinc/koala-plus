package interceptor

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"testing"
)

func TestJwt(t *testing.T) {
	parser := jwt.NewParser()
	parse, err := parser.Parse("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY1MzMzMjUsImlhdCI6MTY4NDk5NzMyNSwiand0VXNlcklkIjoxfQ.fnX3aLgLbPB1dHVx4GRthl58eaJnVCYKZB1YOEVlJbk",
		func(token *jwt.Token) (interface{}, error) {
			fmt.Println(token)
			return []byte("ae0536f9-6450-4606-8e13-5a19ed505da0"), nil
		})
	if err != nil {
		return
	}
	m := parse.Claims.(jwt.MapClaims)

	fmt.Println(m["jwtUserId"])
}
