package controllers

import (
	"api-go-test/models"
	"api-go-test/services"

	// "fmt"
	"net/http"
	"strconv" //konversi string ⇄ angka
	"strings" //fungsi manipulasi string (Trim, Lowercase, dsb).

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Tugasnya melayani permintaan: menambah produk, melihat daftar produk, menghapus produk, dll.
// si pelayan (ProductController) bisa menjalankan pekerjaannya karena dibekali asisten (service) yang mengerjakan bagian teknis ngobrol ke database
type ProductController struct {
	//ProductService (yang tipenya pointer ke services.ProductService), ProductService adalah dependency = hal yang wajib dimiliki ProductController agar ProductController bisa jalan
	ProductService *services.ProductService
}

// constructor function untuk: Membuat ProductService baru (services.NewProductService(db)),
// Mengisi field ProductService di struct ProductController dengan hasil tadi. Mengembalikan pointer *ProductController.
func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{ProductService: services.NewProductService(db)}
}

// memastikan urutan json: meta dulu lalu data
type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type Data struct {
	Content      []models.Product `json:"content"`
	TotalContent int              `json:"total_content"` // jml data di halaman yg tampil setelah pagination (LIMIT dan OFFSET)
	TotalData    int64            `json:"total_data"`    // jml semua data yang match filter slug tersebut di database
	Page         int              `json:"page"`
	Limit        int              `json:"limit"`
}

type Response struct {
	Meta Meta `json:"meta"`
	Data Data `json:"data"`
}

func (pc *ProductController) GetAll(c *gin.Context) {
	// fmt.Println("GetAll c", c)
	// gunakan c.Query agar bisa mendeteksi jika param diberikan tapi kosong (?page=)
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	sort := c.Query("sort")
	order := c.Query("order")
	q := c.Query("q") // search by name

	// default values ketika param kosong atau tidak disertakan
	page := 1
	limit := 10
	// Menghapus semua spasi di awal dan akhir string sort
	if strings.TrimSpace(sort) == "" {
		sort = "id"
	}
	// default order asc jika order kosong
	if strings.TrimSpace(order) == "" {
		order = "asc"
	}
	// parse page jika diberikan
	if pageStr != "" {
		// strconv.Atoi mengubah string dari query parameter (misalnya "1") menjadi integer (3). Atoi = ASCII to Integer.
		// ubah pageStr jadi integer p. Jika konversi berhasil (err == nil) dan nilainya lebih besar dari 0 (p > 0),
		// maka pakai p sebagai nilai page
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// parse limit jika diberikan
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if limit > 100 {
		limit = 100
	}

	// normalize sort & order (service juga punya whitelist)
	sort = strings.ToLower(sort)
	order = strings.ToLower(order)

	products, totalData, err := pc.ProductService.GetPaginated(page, limit, sort, order, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// build response struct - meta dulu, kemudian data
	meta := Meta{
		Message: "success",
		Code:    200,
		Status:  "200",
	}
	data := Data{
		Content:      products,
		TotalContent: len(products),
		TotalData:    totalData,
		Page:         page,
		Limit:        limit,
	}

	resp := Response{
		Meta: meta,
		Data: data,
	}

	c.JSON(http.StatusOK, resp)
}

func (pc *ProductController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := pc.ProductService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (pc *ProductController) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	// Ambil params dari query
	pageStr := c.DefaultQuery("page", "1")
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")
	limitStr := c.DefaultQuery("limit", "10")

	// Konversi page & limit ke int
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	products, totalData, err := pc.ProductService.GetBySlug(slug, page, limit, sort, order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	meta := Meta{
		Message: "success",
		Code:    200,
		Status:  "200",
	}
	data := Data{
		Content:      products,
		TotalContent: len(products),
		TotalData:    totalData,
		Page:         page,
		Limit:        limit,
	}

	resp := Response{
		Meta: meta,
		Data: data,
	}

	c.JSON(http.StatusOK, resp)
}

func (pc *ProductController) Create(c *gin.Context) {
	role := c.GetString("role")
	// fmt.Println("User role:", role)
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := pc.ProductService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (pc *ProductController) Update(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := pc.ProductService.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (pc *ProductController) Delete(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	if err := pc.ProductService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
