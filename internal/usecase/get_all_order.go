package usecase

import (
	"github.com/luanaands/clean-arch/internal/entity"
)

type OrdersDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type GetAllOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewGetAllOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *GetAllOrdersUseCase {
	return &GetAllOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *GetAllOrdersUseCase) Execute() ([]OrdersDTO, error) {
	orders, err := c.OrderRepository.GetAll()
	if err != nil {
		return []OrdersDTO{}, err
	}
	var ordersDTO []OrdersDTO
	for _, order := range orders {
		dto := OrdersDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		}
		ordersDTO = append(ordersDTO, dto)
	}
	return ordersDTO, nil
}
