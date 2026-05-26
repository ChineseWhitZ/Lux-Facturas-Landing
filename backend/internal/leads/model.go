package leads

import (
	"errors"
	"strings"
	"time"
)

type Plan string

const (
	PlanBasic    Plan = "kit_basico"
	PlanStandard Plan = "estandar"
	PlanComplete Plan = "completo"
)

type CustomerType string

const (
	CustomerHasMac     CustomerType = "ya_tiene_mac"
	CustomerNeedsKit   CustomerType = "necesita_kit_completo"
	CustomerWantsQuote CustomerType = "quiere_cotizar"
)

type Status string

const (
	StatusNew       Status = "nuevo"
	StatusContacted Status = "contactado"
	StatusQuoted    Status = "cotizado"
	StatusClosed    Status = "cerrado"
	StatusDiscarded Status = "descartado"
)

type Lead struct {
	ID                    string       `json:"id"`
	CustomerName          string       `json:"customer_name"`
	CompanyName           string       `json:"company_name"`
	DocumentNumber        string       `json:"document_number"`
	Phone                 string       `json:"phone"`
	Email                 string       `json:"email"`
	SelectedPlan          Plan         `json:"selected_plan"`
	CustomerType          CustomerType `json:"customer_type"`
	BusinessCategory      string       `json:"business_category"`
	ApproxProductQuantity int          `json:"approx_product_quantity"`
	City                  string       `json:"city"`
	Message               string       `json:"message"`
	Status                Status       `json:"status"`
	RegisteredAt          time.Time    `json:"registered_at"`
}

type CreateLeadInput struct {
	CustomerName          string       `json:"customer_name"`
	CompanyName           string       `json:"company_name"`
	DocumentNumber        string       `json:"document_number"`
	Phone                 string       `json:"phone"`
	Email                 string       `json:"email"`
	SelectedPlan          Plan         `json:"selected_plan"`
	CustomerType          CustomerType `json:"customer_type"`
	BusinessCategory      string       `json:"business_category"`
	ApproxProductQuantity int          `json:"approx_product_quantity"`
	City                  string       `json:"city"`
	Message               string       `json:"message"`
}

type UpdateStatusInput struct {
	Status Status `json:"status"`
}

func (input CreateLeadInput) Validate() error {
	if strings.TrimSpace(input.CustomerName) == "" {
		return errors.New("customer_name is required")
	}
	if strings.TrimSpace(input.DocumentNumber) == "" {
		return errors.New("document_number is required")
	}
	if strings.TrimSpace(input.Phone) == "" {
		return errors.New("phone is required")
	}
	if strings.TrimSpace(input.Email) == "" {
		return errors.New("email is required")
	}
	if !input.SelectedPlan.Valid() {
		return errors.New("selected_plan is invalid")
	}
	if !input.CustomerType.Valid() {
		return errors.New("customer_type is invalid")
	}
	if input.ApproxProductQuantity < 0 {
		return errors.New("approx_product_quantity cannot be negative")
	}

	return nil
}

func (plan Plan) Valid() bool {
	switch plan {
	case PlanBasic, PlanStandard, PlanComplete:
		return true
	default:
		return false
	}
}

func (customerType CustomerType) Valid() bool {
	switch customerType {
	case CustomerHasMac, CustomerNeedsKit, CustomerWantsQuote:
		return true
	default:
		return false
	}
}

func (status Status) Valid() bool {
	switch status {
	case StatusNew, StatusContacted, StatusQuoted, StatusClosed, StatusDiscarded:
		return true
	default:
		return false
	}
}
