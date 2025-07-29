package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/heinswanhtet/blogora-api/constants"
	"github.com/heinswanhtet/blogora-api/types"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

func init() {
	// Use json field names
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			return ""
		}
		return name
	})
}

type successResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Meta    any    `json:"meta,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type errorResponse struct {
	Status  string `json:"status"`
	Message any    `json:"message"`
}

func GenerateUUID() string {
	return uuid.NewString()
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func WriteError(w http.ResponseWriter, status int, message any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{
		Status:  "fail",
		Message: message,
	})
}

func WriteJSON(w http.ResponseWriter, status int, v any, message string, meta any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(successResponse{
		Status:  "success",
		Message: message,
		Meta:    meta,
		Data:    v,
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// var SECRET []byte = func() []byte {
// 	secret := configs.Envs.SECRET
// 	h := sha256.New()
// 	_, err := h.Write([]byte(secret))
// 	if err != nil {
// 		panic("SHA256 SECRET failed!")
// 	}
// 	sha256SecretByte := h.Sum(nil)
// 	return sha256SecretByte
// }()

func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(constants.SECRET)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(text), nil)

	return fmt.Sprintf("%x:%x", nonce, ciphertext), nil
}

func Decrypt(code string) (string, error) {
	parts := strings.Split(code, ":")

	nonce, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(strings.Join(parts[1:], ":"))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(constants.SECRET)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	text, err := aesgcm.Open(nil, nonce, []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(text), nil
}

func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	validationErrors := err.(validator.ValidationErrors)

	for _, fieldErr := range validationErrors {
		errors[fieldErr.Field()] = getCustomMessage(fieldErr)
	}

	return errors
}

func getCustomMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

func GetOffsetToPaginate(page, pageSize int) int {
	return (page - 1) * pageSize
}

func GenerateMetaPagination(page, pageSize, total int) map[string]int {
	lastPage := math.Ceil(float64(total) / float64(pageSize))

	return map[string]int{
		"current_page": page,
		"last_page":    int(lastPage),
		"per_page":     pageSize,
		"total":        total,
	}
}

type maybe[T any] struct {
	field *T
}

func Maybe[T any](field *T) *maybe[T] {
	return &maybe[T]{
		field: field,
	}
}

func (m *maybe[T]) Else(value T) T {
	if m.field == nil {
		return value
	}
	return *m.field
}

type SetQuery struct {
	Field any
	Col   string
}

func NewSetQuery[T any](field *T, col string) *SetQuery {
	return &SetQuery{
		Field: field,
		Col:   col,
	}
}

func CheckNil(i any) bool {
	v := reflect.ValueOf(i)
	return v.IsNil()
}

func GetSetQuery(id string, table string, dataList *[]*SetQuery) (string, []any, bool) {
	setClauses := []string{}
	args := []any{}

	for _, data := range *dataList {
		if !CheckNil(data.Field) {
			setClauses = append(setClauses, fmt.Sprintf("%v = ?", data.Col))
			args = append(args, data.Field)
		}
	}

	if len(setClauses) == 0 {
		return "", nil, false
	}

	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, time.Now().UTC())

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", table, strings.Join(setClauses, ", "))
	args = append(args, id)

	return query, args, true
}

func GetSetQueryMap(id string, table string, dataList []map[string]any) (string, []any, bool) {
	setClauses := []string{}
	args := []any{}

	for _, data := range dataList {
		if !CheckNil(data["field"]) {
			setClauses = append(setClauses, fmt.Sprintf("%v = ?", data["col"]))
			args = append(args, data["field"])
		}
	}

	if len(setClauses) == 0 {
		return "", nil, false
	}

	setClauses = append(setClauses, "updated_at = ?")
	args = append(args, time.Now().UTC())

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", table, strings.Join(setClauses, ", "))
	args = append(args, id)

	return query, args, true
}

func GetPageAndPageSize(query url.Values) (int, int) {
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(query.Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	return page, pageSize
}

func GetSanitizedQuery(query url.Values, v string, fallback string, allowedList ...string) string {
	allowedMap := make(map[string]bool)
	for _, v := range allowedList {
		allowedMap[v] = true
	}
	field := strings.ToLower(query.Get(v))
	if !allowedMap[field] {
		return fallback
	}
	return field
}

func Contains[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// default list - ["page", "pageSize", "sort_by", "sort_type", "search"]
// Additional values can be provided via arguments
func GetRestOfQuery(query url.Values, list ...string) map[string]string {
	rest := make(map[string]string)
	vs := []string{"page", "pageSize", "sort_by", "sort_type", "search"}
	vs = append(vs, list...)
	for k := range query {
		// isTrue := true
		// for _, s := range vs {
		// 	if k == s {
		// 		isTrue = false
		// 		break
		// 	}
		// }
		// if isTrue {
		// 	rest[k] = query.Get(k)
		// }
		if !Contains(vs, k) {
			rest[k] = query.Get(k)
		}
	}
	return rest
}

func GetSearchQuery(
	validFields []string,
	search string,
	allowedSearchList *[]string,
	otherQuery *map[string]string,
	getTotalQuery string,
) (string, string) {
	var searchQuery string

	if search != "" {
		searchQuery = "AND ("
		list := []string{}
		for _, v := range *allowedSearchList {
			list = append(list, fmt.Sprintf(`%s LIKE "%%%s%%"`, v, search))
		}
		searchQuery = fmt.Sprintf(`%s %s )`, searchQuery, strings.Join(list, " OR "))
		getTotalQuery = fmt.Sprint(getTotalQuery, searchQuery)
	} else if len(*otherQuery) > 0 {
		searchQuery = "AND ("
		list := []string{}
		for k, v := range *otherQuery {
			if !Contains(validFields, k) {
				continue
			}
			if t, err := time.Parse("2006-01-02", v); err == nil {
				tt := t.Add(time.Hour * 24).Format("2006-01-02")
				t := t.Format("2006-01-02")
				list = append(list, fmt.Sprintf(`%s >= "%s" AND %s < "%s"`, k, t, k, tt))
			} else if n, err := strconv.Atoi(v); err == nil {
				list = append(list, fmt.Sprintf(`%s = %v`, k, n))
			} else {
				list = append(list, fmt.Sprintf(`%s LIKE "%%%s%%"`, k, v))
			}
		}
		if len(list) == 0 {
			return "", getTotalQuery
		}
		searchQuery = fmt.Sprintf(`%s %s )`, searchQuery, strings.Join(list, " AND "))
		getTotalQuery = fmt.Sprint(getTotalQuery, searchQuery)
	}

	return searchQuery, getTotalQuery
}

// func GetJSONFieldNames(i any) []string {
// 	t := reflect.TypeOf(i)
// 	if t.Kind() == reflect.Ptr {
// 		t = t.Elem() // to work even if i is passed as ptr
// 	}
// 	var fields []string
// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		tag := field.Tag.Get("json")
// 		if tag != "" && tag != "-" {
// 			fields = append(fields, tag)
// 		}
// 	}
// 	return fields
// }

func RetrieveBearerToken(r *http.Request) (string, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("authorization missing")
	}
	token := strings.Split(authorization, " ")
	if len(token) == 1 {
		return "", fmt.Errorf("invalid token")
	}
	return token[1], nil
}

func CreateJWT(userID string) (string, error) {
	expiration := time.Second * time.Duration(3600*12)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(expiration).Unix(),
		"iat":    time.Now().Unix(),
	})

	tokenString, err := token.SignedString(constants.SECRET)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(constants.SECRET), nil
	})
}

func GetJWTPayload(ctx context.Context) (*types.ContextData, error) {
	payload, ok := ctx.Value(constants.ContextData).(*types.ContextData)
	if !ok {
		return nil, fmt.Errorf("jwt payload failed")
	}
	return payload, nil
}

func GenerateUniqueSlug(title string, exists func(slug string) bool) string {
	base := slugify(title)
	for {
		suffix := secureRandomString(6)
		slug := fmt.Sprintf("%s-%s", base, suffix)
		if !exists(slug) {
			return slug
		}
	}
}

func slugify(title string) string {
	title = strings.ToLower(title)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug := re.ReplaceAllString(title, "-")
	return strings.Trim(slug, "-")
}

// Generate short random alphanumeric string using crypto/rand (stronger than math/rand)
func secureRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // Fail fast in prod
	}
	// Use base32 to keep it URL-safe and alphanumeric
	return strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b))[:length]
}

func FormatUserName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	return s
}
