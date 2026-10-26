package product

// ProductStatus is the lifecycle state of a product.
type ProductStatus int16

const (
	ProductStatusDraft    ProductStatus = 0
	ProductStatusActive   ProductStatus = 1
	ProductStatusArchived ProductStatus = 2
)

// Valid reports whether s is a known product status.
func (s ProductStatus) Valid() bool {
	switch s {
	case ProductStatusDraft, ProductStatusActive, ProductStatusArchived:
		return true
	default:
		return false
	}
}

func (s ProductStatus) String() string {
	switch s {
	case ProductStatusDraft:
		return "draft"
	case ProductStatusActive:
		return "active"
	case ProductStatusArchived:
		return "archived"
	default:
		return "unknown"
	}
}

// ProductType is the kind of book a product represents.
type ProductType int16

const (
	ProductTypeManga      ProductType = 0
	ProductTypeLightNovel ProductType = 1
	ProductTypeNovel      ProductType = 2
	ProductTypeArtbook    ProductType = 3
	ProductTypeMagazine   ProductType = 4
)

// Valid reports whether t is a known product type.
func (t ProductType) Valid() bool {
	switch t {
	case ProductTypeManga, ProductTypeLightNovel, ProductTypeNovel,
		ProductTypeArtbook, ProductTypeMagazine:
		return true
	default:
		return false
	}
}

func (t ProductType) String() string {
	switch t {
	case ProductTypeManga:
		return "manga"
	case ProductTypeLightNovel:
		return "light_novel"
	case ProductTypeNovel:
		return "novel"
	case ProductTypeArtbook:
		return "artbook"
	case ProductTypeMagazine:
		return "magazine"
	default:
		return "unknown"
	}
}
