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
	mu          sync.RWMutex
	tickets     map[string]*model.DeviceFaultTicket
	lines       map[string]*model.ProductionLine
	parts       map[string]*model.Part
	consumptions map[string]*model.PartConsumption
}

func NewTicketRepository() *TicketRepository {
	repo := &TicketRepository{
		tickets:     make(map[string]*model.DeviceFaultTicket),
		lines:       make(map[string]*model.ProductionLine),
		parts:       make(map[string]*model.Part),
		consumptions: make(map[string]*model.PartConsumption),
	}
	repo.initProductionLines()
	repo.initParts()
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

func (r *TicketRepository) initParts() {
	now := time.Now()
	parts := []model.Part{
		{ID: "P001", Name: "主轴轴承", Code: "BRG-001", Category: "传动部件", Unit: "个", UnitPrice: 580.00, Stock: 50, CreatedAt: now},
		{ID: "P002", Name: "伺服电机", Code: "MTR-001", Category: "电气部件", Unit: "台", UnitPrice: 3200.00, Stock: 20, CreatedAt: now},
		{ID: "P003", Name: "接近开关", Code: "SNS-001", Category: "传感器", Unit: "个", UnitPrice: 120.00, Stock: 100, CreatedAt: now},
		{ID: "P004", Name: "PLC控制器", Code: "PLC-001", Category: "电气部件", Unit: "台", UnitPrice: 2800.00, Stock: 15, CreatedAt: now},
		{ID: "P005", Name: "液压油滤芯", Code: "FLT-001", Category: "过滤器", Unit: "个", UnitPrice: 85.00, Stock: 80, CreatedAt: now},
		{ID: "P006", Name: "同步带", Code: "BELT-001", Category: "传动部件", Unit: "条", UnitPrice: 320.00, Stock: 40, CreatedAt: now},
		{ID: "P007", Name: "压力传感器", Code: "SNS-002", Category: "传感器", Unit: "个", UnitPrice: 450.00, Stock: 30, CreatedAt: now},
		{ID: "P008", Name: "接触器", Code: "ELC-001", Category: "电气部件", Unit: "个", UnitPrice: 180.00, Stock: 60, CreatedAt: now},
	}
	for i := range parts {
		r.parts[parts[i].ID] = &parts[i]
	}
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "TKT" + hex.EncodeToString(b)
}

func generatePartID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "PRT" + hex.EncodeToString(b)
}

func generateConsumptionID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "CNS" + hex.EncodeToString(b)
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
	if req.Status == model.StatusClosed {
		now := time.Now()
		ticket.ClosedAt = &now
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

func (r *TicketRepository) CreatePart(req *model.CreatePartRequest) (*model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range r.parts {
		if p.Code == req.Code {
			return nil, errors.New("配件编码已存在")
		}
	}

	now := time.Now()
	part := &model.Part{
		ID:        generatePartID(),
		Name:      req.Name,
		Code:      req.Code,
		Category:  req.Category,
		Unit:      req.Unit,
		UnitPrice: req.UnitPrice,
		Stock:     req.Stock,
		CreatedAt: now,
	}

	r.parts[part.ID] = part
	return part, nil
}

func (r *TicketRepository) GetAllParts(category string) []model.Part {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]model.Part, 0)
	for _, p := range r.parts {
		if category != "" && p.Category != category {
			continue
		}
		result = append(result, *p)
	}
	return result
}

func (r *TicketRepository) GetPartByID(id string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.parts[id]
	if !ok {
		return nil, errors.New("配件不存在")
	}
	return part, nil
}

func (r *TicketRepository) AddConsumption(ticketID string, req *model.AddConsumptionRequest) (*model.PartConsumption, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ticket, ok := r.tickets[ticketID]
	if !ok {
		return nil, errors.New("工单不存在")
	}

	part, ok := r.parts[req.PartID]
	if !ok {
		return nil, errors.New("配件不存在")
	}

	if part.Stock < req.Quantity {
		return nil, errors.New("库存不足")
	}

	part.Stock -= req.Quantity
	totalPrice := float64(req.Quantity) * part.UnitPrice

	consumption := &model.PartConsumption{
		ID:           generateConsumptionID(),
		TicketID:     ticketID,
		PartID:       part.ID,
		PartName:     part.Name,
		PartCode:     part.Code,
		PartCategory: part.Category,
		Quantity:     req.Quantity,
		UnitPrice:    part.UnitPrice,
		TotalPrice:   totalPrice,
		Operator:     req.Operator,
		Remark:       req.Remark,
		CreatedAt:    time.Now(),
	}

	r.consumptions[consumption.ID] = consumption

	if ticket.PartConsumptions == nil {
		ticket.PartConsumptions = make([]model.PartConsumption, 0)
	}
	ticket.PartConsumptions = append(ticket.PartConsumptions, *consumption)
	ticket.TotalMaterialCost += totalPrice
	ticket.UpdatedAt = time.Now()

	return consumption, nil
}

func (r *TicketRepository) GetTicketConsumptions(ticketID string) []model.PartConsumption {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]model.PartConsumption, 0)
	for _, c := range r.consumptions {
		if c.TicketID == ticketID {
			result = append(result, *c)
		}
	}
	return result
}

func (r *TicketRepository) CloseTicketWithConsumptions(ticketID string, req *model.CloseTicketRequest) (*model.DeviceFaultTicket, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ticket, ok := r.tickets[ticketID]
	if !ok {
		return nil, errors.New("工单不存在")
	}

	if req.Consumptions != nil && len(req.Consumptions) > 0 {
		for _, consReq := range req.Consumptions {
			part, ok := r.parts[consReq.PartID]
			if !ok {
				return nil, errors.New("配件ID不存在: " + consReq.PartID)
			}
			if part.Stock < consReq.Quantity {
				return nil, errors.New("配件库存不足: " + part.Name)
			}
		}

		for _, consReq := range req.Consumptions {
			part := r.parts[consReq.PartID]
			part.Stock -= consReq.Quantity
			totalPrice := float64(consReq.Quantity) * part.UnitPrice

			consumption := &model.PartConsumption{
				ID:           generateConsumptionID(),
				TicketID:     ticketID,
				PartID:       part.ID,
				PartName:     part.Name,
				PartCode:     part.Code,
				PartCategory: part.Category,
				Quantity:     consReq.Quantity,
				UnitPrice:    part.UnitPrice,
				TotalPrice:   totalPrice,
				Operator:     consReq.Operator,
				Remark:       consReq.Remark,
				CreatedAt:    time.Now(),
			}

			r.consumptions[consumption.ID] = consumption

			if ticket.PartConsumptions == nil {
				ticket.PartConsumptions = make([]model.PartConsumption, 0)
			}
			ticket.PartConsumptions = append(ticket.PartConsumptions, *consumption)
			ticket.TotalMaterialCost += totalPrice
		}
	}

	ticket.Status = req.Status
	if req.Handler != "" {
		ticket.Handler = req.Handler
	}
	now := time.Now()
	ticket.UpdatedAt = now

	if req.Status == model.StatusClosed {
		ticket.ClosedAt = &now
	}
	if req.Status == model.StatusResolved {
		ticket.ResolvedAt = &now
	}

	return ticket, nil
}

func (r *TicketRepository) GetPartUsageStats(lineID string, faultType model.FaultType, startDate, endDate string) []model.PartUsageStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	statsMap := make(map[string]*model.PartUsageStats)

	for _, cons := range r.consumptions {
		ticket, ok := r.tickets[cons.TicketID]
		if !ok {
			continue
		}

		if lineID != "" && ticket.LineID != lineID {
			continue
		}
		if faultType != "" && ticket.FaultType != faultType {
			continue
		}

		consDate := cons.CreatedAt.Format("2006-01-02")
		if startDate != "" && consDate < startDate {
			continue
		}
		if endDate != "" && consDate > endDate {
			continue
		}

		if statsMap[cons.PartID] == nil {
			statsMap[cons.PartID] = &model.PartUsageStats{
				PartID:       cons.PartID,
				PartName:     cons.PartName,
				PartCode:     cons.PartCode,
				PartCategory: cons.PartCategory,
			}
		}
		statsMap[cons.PartID].TotalUsed += cons.Quantity
		statsMap[cons.PartID].TotalAmount += cons.TotalPrice
		statsMap[cons.PartID].UsageCount++
	}

	result := make([]model.PartUsageStats, 0, len(statsMap))
	for _, s := range statsMap {
		result = append(result, *s)
	}
	return result
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

func GetAllPartCategories() []string {
	return []string{"传动部件", "电气部件", "传感器", "过滤器", "液压元件", "气动元件", "密封件", "紧固件", "其他"}
}
