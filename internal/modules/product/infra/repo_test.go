package infra

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/anton-chornobai/beton.git/internal/db"
	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
	"github.com/google/uuid"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	connStr := os.Getenv("DB_TEST_CONN_STR")
	if connStr == "" {
		log.Fatal("DB_TEST_CONN_STR not set")
	}

	var err error
	testDB, err = db.OpenPostgre(connStr)
	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	defer testDB.Close()

	os.Exit(m.Run())
}

func clearTables(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec("TRUNCATE TABLE product_image, products RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatal("failed to clear tables:", err)
	}
}

func newRepo() *ProductRepository {
	return &ProductRepository{DB: testDB}
}

func TestAdd_MinimalValidProduct(t *testing.T) {
	clearTables(t)

	err := newRepo().Add(context.Background(), &domain.Product{
		Price:  100,
		Title:  "Simple Tile",
		Type:   "tile",
		Status: domain.ProductDisplayed,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestAdd_ValidateFullProduct(t *testing.T) {
	clearTables(t)

	err := newRepo().Add(context.Background(), &domain.Product{
		Price:  100,
		Title:  "Simple Tile",
		Type:   "tile",
		Status: domain.ProductDisplayed,
		Color:  new("white"),
		ImageURLs: []domain.ProductImage{
			{
				ID:  uuid.New(),
				URL: "https://example.com/images/tile1_front.jpg",
			},
			{
				ID:  uuid.New(),
				URL: "https://example.com/images/tile1_angle.jpg",
			},
		},
		Description:   new("This is prpdocyt!"),
		StockQuantity: new(10),
		Weight:        new(200),
		Rating:        new(4),
		Size: &domain.Size{
			Width:  10,
			Height: 10,
		},
	})

	if err != nil {
		t.Fatalf("expected full product to save successfully, got: %v", err)
	}
}

func TestAdd_DuplicateTitleError(t *testing.T) {
	clearTables(t)

	repo := newRepo()
	sharedTitle := "Unique Tile Name"

	p1 := &domain.Product{
		Price:  200,
		Title:  sharedTitle,
		Type:   "brick",
		Status: domain.ProductDisplayed,
	}

	if err := repo.Add(context.Background(), p1); err != nil {
		t.Fatalf("failed to insert initial product: %v", err)
	}

	p2 := &domain.Product{
		Price:  300,
		Title:  sharedTitle,
		Type:   "brick",
		Status: domain.ProductDisplayed,
	}

	err := repo.Add(context.Background(), p2)
	if !errors.Is(err, domain.ErrTitleAlreadyExists) {
		t.Fatalf("expected domain.ErrTitleAlreadyExists, got: %v", err)
	}
}