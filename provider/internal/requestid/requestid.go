package requestid

import "context"

type requestIDKey struct{}

var Key = requestIDKey{}

func Get(ctx context.Context) string {
	id, _ := ctx.Value(Key).(string)

	return id
}
