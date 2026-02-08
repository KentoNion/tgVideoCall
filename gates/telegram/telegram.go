package telegram

import (
	"context"
	"log/slog"
	"time"

	"tgVideoCall/domain/admin"
	"tgVideoCall/domain/admin/models"
	"tgVideoCall/pkg/config"

	tele "gopkg.in/telebot.v3"
)

type Server struct {
	Service Service
	bot     *tele.Bot
	cfg     config.Config
	log     slog.Logger
	ctx     context.Context
}

// NewServer создает новый экземпляр сервера Telegram-бота без регистрации хендлеров
func NewServer(ctx context.Context, log slog.Logger, cfg config.Config, Service admin.Service) *Server {
	return &Server{
		log:     log,
		cfg:     cfg,
		Service: Service,
		bot:     initTelebot(cfg, log),
		ctx:     ctx,
	}
}

func initTelebot(cfg config.Config, log slog.Logger) *tele.Bot {
	const op = "gates.telegram.server.initTelebot"
	//регестрируем бота
	bot, err := tele.NewBot(tele.Settings{
		Token:       cfg.APIKeys.Telegram,
		Synchronous: true,
		Verbose:     false,
		OnError: func(err error, msg tele.Context) {
			if err != nil {
				log.Error(op, "failed to send message", err)
			}
		},
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}
	return bot
}

// RunServer регистрирует все публичные и приватные команды, а также middleware
func (s *Server) RunServer() {
	// публичные команды
	//s.bot.Handle("/start", s.hello)
	//s.bot.Handle("/help", s.help)

	// middleware
	adminGroup := s.bot.Group()
	adminGroup.Use(s.adminMiddleware)

	// приватные команды

	go func(ctx context.Context) {
		time.Sleep(1 * time.Second)
		<-ctx.Done()
		s.bot.Stop()
	}(s.ctx)
	s.bot.Start()
}

func (s *Server) adminMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(teleCtx tele.Context) error {
		const op = "gates.telegram.admin.adminMiddleware"
		s.log.Info(op, "trying to check admin status for user: ", teleCtx.Sender().ID)

		userID := teleCtx.Sender().ID
		admin, err := s.Service.GetAdmin(s.ctx, int(userID))

		if err == models.ErrNotAdmin {
			s.log.Debug(op, "user not found in db or not admin, user_id: ", userID)
			err = teleCtx.Reply("У вас нет прав администратора")
			if err != nil {
				s.log.Error(op, "error sending message", err, "to user: ", teleCtx.Sender().ID)
			}
			return nil // Возвращаем nil чтобы не показывать стектрейс
		}

		if err != nil {
			s.log.Error(op, "error getting admin", err)
			return teleCtx.Reply("Произошла внутренняя ошибка")
		}
		s.log.Debug(op, "admin found: ", admin, "role: ", admin.Role)
		if admin.Role != "admin" && admin.Role != "creator" {
			err = teleCtx.Reply("У вас нет прав администратора")
			if err != nil {
				s.log.Error(op, "error sending message", err, "to user: ", teleCtx.Sender().ID)
			}
			return nil // как при ErrNotAdmin — не показываем стектрейс при отказе в доступе
		}
		return next(teleCtx)
	}
}
