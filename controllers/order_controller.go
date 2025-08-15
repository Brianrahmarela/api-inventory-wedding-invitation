package controllers

import (
	"api-go-invitation/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	OrderService *services.OrderService
}

// Bikin OrderController dan sekaligus kasih dia akses ke DB lewat OrderService
func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{
		OrderService: services.NewOrderService(db),
	}
}

func (oc *OrderController) CreateOrderHandler(c *gin.Context) {
	// Ambil userId dari JWT middleware
	userIDValue, exists := c.Get("userId")
	// fmt.Println("userIDValue:", userIDValue, "exists:", exists)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	//userIDValue.(uint) -> cek apakah userIDValue benar-benar uint (angka tanpa minus).
	//type assertion untuk interface{} → memeriksa dan mengambil nilai asli jika tipe cocok.
	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	// fmt.Println("userID valid!", userID)

	productID := c.PostForm("product_id")
	groomName := c.PostForm("groom_name")
	brideName := c.PostForm("bride_name")
	fmt.Println("productID", productID, "groomName", groomName, "brideName", brideName)

	// Validasi form input
	if productID == "" || groomName == "" || brideName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id, groom_name, and bride_name are required"})
		return
	}

	// Ambil file excel tamu
	file, err := c.FormFile("guest_file")
	// fmt.Println("file", file, "err", err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "guest_file is required"})
		return
	}
	//file.Open() -> buka file upload agar bisa dibaca.
	openedFile, err := file.Open()
	// fmt.Println("openedFile", openedFile, "err", err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read guest file"})
		return
	}
	defer openedFile.Close()
	// Convert string product_id ke uint
	var pid uint
	if pid64, err := strconv.ParseUint(productID, 10, 64); err == nil {
		//convert uint64 keuint32 atau uint64 agar flexible
		pid = uint(pid64)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}
	//Variabel order ini isinya alamat pointer dari *models.Order. Agar menerima objek yang sama yang barusan diisi oleh GORM.
	// Bukan salinan, tapi referensi yg lebih ringan krn tdk duplikasi data order baru
	order, svcErr := oc.OrderService.CreateOrderWithGuests(
		userID,
		pid,
		groomName,
		brideName,
		openedFile,
	)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}
	//gin.H adalah alias untuk map[string]interface{} untuk membentuk objek respons sebelum Gin kirim ke client.
	//meskipun struct order bisa diakses langsung tanpa *, namun di client yg dikirim var order adalah alamat pointer & client tdk bisa akses value krn memorinya di go saja
	//gih.H otomatis jalanin encoding/json, Dia “masuk” lewat alamat pointer (dereference) → baca isi struct dari memori Go.
	//Kemudian dia ubah semua data itu menjadi teks JSON agar bisa dibaca client
	c.JSON(http.StatusCreated, gin.H{
		"message": "order created successfully",
		"data":    order,
	})
}
