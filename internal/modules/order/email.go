package order

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"strings"
	texttemplate "text/template"

	"github.com/808-d/honya-api/pkg/mailer"
)

type orderEmailItem struct {
	Title     string
	Quantity  int32
	UnitPrice string
	Subtotal  string
}

type orderEmailData struct {
	Greeting        string
	OrderID         string
	OrderURL        string
	CreatedAt       string
	PaymentMethod   string
	Items           []orderEmailItem
	Subtotal        string
	Discount        string
	Shipping        string
	Total           string
	ShippingName    string
	ShippingPhone   string
	ShippingAddress string
	Note            string
}

const orderEmailTextTmpl = `Hi {{.Greeting}},

Thanks for your order. We have received order {{.OrderID}}.

Placed: {{.CreatedAt}}
Payment: {{.PaymentMethod}}

Items:
{{range .Items}}- {{.Title}} x{{.Quantity}} — {{.Subtotal}} ({{.UnitPrice}} each)
{{end}}
Subtotal: {{.Subtotal}}
Discount: {{.Discount}}
Shipping: {{.Shipping}}
Total: {{.Total}}

Shipping to:
{{.ShippingName}}
{{.ShippingPhone}}
{{.ShippingAddress}}
{{if .Note}}
Note: {{.Note}}
{{end}}
View your order: {{.OrderURL}}

- Honya
`

const orderEmailHTMLTmpl = `<!DOCTYPE html>
<html>
  <body style="font-family: Arial, Helvetica, sans-serif; color: #222222; line-height: 1.5;">
    <p>Hi {{.Greeting}},</p>
    <p>Thanks for your order. We have received order <strong>{{.OrderID}}</strong>.</p>
    <p>Placed: {{.CreatedAt}}<br>Payment: {{.PaymentMethod}}</p>

    <table style="border-collapse: collapse; width: 100%; max-width: 640px;">
      <thead>
        <tr>
          <th style="text-align: left; padding: 8px 6px; border-bottom: 1px solid #dddddd;">Item</th>
          <th style="text-align: right; padding: 8px 6px; border-bottom: 1px solid #dddddd;">Qty</th>
          <th style="text-align: right; padding: 8px 6px; border-bottom: 1px solid #dddddd;">Unit price</th>
          <th style="text-align: right; padding: 8px 6px; border-bottom: 1px solid #dddddd;">Subtotal</th>
        </tr>
      </thead>
      <tbody>
        {{range .Items}}
        <tr>
          <td style="padding: 8px 6px; border-bottom: 1px solid #eeeeee;">{{.Title}}</td>
          <td style="padding: 8px 6px; border-bottom: 1px solid #eeeeee; text-align: right;">{{.Quantity}}</td>
          <td style="padding: 8px 6px; border-bottom: 1px solid #eeeeee; text-align: right;">{{.UnitPrice}}</td>
          <td style="padding: 8px 6px; border-bottom: 1px solid #eeeeee; text-align: right;">{{.Subtotal}}</td>
        </tr>
        {{end}}
      </tbody>
    </table>

    <p style="margin-top: 16px;">
      Subtotal: {{.Subtotal}}<br>
      Discount: {{.Discount}}<br>
      Shipping: {{.Shipping}}<br>
      <strong>Total: {{.Total}}</strong>
    </p>

    <p>
      <strong>Shipping to</strong><br>
      {{.ShippingName}}<br>
      {{.ShippingPhone}}<br>
      {{.ShippingAddress}}
    </p>

    {{if .Note}}<p><strong>Note:</strong> {{.Note}}</p>{{end}}

    <p>
      <a href="{{.OrderURL}}"
         style="display: inline-block; padding: 10px 18px; background: #111111; color: #ffffff; text-decoration: none; border-radius: 6px;">
        View your order
      </a>
    </p>
    <p>Or copy this link into your browser:<br><a href="{{.OrderURL}}">{{.OrderURL}}</a></p>
    <p>- Honya</p>
  </body>
</html>
`

var (
	orderEmailText = texttemplate.Must(texttemplate.New("order-text").Parse(orderEmailTextTmpl))
	orderEmailHTML = htmltemplate.Must(htmltemplate.New("order-html").Parse(orderEmailHTMLTmpl))
)

func buildOrderConfirmationEmail(contact *UserContact, o *Order, webBaseURL string) (mailer.Message, error) {
	greeting := strings.TrimSpace(contact.DisplayName)
	if greeting == "" {
		greeting = contact.Username
	}

	items := make([]orderEmailItem, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, orderEmailItem{
			Title:     item.ProductTitle,
			Quantity:  item.Quantity,
			UnitPrice: formatOrderMoney(item.UnitPrice),
			Subtotal:  formatOrderMoney(item.Subtotal),
		})
	}

	data := orderEmailData{
		Greeting:        greeting,
		OrderID:         o.PublicID,
		OrderURL:        fmt.Sprintf("%s/orders/%s", strings.TrimRight(webBaseURL, "/"), o.PublicID),
		CreatedAt:       o.CreateTime.Format("Jan 2, 2006 15:04 MST"),
		PaymentMethod:   "Cash on delivery (COD)",
		Items:           items,
		Subtotal:        formatOrderMoney(o.Subtotal),
		Discount:        formatOrderMoney(o.DiscountAmount),
		Shipping:        formatOrderMoney(o.ShippingFee),
		Total:           formatOrderMoney(o.Total),
		ShippingName:    o.ShippingName,
		ShippingPhone:   o.ShippingPhone,
		ShippingAddress: o.ShippingAddress,
		Note:            strings.TrimSpace(o.Note),
	}

	var text, html bytes.Buffer
	if err := orderEmailText.Execute(&text, data); err != nil {
		return mailer.Message{}, err
	}
	if err := orderEmailHTML.Execute(&html, data); err != nil {
		return mailer.Message{}, err
	}

	return mailer.Message{
		To:      contact.Email,
		Subject: "Your Honya order is confirmed",
		Text:    text.String(),
		HTML:    html.String(),
	}, nil
}

func formatOrderMoney(value float64) string {
	return fmt.Sprintf("¥%.2f", value)
}
