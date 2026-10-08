package appcontext

import "context"


type contextKey string

const UserIdKey contextKey = "user_id"

func SetUserId(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, UserIdKey, userId)
}

func GetUserId(ctx context.Context)(int64, bool) {
	val ,ok := ctx.Value(UserIdKey).(int64)

	return  val, ok
}