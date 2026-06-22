package repository

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"device-ticket-service/model"
)

type TicketRepository struct {
	mu     sync.RWMutex
	tickets map[string]*model.DeviceFaultTicket
	lines   map[string]*model.ProductionLine
}

func NewTicketRepository() *TicketRepository {
	repo := &TicketRepository{
		tickets: make(map[string]*model.DeviceFaultTicket),
		lines:   make(map[string]*model.ProductionLine),
	}
	repo.initProductionLines()
	return repo
}

func (r *TicketRepository) initProductionLines() {
	lines := []model.ProductionLine{
		{ID: "L001", Name: "一号装配线", Code: "PL-ASM-001"},
		{ID: "L002", Name: "二号装配线", Code: "PL-ASM-002"},
		{ID: "L003", Name: "焊接生产线", Code: "PL-WLD-001"},
		{ID: "L004", Name: "喷涂生产线", Code: "PL-PNT-001"},
		{ID: "L005", Name: "检测包装线", Code: "PL-PKG-001"},
	}
	for i := range lines {
		r.lines[lines[i].ID] = &lines[i]
	}
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "TKT" + hex.EncodeToString(b)
}

func (r *TicketRepository) GetAllLines() []model.ProductionLine {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]model.ProductionLine, 0, len(r.lines))
	for _, line := range r.lines {
		result = append(result, *line)
	}
	return result
}

func (r *TicketRepository) GetLineByID(lineID string) (*model.ProductionLine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	line, ok := r.lines[lineID]
	if !ok {
		return nil, errors.New("产线不存在")
	}
	return line, nil
}

func (r *TicketRepository) CreateTicket(req *model.CreateTicketRequest) (*model.DeviceFaultTicket, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	line, ok := r.lines[req.LineID]
	if !ok {
		return nil, errors.New("产线编号不存在")
	}

	if !isValidFaultType(req.FaultType) {
		return nil, errors.New("故障类型无效")
	}

	now := time.Now()
	priority := req.Priority
	if priority < 1 {
		priority = 3
	}

	ticket := &model.DeviceFaultTicket{
		ID:          generateID(),
		Title:       req.Title,
		Description: req.Description,
		DeviceID:    req.DeviceID,
		DeviceName:  req.DeviceName,
		LineID:      line.ID,
		LineCode:    line.Code,
		LineName:    line.Name,
		FaultType:   req.FaultType,
		Status:      model.StatusPending,
		Priority:    priority,
		Reporter:    req.Reporter,
		CreatedAt:   now,
		UpdatedAt:   now,
		Tags:        req.Tags,
	}

	r.tickets[ticket.ID] = ticket
	return ticket, nil
}

func (r *TicketRepository) GetTicketByID(id string) (*model.DeviceFaultTicket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ticket, ok := r.tickets[id]
	if !ok {
		return nil, errors.New("工单不存在")
	}
	return ticket, nil
}

func (r *TicketRepository) GetAllTickets(lineID string, faultType model.FaultType, status model.TicketStatus) []model.DeviceFaultTicket {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]model.DeviceFaultTicket, 0)
	for _, t := range r.tickets {
		if lineID != "" && t.LineID != lineID {
			continue
		}
		if faultType != "" && t.FaultType != faultType {
			continue
		}
		if status != "" && t.Status != status {
			continue
		}
		result = append(result, *t)
	}
	return result
}

func (r *TicketRepository) UpdateTicketStatus(id string, req *model.UpdateTicketStatusRequest) (*model.DeviceFaultTicket, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ticket, ok := r.tickets[id]
	if !ok {
		return nil, errors.New("工单不存在")
	}

	ticket.Status = req.Status
	if req.Handler != "" {
		ticket.Handler = req.Handler
	}
	ticket.UpdatedAt = time.Now()

	if req.Status == model.StatusResolved {
		now := time.Now()
		ticket.ResolvedAt = &now
	}

	return ticket, nil
}

func (r *TicketRepository) DeleteTicket(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tickets[id]; !ok {
		return errors.New("工单不存在")
	}
	delete(r.tickets, id)
	return nil
}

func isValidFaultType(ft model.FaultType) bool {
	switch ft {
	case model.FaultTypeMechanical,
		model.FaultTypeElectrical,
		model.FaultTypeSoftware,
		model.FaultTypeSensor,
		model.FaultTypeCommunication,
		model.FaultTypeOther:
		return true
	default:
		return false
	}
}

func GetAllFaultTypes() []map[string]string {
	return []map[string]string{
		{"key": string(model.FaultTypeMechanical), "label": "机械故障"},
		{"key": string(model.FaultTypeElectrical), "label": "电气故障"},
		{"key": string(model.FaultTypeSoftware), "label": "软件故障"},
		{"key": string(model.FaultTypeSensor), "label": "传感器故障"},
		{"key": string(model.FaultTypeCommunication), "label": "通信故障"},
		{"key": string(model.FaultTypeOther), "label": "其他"},
	}
}
