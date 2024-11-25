package external

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

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
		product, err := e.parseProductData(row[1])
		if err != nil {
			return nil, err
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			slog.Error("Error parsing product ID: " + err.Error())
			return nil, err
		}
		product.ID = int64(id)

		products = append(products, product)
	}

	return products, nil
}

func (e *ProductExternal) parseProductData(data string) (product entities.Product, err error) {

	// Регулярное выражение для извлечения названия, срока годности и веса
	re := regexp.MustCompile(`^(.*?)\s+Срок годности\s*:\s*(\d+)\s*суток\s+Вес\s*:\s*([\d,]+)\s*кг\s*$`)
	matches := re.FindStringSubmatch(data)

	if len(matches) != 4 {
		slog.Error("Invalid product data format")
		return product, fmt.Errorf("invalid product data format")
	}

	// Извлечение названия
	product.Name = strings.TrimSpace(matches[1])

	// Извлечение срока годности в секундах
	shelfLifeDays, err := strconv.Atoi(matches[2])
	if err != nil {
		slog.Error("Error parsing shelf life: " + err.Error())
		return product, err
	}
	product.Lifetime = int64(shelfLifeDays * 24 * 60 * 60)

	// Извлечение веса в граммах
	weightKg, err := strconv.ParseFloat(strings.Replace(matches[3], ",", ".", 1), 64)
	if err != nil {
		slog.Error("Error parsing weight: " + err.Error())
		return product, err
	}
	product.Weight = int(weightKg * 1000)

	return product, nil
}

func (e *ProductExternal) SendNewOrderMessage(ctx context.Context, order *entities.Order) error {
	msgText := fmt.Sprintf("Новый заказ от %s, адрес: %s, телефон: %s\n\nСостав заказа:\n\n", order.Customer.Name, order.Customer.Address, order.Customer.Phone)

	for i := range order.Products {
		msgText += fmt.Sprintf("%d. %s - %d шт.", i+1, order.ProductList[i].Name, order.Products[i].Amount)
	}

	msg := tgbotapi.NewMessage(order.Manager.ID, msgText)

	_, err := e.tg.Send(msg)
	if err != nil {
		slog.Error("Failed to send order message: " + err.Error())
		return err
	}

	return nil
}
