package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	if orderUUID == uuid.Nil {
		return errs.ErrInvalidUUID
	}

	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("отменить заказ: %w", err)
	}

	if order.Status == model.OrderStatusCancelled {
		return errs.ErrOrderCancelled
	}
	if order.Status == model.OrderStatusPaid {
		return errs.ErrOrderAlreadyPaid
	}

	// applyUpdate(&order, updateOrder)

	uuids := []string{
		order.GetHullUUID().String(),
		order.GetEngineUUID().String(),
	}
	if order.GetShieldUUID() != nil {
		uuids = append(uuids, order.GetShieldUUID().String())
	}
	if order.GetWeaponUUID() != nil {
		uuids = append(uuids, order.GetWeaponUUID().String())
	}

	err = s.inventoryClient.ReleaseParts(ctx, uuids)
	if err != nil {
		return fmt.Errorf("отменить заказ: отменить резерв деталей: %w", err)
	}

	order.Status = model.OrderStatusCancelled
	order.UpdatedAt = new(time.Now())

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("отменить заказ: %w", err)
	}

	return nil
}

// func applyUpdate(order *model.Order, updateOrder model.Order) {
// 	var engineUUID *uuid.UUID
// 	var hullUUID *uuid.UUID
// 	var shieldUUID *uuid.UUID
// 	var weaponUUID *uuid.UUID
// 	var enginePrice *int64
// 	var hullPrice *int64
// 	var shieldPrice *int64
// 	var weaponPrice *int64

// 	if updateOrder.TransactionUUID != nil {
// 		order.TransactionUUID = updateOrder.TransactionUUID
// 	}
// 	if updateOrder.PaymentMethod != nil {
// 		order.PaymentMethod = updateOrder.PaymentMethod
// 	}
// 	if updateOrder.Status != "" {
// 		order.Status = updateOrder.Status
// 	}
// 	if updateOrder.DeletedAt != nil {
// 		order.DeletedAt = updateOrder.DeletedAt
// 	}

// 	for _, item := range updateOrder.Items {
// 		if item.PartType == model.PartTypeEngine {
// 			engineUUID = &item.PartUUID
// 			enginePrice = &item.Price
// 		}
// 		if item.PartType == model.PartTypeHull {
// 			hullUUID = &item.PartUUID
// 			hullPrice = &item.Price
// 		}
// 		if item.PartType == model.PartTypeShield {
// 			shieldUUID = &item.PartUUID
// 			shieldPrice = &item.Price
// 		}
// 		if item.PartType == model.PartTypeWeapon {
// 			weaponUUID = &item.PartUUID
// 			weaponPrice = &item.Price
// 		}
// 	}

// 	for idx, item := range order.Items {
// 		if engineUUID != nil && item.PartType == model.PartTypeEngine {
// 			order.Items[idx] = model.OrderItem{
// 				PartUUID: *engineUUID,
// 				PartType: model.PartTypeEngine,
// 				Price:    *enginePrice,
// 			}
// 		}
// 		if hullUUID != nil && item.PartType == model.PartTypeHull {
// 			order.Items[idx] = model.OrderItem{
// 				PartUUID: *hullUUID,
// 				PartType: model.PartTypeHull,
// 				Price:    *hullPrice,
// 			}
// 		}
// 		if shieldUUID != nil && item.PartType == model.PartTypeShield {
// 			order.Items[idx] = model.OrderItem{
// 				PartUUID: *hullUUID,
// 				PartType: model.PartTypeShield,
// 				Price:    *shieldPrice,
// 			}
// 		}
// 		if weaponUUID != nil && item.PartType == model.PartTypeWeapon {
// 			order.Items[idx] = model.OrderItem{
// 				PartUUID: *hullUUID,
// 				PartType: model.PartTypeWeapon,
// 				Price:    *weaponPrice,
// 			}
// 		}
// 	}
// }
