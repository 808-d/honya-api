package coupon

// CouponType is how a coupon's value is interpreted.
type CouponType int16

const (
	CouponTypeFixed   CouponType = 0
	CouponTypePercent CouponType = 1
)

// Valid reports whether t is a known coupon type.
func (t CouponType) Valid() bool {
	switch t {
	case CouponTypeFixed, CouponTypePercent:
		return true
	default:
		return false
	}
}

func (t CouponType) String() string {
	switch t {
	case CouponTypeFixed:
		return "fixed"
	case CouponTypePercent:
		return "percent"
	default:
		return "unknown"
	}
}
