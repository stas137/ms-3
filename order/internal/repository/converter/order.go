package repoConverter

import (
	"github.com/stas137/ms-3/order/internal/model"
	"github.com/stas137/ms-3/order/internal/repository/record"
)

func RepoOrderItemsToModelOrderItems(orderItems []record.OrderItem) []model.OrderItem {
	res := make([]model.OrderItem, len(orderItems))
	for idx, item := range orderItems {
		res[idx] = model.OrderItem{
			PartUUID: item.PartUUID,
			PartType: model.PartType(item.PartType),
			Price:    item.Price,
		}
	}
	return res
}

func RepoOrderToModelOrder(order record.Order, orderItems []record.OrderItem) model.Order {
	var paymentMethod *model.PaymentMethod

	if order.PaymentMethod != nil {
		paymentMethod = new(model.PaymentMethod(*order.PaymentMethod))
	}
	return model.Order{
		UUID:            order.UUID,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Items:           RepoOrderItemsToModelOrderItems(orderItems),
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		DeletedAt:       order.DeletedAt,
	}
}

func ModelOrderToRepoOrder(order model.Order) record.Order {
	var paymentMethod *string
	if order.PaymentMethod != nil {
		// model.PaymentMethod — доменный enum, в record он схлопывается в *string
		paymentMethod = new(string(*order.PaymentMethod))
	}
	recordOrder := record.Order{
		UUID:            order.UUID,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          string(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		DeletedAt:       order.DeletedAt,
	}
	return recordOrder
}

func ModelOrderItemsToRepoOrderItems(order model.Order) []record.OrderItem {
	res := make([]record.OrderItem, len(order.Items))
	for idx, item := range order.Items {
		res[idx] = record.OrderItem{
			OrderUUID: order.UUID,
			PartUUID:  item.PartUUID,
			PartType:  string(item.PartType),
			Price:     item.Price,
		}
	}
	return res
}
