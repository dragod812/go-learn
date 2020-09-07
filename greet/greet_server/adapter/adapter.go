package adapter

import "github.com/go-learn/greet/greetpb"

func AdaptStringToGreetResponse(result string) *greetpb.GreetResponse {
	return &greetpb.GreetResponse{
		Result: result,
	}
}
