package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"device-ticket-service/model"
	"device-ticket-service/repository"
)

type TicketHandler struct {
	repo *repository.TicketRepository
}

func NewTicketHandler(repo *repository.TicketRepository) *TicketHandler {
	return &TicketHandler{repo: repo}
}

func writeJSON(w http.ResponseWriter, code int, resp model.ApiResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

func successResponse(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, model.ApiResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func errorResponse(w http.ResponseWriter, httpCode int, msg string) {
	writeJSON(w, httpCode, model.ApiResponse{
		Code:    httpCode,
		Message: msg,
	})
}

func (h *TicketHandler) GetProductionLines(w http.ResponseWriter, r *http.Request) {
	lines := h.repo.GetAllLines()
	successResponse(w, lines)
}

func (h *TicketHandler) GetFaultTypes(w http.ResponseWriter, r *http.Request) {
	types := repository.GetAllFaultTypes()
	successResponse(w, types)
}

func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "请求参数解析失败: "+err.Error())
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		errorResponse(w, http.StatusBadRequest, "工单标题不能为空")
		return
	}
	if strings.TrimSpace(req.DeviceID) == "" {
		errorResponse(w, http.StatusBadRequest, "设备ID不能为空")
		return
	}
	if strings.TrimSpace(req.LineID) == "" {
		errorResponse(w, http.StatusBadRequest, "产线编号不能为空")
		return
	}
	if strings.TrimSpace(req.Reporter) == "" {
		errorResponse(w, http.StatusBadRequest, "报修人不能为空")
		return
	}

	ticket, err := h.repo.CreateTicket(&req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	successResponse(w, ticket)
}

func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tickets/")
	if id == "" {
		errorResponse(w, http.StatusBadRequest, "工单ID不能为空")
		return
	}

	ticket, err := h.repo.GetTicketByID(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	successResponse(w, ticket)
}

func (h *TicketHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	lineID := r.URL.Query().Get("line_id")
	faultType := model.FaultType(r.URL.Query().Get("fault_type"))
	status := model.TicketStatus(r.URL.Query().Get("status"))

	tickets := h.repo.GetAllTickets(lineID, faultType, status)
	successResponse(w, tickets)
}

func (h *TicketHandler) UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tickets/")
	id = strings.TrimSuffix(id, "/status")
	if id == "" {
		errorResponse(w, http.StatusBadRequest, "工单ID不能为空")
		return
	}

	var req model.UpdateTicketStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "请求参数解析失败: "+err.Error())
		return
	}

	if req.Status == "" {
		errorResponse(w, http.StatusBadRequest, "工单状态不能为空")
		return
	}

	ticket, err := h.repo.UpdateTicketStatus(id, &req)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	successResponse(w, ticket)
}

func (h *TicketHandler) DeleteTicket(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tickets/")
	if id == "" {
		errorResponse(w, http.StatusBadRequest, "工单ID不能为空")
		return
	}

	if err := h.repo.DeleteTicket(id); err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	successResponse(w, nil)
}
