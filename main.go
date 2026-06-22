package main

import (
	"log"
	"net/http"

	"device-ticket-service/handler"
	"device-ticket-service/repository"
)

func main() {
	repo := repository.NewTicketRepository()
	h := handler.NewTicketHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/lines", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.GetProductionLines(w, r)
	})

	mux.HandleFunc("/api/fault-types", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.GetFaultTypes(w, r)
	})

	mux.HandleFunc("/api/tickets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateTicket(w, r)
		case http.MethodGet:
			h.ListTickets(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/tickets/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if len(path) > len("/api/tickets/") && path[len("/api/tickets/"):] != "" {
			if len(path) > len("/api/tickets//status") && path[len(path)-len("/status"):] == "/status" {
				if r.Method != http.MethodPut {
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					return
				}
				h.UpdateTicketStatus(w, r)
				return
			}

			switch r.Method {
			case http.MethodGet:
				h.GetTicket(w, r)
			case http.MethodDelete:
				h.DeleteTicket(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			switch r.Method {
			case http.MethodGet:
				h.ListTickets(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	addr := ":8080"
	log.Printf("Server starting on %s ...", addr)
	log.Printf("API Endpoints:")
	log.Printf("  GET    /api/lines           - 获取所有产线列表")
	log.Printf("  GET    /api/fault-types     - 获取所有故障类型")
	log.Printf("  POST   /api/tickets         - 创建设备故障工单")
	log.Printf("  GET    /api/tickets         - 查询工单列表（支持line_id, fault_type, status筛选）")
	log.Printf("  GET    /api/tickets/{id}    - 根据ID查询工单详情")
	log.Printf("  PUT    /api/tickets/{id}/status - 更新工单状态")
	log.Printf("  DELETE /api/tickets/{id}    - 删除工单")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
