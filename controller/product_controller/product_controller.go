package productcontroller

import (
	"be-shop-vision/dto"
	productdto "be-shop-vision/dto/product"
	usecase "be-shop-vision/usecase/product_usecase"
	"be-shop-vision/util"

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

	var images []string

	// Ambil gambar dari ctx.Locals
	rawImages := []interface{}{
		ctx.Locals(middleware.ContextKey("firstImage")),
		ctx.Locals(middleware.ContextKey("secondImage")),
		ctx.Locals(middleware.ContextKey("thirdImage")),
		ctx.Locals(middleware.ContextKey("fourthImage")),
		ctx.Locals(middleware.ContextKey("fifthImage")),
	}

	// Cek apakah ada yang nil
	for _, rawImage := range rawImages {
		if rawImage == nil {
			response.BadRequest(ctx, "Semua gambar harus diunggah", nil)

		}

		// Konversi ke string dan tambahkan ke slice images
		imageStr, ok := rawImage.(string)
		if !ok {
			response.BadRequest(ctx, "Format gambar tidak valid", nil)

		}
		images = append(images, imageStr)
	}

	var payload productdto.CreateProductDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}
	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	storeId := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	err := ctrl.UseCase.CreateProduct(ctx.Context(), payload, storeId, images)
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

	// Pastikan ctx.Locals tidak nil sebelum melakukan type assertion
	getLocal := func(key string) string {
		if val := ctx.Locals(middleware.ContextKey(key)); val != nil {
			if strVal, ok := val.(string); ok {
				return strVal
			}
		}
		return ""
	}

	// Isi Photos dalam payload jika ada perubahan gambar
	payload.Photos = productdto.ProductPhotosDTO{
		FirstImage:  getLocal("firstImage"),
		SecondImage: getLocal("secondImage"),
		ThirdImage:  getLocal("thirdImage"),
		FourthImage: getLocal("fourthImage"),
		FifthImage:  getLocal("fifthImage"),
	}

	// Panggil UseCase untuk update produk
	ctrl.UseCase = ctrl.MakeUseCaseFunction()
	err := ctrl.UseCase.UpdateProduct(ctx.Context(), productId, payload, storeId)
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
// @Tags Skus
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
