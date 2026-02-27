package actionservice

import "context"

type Action interface {
	Name() string
	Validate(params map[string]interface{}) error
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}
