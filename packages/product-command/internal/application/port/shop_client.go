package port

import "context"

type CheckPermissionRequest struct {
	ShopID string
	UserID string
	Action string
}

type CheckPermissionResponse struct {
	IsAllowed bool
	Message   string
}

type ShopClient interface {
	CheckPermission(ctx context.Context, req CheckPermissionRequest) (CheckPermissionResponse, error)
}
