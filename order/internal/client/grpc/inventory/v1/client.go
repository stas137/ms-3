package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stas137/ms-3/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

type client struct {
	inventoryClient inventoryv1.InventoryServiceClient
}

func New(c inventoryv1.InventoryServiceClient) *client {
	return &client{
		inventoryClient: c,
	}
}

func (c *client) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	resp, err := c.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuids,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return nil, errs.ErrPartNotFound
			}
		}
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	modelParts, err := converter.ProtoPartsToModelParts(resp.GetParts())
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	return modelParts, nil
}

func (c *client) ValidateCompatibility(ctx context.Context, hullUUID, engineUUID, shieldUUID, weaponUUID string) error {
	_, err := c.inventoryClient.ValidateCompatibility(ctx, &inventoryv1.ValidateCompatibilityRequest{
		HullUuid:   hullUUID,
		EngineUuid: engineUUID,
		ShieldUuid: shieldUUID,
		WeaponUuid: weaponUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				if st.Message() == "детали несовместимы" {
					return errs.ErrIncompatibleParts
				}
			}
			if st.Code() == codes.InvalidArgument {
				if st.Message() == "тип детали не соответствует слоту корабля" {
					return errs.ErrPartTypeMismatch
				}
			}
		}
		return fmt.Errorf("validate compatibility: %w", err)
	}

	return nil
}

func (c *client) ReserveParts(ctx context.Context, uuids []string) error {
	_, err := c.inventoryClient.ReserveParts(ctx, &inventoryv1.ReservePartsRequest{
		Uuids: uuids,
	})
	if err != nil {
		return fmt.Errorf("зарезервировать детали: %w", err)
	}
	return nil
}

func (c *client) ReleaseParts(ctx context.Context, uuids []string) error {
	_, err := c.inventoryClient.ReleaseParts(ctx, &inventoryv1.ReleasePartsRequest{
		Uuids: uuids,
	})
	if err != nil {
		return fmt.Errorf("отменить резерв деталей: %w", err)
	}
	return nil
}
