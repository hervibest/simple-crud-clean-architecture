package helper

import (
	"simple-crud-clean-architecture/internal/model"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MidtransClient struct {
	midtransServerKey string
	log               *logrus.Logger
}

func NewMidtransClient(viper *viper.Viper, log *logrus.Logger) *MidtransClient {
	midtransServerKey := viper.GetString("app.midtrans_server_key")
	return &MidtransClient{
		midtransServerKey: midtransServerKey,
		log:               log,
	}

}

func (c *MidtransClient) CreateSnapshot(request *model.MidtransSnapshotRequest) (*snap.Response, error) {
	var s snap.Client
	s.New(c.midtransServerKey, midtrans.Sandbox)
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  request.OrderID,
			GrossAmt: int64(request.GrossAmt),
		},
		// CreditCard: &snap.CreditCardDetails{
		// 	Secure: true,
		// },
		CustomerDetail: &midtrans.CustomerDetails{
			// FName: "John",
			// LName: "Doe",
			Email: request.Email,
			// Phone: "081234567890",
		},
	}

	response, err := s.CreateTransaction(req)
	if err != nil {
		c.log.Error("Error creating transaction")
		return nil, err
	}

	return response, nil
}

func (c *MidtransClient) GetMidtransKey() string {
	midtransKey := c.midtransServerKey
	return midtransKey
}
