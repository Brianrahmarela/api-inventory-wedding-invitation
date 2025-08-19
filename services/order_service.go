package services

import (
	"api-go-invitation/models"
	"fmt"
	"mime/multipart"
	"net/url"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type OrderService struct {
	DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{DB: db}
}

func (s *OrderService) CreateOrderWithGuests(
	userID uint,
	productID uint,
	groomName string,
	brideName string,
	guestFile multipart.File,
) (*models.Order, error) {
	//(*models.Order) -> return order yang berhasil dibuat (punya ID, dll)

	tx := s.DB.Begin()
	defer func() {
		//recover() menangkap panic. jika ada, di rollback supaya DB tidak jadi setengah tersimpan.
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Cek stok produk
	var product models.Product
	fmt.Println("models.Product", product)
	if err := tx.First(&product, productID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("product not found")
	}
	if product.Stock <= 0 {
		tx.Rollback()
		return nil, fmt.Errorf("product out of stock")
	}
	fmt.Println("Product after first by productID ->", product)

	// 2. Buat order
	order := models.Order{
		UserID:      userID,
		ProductID:   productID,
		TotalAmount: int64(product.Price),
		Status:      "pending",
		GroomName:   groomName,
		BrideName:   brideName,
	}
	//tx.Create(&order) -> INSERT ke tabel orders, Karena ngirim pointer &order, GORM menuliskan balik
	// value yang dihasilkan DB ke variabel order line 55 di memori:
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 3. Parse Excel tamu
	f, err := excelize.OpenReader(guestFile)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to read excel: %v", err)
	}
	//f.GetRows("Sheet1"): ambil semua baris di sheet bernama “Sheet1” dalam bentuk [][]string.
	//Baris 0 → header (Nama | Partner | Email | Phone). Baris 1..n → data tamu.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to read rows: %v", err)
	}
	// Minimal 1 baris data list tamu (selain header)
	if len(rows) < 2 {
		tx.Rollback()
		return nil, fmt.Errorf("excel file must have at least 1 guest row")
	}

	// Validasi header (baris pertama)
	expectedHeaders := []string{"Nama", "Partner", "Email", "Phone"}
	//rows[0] adalah header. cek jml col min 4,
	if len(rows[0]) < len(expectedHeaders) {
		tx.Rollback()
		return nil, fmt.Errorf("excel header format invalid, expected: %v", expectedHeaders)
	}
	//cek tiap elemen header
	for i, header := range expectedHeaders {
		//strings.TrimSpace(...) -> menghapus spasi di depan dan akhir string, ambil isi header di file Excel, bersihkan spasi, ubah ke lower case.
		//!= strings.ToLower -> Ambil cal header dari expectedHeaders (mis "Nama"), ubah ke lower case. bandingkan "nama" (expected) vs hasil parsing header (harusnya "nama").
		if strings.TrimSpace(strings.ToLower(rows[0][i])) != strings.ToLower(header) {
			tx.Rollback()
			return nil, fmt.Errorf("invalid header in column %d, expected '%s'", i+1, header)
		}
	}

	// Slug nama mempelai
	groomSlug := strings.ReplaceAll(strings.ToLower(groomName), " ", "-")
	brideSlug := strings.ReplaceAll(strings.ToLower(brideName), " ", "-")

	// Loop data tamu
	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		// Pastikan minimal 4 kolom
		if len(row) < 4 {
			tx.Rollback()
			return nil, fmt.Errorf("row %d format invalid, must have columns: Nama, Partner, Email, Phone", i+1)
		}

		guestName := strings.TrimSpace(row[0])
		partnerName := strings.TrimSpace(row[1])
		email := strings.TrimSpace(row[2])
		phone := strings.TrimSpace(row[3])

		if guestName == "" {
			tx.Rollback()
			return nil, fmt.Errorf("row %d: Nama cannot be empty", i+1)
		}

		fullName := guestName
		if partnerName != "" {
			fullName = fmt.Sprintf("%s & %s", guestName, partnerName)
		}

		// 1) encode pakai QueryEscape (meng-encode & -> %26, tapi spasi -> +)
		// 2) ubah + menjadi %20 supaya spasi tampil sebagai %20
		escaped := url.QueryEscape(fullName)
		escaped = strings.ReplaceAll(escaped, "+", "%20")

		link := fmt.Sprintf(
			"https://app.inviteable.id/%s-%s/?to=%s",
			groomSlug,
			brideSlug,
			escaped,
		)

		guest := models.Guest{
			//Field OrderID diisi dengan order.ID (primary key order yang baru dibuat olehtx.Create(&order)). Ini menghubungkan guest ke order (foreign key).
			OrderID:   order.ID,
			Name:      guestName,
			Partner:   partnerName,
			Email:     email,
			Phone:     phone,
			Link:      link,
			CreatedAt: time.Now(),
		}
		//tx.Create(&guest) -> INSERT record guest ke tabel guests melalui transaksi tx
		if err := tx.Create(&guest).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to insert guest at row %d: %v", i+1, err)
		}
	}
	fmt.Println("Product before tx.save", product)
	// 4. Kurangi stok
	product.Stock -= 1
	//tx.Save(&product) -> Jika primary key ada → GORM akan melakukan UPDATE; jika tidak ada → INSERT.
	// Di konteks ini primary key sudah ada sehingga akan UPDATE.
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	fmt.Println("Product after tx.save", product)

	// 5. Commit
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	// &order-> alamat pointer ke order struct (return alamat pointer *models.Order yang dihasilkan).
	return &order, nil
}
