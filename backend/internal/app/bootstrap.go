package app

import (
	"log"

	httpdelivery "logistics-app/backend/internal/delivery/http"
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/infra/db"
	"logistics-app/backend/internal/repository"
	"logistics-app/backend/internal/usecase"

	"github.com/gorilla/mux"
)

func Bootstrap(r *mux.Router) error {
	database, err := db.Connect()
	if err != nil {
		return err
	}
	// AutoMigrate
	if err := database.AutoMigrate(&domain.User{}, &domain.Order{}); err != nil {
		return err
	}
	// Seed simple admin if not exists
	database.Where(domain.User{Email: "admin@example.com"}).Attrs(domain.User{Password: "admin", Role: domain.RoleAdmin, FullName: "Admin", IsActive: true}).FirstOrCreate(&domain.User{})

	orderRepo := repository.NewOrderGormRepo(database)
	orderSvc := usecase.NewOrderService(orderRepo)
	userRepo := repository.NewUserGormRepo(database)
	userSvc := usecase.NewUserService(userRepo)
	h := &httpdelivery.Handler{Orders: orderSvc, Users: userSvc}
	h.Register(r)
	log.Println("Bootstrap completed")
	return nil
}
