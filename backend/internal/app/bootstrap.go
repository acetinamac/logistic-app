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
	if err := database.AutoMigrate(
		&domain.User{},
		&domain.Coordinates{},
		&domain.Address{},
		&domain.PackageType{},
		&domain.Order{},
		&domain.OrderStatusHistory{},
	); err != nil {
		return err
	}

	// Seed simple admin if not exists (store hashed password)
	admin := domain.User{Email: "admin@example.com", FullName: "System Administrator", IsActive: true}
	var existing domain.User

	if err := database.Where(&admin).First(&existing).Error; err != nil {
		// Not found -> create with hashed password
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		database.Create(&domain.User{Email: admin.Email, Password: string(hashed), Role: domain.RoleAdmin, FullName: "System Administrator", IsActive: true})
	} else {
		// If found and password seems not bcrypt (no prefix)
		if len(existing.Password) > 0 && !(len(existing.Password) > 4 && existing.Password[0] == '$') {
			hashed, _ := bcrypt.GenerateFromPassword([]byte(existing.Password), bcrypt.DefaultCost)
			existing.Password = string(hashed)
			database.Save(&existing)
		}
	}

	// Seed default package types if not present
	var count int64
	database.Model(&domain.PackageType{}).Count(&count)
	if count == 0 {
		// S, M, L
		pts := []domain.PackageType{
			{SizeCode: domain.PackageS, MaxWeightKg: 5.00, Description: "Small package - up to 5 kg", IsActive: true},
			{SizeCode: domain.PackageM, MaxWeightKg: 15.00, Description: "Medium package - up to 15 kg", IsActive: true},
			{SizeCode: domain.PackageL, MaxWeightKg: 25.00, Description: "Large package - up to 25 kg", IsActive: true},
		}
		for _, p := range pts {
			_ = database.Create(&p).Error
		}
	}

	orderRepo := repository.NewOrderGormRepo(database)
	userRepo := repository.NewUserGormRepo(database)
	userSvc := usecase.NewUserService(userRepo)
	ptRepo := repository.NewPackageTypeGormRepo(database)
	ptSvc := usecase.NewPackageTypeService(ptRepo)
	orderSvc := usecase.NewOrderService(orderRepo, ptSvc)
	addrRepo := repository.NewAddressGormRepo(database)
	addrSvc := usecase.NewAddressService(addrRepo)
	h := &httpdelivery.Handler{Orders: orderSvc, Users: userSvc, PackageTypes: ptSvc, Addresses: addrSvc}
	h.Register(r)
	log.Println("Bootstrap completed")
	return nil
}
