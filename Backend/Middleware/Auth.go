package Middleware

import (
	"elucid503/SproutLoggingRelay/Models"
	"elucid503/SproutLoggingRelay/Types"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var AccountCache map[string]Models.SproutAccount = make(map[string]Models.SproutAccount) // Will hold recent accounts so we don't have to make costly external requests every time

func CheckAccountFromCache(Token string) (Models.SproutAccount, bool) {

	Account, Exists := AccountCache[Token]

	return Account, Exists

}

func AddAccountToCache(Token string, Account Models.SproutAccount) {

	AccountCache[Token] = Account

	RemoveAccountFromCacheAfter(Token, 60*60) // 1 hour

}

func RemoveAccountFromCacheAfter(Token string, Timeout int) {

	go func() {

		// goroutine to run in background after timeout

		<-time.After(time.Second * time.Duration(Timeout)) // write to channel; everything after this will run after Timeout

		delete(AccountCache, Token)

	}()
}

func Authorize() gin.HandlerFunc {

	return func(GinContext *gin.Context) {

		Account, AuthorizedError := AuthorizeFromRequest(GinContext)

		if AuthorizedError != nil || !(slices.Contains(Account.Flags, "Developer") || slices.Contains(Account.Flags, "LoggingAccess")) {

			GinContext.JSON(401, Types.Response{

				Message: "Unauthorized",
			})

			GinContext.Abort() // Stop the request from continuing

			return

		}

		GinContext.Set("Account", Account)
		GinContext.Next()

	}

}

func AuthorizeFromRequest(GinContext *gin.Context) (*Models.SproutAccount, error) {

	Token, Error := GinContext.Cookie("Sprout-JWT")

	if Error != nil || Token == "" {

		// Try auth header

		Token = GinContext.GetHeader("Authorization")

		if strings.Contains(Token, "Bearer") {

			Token = strings.Split(Token, " ")[1] // Bearer token

		}

	}

	if Token == "" {

		return nil, errors.New("No token provided")

	}

	if Account, Exists := CheckAccountFromCache(Token); Exists {

		return &Account, nil

	}

	// Get account

	Account, Err := Models.GetSproutAccountByToken(Token)

	if Err != nil {

		return nil, errors.New("Failed to get account")

	}

	// Was successful, so add to cache

	AddAccountToCache(Token, Account)

	return &Account, nil

}
