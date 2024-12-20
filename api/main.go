package main

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gitea.nichijou.dev/johynpapin/entrelac-server/auth"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/fogleman/gg"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/golang-jwt/jwt/v4"
	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/checkout/session"
	"github.com/stripe/stripe-go/v73/customer"
	"github.com/stripe/stripe-go/webhook"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
)

//go:embed gift.png
var giftImageBytes []byte

//go:embed robotomono.ttf
var robotoMonoTTF []byte
var giftImage image.Image
var codeFace font.Face
var firstNameFace font.Face
var sharesFace font.Face

func sendConfirmAccountEmail(mg mailgun.Mailgun, recipient, token string) error {
	sender := "no-reply@entrelac.coop"
	subject := "Confirmer votre compte Entrelac.coop"
	body := ""

	message := mg.NewMessage(sender, subject, body, recipient)
	message.SetTemplate("confirm-account")
	err := message.AddTemplateVariable("token", token)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err = mg.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func sendResetAccountEmail(mg mailgun.Mailgun, recipient, token string) error {
	sender := "no-reply@entrelac.coop"
	subject := "Réinitialiser votre mot de passe Entrelac.coop"
	body := ""

	message := mg.NewMessage(sender, subject, body, recipient)
	message.SetTemplate("reset-account")
	err := message.AddTemplateVariable("token", token)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err = mg.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID            string  `bun:"id,pk,type:uuid,default:gen_new_uuid()" json:"id"`
	Admin         bool    `bun:"admin,notnull,default:false" json:"admin"`
	Confirmed     bool    `bun:"confirmed,notnull,default:false" json:"confirmed"`
	ConfirmToken  *string `bun:"confirm_token"`
	ResetToken    *string `bun:"reset_token"`
	Email         string  `bun:"email,unique,notnull" json:"email"`
	Password      string  `bun:"password,notnull"`
	PhoneNumber   string  `bun:"phone_number,notnull" json:"phoneNumber"`
	FirstName     string  `bun:"first_name,notnull" json:"firstName"`
	LastName      string  `bun:"last_name,notnull" json:"lastName"`
	Address       string  `bun:"address,notnull" json:"address"`
	PostalCode    string  `bun:"postal_code,notnull" json:"postalCode"`
	City          string  `bun:"city,notnull" json:"city"`
	Country       string  `bun:"country,notnull" json:"country"`
	Category      string  `bun:"category,notnull" json:"category"`
	Reason        *string `bun:"reason" json:"reason"`
	Customer      string  `bun:"customer,notnull" json:"customer"`
	IdentityFront *string `bun:"identity_front" json:"identityFront"`
	IdentityBack  *string `bun:"identity_back" json:"identityBack"`
	AddressProof  *string `bun:"address_proof" json:"addressProof"`
	Accepted      bool    `bun:"accepted,notnull,default:false" json:"accepted"`
	InitialShares uint    `bun:"initial_shares,notnull,default:0" json:"initialShares"`

	Payments []*Payment `bun:"rel:has-many,join:id=user_id"`
}

type Payment struct {
	bun.BaseModel `bun:"table:payments"`

	ID            string    `bun:"id,pk,type:uuid,default:gen_new_uuid()"`
	StripeEventID string    `bun:"stripe_event_id,unique,notnull"`
	Shares        uint      `bun:"shares,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	UserID        string    `bun:"user_id,notnull"`
	GiftID        *string   `bun:"gift_id,unique"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

type Gift struct {
	bun.BaseModel `bun:"table:gifts"`

	ID              string  `bun:"id,pk,type:uuid,default:gen_new_uuid()"`
	Code            string  `bun:"code,unique,notnull"`
	ClaimedByUserID *string `bun:"claimed_by_user_id"`

	Payment *Payment `bun:"rel:has-one,join:id=gift_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type CreateUserRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required"`
	PhoneNumber string  `json:"phone_number" binding:"required"`
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	PostalCode  string  `json:"postal_code" binding:"required"`
	City        string  `json:"city" binding:"required"`
	Country     string  `json:"country" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Reason      *string `json:"reason" binding:"required_unless=Category supporters"`
}

type CreateTokenRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ConfirmUserRequest struct {
	Email string `json:"email" binding:"required"`
	Token string `json:"token" binding:"required"`
}

type StartConfirmUserRequest struct {
	Email string `json:"email" binding:"required"`
}

type ResetUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

type StartResetUserRequest struct {
	Email string `json:"email" binding:"required"`
}

type UseGiftCodeRequest struct {
	GiftCode string `json:"gift_code" binding:"required"`
}

type AdminCSVGetUsersItem struct {
	ID          string  `json:"id"`
	Confirmed   bool    `json:"confirmed"`
	Accepted    bool    `json:"accepted"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Address     string  `json:"address"`
	PostalCode  string  `json:"postalCode"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Category    string  `json:"category"`
	Reason      *string `json:"reason"`
	Shares      uint    `json:"shares"`
}

func (item *AdminCSVGetUsersItem) EncodeCSV() []string {
	var confirmed, accepted, reason string

	if item.Confirmed {
		confirmed = "true"
	} else {
		confirmed = "false"
	}

	if item.Accepted {
		accepted = "true"
	} else {
		accepted = "false"
	}

	if item.Reason != nil {
		reason = *item.Reason
	}

	return []string{
		item.ID,
		confirmed,
		accepted,
		item.Email,
		item.PhoneNumber,
		item.FirstName,
		item.LastName,
		item.Address,
		item.PostalCode,
		item.City,
		item.Country,
		item.Category,
		reason,
		strconv.Itoa(int(item.Shares)),
	}
}

type AdminGetUsersResponseItem struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Accepted  bool   `json:"accepted"`
	Category  string `json:"category"`
	Shares    uint   `json:"shares"`
}

type AdminGetUserResponse struct {
	ID            string  `json:"id"`
	Confirmed     bool    `json:"confirmed"`
	Admin         bool    `json:"admin"`
	Email         string  `json:"email"`
	FirstName     string  `json:"firstName"`
	LastName      string  `json:"lastName"`
	PhoneNumber   string  `json:"phoneNumber"`
	Address       string  `json:"address"`
	PostalCode    string  `json:"postalCode"`
	City          string  `json:"city"`
	Country       string  `json:"country"`
	Category      string  `json:"category"`
	Reason        *string `json:"reason"`
	IdentityFront *string `json:"identityFront"`
	IdentityBack  *string `json:"identityBack"`
	AddressProof  *string `json:"addressProof"`
	Shares        uint    `json:"shares"`
}

type UploadDocumentsForm struct {
	IdentityFront *multipart.FileHeader `form:"identity_front" binding:"required"`
	IdentityBack  *multipart.FileHeader `form:"identity_back"`
	AddressProof  *multipart.FileHeader `form:"address_proof" binding:"required"`
}

type CreateCheckoutSessionRequest struct {
	Quantity uint `json:"quantity" binding:"required"`
	Gift     bool `json:"gift"`
}

var Migrations = migrate.NewMigrations()

//go:embed migrations/*.sql
var sqlMigrations embed.FS

func generateGiftCard(firstName string, code string, shares uint) image.Image {
	codeX := 1149.0 + 30.0
	codeY := 1072.0 + 193.0 - 46

	firstNameX := 480.0 + 20.0
	firstNameY := 1000.0

	sharesX := 624.0
	sharesY := 791.0

	dc := gg.NewContextForImage(giftImage)

	dc.SetFontFace(codeFace)
	dc.SetHexColor("074884")
	dc.DrawString(code, codeX, codeY)

	dc.SetFontFace(firstNameFace)
	dc.SetHexColor("074884")
	dc.DrawString(firstName, firstNameX, firstNameY)

	dc.SetFontFace(sharesFace)
	dc.SetHexColor("074884")
	dc.DrawString(strconv.Itoa(int(shares)), sharesX, sharesY)

	return dc.Image()
}

func init() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		log.Fatal(err)
	}

	var err error
	giftImage, err = png.Decode(bytes.NewReader(giftImageBytes))
	if err != nil {
		log.Fatal(err)
	}

	codeFont, err := truetype.Parse(robotoMonoTTF)
	if err != nil {
		log.Fatal(err)
	}

	codeFace = truetype.NewFace(codeFont, &truetype.Options{Size: 132})

	firstNameFont, err := truetype.Parse(gobold.TTF)
	if err != nil {
		log.Fatal(err)
	}

	firstNameFace = truetype.NewFace(firstNameFont, &truetype.Options{Size: 56})

	sharesFace = truetype.NewFace(firstNameFont, &truetype.Options{Size: 48})
}

func main() {
	_ = godotenv.Load()

	stripeKey := os.Getenv("STRIPE_KEY")
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	stripePrice := os.Getenv("STRIPE_PRICE")
	mailgunDomain := os.Getenv("MAILGUN_DOMAIN")
	mailgunKey := os.Getenv("MAILGUN_KEY")
	mailgunAPIBase := os.Getenv("MAILGUN_API_BASE")
	dataPath := os.Getenv("DATA_PATH")
	dsn := os.Getenv("DSN")
	appBaseURL := os.Getenv("APP_BASE_URL")
	key := []byte(os.Getenv("KEY"))

	if err := os.MkdirAll(filepath.Join(dataPath, "uploads"), os.ModePerm); err != nil {
		log.Fatalf("error creating uploads directory: %v", err)
	}

	stripe.Key = stripeKey

	mg := mailgun.NewMailgun(mailgunDomain, mailgunKey)
	mg.SetAPIBase(mailgunAPIBase)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	migrator := migrate.NewMigrator(db, Migrations)

	if err := migrator.Init(context.Background()); err != nil {
		log.Fatalf("error init migrator: %v", err)
	}

	group, err := migrator.Migrate(context.Background())
	if err != nil {
		log.Fatalf("error migrating: %v", err)
	}
	if group.IsZero() {
		log.Printf("there are no new migrations to run (database is up to date)")
	} else {
		log.Printf("migrated to %s", group)
	}

	r := gin.Default()

	if gin.Mode() == gin.ReleaseMode {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              "https://eb983852916d4ca6937de151edcc4655@o1422105.ingest.sentry.io/6768722",
			TracesSampleRate: 1.0,
		}); err != nil {
			log.Fatalf("error initializing Sentry: %v", err)
		}
		defer sentry.Flush(2 * time.Second)

		r.Use(sentrygin.New(sentrygin.Options{}))
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	r.Use(cors.New(corsConfig))

	r.POST("/tokens", func(c *gin.Context) {
		var json CreateTokenRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error(), "bad-request"})
			return
		}

		user := new(User)
		if err := db.NewSelect().Model(user).Where("email = ?", json.Email).Scan(c); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusBadRequest, ErrorResponse{"No user exists with this email address.", "email-unknown"})
				return
			}

			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		if !auth.CheckPassword(json.Password, user.Password) {
			c.JSON(http.StatusBadRequest, ErrorResponse{"This password is invalid.", "password-invalid"})
			return
		}

		if !user.Confirmed {
			c.JSON(http.StatusUnauthorized, ErrorResponse{"This account is not confirmed.", "not-confirmed"})
			return
		}

		token, err := auth.NewToken(key, user.ID, user.Admin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	r.POST("/users", func(c *gin.Context) {
		var json CreateUserRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error(), "bad-request"})
			return
		}

		exists, err := db.NewSelect().Table("users").Where("email = ?", json.Email).Exists(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, ErrorResponse{"An user using this email address already exists.", "email-used"})
			return
		}

		passwordHash, err := auth.HashPassword(json.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		params := &stripe.CustomerParams{
			Email: stripe.String(json.Email),
			Name:  stripe.String(json.FirstName + " " + json.LastName),
		}

		stripeCustomer, err := customer.New(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		token := auth.NewConfirmToken()

		user := &User{
			ConfirmToken:  &token,
			Email:         json.Email,
			Password:      passwordHash,
			PhoneNumber:   json.PhoneNumber,
			FirstName:     json.FirstName,
			LastName:      json.LastName,
			Address:       json.Address,
			PostalCode:    json.PostalCode,
			City:          json.City,
			Country:       json.Country,
			Category:      json.Category,
			Reason:        json.Reason,
			Customer:      stripeCustomer.ID,
			Accepted:      false,
			InitialShares: 0,
		}

		_, err = db.NewInsert().Model(user).Exec(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		err = sendConfirmAccountEmail(mg, json.Email, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	r.POST("/users/confirm", func(c *gin.Context) {
		var json ConfirmUserRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error(), "bad-request"})
			return
		}

		user := new(User)
		if err := db.NewSelect().Model(user).Where("email = ?", json.Email).Scan(c); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusBadRequest, ErrorResponse{"No user exists with this email address.", "email-unknown"})
				return
			}

			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		if user.ConfirmToken == nil || json.Token != *user.ConfirmToken {
			c.JSON(http.StatusUnauthorized, ErrorResponse{"This token is invalid.", "bad-token"})
			return
		}

		update := &User{ID: user.ID, Confirmed: true, ConfirmToken: nil}
		_, err = db.NewUpdate().Model(update).Column("confirmed", "confirm_token").WherePK().Exec(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
		}

		token, err := auth.NewToken(key, user.ID, user.Admin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	r.POST("/users/confirm/start", func(c *gin.Context) {
		var json StartConfirmUserRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error(), "bad-request"})
			return
		}

		user := new(User)
		if err := db.NewSelect().Model(user).Where("email = ?", json.Email).Scan(c); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusBadRequest, ErrorResponse{"No user exists with this email address.", "email-unknown"})
				return
			}

			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		if user.Confirmed {
			c.JSON(http.StatusUnauthorized, ErrorResponse{"The user is already confirmed.", "confirmed"})
			return
		}

		token := auth.NewConfirmToken()

		update := &User{ID: user.ID, ConfirmToken: &token}
		_, err = db.NewUpdate().Model(update).Column("confirm_token").WherePK().Exec(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
		}

		err = sendConfirmAccountEmail(mg, json.Email, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	r.POST("/users/reset", func(c *gin.Context) {
		var json ResetUserRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error(), "bad-request"})
			return
		}

		user := new(User)
		if err := db.NewSelect().Model(user).Where("email = ?", json.Email).Scan(c); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusBadRequest, ErrorResponse{"No user exists with this email address.", "email-unknown"})
				return
			}

			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		if user.ResetToken == nil || json.Token != *user.ResetToken {
			c.JSON(http.StatusUnauthorized, ErrorResponse{"This token is invalid.", "bad-token"})
			return
		}

		passwordHash, err := auth.HashPassword(json.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		update := &User{ID: user.ID, Confirmed: true, ResetToken: nil, ConfirmToken: nil, Password: passwordHash}
		_, err = db.NewUpdate().Model(update).Column("confirmed", "reset_token", "confirm_token", "password").WherePK().Exec(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
		}

		token, err := auth.NewToken(key, user.ID, user.Admin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	r.POST("/users/reset/start", func(c *gin.Context) {
		var json StartResetUserRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error(), "bad-request"})
			return
		}

		user := new(User)
		if err := db.NewSelect().Model(user).Where("email = ?", json.Email).Scan(c); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusBadRequest, ErrorResponse{"No user exists with this email address.", "email-unknown"})
				return
			}

			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		token := auth.NewResetToken()

		update := &User{ID: user.ID, ResetToken: &token}
		_, err = db.NewUpdate().Model(update).Column("reset_token").WherePK().Exec(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
		}

		err = sendResetAccountEmail(mg, json.Email, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/gifts/:giftID", func(c *gin.Context) {
		giftID := c.Param("giftID")

		gift := new(Gift)
		if err := db.NewSelect().Model(gift).Where("gift.id = ?", giftID).Relation("Payment").Relation("Payment.User").Scan(c); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, ErrorResponse{"Gift not found.", "not-found"})
				return
			}

			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		image := generateGiftCard(gift.Payment.User.FirstName, gift.Code, gift.Payment.Shares)

		c.Header("Content-Type", "image/png")
		c.Header("Content-Disposition", `attachment; filename="cadeau.png"`)

		if err := png.Encode(c.Writer, image); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		c.Status(http.StatusOK)
	})

	authorized := r.Group("/", auth.Middleware(key))

	authorized.POST("/users/me/use-gift-code", func(c *gin.Context) {
		var json UseGiftCodeRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error(), "bad-request"})
			return
		}

		log.Println(json)

		gift := new(Gift)
		if err := db.NewSelect().Model(gift).Where("code = ?", json.GiftCode).Scan(c); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, ErrorResponse{"Gift not found.", "not-found"})
				return
			}

			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		if gift.ClaimedByUserID != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{"Gift Code already claimed.", "gift-code-claimed"})
			return
		}

		userID := c.GetString("userID")
		giftUpdate := &Gift{ID: gift.ID, ClaimedByUserID: &userID}
		_, err = db.NewUpdate().Model(giftUpdate).Column("claimed_by_user_id").WherePK().Exec(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}
	})

	authorized.GET("/users/me", func(c *gin.Context) {
		userID := c.GetString("userID")

		user := new(User)
		err := db.NewSelect().Model(user).Where("id = ?", userID).Scan(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		var shares uint
		err = db.NewSelect().Table("payments").Join("LEFT JOIN gifts").JoinOn("gifts.id = payments.gift_id").ColumnExpr("COALESCE(SUM(shares), 0)").Where("payments.user_id = ?", userID).WhereOr("gifts.claimed_by_user_id = ?", userID).Scan(c, &shares)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		shares += user.InitialShares

		mustUploadDocuments := user.IdentityFront == nil || user.AddressProof == nil

		c.JSON(http.StatusOK, gin.H{
			"email":               user.Email,
			"mustUploadDocuments": mustUploadDocuments,
			"shares":              shares,
		})
	})

	authorized.POST("/users/me/documents", func(c *gin.Context) {
		userID := c.GetString("userID")

		var form UploadDocumentsForm
		if err := c.ShouldBindWith(&form, binding.FormMultipart); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error(), "bad-request"})
			return
		}

		if err := os.MkdirAll(filepath.Join(dataPath, "uploads", userID), os.ModePerm); err != nil {
			log.Fatalf("error creating user uploads directory: %v", err)
		}

		identityFrontKey := uuid.New().String()
		err := c.SaveUploadedFile(form.IdentityFront, filepath.Join(dataPath, "uploads", userID, identityFrontKey))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		var identityBack *string
		if form.IdentityBack != nil {
			identityBackKey := uuid.New().String()
			identityBack = &identityBackKey
			err = c.SaveUploadedFile(form.IdentityBack, filepath.Join(dataPath, "uploads", userID, identityBackKey))
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
				return
			}
		}

		addressProofKey := uuid.New().String()
		err = c.SaveUploadedFile(form.AddressProof, filepath.Join(dataPath, "uploads", userID, addressProofKey))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		user := &User{ID: userID, IdentityFront: &identityFrontKey, IdentityBack: identityBack, AddressProof: &addressProofKey}
		_, err = db.NewUpdate().Model(user).Column("identity_front", "identity_back", "address_proof").WherePK().Exec(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	authorized.POST("/users/me/checkout/sessions", func(c *gin.Context) {
		var json CreateCheckoutSessionRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error(), "bad-request"})
			return
		}

		userID := c.GetString("userID")

		user := new(User)
		err := db.NewSelect().Model(user).Where("id = ?", userID).Scan(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		params := &stripe.CheckoutSessionParams{
			Customer: stripe.String(user.Customer),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(stripePrice),
					Quantity: stripe.Int64(int64(json.Quantity)),
				},
			},
			Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
			SuccessURL: stripe.String(appBaseURL + "payment/success"),
			CancelURL:  stripe.String(appBaseURL + "payment/cancel"),
			Params: stripe.Params{
				Metadata: map[string]string{
					"shares": strconv.Itoa(int(json.Quantity)),
					"userID": userID,
				},
			},
		}

		log.Println(json)

		if json.Gift {
			gift := &Gift{
				Code: gofakeit.Regex("[ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789]{8}"),
			}
			_, err := db.NewInsert().Model(gift).Returning("id").Exec(c)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}

			params.SuccessURL = stripe.String(*params.SuccessURL + "?giftID=" + gift.ID)
			params.Params.Metadata["giftID"] = gift.ID
		}

		s, err := session.New(params)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"url": s.URL,
		})
	})

	admin := authorized.Group("/admin", auth.AdminMiddleware())

	admin.GET("/csv/users", func(c *gin.Context) {
		users := make([]AdminCSVGetUsersItem, 0)
		if err := db.NewRaw("SELECT u.id, u.confirmed, u.accepted, u.email, u.phone_number, u.first_name, u.last_name, u.address, u.postal_code, u.city, u.country, u.category, u.reason, COALESCE(SUM(p.shares), 0) + u.initial_shares AS shares FROM users AS u LEFT JOIN payments AS p ON u.id = p.user_id GROUP BY u.id, u.confirmed, u.accepted, u.email, u.phone_number, u.first_name, u.last_name, u.address, u.postal_code, u.city, u.country, u.category, u.reason ORDER BY u.email ASC").Scan(c, &users); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		buf := new(bytes.Buffer)
		w := csv.NewWriter(buf)

		for _, user := range users {
			if err := w.Write(user.EncodeCSV()); err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
				return
			}
		}

		now := time.Now()
		fileName := now.Format("02-01-2006-15-04") + ".csv"

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="` + fileName + `"`,
		}

		c.DataFromReader(http.StatusOK, int64(buf.Len()), "text/csv", buf, extraHeaders)
		return
	})

	admin.GET("/users", func(c *gin.Context) {
		users := make([]AdminGetUsersResponseItem, 0)
		if err := db.NewRaw("SELECT u.id, u.email, u.first_name, u.last_name, u.accepted, u.category, COALESCE(SUM(p.shares), 0) + u.initial_shares AS shares FROM users AS u LEFT JOIN payments AS p ON u.id = p.user_id GROUP BY u.id, u.email, u.first_name, u.last_name, u.accepted, u.category ORDER BY u.email ASC").Scan(c, &users); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		c.JSON(http.StatusOK, users)
		return
	})

	admin.GET("/users/:userID", func(c *gin.Context) {
		userID := c.Param("userID")

		user := new(User)
		if err := db.NewSelect().Model(user).Where("id = ?", userID).Scan(c); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusBadRequest, ErrorResponse{"No user exists with this ID.", "id-unknown"})
				return
			}

			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		var shares uint
		err = db.NewSelect().Table("payments").ColumnExpr("COALESCE(SUM(shares), 0)").Where("user_id = ?", userID).Scan(c, &shares)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{"Internal server error.", "internal"})
			return
		}

		shares += user.InitialShares

		response := &AdminGetUserResponse{
			ID:            user.ID,
			Confirmed:     user.Confirmed,
			Admin:         user.Admin,
			Email:         user.Email,
			PhoneNumber:   user.PhoneNumber,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Address:       user.Address,
			PostalCode:    user.PostalCode,
			City:          user.City,
			Country:       user.Country,
			Category:      user.Category,
			Reason:        user.Reason,
			IdentityFront: user.IdentityFront,
			IdentityBack:  user.IdentityBack,
			AddressProof:  user.AddressProof,
			Shares:        shares,
		}

		c.JSON(http.StatusOK, response)
		return
	})

	admin.GET("/users/:userID/documents/:documentID", func(c *gin.Context) {
		userID := c.Param("userID")
		documentID := c.Param("documentID")

		documentPath := filepath.Join(dataPath, "uploads", userID, documentID)

		c.File(documentPath)
	})

	r.POST("/stripe/webhook", func(c *gin.Context) {
		const MaxBodyBytes = int64(65536)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{})
			return
		}

		event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), stripeWebhookSecret)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		if event.Type == "checkout.session.completed" {
			var session stripe.CheckoutSession
			err := json.Unmarshal(event.Data.Raw, &session)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{})
				return
			}

			userID := session.Metadata["userID"]
			giftID := session.Metadata["giftID"]

			shares, err := strconv.Atoi(session.Metadata["shares"])
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}

			payment := &Payment{
				StripeEventID: event.ID,
				UserID:        userID,
				CreatedAt:     time.Unix(event.Created, 0),
				Shares:        uint(shares),
			}

			if giftID != "" {
				log.Println("giftID is not empty")
				payment.GiftID = &giftID
			}
			_, err = db.NewInsert().Model(payment).Exec(c)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	r.Run()
}
