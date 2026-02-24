package service

import (
	"context"

	"github.com/luanaands/clean-arch/internal/infra/grpc/pb"
	"github.com/luanaands/clean-arch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase  usecase.CreateOrderUseCase
	GetAllOrdersUseCase usecase.GetAllOrdersUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, getAllOrdersUseCase usecase.GetAllOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase:  createOrderUseCase,
		GetAllOrdersUseCase: getAllOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) GetAllOrders(ctx context.Context, in *pb.Blank) (*pb.OrderList, error) {
	output, err := s.GetAllOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}
	var orders []*pb.GetOrderResponse
	for _, order := range output {
		orders = append(orders, &pb.GetOrderResponse{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
	}
	return &pb.OrderList{Orders: orders}, nil
}
