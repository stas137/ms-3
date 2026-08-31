package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
	"github.com/stas137/ms-3/order/internal/service/input"
)

// CreateOrder реализует операцию createOrder
// POST /api/v1/orders
func (s *service) Create(ctx context.Context, in input.CreateOrderInput) (model.Order, error) {
	if in.EngineUUID.String() == "" {
		return model.Order{}, errs.ErrInvalidUUID
	}
	if in.EngineUUID == uuid.Nil {
		return model.Order{}, errs.ErrInvalidUUID
	}
	if _, err := uuid.Parse(in.EngineUUID.String()); err != nil {
		return model.Order{}, errs.ErrInvalidUUID
	}
	if in.HullUUID.String() == "" {
		return model.Order{}, errs.ErrInvalidUUID
	}
	if in.HullUUID == uuid.Nil {
		return model.Order{}, errs.ErrInvalidUUID
	}
	if _, err := uuid.Parse(in.HullUUID.String()); err != nil {
		return model.Order{}, errs.ErrInvalidUUID
	}

	uuids := []string{
		in.HullUUID.String(),
		in.EngineUUID.String(),
	}
	if in.WeaponUUID != nil {
		uuids = append(uuids, in.WeaponUUID.String())
	}
	if in.ShieldUUID != nil {
		uuids = append(uuids, in.ShieldUUID.String())
	}

	// inventory service
	listParts, err := s.inventoryClient.ListParts(
		ctx,
		uuids,
	)
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			// listParts = defaultParts // defaultParts
			return model.Order{}, fmt.Errorf("создать заказ: (деталь не найдена): %w", err)
		} else {
			return model.Order{}, fmt.Errorf("создать заказ (получить детали): %w", err)
		}
	}

	for _, part := range listParts {
		if part.StockQuantity <= 0 {
			return model.Order{}, fmt.Errorf("создать заказ: деталь %s %w", part.Name, errs.ErrOutOfStock)
		}
	}

	hullUUID := in.HullUUID.String()
	engineUUID := in.EngineUUID.String()

	var shieldUUID string
	if in.ShieldUUID != nil {
		shieldUUID = in.ShieldUUID.String()
	}

	var weaponUUID string
	if in.WeaponUUID != nil {
		weaponUUID = in.WeaponUUID.String()
	}

	err = s.inventoryClient.ValidateCompatibility(ctx, hullUUID, engineUUID, shieldUUID, weaponUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("создать заказ: validate compatibility %w", err)
	}

	orderUUID := uuid.New()
	createOrder := model.Order{
		UUID:      orderUUID,
		Items:     getOrderItems(listParts),
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}

	err = s.inventoryClient.ReserveParts(ctx, []string{hullUUID, engineUUID, shieldUUID, weaponUUID})
	if err != nil {
		return model.Order{}, fmt.Errorf("cоздать заказ: зарезервировать детали: %w", err)
	}

	err = s.orderRepo.Create(ctx, createOrder, getOrderItems(listParts))
	if err != nil {
		return model.Order{}, fmt.Errorf("создать заказ: %w", err)
	}

	return createOrder, nil
}

func getOrderItems(parts []model.Part) []model.OrderItem {
	orderItems := make([]model.OrderItem, len(parts))
	for idx, part := range parts {
		orderItems[idx] = model.OrderItem{
			PartUUID: part.UUID,
			PartType: part.PartType,
			Price:    part.Price,
		}
	}
	return orderItems
}
