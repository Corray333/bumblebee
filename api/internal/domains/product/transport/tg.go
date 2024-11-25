package transport

import (
	"context"
	"log/slog"

	"github.com/Corray333/bumblebee/internal/domains/product/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ManagerStateIdle = iota
	ManagerStateWaitingEmail
)

func (t *ProductTransport) RunTelegram() {
	t.tg.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.tg.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go func() {
			err := t.parseMessage(context.Background(), update.Message)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, произошла ошибка. Попробуйте позже.")
				_, err := t.tg.Send(msg)
				if err != nil {
					slog.Error("Failed to send error message: " + err.Error())
				}
			}
		}()
	}
}

func (t *ProductTransport) parseMessage(ctx context.Context, msg *tgbotapi.Message) error {
	if !msg.Chat.IsPrivate() {
		return nil
	}

	manager, err := t.service.GetManagerByID(ctx, &entities.Manager{ID: msg.From.ID})
	if err != nil {
		return err
	}

	if msg.Text == "/start" {
		msg := tgbotapi.NewMessage(msg.Chat.ID, "Привет. Чтобы добавить себя в список менеджеров, поделитесь своим номером телефона. Просто нажмите кнопку ниже⬇")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonContact("Поделиться номером телефона"),
			),
		)

		_, err := t.tg.Send(msg)
		if err != nil {
			slog.Error("Failed to send message: " + err.Error())
			return err
		}
	}

	if manager == nil || manager.State == ManagerStateIdle {
		if msg.Contact != nil {
			manager := &entities.Manager{
				ID:    msg.From.ID,
				Phone: msg.Contact.PhoneNumber,
				State: ManagerStateWaitingEmail,
			}
			err := t.service.SetManager(ctx, manager)
			if err != nil {
				return err
			}

			msg := tgbotapi.NewMessage(msg.Chat.ID, "Спасибо, теперь отправьте свою почту📧")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err = t.tg.Send(msg)
			if err != nil {
				slog.Error("Failed to send message: " + err.Error())
				return err
			}
		}
	} else {
		if msg.Text != "" {
			switch manager.State {
			case ManagerStateWaitingEmail:
				manager.Email = msg.Text
				manager.State = ManagerStateIdle
				err := t.service.SetManager(ctx, manager)
				if err != nil {
					return err
				}

				msg := tgbotapi.NewMessage(msg.Chat.ID, "Спасибо, теперь вы в списке менеджеров👍")
				_, err = t.tg.Send(msg)
				if err != nil {
					slog.Error("Failed to send message: " + err.Error())
					return err
				}
			}
		}
	}
	return nil
}
