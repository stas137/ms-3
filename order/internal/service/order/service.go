package order

type service struct {
	orderRepo       OrderRepository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	txManager       TxManager
}

func NewService(
	orderRepo OrderRepository,
	inventoryClient InventoryClient,
	paymentClient PaymentClient,
	txManager TxManager,
) *service {
	return &service{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		txManager:       txManager,
	}
}
