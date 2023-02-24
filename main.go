package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

func GenerateRandomString() string {
	rand.Seed(time.Now().Unix())

	str := "AsDfGhBvCX123456MnBp"

	shuff := []rune(str)

	// Shuffling the string
	rand.Shuffle(len(shuff), func(i, j int) {
		shuff[i], shuff[j] = shuff[j], shuff[i]
	})

	// Displaying the random string
	// fmt.Println(string(shuff))
	return string(shuff)
}

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))
	// Endpoint to retrieve the user's location and insert it into the database
	e.POST("/create", func(c echo.Context) error {
		// Retrieve the user's location
		var s = snap.Client{}
		s.New("**your Midtrans Server Key**", midtrans.Sandbox)
		// Use to midtrans.Production if you want Production Environment (accept real transaction).
		// 2. Initiate Snap request param
		orderID := GenerateRandomString()
		orderID = "FAUZISHOP-ORDER-ID-" + orderID
		req := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: 25000,
			},
			CreditCard: &snap.CreditCardDetails{
				Secure: true,
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName:    "Fauzi",
				LName:    "Sofyan",
				Email:    "fauzi@example.com",
				Phone:    "085123123",
				ShipAddr: &midtrans.CustomerAddress{Address: "Palad Jaya 1 Street,Bandung,40124,West Java, Indonesia "},
			},
		}
		snapResp, _ := s.CreateTransaction(req)
		if snapResp.RedirectURL == "" {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Payment Error, error when creating transaction payment to midtrans"})
		}
		log.Println(snapResp)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"name":             "Fauzi Sofyan",
			"email":            "fauzi@example.com",
			"phone":            "085123123",
			"address":          "Palad Jaya 1 Street,Bandung,40124,West Java, Indonesia ",
			"payment_url":      snapResp.RedirectURL,
			"transaction_code": orderID,
			"message":          "create payment success",
		})
	})
	e.PUT("/update", func(c echo.Context) error {
		var z = coreapi.Client{}
		z.New("SB-Mid-server-nP8oOrzwnFwp8UTSeDXEhm7v", midtrans.Sandbox)
		// cekStatus, _ := z.CheckTransaction("your transaction code or token transaction here")
		cekStatus, _ := z.CheckTransaction("FAUZISHOP-ORDER-ID-fhpBX6ABGsv32DM15Cn4")
		log.Println(cekStatus.TransactionStatus)
		log.Println(cekStatus.StatusMessage)
		if cekStatus.TransactionStatus != "settlement" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "payment not complete, please complete the payment first"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "check payment status success"})
	})

	// Start the server
	if err := e.Start(":8000"); err != nil {
		log.Println(err.Error())
	}
}
