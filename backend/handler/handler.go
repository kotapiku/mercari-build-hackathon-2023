package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kotapiku/mercari-build-hackathon-2023/backend/db"
	"github.com/kotapiku/mercari-build-hackathon-2023/backend/domain"
	"github.com/kotapiku/mercari-build-hackathon-2023/backend/service"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

var (
	logFile = getEnv("LOGFILE", "access.log")
)

type InitializeResponse struct {
	Message string `json:"message"`
}

type registerRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type getUserItemsResponse struct {
	ID           int32             `json:"id"`
	Name         string            `json:"name"`
	Price        int64             `json:"price"`
	CategoryName string            `json:"category_name"`
	Status       domain.ItemStatus `json:"status"`
}

type getItemsResponse struct {
	ID           int32  `json:"id"`
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	CategoryName string `json:"category_name"`
}

type getItemResponse struct {
	ID           int32             `json:"id"`
	Name         string            `json:"name"`
	CategoryID   int64             `json:"category_id"`
	CategoryName string            `json:"category_name"`
	UserID       int64             `json:"user_id"`
	Price        int64             `json:"price"`
	Description  string            `json:"description"`
	Status       domain.ItemStatus `json:"status"`
}

type getCategoriesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type sellRequest struct {
	ItemID int32 `json:"item_id"`
}

type addItemRequest struct {
	Name        string `form:"name"`
	CategoryID  int64  `form:"category_id"`
	Price       int64  `form:"price"`
	Description string `form:"description"`
}

type addItemResponse struct {
	ID int64 `json:"id"`
}

type addBalanceRequest struct {
	Balance int64 `json:"balance"`
}

type getBalanceResponse struct {
	Balance int64 `json:"balance"`
}

type editItemRequest struct {
	Name        string `json:"name"`
	CategoryID  int64  `json:"category_id"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
}

type LoginRequestByID struct {
	UserID   int64  `json:"user_id"`
	Password string `json:"password"`
}

type LoginRequestByName struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type loginResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type description struct {
	ItemName    string `json:"item_name"`
	Description string `json:"description"`
}

type DescriptionResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type DescriptionRequest struct {
	Model     string                       `json:"model"`
	Messages  []*DescriptionRequestMessage `json:"messages"`
	MaxTokens int                          `json:"max_tokens"`
}

type DescriptionRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Message      *DescriptionRequestMessage `json:"message"`
	FinishReason string                     `json:"finish_reason"`
	Index        int                        `json:"index"`
}

type Response struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Usage   *Usage    `json:"usage"`
	Choices []*Choice `json:"choices"`
}

type Handler struct {
	DB           *sql.DB
	UserRepo     db.UserRepository
	ItemRepo     db.ItemRepository
	LoginService service.LoginService
}

func (h *Handler) Initialize(c echo.Context) error {
	err := os.Truncate(logFile, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "Failed to truncate access log"))
	}

	err = db.Initialize(c.Request().Context(), h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "Failed to initialize"))
	}

	return c.JSON(http.StatusOK, InitializeResponse{Message: "Success"})
}

func (h *Handler) AccessLog(c echo.Context) error {
	return c.File(logFile)
}

func isValidName(name string) bool {
	// ユーザー名, アイテム名に使用できるか
	return name != ""
}

func isValidPassword(password string) bool {
	// パスワードに使用できる文字の正規表現パターン
	pattern := "^[a-zA-Z0-9!@#$%^&*]+$"
	reg := regexp.MustCompile(pattern)
	return password != "" && reg.MatchString(password)
}

func (h *Handler) Register(c echo.Context) error {
	req := new(registerRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validation
	if !isValidName(req.Name) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid username"))
	}
	if !isValidPassword(req.Password) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid password"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userID, err := h.UserRepo.AddUser(c.Request().Context(), domain.User{Name: req.Name, Password: string(hash)})
	if err != nil {
		if err == db.ErrConflict {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, registerResponse{ID: userID, Name: req.Name})
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(LoginRequestByID)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validation
	if !isValidPassword(req.Password) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid password"))
	}

	user, encodedToken, err := h.LoginService.LoginByID(ctx, req.UserID, req.Password)
	if err != nil {
		if err == service.ErrMismatchPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, loginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: encodedToken,
	})
}

func (h *Handler) LoginByName(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(LoginRequestByName)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validation
	if !isValidName(req.UserName) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid username"))
	}
	if !isValidPassword(req.Password) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid password"))
	}

	user, encodedToken, err := h.LoginService.LoginByName(ctx, req.UserName, req.Password)
	if err != nil {
		if err == service.ErrMismatchPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, loginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: encodedToken,
	})
}

func (h *Handler) AddItem(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(addItemRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// validation
	if file.Size > 1<<20 {
		return echo.NewHTTPError(http.StatusBadRequest, "file size is too large (> 1MB)")
	}
	if req.Price <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "price must be greater than 0")
	}
	if !isValidName(req.Name) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid name")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	var dest []byte
	blob := bytes.NewBuffer(dest)

	if _, err := io.Copy(blob, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, err = h.ItemRepo.GetCategory(ctx, req.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid categoryID")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	itemID, err := h.ItemRepo.AddItem(c.Request().Context(), domain.Item{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		UserID:      userID,
		Price:       req.Price,
		Description: req.Description,
		Image:       blob.Bytes(),
		Status:      domain.ItemStatusInitial,
	})
	if err != nil {
		if err == db.ErrConflict {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, addItemResponse{ID: itemID})
}

func (h *Handler) Sell(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(sellRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	item, err := h.ItemRepo.GetItem(ctx, req.ItemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	userID, err := GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if item.UserID != userID {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "can not sell other's item")
	}
	if item.Status != domain.ItemStatusInitial {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "invalid item status")
	}

	if err := h.ItemRepo.UpdateItemStatus(ctx, item.ID, domain.ItemStatusOnSale); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "successful")
}

func (h *Handler) getItems(c echo.Context, onSaleOnly bool) error {
	ctx := c.Request().Context()

	items, err := h.ItemRepo.GetItems(ctx, onSaleOnly)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	itemsRsp := make([]getItemResponse, 0, len(items))
	for _, item := range items {
		itemsRsp = append(itemsRsp, getItemResponse{
			ID:           item.Item.ID,
			Name:         item.Item.Name,
			CategoryID:   item.Category.ID,
			CategoryName: item.Category.Name,
			UserID:       item.Item.UserID,
			Price:        item.Item.Price,
			Description:  item.Item.Description,
			Status:       item.Item.Status,
		})
	}
	return c.JSON(http.StatusOK, itemsRsp)
}

func (h *Handler) GetOnSaleItems(c echo.Context) error {
	return h.getItems(c, true)
}

func (h *Handler) GetOnSaleSoldOutItems(c echo.Context) error {
	return h.getItems(c, false)
}

func (h *Handler) GetItem(c echo.Context) error {
	ctx := c.Request().Context()

	itemID, err := strconv.ParseInt(c.Param("itemID"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	item, err := h.ItemRepo.GetItem(ctx, int32(itemID))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	category, err := h.ItemRepo.GetCategory(ctx, item.CategoryID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK,
		getItemResponse{
			ID:           item.ID,
			Name:         item.Name,
			CategoryID:   item.CategoryID,
			CategoryName: category.Name,
			UserID:       item.UserID,
			Price:        item.Price,
			Description:  item.Description,
			Status:       item.Status,
		})
}

func (h *Handler) GetUserItems(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid userID type")
	}

	items, err := h.ItemRepo.GetItemsByUserID(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	var res []getUserItemsResponse
	for _, item := range items {
		cats, err := h.ItemRepo.GetCategories(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		for _, cat := range cats {
			if cat.ID == item.CategoryID {
				res = append(res, getUserItemsResponse{ID: item.ID, Name: item.Name, Price: item.Price, CategoryName: cat.Name, Status: item.Status})
			}
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetCategories(c echo.Context) error {
	ctx := c.Request().Context()

	cats, err := h.ItemRepo.GetCategories(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	res := make([]getCategoriesResponse, len(cats))
	for i, cat := range cats {
		res[i] = getCategoriesResponse{ID: cat.ID, Name: cat.Name}
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetImage(c echo.Context) error {
	ctx := c.Request().Context()

	itemID, err := strconv.ParseInt(c.Param("itemID"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid itemID type")
	}

	data, err := h.ItemRepo.GetItemImage(ctx, int32(itemID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Blob(http.StatusOK, "image/jpeg", data)
}

func DescriptRequestMessage(itemName string, description string) *DescriptionRequestMessage {
	content := fmt.Sprintf(
		`Product Name: %s
		Current Description: %s

		Please generate a more appealing description under 15 words based on this information.
	`, itemName, description)
	return &DescriptionRequestMessage{
		Role:    "user",
		Content: content,
	}
}

func DescriptRequest(itemName string, description string, maxTokens int) *DescriptionRequest {
	messages := []*DescriptionRequestMessage{DescriptRequestMessage(itemName, description)}
	return &DescriptionRequest{
		Model:     "gpt-3.5-turbo",
		Messages:  messages,
		MaxTokens: maxTokens,
	}
}

func (h *Handler) DescriptionHelper(c echo.Context) error {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		echo.NewHTTPError(http.StatusInternalServerError, errors.New("API_KEY is not set"))
	}
	req := new(description)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// create api request from user input
	new_request := DescriptRequest(req.ItemName, req.Description, 20)
	data, err := json.Marshal(new_request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	reqGpt, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(data))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// set api key
	reqGpt.Header.Set("Content-Type", "application/json")
	reqGpt.Header.Set("Authorization", "Bearer "+apiKey)

	// send api request
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	res, err := client.Do(reqGpt)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	// responseをfrontに送る
	return c.JSON(http.StatusOK, response.Choices[0].Message.Content)
}

func (h *Handler) Search(c echo.Context) error {
	ctx := c.Request().Context()

	itemName := c.QueryParam("name")
	items, err := h.ItemRepo.SearchItem(ctx, itemName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	itemsRsp := make([]getItemResponse, 0, len(items))
	for _, item := range items {
		itemsRsp = append(itemsRsp, getItemResponse{
			ID:           item.Item.ID,
			Name:         item.Item.Name,
			CategoryID:   item.Category.ID,
			CategoryName: item.Category.Name,
			UserID:       item.Item.UserID,
			Price:        item.Item.Price,
			Description:  item.Item.Description,
			Status:       item.Item.Status,
		})
	}
	return c.JSON(http.StatusOK, itemsRsp)
}

func (h *Handler) AddBalance(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(addBalanceRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.Balance <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("balance must be positive"))
	}

	userID, err := GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	user, err := h.UserRepo.GetUser(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if err := h.UserRepo.UpdateBalance(ctx, userID, user.Balance+req.Balance); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "successful")
}

func (h *Handler) GetBalance(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	user, err := h.UserRepo.GetUser(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, getBalanceResponse{Balance: user.Balance})
}

func (h *Handler) Purchase(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	itemID, err := strconv.ParseInt(c.Param("itemID"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	item, err := h.ItemRepo.GetItem(ctx, int32(itemID))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// user, sellerの取得
	user, err := h.UserRepo.GetUser(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	sellerID := item.UserID
	seller, err := h.UserRepo.GetUser(ctx, sellerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// 売買が成立するかどうかの判定
	if item.Status != domain.ItemStatusOnSale {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "item is not on sale")
	}
	if user.Balance < item.Price {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "balance is not enough")
	}
	if userID == sellerID {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "can not buy own items")
	}

	// 売買
	if err := h.UserRepo.UpdateBalance(ctx, userID, user.Balance-item.Price); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := h.UserRepo.UpdateBalance(ctx, sellerID, seller.Balance+item.Price); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := h.ItemRepo.UpdateItemStatus(ctx, int32(itemID), domain.ItemStatusSoldOut); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "successful")
}

func (h *Handler) EditItem(c echo.Context) error {
	ctx := c.Request().Context()
	var req editItemRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	// convert string ID to int32
	id, err := strconv.ParseInt(c.Param("itemID"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	item, err := h.ItemRepo.GetItem(ctx, int32(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	// check if the user is the owner of the item
	userID, err := GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if item.UserID != userID {
		return echo.NewHTTPError(http.StatusPreconditionFailed, "can not update other's item")
	}
	// check category exists
	_, err = h.ItemRepo.GetCategory(ctx, req.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid categoryID")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.ItemRepo.EditItem(ctx, int32(id), req.Name, req.CategoryID, req.Price, req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, "successful")
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetUserID(c echo.Context) (int64, error) {
	user := c.Get("user").(*jwt.Token)
	if user == nil {
		return -1, fmt.Errorf("invalid token")
	}
	claims := user.Claims.(*service.JwtCustomClaims)
	if claims == nil {
		return -1, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}
