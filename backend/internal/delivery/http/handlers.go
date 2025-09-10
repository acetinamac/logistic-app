package http

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"logistics-app/backend/internal/domain"
	"logistics-app/backend/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type Handler struct {
	Orders       *usecase.OrderService
	Users        *usecase.UserService
	PackageTypes *usecase.PackageTypeService
	Addresses    *usecase.AddressService
}

type claims struct {
	UserID uint        `json:"uid"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

func jwtKey() []byte { return []byte(getenv("JWT_SECRET", "dev_secret")) }
func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/api/login", h.Login).Methods(http.MethodPost)
	// Users
	r.HandleFunc("/api/users", h.RegisterUser).Methods(http.MethodPost)
	r.HandleFunc("/api/users/{id}", h.DeleteUser).Methods(http.MethodDelete)
	// Package Types
	r.HandleFunc("/api/package-types", h.ListPackageTypes).Methods(http.MethodGet)
	r.HandleFunc("/api/package-types/{id}/active", h.SetPackageTypeActive).Methods(http.MethodPatch)
	// Addresses
	r.HandleFunc("/api/addresses", h.CreateAddress).Methods(http.MethodPost)
	r.HandleFunc("/api/addresses", h.ListAddresses).Methods(http.MethodGet)
	r.HandleFunc("/api/addresses/{id}", h.GetAddress).Methods(http.MethodGet)
	r.HandleFunc("/api/addresses/{id}", h.UpdateAddress).Methods(http.MethodPut)
	r.HandleFunc("/api/addresses/{id}", h.DeleteAddress).Methods(http.MethodDelete)
	r.HandleFunc("/api/addresses/{id}/active", h.SetAddressActive).Methods(http.MethodPatch)
	// Orders
	r.HandleFunc("/api/orders", h.CreateOrder).Methods(http.MethodPost)
	r.HandleFunc("/api/orders", h.MyOrders).Methods(http.MethodGet)
	r.HandleFunc("/api/orders/{id}/status", h.UpdateStatus).Methods(http.MethodPatch)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200); _, _ = w.Write([]byte("ok")) }).Methods(http.MethodGet)
}

// Login godoc
// @Summary Login endpoint
// @Description Authenticates user and returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string} true "Login credentials"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Not authorized"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	u, err := h.Users.Authenticate(body.Email, body.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	now := time.Now()
	cl := &claims{
		UserID: u.ID,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 2)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	s, err := token.SignedString(jwtKey())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"token": s})
}

// RegisterUser godoc
// @Summary Register new user
// @Description Creates a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body object{email=string,password=string,full_name=string,phone=string,role=string} true "User registration details"
// @Success 201 {object} domain.User "Created user"
// @Failure 400 {string} string "Bad request"
// @Router /users [post]
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string      `json:"email"`
		Password string      `json:"password"`
		FullName string      `json:"full_name"`
		Phone    string      `json:"phone"`
		Role     domain.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	u, err := h.Users.Register(body.Email, body.Password, body.FullName, body.Phone, body.Role)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(201)
	_ = json.NewEncoder(w).Encode(u)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Deletes a user account (admin or own account only)
// @Tags users
// @Param id path integer true "User ID"
// @Success 204 "No content"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	uid, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	if role != domain.RoleAdmin && uid != id {
		http.Error(w, "forbidden", 403)
		return
	}
	if err := h.Users.Delete(id); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(204)
}

// CreateOrder godoc
// @Summary Create new order
// @Description Creates a new order for authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Param order body domain.Order true "Order details"
// @Success 201 {object} domain.Order "Created order"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Security BearerAuth
// @Router /orders [post]
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	uid, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	if role != domain.RoleClient && role != domain.RoleAdmin {
		http.Error(w, "forbidden", 403)
		return
	}
	var o domain.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	o.CustomerID = uid
	o.CreatedBy = uid
	o.UpdatedBy = &uid
	if err := h.Orders.Create(&o); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(201)
	_ = json.NewEncoder(w).Encode(o)
}

// Orders list godoc
// @Summary List orders
// @Description Clients see only their orders. Admins can see all orders by passing ?all=1.
// @Tags orders
// @Produce json
// @Param all query string false "If set to 1 and requester is admin, returns all orders; otherwise returns only own orders"
// @Success 200 {array} domain.Order "List of orders"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /orders [get]
func (h *Handler) MyOrders(w http.ResponseWriter, r *http.Request) {
	uid, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	if role != domain.RoleClient && role != domain.RoleAdmin {
		http.Error(w, "forbidden", 403)
		return
	}
	var (
		orders []domain.Order
		err    error
	)
	if role == domain.RoleAdmin && r.URL.Query().Get("all") == "1" {
		orders, err = h.Orders.FindAll()
	} else {
		orders, err = h.Orders.FindByCustomer(uid)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(orders)
}

// UpdateStatus godoc
// @Summary Update order status
// @Description Updates the status of an order (admin only)
// @Tags orders
// @Accept json
// @Param id path integer true "Order ID"
// @Param status body object{status=string} true "New status"
// @Success 204 "No content"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Security BearerAuth
// @Router /orders/{id}/status [patch]
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	uid, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	if role != domain.RoleAdmin {
		http.Error(w, "forbidden", 403)
		return
	}
	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	var body struct {
		Status domain.OrderStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := h.Orders.UpdateStatus(uint(id64), body.Status, uid); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(204)
}

// ListPackageTypes godoc
// @Summary List package types
// @Description Returns package types. If ?all=1 and requester is admin, includes inactive; otherwise only active.
// @Tags package_types
// @Produce json
// @Param all query string false "If set to 1 and requester is admin, returns active and inactive"
// @Success 200 {array} domain.PackageType
// @Failure 401 {string} string "Unauthorized"
// @Security BearerAuth
// @Router /package-types [get]
func (h *Handler) ListPackageTypes(w http.ResponseWriter, r *http.Request) {
	_, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	includeInactive := false
	if role == domain.RoleAdmin && r.URL.Query().Get("all") == "1" {
		includeInactive = true
	}
	list, err := h.PackageTypes.List(includeInactive)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(list)
}

// SetPackageTypeActive godoc
// @Summary Set PackageType active status
// @Description Admin only. Sets is_active true/false for a PackageType
// @Tags package_types
// @Accept json
// @Param id path integer true "PackageType ID"
// @Param request body object{active=boolean} true "Desired active state"
// @Success 204 "No content"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Security BearerAuth
// @Router /package-types/{id}/active [patch]
func (h *Handler) SetPackageTypeActive(w http.ResponseWriter, r *http.Request) {
	_, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	if role != domain.RoleAdmin {
		http.Error(w, "forbidden", 403)
		return
	}
	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	var body struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := h.PackageTypes.ToggleActive(uint(id64), body.Active); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(204)
}

// Addresses Handlers
// CreateAddress
// @Summary Create address (with coordinates)
// @Description Creates coordinates first if provided, then address; CustomerID is set from JWT. Clients create their own; admins can also create for themselves only in this endpoint.
// @Tags addresses
// @Accept json
// @Produce json
// @Param request body object{street=string,exterior_number=string,interior_number=string,neighborhood=string,postal_code=string,city=string,state=string,country=string,coordinates=object{latitude=number,longitude=number}} true "Address with optional coordinates"
// @Success 201 {object} domain.Address
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Security BearerAuth
// @Router /addresses [post]
func (h *Handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	uid, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	if role != domain.RoleClient && role != domain.RoleAdmin {
		http.Error(w, "forbidden", 403)
		return
	}
	var req usecase.AddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	addr, _, err := h.Addresses.Create(uid, req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(201)
	_ = json.NewEncoder(w).Encode(addr)
}

// ListAddresses
// @Summary List addresses
// @Description Clients see only their addresses (active by default). Admins can pass ?all=1 to see all, and ?include_inactive=1 to include inactive.
// @Tags addresses
// @Produce json
// @Param all query string false "Admin only: if set to 1, list all users' addresses"
// @Param include_inactive query string false "Admin only: if set to 1, include inactive addresses"
// @Success 200 {array} domain.Address
// @Failure 401 {string} string "Unauthorized"
// @Security BearerAuth
// @Router /addresses [get]
func (h *Handler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	uid, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	includeInactive := role == domain.RoleAdmin && r.URL.Query().Get("include_inactive") == "1"
	all := role == domain.RoleAdmin && r.URL.Query().Get("all") == "1"
	list, err := h.Addresses.List(uid, role, includeInactive, all)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(list)
}

// GetAddress
// @Summary Get single address
// @Tags addresses
// @Produce json
// @Param id path integer true "Address ID"
// @Success 200 {object} domain.Address
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Security BearerAuth
// @Router /addresses/{id} [get]
func (h *Handler) GetAddress(w http.ResponseWriter, r *http.Request) {
	uid, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	isAdmin := role == domain.RoleAdmin
	a, err := h.Addresses.Get(uid, isAdmin, uint(id64))
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	_ = json.NewEncoder(w).Encode(a)
}

// UpdateAddress
// @Summary Update address (and coordinates)
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path integer true "Address ID"
// @Param request body object{street=string,exterior_number=string,interior_number=string,neighborhood=string,postal_code=string,city=string,state=string,country=string,is_active=boolean,coordinates=object{latitude=number,longitude=number}} true "Address update"
// @Success 200 {object} domain.Address
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Security BearerAuth
// @Router /addresses/{id} [put]
func (h *Handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	uid, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	var req usecase.AddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	addr, _, err := h.Addresses.Update(uid, role == domain.RoleAdmin, uint(id64), req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(addr)
}

// DeleteAddress
// @Summary Delete address
// @Description Deletes address only if it is not referenced by any order. Owner or admin only.
// @Tags addresses
// @Param id path integer true "Address ID"
// @Success 204 "No content"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Security BearerAuth
// @Router /addresses/{id} [delete]
func (h *Handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	uid, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	if err := h.Addresses.Delete(uid, role == domain.RoleAdmin, uint(id64)); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(204)
}

// SetAddressActive
// @Summary Set Address active status
// @Tags addresses
// @Accept json
// @Param id path integer true "Address ID"
// @Param request body object{active=boolean} true "Desired active state"
// @Success 204 "No content"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Security BearerAuth
// @Router /addresses/{id}/active [patch]
func (h *Handler) SetAddressActive(w http.ResponseWriter, r *http.Request) {
	uid, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	var body struct {
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := h.Addresses.ToggleActive(uid, role == domain.RoleAdmin, uint(id64), body.Active); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(204)
}

func auth(r *http.Request) (uint, domain.Role, bool) {
	h := r.Header.Get("Authorization")
	if h == "" || !strings.HasPrefix(h, "Bearer ") {
		return 0, "", false
	}
	tokStr := strings.TrimPrefix(h, "Bearer ")
	token, err := jwt.ParseWithClaims(tokStr, &claims{}, func(t *jwt.Token) (interface{}, error) { return jwtKey(), nil })
	if err != nil || !token.Valid {
		return 0, "", false
	}
	cl, ok := token.Claims.(*claims)
	if !ok {
		return 0, "", false
	}
	return cl.UserID, cl.Role, true
}
