package productcontroller

import (
	"archive/zip"
	"be-shop-vision/dto"
	productdto "be-shop-vision/dto/product"
	usecase "be-shop-vision/usecase/product_usecase"
	"be-shop-vision/util"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/middleware"
	"github.com/susatyo441/go-ta-utils/parser"
	"github.com/susatyo441/go-ta-utils/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type makeProductUseCaseFunc = func() usecase.IProductUseCase

type ProductController struct {
	UseCase             usecase.IProductUseCase
	MakeUseCaseFunction makeProductUseCaseFunc
}

func MakeProductController(makeUseCaseFunc makeProductUseCaseFunc) *ProductController {
	return &ProductController{MakeUseCaseFunction: makeUseCaseFunc}
}

// CreateProduct godoc
// @Summary Create Product
// @Description Create Product
// @Tags Product
// @Produce  json
// @Router /api/product [post]
// @Param payload body productdto.CreateProductDTO true "Payload to create"
// @Security BearerAuth
func (ctrl *ProductController) CreateProduct(ctx *fiber.Ctx) error {

	// Ambil semua file dari form
	files := map[string]*multipart.FileHeader{}
	attributes := []string{"image1", "image2", "image3", "image4", "image5"}
	for _, attr := range attributes {
		fileHeader, err := ctx.FormFile(attr)
		if err == nil { // File ada
			files[attr] = fileHeader
		}
	}

	// Validasi jumlah file harus pas 5
	if len(files) != 5 {
		return response.BadRequest(ctx, "Jumlah foto harus 5", nil)
	}

	var payload productdto.CreateProductDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}
	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	storeId := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	err := ctrl.UseCase.CreateProduct(ctx.Context(), payload, storeId, files)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Created(ctx, "Successfully create product", nil)
}

// BulkDeleteProducts godoc
// @Summary Bulk Delete Products
// @Description Bulk Delete Products
// @Tags Product
// @Produce  json
// @Router /api/product [delete]
// @Param payload body dto.ArrayOfIdDTO true "Payload to bulk delete"
// @Security BearerAuth
func (ctrl *ProductController) BulkDeleteProducts(ctx *fiber.Ctx) error {
	var payload dto.ArrayOfIdDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}

	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	storeId := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	err := ctrl.UseCase.BulkDeleteProducts(ctx.Context(), payload, storeId)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Success(ctx, "Successfully delete products", nil)
}

// UpdateProduct godoc
// @Summary Update Products
// @Description Update Products
// @Tags Product
// @Produce  json
// @Router /api/product/{productId} [put]
// @Param productId path string true "product ID"
// @Param payload body productdto.UpdateProductDTO true "Payload to update"
// @Security BearerAuth
func (ctrl *ProductController) UpdateProduct(ctx *fiber.Ctx) error {
	// Ambil productId dari parameter
	productId, paramErr := functions.ParamToObjectID(ctx, "productId")
	if paramErr != nil {
		return response.BadRequest(ctx, "Invalid product id format", nil)
	}

	// Parsing payload
	var payload productdto.UpdateProductDTO
	if err := ctx.BodyParser(&payload); err != nil {
		return response.BadRequest(ctx, "Invalid request body", nil)
	}

	// Validasi payload
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}

	// Ambil storeId dari context
	storeId, ok := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)
	if !ok {
		return response.BadRequest(ctx, "Invalid store ID", nil)
	}

	files := map[string]*multipart.FileHeader{}
	attributes := []string{"image1", "image2", "image3", "image4", "image5"}
	for _, attr := range attributes {
		if fileHeader, err := ctx.FormFile(attr); err == nil {
			files[attr] = fileHeader
		}
	}

	// Panggil UseCase untuk update produk
	ctrl.UseCase = ctrl.MakeUseCaseFunction()
	err := ctrl.UseCase.UpdateProduct(ctx.Context(), productId, payload, storeId, files)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Success(ctx, "Successfully updated product", nil)
}

// UpdateProduct godoc
// @Summary Update Products Stock
// @Description Update Products Stock
// @Tags Product
// @Produce  json
// @Router /api/product/{productId}/stock [put]
// @Param productId path string true "product ID"
// @Param payload body productdto.UpdateProductStockDTO true "Payload to update stock"
// @Security BearerAuth
func (ctrl *ProductController) UpdateProductStock(ctx *fiber.Ctx) error {
	// Ambil productId dari parameter
	productId, paramErr := functions.ParamToObjectID(ctx, "productId")
	if paramErr != nil {
		return response.BadRequest(ctx, "Invalid product id format", nil)
	}

	// Parsing payload
	var payload productdto.UpdateProductStockDTO
	if err := ctx.BodyParser(&payload); err != nil {
		return response.BadRequest(ctx, "Invalid request body", nil)
	}

	// Validasi payload
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}

	// Ambil storeId dari context
	storeId, ok := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)
	if !ok {
		return response.BadRequest(ctx, "Invalid store ID", nil)
	}

	// Panggil UseCase untuk update produk
	ctrl.UseCase = ctrl.MakeUseCaseFunction()
	err := ctrl.UseCase.UpdateProductStock(ctx.Context(), productId, payload, storeId)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Success(ctx, "Successfully updated product", nil)
}

// GetProductDetail godoc
// @Summary Get Product Detail
// @Description Get Product Detail
// @Tags Product
// @Produce  json
// @Router /api/product/{productId} [get]
// @Param productId path string true "product ID"
// @Security BearerAuth
func (ctrl *ProductController) GetProductDetail(ctx *fiber.Ctx) error {
	productID, err := functions.ParamToObjectID(ctx, "productId")
	if err != nil {
		return response.BadRequest(ctx, "Invalid product ID format", nil)
	}

	storeID := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	productDetail, errUseCase := ctrl.UseCase.GetProductDetail(ctx.Context(), productID, storeID)
	if errUseCase != nil {
		return response.SendResponse(ctx, errUseCase.Code, nil, errUseCase.Message)
	}

	return response.Success(ctx, "Successfully retrieved product details", productDetail)
}

// GetProductList godoc
// @Summary Get Products list
// @Description Get Products list
// @Tags Product
// @Produce  json
// @Router /api/product [get]
// @Param q query dto.PaginationQuery false "Query"
// @Security BearerAuth
func (c *ProductController) GetProductList(ctx *fiber.Ctx) error {

	// parse the request query into a `GetDeviceNamesQuery`
	rawQuery := ctx.Queries()
	query, _ := parser.ParseQuery[dto.PaginationQuery](rawQuery)

	// validate parsed query
	if err := util.ValidateStruct(query); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}

	storeID := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	// initialize the use case with logic specific to the company
	c.UseCase = c.MakeUseCaseFunction()

	// get device names using the validated query
	result, err := c.UseCase.GetProductList(ctx.UserContext(), *query, storeID)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	// return a success response with the fetched device names and a 200 status
	return response.Success(ctx, "successfully get product list", result)
}

// Export Database godoc
// @Summary Export Database
// @Description Export All Database
// @Tags Export
// @Produce  json
// @Router /api/export-all [get]
// @Security BearerAuth
func (c *ProductController) ExportAll(ctx *fiber.Ctx) error {
	c.UseCase = c.MakeUseCaseFunction()
	products, categories, productPhotos, transactions, stores, users, err := c.UseCase.ExportAllData(ctx.UserContext())
	if err != nil {
		return response.InternalServerError(ctx, err.Error(), nil)

	}

	// Create buffer untuk zip file
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Fungsi untuk menambahkan file ke zip
	addToZip := func(filename string, data interface{}) error {
		// Data harus berupa slice
		slice, ok := data.([]interface{})
		if !ok {
			// Coba pakai refleksi sebagai alternatif umum untuk semua slice
			val := reflect.ValueOf(data)
			if val.Kind() != reflect.Slice {
				return errors.New("data must be a slice")
			}
			// Konversi slice ke []interface{}
			slice = make([]interface{}, val.Len())
			for i := 0; i < val.Len(); i++ {
				slice[i] = val.Index(i).Interface()
			}
		}

		var buf bytes.Buffer
		buf.WriteString("[\n")

		for i, item := range slice {
			var jsonBuf bytes.Buffer
			writer, err := bsonrw.NewExtJSONValueWriter(&jsonBuf, true, false)
			if err != nil {
				return err
			}
			enc, err := bson.NewEncoder(writer)
			if err != nil {
				return err
			}
			if err := enc.Encode(item); err != nil {
				return err
			}
			buf.Write(jsonBuf.Bytes())
			if i < len(slice)-1 {
				buf.WriteString(",\n")
			}
		}

		buf.WriteString("\n]")

		writer, err := zipWriter.Create(filename)
		if err != nil {
			return err
		}

		_, err = writer.Write(buf.Bytes())
		return err
	}

	// Tambahkan semua data ke zip
	if err := addToZip("products.json", products); err != nil {
		return response.InternalServerError(ctx, err.Error(), nil)
	}

	if err := addToZip("categories.json", categories); err != nil {
		return response.InternalServerError(ctx, err.Error(), nil)
	}

	if err := addToZip("product_photos.json", productPhotos); err != nil {
		return response.InternalServerError(ctx, err.Error(), nil)
	}

	if err := addToZip("transactions.json", transactions); err != nil {
		return response.InternalServerError(ctx, err.Error(), nil)
	}

	if err := addToZip("stores.json", stores); err != nil {
		return response.InternalServerError(ctx, err.Error(), nil)
	}

	if err := addToZip("users.json", users); err != nil {
		return response.InternalServerError(ctx, err.Error(), nil)
	}

	// Tutup zip writer
	if err := zipWriter.Close(); err != nil {
		return response.InternalServerError(ctx, err.Error(), nil)
	}

	// Set header untuk response
	ctx.Set("Content-Type", "application/zip")
	ctx.Set("Content-Disposition", "attachment; filename=export_"+time.Now().Format("20060102150405")+".zip")

	return ctx.Send(buf.Bytes())
}

// ExportAllProductPhotos godoc
// @Summary Export All Product Photos
// @Description Export All Product Photos
// @Tags Export
// @Produce  json
// @Router /api/product/export-photos [get]
// @Security BearerAuth
func (c *ProductController) ExportAllProductPhotos(ctx *fiber.Ctx) error {
	storeID := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)
	folderPath := fmt.Sprintf("../acts-files/%s/product", storeID.Hex()) // Disarankan menggunakan path relatif dari root proyek

	// 2. Validasi apakah folder ada dan merupakan sebuah direktori
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return response.NotFound(ctx, fmt.Sprintf("Folder untuk toko ini tidak ditemukan: %s", folderPath), nil)
	}
	if !info.IsDir() {
		return response.BadRequest(ctx, fmt.Sprintf("Path yang ditemukan bukan sebuah folder: %s", folderPath), nil)
	}

	// 3. Atur HTTP Header untuk respons file download
	zipFileName := fmt.Sprintf("semua-foto-%s.zip", storeID.Hex())
	ctx.Set("Content-Type", "application/zip")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFileName))

	// 4. Buat zip.Writer yang akan menulis langsung ke response body Fiber
	// INI PERUBAHAN UTAMA: Gunakan ctx.Response().BodyWriter()
	zipWriter := zip.NewWriter(ctx.Response().BodyWriter())
	defer zipWriter.Close()

	// 5. Jelajahi folder secara rekursif menggunakan filepath.Walk
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(folderPath, path)
		if err != nil {
			return err
		}
		header.Name = strings.ReplaceAll(header.Name, "\\", "/")
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		fileToZip, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		_, err = io.Copy(writer, fileToZip)
		return err
	})

	if err != nil {
		// Error saat zipping, kembalikan response error server
		// Perhatikan bahwa header mungkin sudah terkirim, klien mungkin menerima file korup.
		// Ini adalah batasan dari streaming HTTP.
		return response.InternalServerError(ctx, "Gagal saat memproses file zip: "+err.Error(), nil)
	}

	// PENTING: Kembalikan nil untuk menandakan ke Fiber bahwa response berhasil ditangani
	return nil
}
