package inventory

import inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer
	partService PartService
}

func NewApi(partService PartService) *api {
	return &api{
		partService: partService,
	}
}
