package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

// ProductHandler is a struct to handle product requests
type ProductHandler struct {
	productService         services.ProductService         // ProductService is a service for product operations
	deliveryService        services.DeliveryService        // DeliveryService is a service for delivery operations
	productCategoryService services.ProductCategoryService // ProductCategoryService is a service for product category operations
	companyService         services.CompanyService         // CompanyService is a service for company operations
	userCompanyService     services.UserCompanyService     // UserCompanyService is a service for user company operations
	LogService             services.LogService             // LogService is a service for logging operations
}

// NewProductHandler creates a new ProductHandler
// productService: ProductService is a service for product operations
// deliveryService: DeliveryService is a service for delivery operations
// productCategoryService: ProductCategoryService is a service for product category operations
// companyService: CompanyService is a service for company operations
// userCompanyService: UserCompanyService is a service for user company operations
// returns a new ProductHandler
func NewProductHandler(
	productService services.ProductService,
	deliveryService services.DeliveryService,
	productCategoryService services.ProductCategoryService,
	companyService services.CompanyService,
	userCompanyService services.UserCompanyService,
	logService services.LogService,
) *ProductHandler {
	return &ProductHandler{
		productService:         productService,
		deliveryService:        deliveryService,
		productCategoryService: productCategoryService,
		companyService:         companyService,
		userCompanyService:     userCompanyService,
		LogService:             logService,
	}
}

// CreateProductRequest is a struct to bind request body for creating a product
type CreateProductRequest struct {
	// Name is the name of the product
	// example: Apple
	Name string `gorm:"size:255;not null;column:name"`

	// Weight is the weight of the product
	// example: 0.5
	Weight float64 `gorm:"not null;column:weight"`

	// ProductType is the type of the product
	// example: Fruits
	ProductType string `json:"product_type" example:"Fruits | Vegetables | Frozen Foods | Dairy Products | Meat"`

	// DeliveryID is the ID of the delivery
	// example: 1
	DeliveryID uint `gorm:"not null;column:delivery_id"`
}

// UpdateProductRequest is a struct to bind request body for updating a product
type UpdateProductRequest struct {
	// Name is the name of the product
	// example: Apple
	Name string `gorm:"size:255;not null;column:name"`

	// Weight is the weight of the product
	// example: 0.5
	Weight float64 `gorm:"not null;column:weight"`

	// ProductType is the type of the product
	// example: Fruits
	ProductType string `json:"product_type"`
}

// CreateProduct godoc
// @Summary      Create a product
// @Description  Create a product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        request body CreateProductRequest true "CreateProductRequest"
// @Security     BearerAuth
// @Router       /products [post]
func (h *ProductHandler) AddProduct(c *gin.Context) {
	var productRequest CreateProductRequest
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "create_product", "Failed to create product: "+err.Error(), false)

		return
	}

	delivery, err := h.deliveryService.GetByID(context.Background(), productRequest.DeliveryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		h.logAction(c, "create_product", "Failed to create product: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: delivery.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		h.logAction(c, "create_product", "Failed to create product: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	productCategory, err := h.productCategoryService.Where(context.Background(), &models.ProductCategory{
		Name: productRequest.ProductType,
	})
	if err != nil || len(productCategory) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product category not found"})
		h.logAction(c, "create_product", "Failed to create product: "+err.Error(), false)

		return
	}

	if productRequest.Weight <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid weight"})
		return
	}

	product := &models.Product{
		Name:              productRequest.Name,
		Weight:            productRequest.Weight,
		ProductCategoryID: productCategory[0].ID,
		DeliveryID:        productRequest.DeliveryID,
	}

	if err := h.productService.Create(context.Background(), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_product", "Failed to create product: "+err.Error(), false)

		return
	}

	productDTO := &dtos.ProductDTO{}
	if err = dtoMapper.Map(productDTO, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "create_product", "Failed to create product: "+err.Error(), false)

		return
	}

	h.logAction(c, "create_product", "Product '"+productRequest.Name+"' created successfully", true)

	c.JSON(http.StatusOK, productDTO)
}

// GetProduct godoc
// @Summary      Get a product
// @Description  Get a product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        product_id path int true "Product ID"
// @Security     BearerAuth
// @Router       /products/{product_id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID format"})
		h.logAction(c, "get_product", "Failed to fetch product: "+err.Error(), false)

		return
	}

	product, err := h.productService.GetByID(context.Background(), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_product", "Failed to fetch product: "+err.Error(), false)

		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "get_product", "Failed to fetch product: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: product.Delivery.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		h.logAction(c, "get_product", "Failed to fetch product: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		h.logAction(c, "get_product", "Failed to fetch product: "+err.Error(), false)

		return
	}

	productDTO := &dtos.ProductDTO{}
	if err = dtoMapper.Map(productDTO, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "get_product", "Failed to fetch product: "+err.Error(), false)

		return
	}

	h.logAction(c, "get_product", fmt.Sprintf("Product %d fetched successfully", productID), true)

	c.JSON(http.StatusOK, productDTO)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update a product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        product_id path int true "Product ID"
// @Param        request body UpdateProductRequest true "UpdateProductRequest"
// @Security     BearerAuth
// @Router       /products/{product_id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var productRequest UpdateProductRequest
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.logAction(c, "update_product", "Failed to update product: "+err.Error(), false)

		return
	}

	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID format"})
		h.logAction(c, "update_product", "Failed to update product: "+err.Error(), false)

		return
	}

	product, err := h.productService.GetByID(context.Background(), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_product", "Failed to update product: "+err.Error(), false)

		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "update_product", "Failed to update product: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: product.Delivery.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		h.logAction(c, "update_product", "Failed to update product: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if productRequest.ProductType != "" {
		productCategory, err := h.productCategoryService.Where(context.Background(), &models.ProductCategory{
			Name: productRequest.ProductType,
		})
		if err != nil || len(productCategory) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product category not found"})
			return
		}
		product.ProductCategoryID = productCategory[0].ID
	}
	if productRequest.Name != "" {
		product.Name = productRequest.Name
	}
	if productRequest.Weight != 0 {
		product.Weight = productRequest.Weight
	}

	product.Delivery = models.Delivery{}
	product.ProductCategory = models.ProductCategory{}

	if err := h.productService.Update(context.Background(), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_product", "Failed to update product: "+err.Error(), false)

		return
	}

	productDTO := &dtos.ProductDTO{}
	if err = dtoMapper.Map(productDTO, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "update_product", "Failed to update product: "+err.Error(), false)

		return
	}

	h.logAction(c, "update_product", fmt.Sprintf("Product %d updated successfully", productID), true)

	c.JSON(http.StatusOK, productDTO)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete a product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        product_id path int true "Product ID"
// @Security     BearerAuth
// @Router       /products/{product_id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID format"})
		h.logAction(c, "delete_product", "Failed to delete product: "+err.Error(), false)

		return
	}

	product, err := h.productService.GetByID(context.Background(), uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "delete_product", "Failed to delete product: "+err.Error(), false)

		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		h.logAction(c, "delete_product", "Failed to delete product: "+err.Error(), false)

		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: product.Delivery.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		h.logAction(c, "delete_product", "Failed to delete product: "+err.Error(), false)

		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := h.productService.Delete(context.Background(), uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logAction(c, "delete_product", "Failed to delete product: "+err.Error(), false)

		return
	}

	h.logAction(c, "delete_product", fmt.Sprintf("Product %d deleted successfully", productID), true)

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (h *ProductHandler) logAction(c *gin.Context, actionType, description string, success bool) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		userID = nil
	}
	h.LogService.LogAction(userID, actionType, description, success)
}
