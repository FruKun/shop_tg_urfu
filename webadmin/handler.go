package webadmin

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/FruKun/shop_tg_bot_urfu/config"
	"github.com/FruKun/shop_tg_bot_urfu/db"
	"github.com/FruKun/shop_tg_bot_urfu/models"
)

type Handler struct {
	storage *db.Database
	config  *config.Config
}

func New(storage *db.Database, cfg *config.Config) *Handler {
	return &Handler{storage: storage, config: cfg}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/dashboard", h.basicAuth(h.Dashboard))
	mux.HandleFunc("/admin/add", h.basicAuth(h.AddForm))
	mux.HandleFunc("/admin/add/product", h.basicAuth(h.AddProduct))
	mux.HandleFunc("/admin/delete", h.basicAuth(h.DeleteProduct))

	fileServer := http.FileServer(http.Dir(h.config.UploadDir))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", fileServer))
}

func (h *Handler) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != h.config.AdminLogin || pass != h.config.AdminPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	products, err := h.storage.GetAllProduct()
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка загрузки товаров", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.New("dashboard").Parse(dashboardTemplate))
	tmpl.Execute(w, products)
}

func (h *Handler) AddForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("add").Parse(addTemplate))
	tmpl.Execute(w, nil)
}

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Файл слишком большой (макс. 10 МБ)", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "некорректная цена", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		http.Error(w, "некорректное количество", http.StatusBadRequest)
		return
	}

	var imageUrl string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			http.Error(w, "Допустимы только изображения (jpg, jpeg, png)", http.StatusBadRequest)
			return
		}

		fileName := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), rand.Int31n(999999), ext)
		filePath := filepath.Join(h.config.UploadDir, fileName)

		if err := os.MkdirAll(h.config.UploadDir, 0755); err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Ошибка копирования файла", http.StatusInternalServerError)
			return
		}

		imageUrl = "/uploads/" + fileName

	}

	product := &models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		ImageUrl:    imageUrl,
	}

	if err := h.storage.CreateProduct(product); err != nil {
		http.Error(w, "Ошибка создания товара", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	product, err := h.storage.GetProduct(id)
	if err == nil && product.ImageUrl != "" {
		imgPath := filepath.Join(h.config.UploadDir, filepath.Base(product.ImageUrl))
		os.Remove(imgPath)
	}

	if err := h.storage.DeleteProduct(id); err != nil {
		http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
