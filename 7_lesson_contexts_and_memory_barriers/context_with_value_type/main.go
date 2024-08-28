package main

import (
	"context"
	"fmt"
)

func main() {
	{
		ctx := context.WithValue(context.Background(), "key", "value1")
		ctx = context.WithValue(ctx, "key", "value2")

		fmt.Println("string =", ctx.Value("key").(string))
	}
	{
		type key1 string
		type key2 string
		const k1 key1 = "key"
		const k2 key2 = "key"

		ctx := context.WithValue(context.Background(), k1, "value1")
		ctx = context.WithValue(ctx, k2, "value2")

		fmt.Println("key1 =", ctx.Value(k1).(string))
		fmt.Println("key2 =", ctx.Value(k2).(string))
	}
}
