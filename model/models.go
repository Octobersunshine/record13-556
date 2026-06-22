package model

import "time"

type FaultType string

const (
	FaultTypeMechanical  FaultType = "机械故障"
	FaultTypeElectrical  FaultType = "电气故障"
	FaultTypeSoftware    FaultType = "软件故障"
	FaultTypeSensor      FaultType = "传感器故障"
	FaultTypeCommunication FaultType = "通信故障"
	FaultTypeOther       FaultType = "其他"
)

type TicketStatus string

const (
	StatusPending   TicketStatus = "待处理"
	StatusProcessing TicketStatus = "处理中"
	StatusResolved  TicketStatus = "已解决"
	StatusClosed    TicketStatus = "已关闭"
)

type ProductionLine struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type Part struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	Category  string  `json:"category"`
	Unit      string  `json:"unit"`
	UnitPrice float64 `json:"unit_price"`
	Stock     int     `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}

type PartConsumption struct {
	ID           string    `json:"id"`
	TicketID     string    `json:"ticket_id"`
	PartID       string    `json:"part_id"`
	PartName     string    `json:"part_name"`
	PartCode     string    `json:"part_code"`
	PartCategory string    `json:"part_category"`
	Quantity     int       `json:"quantity"`
	UnitPrice    float64   `json:"unit_price"`
	TotalPrice   float64   `json:"total_price"`
	Operator     string    `json:"operator"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
}

type PartUsageStats struct {
	PartID       string  `json:"part_id"`
	PartName     string  `json:"part_name"`
	PartCode     string  `json:"part_code"`
	PartCategory string  `json:"part_category"`
	TotalUsed    int     `json:"total_used"`
	TotalAmount  float64 `json:"total_amount"`
	UsageCount   int     `json:"usage_count"`
}

type CategoryCostBreakdown struct {
	Category    string                 `json:"category"`
	TotalAmount float64                `json:"total_amount"`
	TotalItems  int                    `json:"total_items"`
	Items       []PartConsumptionDetail `json:"items"`
}

type PartConsumptionDetail struct {
	ID           string  `json:"id"`
	PartID       string  `json:"part_id"`
	PartName     string  `json:"part_name"`
	PartCode     string  `json:"part_code"`
	PartCategory string  `json:"part_category"`
	Quantity     int     `json:"quantity"`
	Unit         string  `json:"unit"`
	UnitPrice    float64 `json:"unit_price"`
	TotalPrice   float64 `json:"total_price"`
	Operator     string  `json:"operator"`
	Remark       string  `json:"remark"`
	CreatedAt    string  `json:"created_at"`
}

type TicketCostDetail struct {
	TicketID         string               `json:"ticket_id"`
	TicketTitle      string               `json:"ticket_title"`
	DeviceID         string               `json:"device_id"`
	DeviceName       string               `json:"device_name"`
	LineID           string               `json:"line_id"`
	LineName         string               `json:"line_name"`
	FaultType        FaultType            `json:"fault_type"`
	Status           TicketStatus         `json:"status"`
	Handler          string               `json:"handler"`
	TotalItemCount   int                  `json:"total_item_count"`
	TotalPartCount   int                  `json:"total_part_count"`
	TotalCost        float64              `json:"total_cost"`
	CategoryBreakdown []CategoryCostBreakdown `json:"category_breakdown"`
	CostPerCategory  map[string]float64   `json:"cost_per_category"`
	PercentagePerCategory map[string]float64 `json:"percentage_per_category"`
	AvgItemCost      float64              `json:"avg_item_cost"`
	MaxSingleItemCost float64             `json:"max_single_item_cost"`
	MaxCostPartName  string               `json:"max_cost_part_name"`
	RepairDuration   string               `json:"repair_duration"`
	ClosedAt         *string              `json:"closed_at,omitempty"`
}

type DeviceFaultTicket struct {
	ID                string            `json:"id"`
	Title             string            `json:"title"`
	Description       string            `json:"description"`
	DeviceID          string            `json:"device_id"`
	DeviceName        string            `json:"device_name"`
	LineID            string            `json:"line_id"`
	LineCode          string            `json:"line_code"`
	LineName          string            `json:"line_name"`
	FaultType         FaultType         `json:"fault_type"`
	Status            TicketStatus      `json:"status"`
	Priority          int               `json:"priority"`
	Reporter          string            `json:"reporter"`
	Handler           string            `json:"handler,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	ResolvedAt        *time.Time        `json:"resolved_at,omitempty"`
	ClosedAt          *time.Time        `json:"closed_at,omitempty"`
	Tags              []string          `json:"tags,omitempty"`
	PartConsumptions  []PartConsumption `json:"part_consumptions,omitempty"`
	TotalMaterialCost float64           `json:"total_material_cost"`
}

type CreateTicketRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DeviceID    string    `json:"device_id" binding:"required"`
	DeviceName  string    `json:"device_name"`
	LineID      string    `json:"line_id" binding:"required"`
	FaultType   FaultType `json:"fault_type" binding:"required"`
	Priority    int       `json:"priority"`
	Reporter    string    `json:"reporter" binding:"required"`
	Tags        []string  `json:"tags"`
}

type UpdateTicketStatusRequest struct {
	Status  TicketStatus `json:"status" binding:"required"`
	Handler string       `json:"handler"`
}

type CreatePartRequest struct {
	Name      string  `json:"name" binding:"required"`
	Code      string  `json:"code" binding:"required"`
	Category  string  `json:"category"`
	Unit      string  `json:"unit" binding:"required"`
	UnitPrice float64 `json:"unit_price"`
	Stock     int     `json:"stock"`
}

type AddConsumptionRequest struct {
	PartID   string `json:"part_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Operator string `json:"operator" binding:"required"`
	Remark   string `json:"remark"`
}

type CloseTicketRequest struct {
	Status       TicketStatus             `json:"status" binding:"required"`
	Handler      string                   `json:"handler"`
	Consumptions []AddConsumptionRequest  `json:"consumptions"`
	Remark       string                   `json:"remark"`
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
