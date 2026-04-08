package application

import (
	"context"
	"testing"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

type MockRepo struct {
	AddCalled        bool
	LastAddedProduct *domain.Product

	GetWithLimitCalled bool
	LastLimit          int
	GetWithLimitResult []domain.Product
	GetWithLimitErr    error

	GetByIDCalled bool
	LastGetID     int
	GetByIDResult *domain.Product
	GetByIDErr    error

	DeleteByIDCalled bool
	LastDeleteID     int
	DeleteByIDErr    error

	UpdateCalled  bool
	LastUpdateID  int
	LastUpdateReq domain.ProductUpdate
	UpdateErr     error
}

func (m *MockRepo) Add(ctx context.Context, product *domain.Product) error {
	m.AddCalled = true
	m.LastAddedProduct = product
	return nil
}

func (m *MockRepo) GetProducts(ctx context.Context, limit int, status *domain.ProductStatus) ([]domain.Product, error) {
	m.GetWithLimitCalled = true
	m.LastLimit = limit
	if m.GetWithLimitErr != nil {
		return nil, m.GetWithLimitErr
	}
	return m.GetWithLimitResult, nil
}

func (m *MockRepo) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	m.GetByIDCalled = true
	m.LastGetID = id
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	return m.GetByIDResult, nil
}

func (m *MockRepo) DeleteByID(ctx context.Context, id int) error {
	m.DeleteByIDCalled = true
	m.LastDeleteID = id
	return m.DeleteByIDErr
}

func (m *MockRepo) Update(ctx context.Context, id int, req domain.ProductUpdate) error {
	m.UpdateCalled = true
	m.LastUpdateID = id
	m.LastUpdateReq = req
	return m.UpdateErr
}

func TestProductAdd(t *testing.T) {
	mockRepo := &MockRepo{}
	service, _ := NewProductService(mockRepo)
	ctx := context.Background()

	color := "white"
	stock := 120
	weight := 100
	rating := 5
	size := &domain.Size{Width: 10, Height: 20}

	input := domain.Product{
		Price:         100,
		Title:         "product1",
		Type:          "tile",
		Color:         &color,
		Status:        domain.ProductArchived,
		ImageURL:      nil,
		Description:   nil,
		StockQuantity: &stock,
		Weight:        &weight,
		Rating:        &rating,
		Size:          size,
	}

	err := service.Add(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mockRepo.AddCalled {
		t.Fatal("expected Add to be called on repo")
	}

	if mockRepo.LastAddedProduct.Title != input.Title {
		t.Fatalf("expected Title %q, got %q", input.Title, mockRepo.LastAddedProduct.Title)
	}
}
