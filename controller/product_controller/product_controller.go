package productcontroller

import (
	"be-shop-vision/dto"
	productdto "be-shop-vision/dto/product"
	usecase "be-shop-vision/usecase/product_usecase"
	"be-shop-vision/util"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/middleware"
	"github.com/susatyo441/go-ta-utils/parser"
	"github.com/susatyo441/go-ta-utils/response"
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
// @Router /product [post]
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
// @Router /product [delete]
// @Param payload body dto.ArrayOfIdDTO true "Payload to delete"
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
// @Router /product/{productId} [put]
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

// GetProductDetail godoc
// @Summary Get Product Detail
// @Description Get Product Detail
// @Tags Product
// @Produce  json
// @Router /product/{productId} [get]
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
// @Router /product [get]
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
