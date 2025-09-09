package app

import (
	"log"

	httpdelivery "logistics-app/backend/internal/delivery/http"
	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/infra/db"
	"logistics-app/backend/internal/repository"
	"logistics-app/backend/internal/usecase"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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
	// Seed simple admin if not exists (store hashed password)
	admin := domain.User{Email: "admin@example.com", FullName: "Admin", IsActive: true}
	var existing domain.User
	if err := database.Where(&admin).First(&existing).Error; err != nil {
		// Not found -> create with hashed password
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		database.Create(&domain.User{Email: admin.Email, Password: string(hashed), Role: domain.RoleAdmin, FullName: "Admin", IsActive: true})
	} else {
		// If found and password seems not bcrypt (no prefix), rehash simple case (best-effort)
		if len(existing.Password) > 0 && !(len(existing.Password) > 4 && existing.Password[0] == '$') {
			hashed, _ := bcrypt.GenerateFromPassword([]byte(existing.Password), bcrypt.DefaultCost)
			existing.Password = string(hashed)
			database.Save(&existing)
		}
	}

	orderRepo := repository.NewOrderGormRepo(database)
	orderSvc := usecase.NewOrderService(orderRepo)
	userRepo := repository.NewUserGormRepo(database)
	userSvc := usecase.NewUserService(userRepo)
	h := &httpdelivery.Handler{Orders: orderSvc, Users: userSvc}
	h.Register(r)
	log.Println("Bootstrap completed")
	return nil
}
