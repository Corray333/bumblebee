package external

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/Corray333/bumblebee/internal/domains/product/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
)

type ProductExternal struct {
	tg *tgbotapi.BotAPI
}

func New(tg *tgbotapi.BotAPI) *ProductExternal {
	return &ProductExternal{
		tg: tg,
	}
}

func (e *ProductExternal) GetProducts() (products []entities.Product, err error) {
	f, err := excelize.OpenFile(viper.GetString("price_list_path"))
	if err != nil {
		slog.Error("Error opening price list file: " + err.Error())
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Лист1")
	if err != nil {
		slog.Error("Error reading rows from price list: " + err.Error())
		return nil, err
	}

	for _, row := range rows {
		id, err := strconv.Atoi(row[0])
		if err != nil {
			slog.Error("Error parsing product ID: " + err.Error())
			return nil, err
		}
		product := entities.Product{
			ID:          int64(id),
			Description: row[1],
		}
		products = append(products, product)
	}

	return products, nil
}

func (e *ProductExternal) SendNewOrderMessage(ctx context.Context, order *entities.Order) error {
	msgText := fmt.Sprintf("Заказ #%d от %s, адрес: %s, телефон: %s\n\nСостав заказа:\n\n", order.ID, order.Customer.Name, order.Customer.Address, order.Customer.Phone)

	for i := range order.Products {
		msgText += fmt.Sprintf("%d. %s - %d шт.", i+1, order.ProductList[i].Description, order.Products[i].Amount)
	}

	msg := tgbotapi.NewMessage(order.Manager.ID, msgText)

	_, err := e.tg.Send(msg)
	if err != nil {
		slog.Error("Failed to send order message: " + err.Error())
		return err
	}

	return nil
}
