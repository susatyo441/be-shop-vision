package usercontroller

import (
	userdto "be-shop-vision/dto/user"
	usecase "be-shop-vision/usecase/user_usecase"
	"be-shop-vision/util"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/susatyo441/go-ta-utils/response"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type makeUserUseCaseFunc = func() usecase.IUserUseCase

type UserController struct {
	UseCase             usecase.IUserUseCase
	MakeUseCaseFunction makeUserUseCaseFunc
}

func MakeUserController(makeUseCaseFunc makeUserUseCaseFunc) *UserController {
	return &UserController{MakeUseCaseFunction: makeUseCaseFunc}
}

// --- Konfigurasi Google OAuth ---
// Taruh ini di level package agar tidak dibuat ulang terus-menerus.
var googleOauthConfig *oauth2.Config

// Gunakan fungsi init() untuk setup konfigurasi sekali saat aplikasi dimulai.
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),     // Ambil dari .env
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"), // Ambil dari .env
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),  // Ambil dari .env
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// Metode baru untuk memulai alur login Google
func (ctrl *UserController) GoogleLogin(ctx *fiber.Ctx) error {
	// Buat state acak untuk proteksi CSRF dan simpan di cookie
	expiration := time.Now().Add(20 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &fiber.Cookie{
		Name:    "oauthstate",
		Value:   state,
		Expires: expiration,
	}
	ctx.Cookie(cookie)

	// Arahkan user ke halaman login Google
	url := googleOauthConfig.AuthCodeURL(state)
	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

// Metode baru untuk menangani callback dari Google
func (ctrl *UserController) GoogleCallback(ctx *fiber.Ctx) error {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// 1. Validasi state CSRF
	oauthState := ctx.Cookies("oauthstate")
	if ctx.Query("state") != oauthState {
		return response.BadRequest(ctx, "Invalid oauth state", nil)
	}

	// 2. Tukarkan authorization code dengan token
	code := ctx.Query("code")
	token, err := googleOauthConfig.Exchange(ctx.Context(), code)
	if err != nil {
		return response.InternalServerError(ctx, "Failed to exchange token", err)
	}

	// 3. Dapatkan informasi user dari Google menggunakan token
	oauth2Service, err := googleOauth2.NewService(ctx.Context(), option.WithTokenSource(googleOauthConfig.TokenSource(ctx.Context(), token)))
	if err != nil {
		return response.InternalServerError(ctx, "Failed to create google service", err)
	}
	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return response.InternalServerError(ctx, "Failed to get user info", err)
	}

	// 4. Panggil Use Case
	ctrl.UseCase = ctrl.MakeUseCaseFunction()
	loginResponse, errUseCase := ctrl.UseCase.LoginGoogleCallback(ctx.Context(), userInfo)
	if errUseCase != nil {
		return response.SendResponse(ctx, errUseCase.Code, nil, errUseCase.Message)
	}
	// --- PERUBAHAN DIMULAI DI SINI ---

	// 1. Ubah seluruh objek loginResponse menjadi JSON
	jsonData, err := json.Marshal(loginResponse)
	if err != nil {
		return response.InternalServerError(ctx, "Failed to serialize response", err)
	}

	// 2. Encode JSON menjadi string Base64 yang aman untuk URL
	encodedData := base64.URLEncoding.EncodeToString(jsonData)

	// 3. Buat URL redirect dengan data yang sudah di-encode
	frontendURL := os.Getenv("FRONTEND_URL")

	redirectURL := fmt.Sprintf("%s/auth/callback?data=%s", frontendURL, encodedData)

	// 4. Arahkan pengguna ke URL tersebut
	return ctx.Redirect(redirectURL, http.StatusTemporaryRedirect)
}

// RegisterUser godoc
// @Summary Register User
// @Description Register User
// @Tags User
// @Produce  json
// @Router /api/user/register [post]
// @Param payload body userdto.RegisterUserDTO true "Payload to create"
// @Security BearerAuth
func (ctrl *UserController) RegisterUser(ctx *fiber.Ctx) error {

	var payload userdto.RegisterUserDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}
	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	files := map[string]*multipart.FileHeader{}
	attributes := "profile_picture"

	fileHeader, err := ctx.FormFile(attributes)
	if err == nil { // File ada
		files[attributes] = fileHeader
	}

	errUseCase := ctrl.UseCase.RegisterUser(ctx.Context(), payload, files)
	if errUseCase != nil {
		return response.SendResponse(ctx, errUseCase.Code, nil, errUseCase.Message)
	}

	return response.Created(ctx, "Successfully register user", nil)
}

// LoginUser godoc
// @Summary Login User
// @Description Login User
// @Tags User
// @Produce  json
// @Router /api/user/login [post]
// @Param payload body userdto.LoginUserDTO true "Payload to login"
// @Security BearerAuth
func (ctrl *UserController) LoginUser(c *fiber.Ctx) error {
	var body userdto.LoginUserDTO

	// Parse request body
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validasi input
	if err := util.ValidateStruct(body); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	// Panggil usecase untuk login
	userData, err := ctrl.UseCase.LoginUser(c.Context(), body)
	if err != nil {
		return response.SendResponse(c, err.Code, nil, err.Message)
	}

	// Simpan session user ID
	c.Locals("session", userData.Token)

	return response.Success(c, "Login berhasil", userData)
}
