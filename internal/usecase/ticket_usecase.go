package usecase

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"strings"
	"time"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type TicketUsecase interface {
	GetMyTickets(userID string) ([]dto.MyTicketListResponse, error)
	GetMyTicketDetail(userID, ticketID string) (*dto.MyTicketDetailResponse, error)
	DownloadMyTicketPDF(userID, ticketID string) ([]byte, string, error)
	GetAttendeesByEventID(eventID string, search string, page, limit int) ([]dto.AttendeeResponse, int64, error)
	CheckInTicket(eventID, ticketCode string) (*dto.CheckInResponse, error)
}

type ticketUsecase struct {
	ticketRepo repository.TicketRepository
}

func NewTicketUsecase(
	ticketRepo repository.TicketRepository,
) TicketUsecase {
	return &ticketUsecase{
		ticketRepo: ticketRepo,
	}
}

func (u *ticketUsecase) GetMyTickets(userID string) ([]dto.MyTicketListResponse, error) {
	tickets, err := u.ticketRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to fetch tickets")
	}

	result := make([]dto.MyTicketListResponse, 0, len(tickets))
	for _, ticket := range tickets {
		item := dto.MyTicketListResponse{
			ID:         ticket.ID,
			TicketCode: ticket.TicketCode,
			QRCode:     ticket.QRCode,
			HolderName: ticket.HolderName,
			OrderID:    ticket.OrderID,
			IsUsed:     ticket.IsUsed,
			CreatedAt:  ticket.CreatedAt,
		}

		item.Event.ID = ticket.Order.Event.ID
		item.Event.Title = ticket.Order.Event.Title
		item.Event.Slug = ticket.Order.Event.Slug
		item.Event.Venue = ticket.Order.Event.Venue
		item.Event.Address = ticket.Order.Event.Address
		item.Event.Time = ticket.Order.Event.Time
		item.Event.StartDate = ticket.Order.Event.StartDate
		item.Event.EndDate = ticket.Order.Event.EndDate
		item.Event.BannerURL = ticket.Order.Event.BannerURL

		result = append(result, item)
	}

	return result, nil
}

func (u *ticketUsecase) GetMyTicketDetail(userID, ticketID string) (*dto.MyTicketDetailResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user session")
	}

	ticket, err := u.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	if ticket.Order.UserID != parsedUserID {
		return nil, errors.New("unauthorized")
	}

	if ticket.Order.Status != "paid" {
		return nil, errors.New("ticket not found")
	}

	response := &dto.MyTicketDetailResponse{
		ID:             ticket.ID,
		TicketCode:     ticket.TicketCode,
		QRCode:         ticket.QRCode,
		HolderName:     ticket.HolderName,
		IdentityType:   ticket.IdentityType,
		IdentityNumber: ticket.IdentityNumber,
		HolderPhone:    ticket.HolderPhone,
		HolderEmail:    ticket.HolderEmail,
		Notes:          ticket.Notes,
		IsUsed:         ticket.IsUsed,
		UsedAt:         ticket.UsedAt,
		CreatedAt:      ticket.CreatedAt,
	}

	response.Order.OrderID = ticket.Order.ID
	response.Order.EventName = ticket.Order.Event.Title
	response.Order.UnitPrice = ticket.Order.UnitPrice
	response.Order.OrderStatus = ticket.Order.Status

	return response, nil
}

func (u *ticketUsecase) DownloadMyTicketPDF(userID, ticketID string) ([]byte, string, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, "", errors.New("invalid user session")
	}

	ticket, err := u.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, "", errors.New("ticket not found")
	}

	if ticket.Order.UserID != parsedUserID {
		return nil, "", errors.New("unauthorized")
	}

	if ticket.Order.Status != "paid" {
		return nil, "", errors.New("ticket not found")
	}

	pdfBytes, err := generateTicketPDF(ticket)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate ticket pdf: %w", err)
	}

	filename := fmt.Sprintf("ticket-%s.pdf", sanitizeFileName(ticket.TicketCode))
	return pdfBytes, filename, nil
}

func (u *ticketUsecase) GetAttendeesByEventID(eventID string, search string, page, limit int) ([]dto.AttendeeResponse, int64, error) {
	offset := (page - 1) * limit

	tickets, total, err := u.ticketRepo.FindByEventID(eventID, search, limit, offset)
	if err != nil {
		return nil, 0, errors.New("failed to fetch attendees")
	}

	responses := make([]dto.AttendeeResponse, 0, len(tickets))
	for _, ticket := range tickets {
		responses = append(responses, dto.ToAttendeeResponse(ticket))
	}

	return responses, total, nil
}

func (u *ticketUsecase) CheckInTicket(eventID, ticketCode string) (*dto.CheckInResponse, error) {
	ticket, err := u.ticketRepo.FindByCodeAndEventID(ticketCode, eventID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	if ticket.IsUsed {
		return nil, errors.New("ticket already checked in")
	}

	now := time.Now()
	ticket.IsUsed = true
	ticket.UsedAt = &now

	if err := u.ticketRepo.Update(ticket); err != nil {
		return nil, errors.New("failed to check in ticket")
	}

	usedAtStr := now.Format(time.RFC3339)
	return &dto.CheckInResponse{
		TicketCode: ticket.TicketCode,
		HolderName: ticket.HolderName,
		IsUsed:     ticket.IsUsed,
		UsedAt:     &usedAtStr,
	}, nil
}

func generateTicketPDF(ticket *model.Ticket) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(12, 12, 12)
	pdf.SetAutoPageBreak(true, 12)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(24, 24, 24)
	pdf.CellFormat(0, 10, toPDFSafeText("E-TICKET"), "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(0, 7, toPDFSafeText(ticket.Order.Event.Title), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetDrawColor(210, 210, 210)
	pdf.Rect(12, 32, 186, 112, "D")

	leftX := 16.0
	rightX := 122.0
	startY := 38.0

	pdf.SetTextColor(30, 30, 30)
	pdf.SetFont("Arial", "B", 11)
	pdf.Text(leftX, startY, toPDFSafeText("Informasi Tiket"))

	pdf.SetFont("Arial", "", 10)
	currentY := startY + 6

	rows := [][2]string{
		{toPDFSafeText("Nama Pemegang"), toPDFSafeText(ticket.HolderName)},
		{toPDFSafeText("Kode Tiket"), toPDFSafeText(ticket.TicketCode)},
		{toPDFSafeText("Order ID"), toPDFSafeText(ticket.Order.ID.String())},
		{toPDFSafeText("Email"), toPDFSafeText(ticket.HolderEmail)},
		{toPDFSafeText("Nomor HP"), toPDFSafeText(ticket.HolderPhone)},
		{toPDFSafeText("Identitas"), toPDFSafeText(fmt.Sprintf("%s - %s", ticket.IdentityType, ticket.IdentityNumber))},
		{toPDFSafeText("Venue"), toPDFSafeText(safeString(ticket.Order.Event.Venue))},
		{toPDFSafeText("Alamat"), toPDFSafeText(safeString(ticket.Order.Event.Address))},
		{toPDFSafeText("Tanggal Event"), toPDFSafeText(formatEventDate(ticket.Order.Event.StartDate, ticket.Order.Event.EndDate))},
		{toPDFSafeText("Waktu"), toPDFSafeText(safeString(ticket.Order.Event.Time))},
		{toPDFSafeText("Harga"), toPDFSafeText(formatCurrency(ticket.Order.UnitPrice))},
		{toPDFSafeText("Status"), toPDFSafeText(ticketStatusLabel(ticket))},
	}

	for _, row := range rows {
		currentY = writeRow(pdf, leftX, currentY, row[0], row[1], 92)
	}

	pdf.SetFont("Arial", "B", 11)
	pdf.Text(rightX, startY, toPDFSafeText("Barcode"))

	barcodePNG, err := generateCode128PNG(ticket.TicketCode)
	if err != nil {
		return nil, err
	}

	barcodeOptions := gofpdf.ImageOptions{ImageType: "PNG"}
	pdf.RegisterImageOptionsReader("ticket-barcode", barcodeOptions, bytes.NewReader(barcodePNG))
	pdf.ImageOptions("ticket-barcode", rightX, startY+4, 70, 18, false, barcodeOptions, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(90, 90, 90)
	pdf.SetXY(rightX, startY+24)
	pdf.CellFormat(70, 5, toPDFSafeText(ticket.TicketCode), "", 1, "C", false, 0, "")

	if ticket.QRCode != nil && *ticket.QRCode != "" {
		qrPNG, decodeErr := decodeBase64ImageToPNG(*ticket.QRCode)
		if decodeErr == nil {
			qrOptions := gofpdf.ImageOptions{ImageType: "PNG"}
			pdf.SetFont("Arial", "B", 11)
			pdf.SetTextColor(30, 30, 30)
			pdf.Text(rightX, startY+37, toPDFSafeText("QR Code"))
			pdf.RegisterImageOptionsReader("ticket-qr", qrOptions, bytes.NewReader(qrPNG))
			pdf.ImageOptions("ticket-qr", rightX+12, startY+41, 45, 45, false, qrOptions, 0, "")
		}
	}

	pdf.SetFont("Arial", "I", 9)
	pdf.SetTextColor(110, 110, 110)
	pdf.SetY(148)
	pdf.MultiCell(0, 5, toPDFSafeText("Tunjukkan barcode atau QR code ini saat check-in. Tiket hanya berlaku satu kali sesuai kode tiket."), "", "L", false)

	var out bytes.Buffer
	if err := pdf.Output(&out); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func writeRow(pdf *gofpdf.Fpdf, x, y float64, label, value string, valueWidth float64) float64 {
	pdf.SetXY(x, y)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(70, 70, 70)
	pdf.CellFormat(34, 5, label, "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(30, 30, 30)
	pdf.MultiCell(valueWidth, 5, value, "", "L", false)

	if pdf.GetY() < y+5 {
		return y + 5
	}
	return pdf.GetY()
}

func generateCode128PNG(text string) ([]byte, error) {
	code, err := code128.Encode(text)
	if err != nil {
		return nil, err
	}

	scaled, err := barcode.Scale(code, 440, 100)
	if err != nil {
		return nil, err
	}

	return encodePNG8Bit(scaled)
}

func decodeBase64ImageToPNG(raw string) ([]byte, error) {
	parts := strings.SplitN(raw, ",", 2)
	encoded := raw
	if len(parts) == 2 {
		encoded = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		return nil, err
	}

	return encodePNG8Bit(img)
}

func encodePNG8Bit(src image.Image) ([]byte, error) {
	bounds := src.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, src, bounds.Min, draw.Src)

	var out bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestSpeed}
	if err := encoder.Encode(&out, rgba); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func formatCurrency(value float64) string {
	return fmt.Sprintf("Rp %.0f", value)
}

func formatEventDate(startDate, endDate *time.Time) string {
	if startDate == nil && endDate == nil {
		return "-"
	}

	if startDate != nil && endDate != nil {
		return startDate.Format("02 Jan 2006") + " - " + endDate.Format("02 Jan 2006")
	}

	if startDate != nil {
		return startDate.Format("02 Jan 2006")
	}

	return endDate.Format("02 Jan 2006")
}

func ticketStatusLabel(ticket *model.Ticket) string {
	if ticket.IsUsed {
		if ticket.UsedAt != nil {
			return "Sudah check-in (" + ticket.UsedAt.Format("02 Jan 2006 15:04") + ")"
		}
		return "Sudah check-in"
	}

	return "Belum check-in"
}

func safeString(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "-"
	}
	return *value
}

func sanitizeFileName(value string) string {
	replacer := strings.NewReplacer("/", "-", "\\", "-", " ", "-", ":", "-", "*", "", "?", "", "\"", "", "<", "", ">", "", "|", "")
	return strings.ToLower(replacer.Replace(strings.TrimSpace(value)))
}

func toPDFSafeText(value string) string {
	if value == "" {
		return "-"
	}

	var builder strings.Builder
	builder.Grow(len(value))

	for _, r := range value {
		switch {
		case r == '\n' || r == '\r' || r == '\t':
			builder.WriteRune(' ')
		case r < 32:
			continue
		case r > 255:
			builder.WriteRune('?')
		default:
			builder.WriteRune(r)
		}
	}

	sanitized := strings.TrimSpace(builder.String())
	if sanitized == "" {
		return "-"
	}

	return sanitized
}
