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

type DeviceFaultTicket struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	DeviceID    string       `json:"device_id"`
	DeviceName  string       `json:"device_name"`
	LineID      string       `json:"line_id"`
	LineCode    string       `json:"line_code"`
	LineName    string       `json:"line_name"`
	FaultType   FaultType    `json:"fault_type"`
	Status      TicketStatus `json:"status"`
	Priority    int          `json:"priority"`
	Reporter    string       `json:"reporter"`
	Handler     string       `json:"handler,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	ResolvedAt  *time.Time   `json:"resolved_at,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
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

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
