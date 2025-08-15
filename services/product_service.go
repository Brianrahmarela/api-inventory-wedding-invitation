package services

import (
	"api-go-invitation/models"
	"api-go-invitation/utils"
	"errors" //Package bawaan Go untuk bikin error manual
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// type ProductService struct → Struct yang dipakai untuk menyimpan dependency.
type ProductService struct {
	DB *gorm.DB
	// DB *gorm.DB → Pointer ke koneksi database GORM. Jadi semua method di struct ini bisa akses DB
}

// Fungsi constructor → Biar gampang bikin ProductService baru dengan DB yang udah terkoneksi.
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db} // Buat struct baru dan set field DB dengan parameter db.
}

// func (ps *ProductService) → Method ini “nempel” ke ProductService.
func (ps *ProductService) GetAll() ([]models.Product, error) {
	var products []models.Product
	//ps.DB.Find(&products) → Query SELECT * FROM products.
	if err := ps.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetPaginated returns products for a given page & limit with sorting and optional q filter, plus total count
func (ps *ProductService) GetPaginated(page, limit int, sort, order, q string) ([]models.Product, int64, error) {
	fmt.Println("GetPaginated called with page:", page, "limit:", limit, "sort:", sort, "order:", order, "q:", q)
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	fmt.Println("Offset:", offset)

	// whitelist kolom yang boleh di-sort untuk mencegah SQL injection
	allowedSort := map[string]bool{
		"id":         true,
		"name":       true,
		"price":      true,
		"stock":      true,
		"created_at": true,
		"updated_at": true,
	}

	sort = strings.ToLower(sort)
	//Kalau kolom sort tidak ada di daftar → default jadi "id"
	if !allowedSort[sort] {
		sort = "id"
	}

	order = strings.ToLower(order)
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	orderClause := fmt.Sprintf("%s %s", sort, order)
	fmt.Println("Order clause:", orderClause)

	// build base query
	dbQuery := ps.DB.Model(&models.Product{})
	fmt.Println("dbQuery query before:", dbQuery)
	if q != "" {
		like := "%" + q + "%"
		dbQuery = dbQuery.Where("name LIKE ?", like)
	}
	fmt.Println("dbQuery query after:", dbQuery)

	var total int64
	// Count(&total) → Hitung jumlah baris yang cocok dengan filter.
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []models.Product
	//.Order(orderClause) → Urutkan sesuai sort & order.
	//.Limit(limit) → Ambil sejumlah limit.
	//.Offset(offset) → Lewati offset data.
	//.Find(&products) → Eksekusi query.
	if err := dbQuery.Order(orderClause).Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	//products -> Data produk yang diambil (per halaman).
	//total -> jumlah semua data yang cocok (tanpa limit & offset)
	return products, total, nil
}

func (ps *ProductService) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := ps.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (ps *ProductService) GetBySlug(slug string, page, limit int, sort, order string) ([]models.Product, int64, error) {
	// Default handling untuk sort/order sudah di controller
	if sort == "" {
		sort = "id"
	}
	if order == "" {
		order = "asc"
	}

	offset := (page - 1) * limit

	var products []models.Product
	query := ps.DB.Where("slug LIKE ?", "%"+slug+"%").
		Order(fmt.Sprintf("%s %s", sort, order)).
		Limit(limit).
		Offset(offset)

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := ps.DB.Model(&models.Product{}).Where("slug LIKE ?", "%"+slug+"%").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (ps *ProductService) Create(req *models.CreateProductRequest) (*models.Product, error) {
	slug := utils.GenerateSlug(req.Name)

	// pastikan slug unik
	var exist models.Product
	if err := ps.DB.Where("slug = ?", slug).First(&exist).Error; err == nil {
		return nil, errors.New("product with this name already exists")
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		Slug:        slug,
	}
	if err := ps.DB.Create(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (ps *ProductService) Update(id uint, req *models.UpdateProductRequest) (*models.Product, error) {
	product, err := ps.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		product.Name = req.Name
		product.Slug = utils.GenerateSlug(req.Name)
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.Stock != 0 {
		product.Stock = req.Stock
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if err := ps.DB.Save(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (ps *ProductService) Delete(id uint) error {
	if err := ps.DB.Delete(&models.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}
