package repoConverter

import (
	"github.com/google/uuid"

	"github.com/stas137/ms-3/order/internal/model"
	"github.com/stas137/ms-3/order/internal/repository/record"
)

func getItems(order record.Order) []model.OrderItem {
	var orderItems []model.OrderItem

	orderItems = append(orderItems, model.OrderItem{
		PartUUID: order.EngineUUID,
		PartType: model.PartTypeEngine,
		Price:    order.EnginePrice,
	})
	orderItems = append(orderItems, model.OrderItem{
		PartUUID: order.HullUUID,
		PartType: model.PartTypeHull,
		Price:    order.HullPrice,
	})

	if order.ShieldUUID != nil {
		orderItems = append(orderItems, model.OrderItem{
			PartUUID: *order.ShieldUUID,
			PartType: model.PartTypeShield,
			Price:    *order.ShieldPrice,
		})
	}
	if order.WeaponUUID != nil {
		orderItems = append(orderItems, model.OrderItem{
			PartUUID: *order.WeaponUUID,
			PartType: model.PartTypeWeapon,
			Price:    *order.WeaponPrice,
		})
	}

	return orderItems
}

func RepoOrderToModelOrder(order record.Order) model.Order {
	return model.Order{
		UUID:            order.OrderUUID,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   (*model.PaymentMethod)(order.PaymentMethod),
		Items:           getItems(order),
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func ModelOrderToRepoOrder(order model.Order) record.Order {
	var engineUUID uuid.UUID
	var hullUUID uuid.UUID
	var shieldUUID *uuid.UUID
	var weaponUUID *uuid.UUID

	var enginePrice int64
	var hullPrice int64
	var shieldPrice *int64
	var weaponPrice *int64

	var paymentMethod *string

	for _, item := range order.Items {
		if item.PartType == model.PartTypeHull {
			hullUUID = item.PartUUID
			hullPrice = item.Price
		}
		if item.PartType == model.PartTypeEngine {
			engineUUID = item.PartUUID
			enginePrice = item.Price
		}
		if item.PartType == model.PartTypeShield {
			shieldUUID = &item.PartUUID
			shieldPrice = &item.Price
		}
		if item.PartType == model.PartTypeWeapon {
			weaponUUID = &item.PartUUID
			weaponPrice = &item.Price
		}
	}

	if order.PaymentMethod != nil {
		paymentMethod = (*string)(order.PaymentMethod)
	}

	recordOrder := record.Order{
		OrderUUID:       order.UUID,
		EngineUUID:      engineUUID,
		HullUUID:        hullUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		EnginePrice:     enginePrice,
		HullPrice:       hullPrice,
		ShieldPrice:     shieldPrice,
		WeaponPrice:     weaponPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          string(order.Status),
		TotalPrice:      order.TotalPrice(),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		DeletedAt:       order.DeletedAt,
	}

	return recordOrder
}
