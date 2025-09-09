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
	Orders *usecase.OrderService
	Users  *usecase.UserService
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
	// Orders
	r.HandleFunc("/api/orders", h.CreateOrder).Methods(http.MethodPost)
	r.HandleFunc("/api/orders", h.MyOrders).Methods(http.MethodGet)
	r.HandleFunc("/api/admin/orders", h.AllOrders).Methods(http.MethodGet)
	r.HandleFunc("/api/admin/orders/{id}/status", h.UpdateStatus).Methods(http.MethodPatch)
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
// @Param request body object{email=string,password=string,role=string} true "User registration details"
// @Success 201 {object} domain.User "Created user"
// @Failure 400 {string} string "Bad request"
// @Router /users [post]
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string      `json:"email"`
		Password string      `json:"password"`
		Role     domain.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	u, err := h.Users.Register(body.Email, body.Password, body.Role)
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
	if err := h.Orders.Create(&o); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(201)
	_ = json.NewEncoder(w).Encode(o)
}

// MyOrders godoc
// @Summary Get user orders
// @Description Returns orders for authenticated user
// @Tags orders
// @Produce json
// @Param all query string false "Get all orders (admin only)"
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

// AllOrders godoc
// @Summary Get all orders
// @Description Returns all orders (admin only)
// @Tags admin
// @Produce json
// @Success 200 {array} domain.Order "List of all orders"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /admin/orders [get]
func (h *Handler) AllOrders(w http.ResponseWriter, r *http.Request) {
	_, role, ok := auth(r)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}
	if role != domain.RoleAdmin {
		http.Error(w, "forbidden", 403)
		return
	}
	orders, err := h.Orders.FindAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(orders)
}

// UpdateStatus godoc
// @Summary Update order status
// @Description Updates the status of an order (admin only)
// @Tags admin
// @Accept json
// @Param id path integer true "Order ID"
// @Param status body object{status=string} true "New status"
// @Success 204 "No content"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Security BearerAuth
// @Router /admin/orders/{id}/status [patch]
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
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
		Status domain.OrderStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := h.Orders.UpdateStatus(uint(id64), body.Status); err != nil {
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
