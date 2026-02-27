package application_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
// 	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
// )

// type FakeHasher struct {}


// func (f *FakeHasher) HashPassword(password string) (string, error) {
//     return "fake-hash-" + password, nil
// }

// func (f *FakeHasher) CompareHashAndPassword(hashed, password string) error {
//     if hashed != "fake-hash-"+password {
//         return errors.New("password mismatch")
//     }
//     return nil
// }

// type FakeTokenManager struct{}

// func (f *FakeTokenManager) GenerateToken(id, role string) (string, error) {
//     return "fake-token", nil
// }

// type FakeRepo struct{}

// func (f *FakeRepo) SignupByEmail(ctx context.Context, user domain.User) error {
//     return nil
// }
// func TestUserService_SignupByEmail_Success(t *testing.T) {
//     repo := &FakeRepo{}
//     tokenManager := &FakeTokenManager{}
//     hasher := &FakeHasher{}

//     service := application.NewUserService(repo, tokenManager, hasher)

//     token, err := service.SignupByEmail(
//         context.Background(),
//         "test@example.com",
//         "Password123",
//     )

//     if err != nil {
//         t.Fatalf("expected no error, got %v", err)
//     }

//     if token != "fake-token" {
//         t.Fatalf("expected fake-token, got %s", token)
//     }
// }