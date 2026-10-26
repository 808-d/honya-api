package order

// OrderStatus is the fulfillment lifecycle of an order.
type OrderStatus int16

const (
	OrderStatusPending    OrderStatus = 0
	OrderStatusConfirmed  OrderStatus = 1
	OrderStatusProcessing OrderStatus = 2
	OrderStatusShipping   OrderStatus = 3
	OrderStatusCompleted  OrderStatus = 4
	OrderStatusCancelled  OrderStatus = 5
	OrderStatusRefunded   OrderStatus = 6
)

// Valid reports whether s is a known order status.
func (s OrderStatus) Valid() bool {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusProcessing,
		OrderStatusShipping, OrderStatusCompleted, OrderStatusCancelled,
		OrderStatusRefunded:
		return true
	default:
		return false
	}
}

// CanTransitionTo reports whether the state machine allows s -> next.
func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	switch s {
	case OrderStatusPending:
		return next == OrderStatusConfirmed || next == OrderStatusCancelled
	case OrderStatusConfirmed:
		return next == OrderStatusProcessing || next == OrderStatusCancelled
	case OrderStatusProcessing:
		return next == OrderStatusShipping || next == OrderStatusCancelled
	case OrderStatusShipping:
		return next == OrderStatusCompleted
	case OrderStatusCompleted:
		return next == OrderStatusRefunded
	default:
		return false
	}
}

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "pending"
	case OrderStatusConfirmed:
		return "confirmed"
	case OrderStatusProcessing:
		return "processing"
	case OrderStatusShipping:
		return "shipping"
	case OrderStatusCompleted:
		return "completed"
	case OrderStatusCancelled:
		return "cancelled"
	case OrderStatusRefunded:
		return "refunded"
	default:
		return "unknown"
	}
}

// PaymentStatus is the payment state of an order, independent of fulfillment.
type PaymentStatus int16

const (
	PaymentStatusUnpaid   PaymentStatus = 0
	PaymentStatusPaid     PaymentStatus = 1
	PaymentStatusFailed   PaymentStatus = 2
	PaymentStatusRefunded PaymentStatus = 3
)

// Valid reports whether s is a known payment status.
func (s PaymentStatus) Valid() bool {
	switch s {
	case PaymentStatusUnpaid, PaymentStatusPaid, PaymentStatusFailed,
		PaymentStatusRefunded:
		return true
	default:
		return false
	}
}

// CanTransitionTo reports whether the payment state machine allows s -> next.
func (s PaymentStatus) CanTransitionTo(next PaymentStatus) bool {
	switch s {
	case PaymentStatusUnpaid:
		return next == PaymentStatusPaid || next == PaymentStatusFailed
	case PaymentStatusFailed:
		return next == PaymentStatusPaid
	case PaymentStatusPaid:
		return next == PaymentStatusRefunded
	default:
		return false
	}
}

func (s PaymentStatus) String() string {
	switch s {
	case PaymentStatusUnpaid:
		return "unpaid"
	case PaymentStatusPaid:
		return "paid"
	case PaymentStatusFailed:
		return "failed"
	case PaymentStatusRefunded:
		return "refunded"
	default:
		return "unknown"
	}
}

// PaymentMethod is how an order is paid for.
type PaymentMethod int16

const (
	PaymentMethodCOD PaymentMethod = 0
)

// Valid reports whether m is a supported payment method.
func (m PaymentMethod) Valid() bool {
	return m == PaymentMethodCOD
}

func (m PaymentMethod) String() string {
	switch m {
	case PaymentMethodCOD:
		return "cod"
	default:
		return "unknown"
	}
}
