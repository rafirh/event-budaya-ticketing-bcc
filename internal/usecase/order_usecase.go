package usecase

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"
	"event-budaya-ticketing-bcc/pkg/email"
	"event-budaya-ticketing-bcc/pkg/payment"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

const serviceFeePerTicket = 2000.0

type OrderUsecase interface {
	CreateTicketOrder(userID string, req *dto.CreateTicketOrderRequest) (*dto.CreateTicketOrderResponse, error)
	HandleMidtransWebhook(req *dto.MidtransWebhookRequest) error
	GetMyOrders(userID string) ([]dto.MyOrderResponse, error)
	GetMyOrderDetail(userID, orderID string) (*dto.MyOrderDetailResponse, error)
}

type orderUsecase struct {
	userRepo        repository.UserRepository
	eventRepo       repository.EventRepository
	orderRepo       repository.OrderRepository
	ticketRepo      repository.TicketRepository
	paymentRepo     repository.PaymentRepository
	walletRepo      repository.PromoterWalletRepository
	transactionRepo repository.WalletTransactionRepository
	eventUsecase    EventUsecase
	mailSender      email.Sender
	midtransClient  *payment.Client
	midtransServer  string
}

func NewOrderUsecase(
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
	orderRepo repository.OrderRepository,
	ticketRepo repository.TicketRepository,
	paymentRepo repository.PaymentRepository,
	walletRepo repository.PromoterWalletRepository,
	transactionRepo repository.WalletTransactionRepository,
	eventUsecase EventUsecase,
	mailSender email.Sender,
	midtransClient *payment.Client,
	midtransServer string,
) OrderUsecase {
	return &orderUsecase{
		userRepo:        userRepo,
		eventRepo:       eventRepo,
		orderRepo:       orderRepo,
		ticketRepo:      ticketRepo,
		paymentRepo:     paymentRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		eventUsecase:    eventUsecase,
		mailSender:      mailSender,
		midtransClient:  midtransClient,
		midtransServer:  midtransServer,
	}
}

func (u *orderUsecase) CreateTicketOrder(userID string, req *dto.CreateTicketOrderRequest) (*dto.CreateTicketOrderResponse, error) {
	if u.midtransClient == nil || u.midtransServer == "" {
		return nil, errors.New("midtrans is not configured")
	}

	event, err := u.eventRepo.FindByID(req.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	ticketCount := len(req.Tickets)
	if ticketCount < 1 {
		return nil, errors.New("at least one ticket is required")
	}

	if event.Sold >= event.Quota {
		return nil, errors.New("ticket is sold out")
	}

	if (event.Sold + ticketCount) > event.Quota {
		remaining := event.Quota - event.Sold
		return nil, fmt.Errorf("only %d tickets available", remaining)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user session")
	}

	unitPrice := event.Price
	serviceFeeTotal := float64(ticketCount) * serviceFeePerTicket
	totalPrice := (unitPrice * float64(ticketCount)) + serviceFeeTotal

	order := &model.Order{
		UserID:     parsedUserID,
		EventID:    event.ID,
		Quantity:   ticketCount,
		UnitPrice:  unitPrice,
		ServiceFee: serviceFeeTotal,
		TotalPrice: totalPrice,
		Status:     "pending",
	}
	if err := u.orderRepo.Create(order); err != nil {
		return nil, errors.New("failed to create order")
	}

	tickets := make([]model.Ticket, 0, ticketCount)
	for _, t := range req.Tickets {
		notes := ""
		if t.Notes != nil {
			notes = strings.TrimSpace(*t.Notes)
		}

		ticketCode := generateTicketCode(order.ID)
		qrCodeBase64, err := generateTicketQRCodeBase64(ticketCode)
		if err != nil {
			return nil, errors.New("failed to generate ticket qr code")
		}

		ticket := model.Ticket{
			OrderID:        order.ID,
			TicketCode:     ticketCode,
			QRCode:         &qrCodeBase64,
			HolderName:     t.HolderName,
			IdentityType:   t.IdentityType,
			IdentityNumber: t.IdentityNumber,
			HolderPhone:    t.HolderPhone,
			HolderEmail:    t.HolderEmail,
			Notes:          notes,
			IsUsed:         false,
		}
		tickets = append(tickets, ticket)
	}

	if err := u.ticketRepo.CreateBatch(tickets); err != nil {
		return nil, errors.New("failed to create tickets")
	}

	paymentData := &model.Payment{
		OrderID: order.ID,
		Amount:  totalPrice,
		Status:  "waiting",
	}
	gateway := "midtrans"
	paymentData.PaymentGateway = &gateway
	if err := u.paymentRepo.Create(paymentData); err != nil {
		return nil, errors.New("failed to create payment")
	}

	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	snapResp, err := u.midtransClient.CreateSnapTransaction(payment.SnapTransactionRequest{
		TransactionDetails: payment.TransactionDetails{
			OrderID:  order.ID.String(),
			GrossAmt: int64(totalPrice),
		},
		CustomerDetails: &payment.CustomerDetails{
			FirstName: user.Name,
			Email:     user.Email,
			Phone:     derefString(user.Phone),
		},
	})
	if err != nil {
		return nil, errors.New("failed to create midtrans transaction")
	}

	paymentURL := snapResp.RedirectURL
	paymentData.PaymentURL = &paymentURL
	if err := u.paymentRepo.Update(paymentData); err != nil {
		return nil, errors.New("failed to update payment url")
	}

	event.Sold += ticketCount
	if event.Sold > event.Quota {
		return nil, errors.New("failed to reserve event tickets")
	}
	if err := u.eventRepo.Update(event); err != nil {
		return nil, errors.New("failed to reserve event tickets")
	}

	return &dto.CreateTicketOrderResponse{
		OrderID:         order.ID,
		EventID:         order.EventID,
		TicketCount:     ticketCount,
		UnitPrice:       unitPrice,
		ServiceFee:      serviceFeePerTicket,
		ServiceFeeTotal: serviceFeeTotal,
		TotalPrice:      totalPrice,
		PaymentStatus:   paymentData.Status,
		PaymentToken:    snapResp.Token,
		PaymentURL:      snapResp.RedirectURL,
	}, nil
}

func (u *orderUsecase) HandleMidtransWebhook(req *dto.MidtransWebhookRequest) error {
	if strings.HasPrefix(req.OrderID, "EVT-") {
		return u.eventUsecase.HandleEventPaymentWebhook(req)
	}

	if u.midtransServer == "" {
		return errors.New("midtrans server key is not configured")
	}

	if req.OrderID == "" || req.StatusCode == "" || req.GrossAmount == "" || req.SignatureKey == "" {
		return errors.New("invalid midtrans payload")
	}

	if !isValidMidtransSignature(req, u.midtransServer) {
		return errors.New("invalid midtrans signature")
	}

	order, err := u.orderRepo.FindByID(req.OrderID)
	if err != nil {
		return errors.New("order not found")
	}

	paymentData, err := u.paymentRepo.FindByOrderID(req.OrderID)
	if err != nil {
		return errors.New("payment not found")
	}

	paymentStatus, orderStatus := mapMidtransStatus(req.TransactionStatus, req.FraudStatus)
	isStatusChanged := paymentData.Status != paymentStatus || order.Status != orderStatus

	paymentData.Status = paymentStatus
	method := req.PaymentType
	if method != "" {
		paymentData.PaymentMethod = &method
	}
	if paymentStatus == "success" {
		now := time.Now()
		paymentData.PaidAt = &now
	}
	if err := u.paymentRepo.Update(paymentData); err != nil {
		return errors.New("failed to update payment")
	}

	wasPaid := order.Status == "paid"
	wasPending := order.Status == "pending"
	order.Status = orderStatus
	if err := u.orderRepo.Update(order); err != nil {
		return errors.New("failed to update order")
	}

	if isStatusChanged {
		if notifyErr := u.sendPaymentStatusNotification(order, paymentData, req.TransactionStatus, req.TransactionID); notifyErr != nil {
			log.Printf("warning: failed to send payment notification email for order %s: %v", order.ID.String(), notifyErr)
		}
	}

	if orderStatus == "paid" && !wasPaid {
		event, err := u.eventRepo.FindByID(order.EventID.String())
		if err != nil {
			return errors.New("event not found")
		}

		promoterID := event.PromoterID.String()
		wallet, err := u.walletRepo.FindByPromoterID(promoterID)
		if err != nil {
			wallet = &model.PromotorWallet{
				PromoterID: event.PromoterID,
				Balance:    0,
			}
			if err := u.walletRepo.Create(wallet); err != nil {
				return errors.New("failed to create promoter wallet")
			}
		}

		commission := order.TotalPrice - order.ServiceFee
		wallet.Balance += commission
		if err := u.walletRepo.Update(wallet); err != nil {
			return errors.New("failed to update wallet balance")
		}

		desc := fmt.Sprintf("Commission from order %s", order.ID.String())
		transaction := &model.WalletTransaction{
			WalletID:    wallet.ID,
			Type:        "TICKET_COMMISSION",
			Direction:   "IN",
			Amount:      commission,
			ReferenceID: &order.ID,
			Description: &desc,
		}
		if err := u.transactionRepo.Create(transaction); err != nil {
			return errors.New("failed to create wallet transaction")
		}

		return nil
	}

	if orderStatus == "cancelled" && wasPending {
		event, err := u.eventRepo.FindByID(order.EventID.String())
		if err != nil {
			return errors.New("event not found")
		}
		event.Sold -= order.Quantity
		if event.Sold < 0 {
			event.Sold = 0
		}
		if err := u.eventRepo.Update(event); err != nil {
			return errors.New("failed to restore event tickets")
		}
	}

	return nil
}

func (u *orderUsecase) GetMyOrders(userID string) ([]dto.MyOrderResponse, error) {
	orders, err := u.orderRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to fetch orders")
	}

	result := make([]dto.MyOrderResponse, 0, len(orders))
	for _, order := range orders {
		paymentData, err := u.paymentRepo.FindByOrderID(order.ID.String())
		if err != nil {
			continue
		}

		result = append(result, dto.MyOrderResponse{
			OrderID:       order.ID,
			EventID:       order.EventID,
			EventName:     order.Event.Title,
			TicketCount:   order.Quantity,
			TotalPrice:    order.TotalPrice,
			PaymentStatus: paymentData.Status,
			PaymentURL:    paymentData.PaymentURL,
			OrderStatus:   order.Status,
			CreatedAt:     order.CreatedAt,
		})
	}

	return result, nil
}

func (u *orderUsecase) GetMyOrderDetail(userID, orderID string) (*dto.MyOrderDetailResponse, error) {
	order, err := u.orderRepo.FindByIDWithRelations(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.UserID.String() != userID {
		return nil, errors.New("unauthorized")
	}

	paymentData, err := u.paymentRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	tickets, err := u.ticketRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, errors.New("failed to fetch tickets")
	}

	ticketDetails := make([]dto.TicketDetail, 0, len(tickets))
	for _, ticket := range tickets {
		ticketDetails = append(ticketDetails, dto.TicketDetail{
			ID:             ticket.ID,
			TicketCode:     ticket.TicketCode,
			HolderName:     ticket.HolderName,
			IdentityType:   ticket.IdentityType,
			IdentityNumber: ticket.IdentityNumber,
			HolderPhone:    ticket.HolderPhone,
			HolderEmail:    ticket.HolderEmail,
			Notes:          ticket.Notes,
			IsUsed:         ticket.IsUsed,
			UsedAt:         ticket.UsedAt,
		})
	}

	return &dto.MyOrderDetailResponse{
		OrderID:     order.ID,
		EventID:     order.EventID,
		EventName:   order.Event.Title,
		TicketCount: order.Quantity,
		UnitPrice:   order.UnitPrice,
		ServiceFee:  order.ServiceFee,
		TotalPrice:  order.TotalPrice,
		OrderStatus: order.Status,
		CreatedAt:   order.CreatedAt,
		Tickets:     ticketDetails,
		Payment: dto.PaymentDetail{
			PaymentMethod:  paymentData.PaymentMethod,
			PaymentGateway: paymentData.PaymentGateway,
			Amount:         paymentData.Amount,
			Status:         paymentData.Status,
			PaymentURL:     paymentData.PaymentURL,
			PaidAt:         paymentData.PaidAt,
		},
	}, nil
}

func mapMidtransStatus(transactionStatus, fraudStatus string) (paymentStatus, orderStatus string) {
	switch transactionStatus {
	case "settlement":
		return "success", "paid"
	case "capture":
		if fraudStatus == "accept" {
			return "success", "paid"
		}
		return "waiting", "pending"
	case "pending":
		return "waiting", "pending"
	case "deny", "cancel", "expire", "failure":
		return "failed", "cancelled"
	default:
		return "waiting", "pending"
	}
}

func isValidMidtransSignature(req *dto.MidtransWebhookRequest, serverKey string) bool {
	raw := req.OrderID + req.StatusCode + req.GrossAmount + serverKey
	sum := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(sum[:])
	return strings.EqualFold(expected, req.SignatureKey)
}

func generateTicketCode(orderID uuid.UUID) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("TIX-%s-%d", strings.ToUpper(orderID.String()[0:8]), timestamp)
}

func generateTicketQRCodeBase64(ticketCode string) (string, error) {
	pngData, err := qrcode.Encode(ticketCode, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(pngData)
	return "data:image/png;base64," + encoded, nil
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (u *orderUsecase) sendPaymentStatusNotification(order *model.Order, paymentData *model.Payment, transactionStatus, transactionID string) error {
	if u.mailSender == nil {
		return nil
	}

	user, err := u.userRepo.FindByID(order.UserID.String())
	if err != nil {
		return errors.New("failed to send payment notification email")
	}

	event, err := u.eventRepo.FindByID(order.EventID.String())
	if err != nil {
		return errors.New("failed to send payment notification email")
	}

	eventDate := "-"
	if event.StartDate != nil {
		eventDate = event.StartDate.Format("02 Jan 2006")
		if event.EndDate != nil {
			eventDate = eventDate + " - " + event.EndDate.Format("02 Jan 2006")
		}
	}

	eventTime := "-"
	if event.Time != nil && strings.TrimSpace(*event.Time) != "" {
		eventTime = strings.TrimSpace(*event.Time)
	}

	venue := "-"
	if event.Venue != nil && strings.TrimSpace(*event.Venue) != "" {
		venue = strings.TrimSpace(*event.Venue)
	}

	method := "-"
	if paymentData.PaymentMethod != nil && strings.TrimSpace(*paymentData.PaymentMethod) != "" {
		method = strings.TrimSpace(*paymentData.PaymentMethod)
	}

	paymentURL := "-"
	if paymentData.PaymentURL != nil && strings.TrimSpace(*paymentData.PaymentURL) != "" {
		paymentURL = strings.TrimSpace(*paymentData.PaymentURL)
	}

	subject := "Pembayaran Tiket Diproses"
	header := "Update Pembayaran Tiket"
	statusText := "Status pembayaran kamu saat ini sedang diproses."

	if order.Status == "paid" {
		subject = "Pembayaran Berhasil - E-Ticket Siap Digunakan"
		header = "Pembayaran Berhasil"
		statusText = "Pembayaran kamu sudah kami terima. E-ticket siap dipakai saat check-in acara."
	} else if order.Status == "cancelled" {
		subject = "Pembayaran Gagal / Dibatalkan"
		header = "Pembayaran Tidak Berhasil"
		statusText = "Pembayaran kamu tidak berhasil atau dibatalkan. Kamu bisa melakukan pemesanan ulang."
	}

	body := fmt.Sprintf(`
		<p>Halo %s,</p>
		<p>%s</p>
		<h3 style="margin-bottom:8px;">%s</h3>
		<table style="border-collapse:collapse; width:100%%; max-width:640px;">
			<tr><td style="padding:6px 0; width:220px;"><strong>Order ID</strong></td><td style="padding:6px 0;">%s</td></tr>
			<tr><td style="padding:6px 0;"><strong>Transaction ID</strong></td><td style="padding:6px 0;">%s</td></tr>
			<tr><td style="padding:6px 0;"><strong>Event</strong></td><td style="padding:6px 0;">%s</td></tr>
			<tr><td style="padding:6px 0;"><strong>Tanggal Event</strong></td><td style="padding:6px 0;">%s</td></tr>
			<tr><td style="padding:6px 0;"><strong>Waktu Event</strong></td><td style="padding:6px 0;">%s</td></tr>
			<tr><td style="padding:6px 0;"><strong>Lokasi</strong></td><td style="padding:6px 0;">%s</td></tr>
			<tr><td style="padding:6px 0;"><strong>Jumlah Tiket</strong></td><td style="padding:6px 0;">%d tiket</td></tr>
			<tr><td style="padding:6px 0;"><strong>Total Pembayaran</strong></td><td style="padding:6px 0;">%s</td></tr>
			<tr><td style="padding:6px 0;"><strong>Metode Pembayaran</strong></td><td style="padding:6px 0;">%s</td></tr>
			<tr><td style="padding:6px 0;"><strong>Status Midtrans</strong></td><td style="padding:6px 0;">%s</td></tr>
			<tr><td style="padding:6px 0;"><strong>Status Order</strong></td><td style="padding:6px 0;">%s</td></tr>
			<tr><td style="padding:6px 0;"><strong>Link Pembayaran</strong></td><td style="padding:6px 0;">%s</td></tr>
		</table>
		<p style="margin-top:16px;">Terima kasih sudah menggunakan layanan kami.</p>
	`,
		html.EscapeString(user.Name),
		html.EscapeString(statusText),
		html.EscapeString(header),
		html.EscapeString(order.ID.String()),
		html.EscapeString(defaultIfEmpty(transactionID, "-")),
		html.EscapeString(event.Title),
		html.EscapeString(eventDate),
		html.EscapeString(eventTime),
		html.EscapeString(venue),
		order.Quantity,
		html.EscapeString(formatRupiah(order.TotalPrice)),
		html.EscapeString(method),
		html.EscapeString(defaultIfEmpty(transactionStatus, "-")),
		html.EscapeString(order.Status),
		html.EscapeString(paymentURL),
	)

	if err := u.mailSender.Send(user.Email, subject, body); err != nil {
		return errors.New("failed to send payment notification email")
	}

	return nil
}

func formatRupiah(amount float64) string {
	return fmt.Sprintf("Rp %.0f", amount)
}

func defaultIfEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
