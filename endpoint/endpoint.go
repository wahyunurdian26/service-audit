package endpoint

import (
	"context"
	"math"

	"github.com/go-kit/kit/endpoint"
	"github.com/wahyunurdian26/service-audit/model"
	"github.com/wahyunurdian26/service-audit/service"
)

type AuditEndpoints struct {
	RecordEventEndpoint     endpoint.Endpoint
	GetListAuditLogEndpoint endpoint.Endpoint
}

func MakeEndpoints(s service.AuditService) AuditEndpoints {
	return AuditEndpoints{
		RecordEventEndpoint:     makeRecordEventEndpoint(s),
		GetListAuditLogEndpoint: makeGetListAuditLogEndpoint(s),
	}
}

func makeRecordEventEndpoint(s service.AuditService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*model.AuditLog)
		err := s.RecordEvent(ctx, req)
		return nil, err
	}
}

type GetListAuditLogRequest struct {
	Page     int    `json:"page"`
	Show     int    `json:"show"`
	ItemID   string `json:"item_id"`
	ItemType string `json:"item_type"`
}

type GetListAuditLogResponse struct {
	Meta Meta        `json:"meta"`
	Data []AuditData `json:"data"`
}

type Meta struct {
	Page      int `json:"page"`
	Show      int `json:"show"`
	TotalPage int `json:"total_page"`
}

type AuditData struct {
	ItemID        string      `json:"item_id"`
	Event         string      `json:"event"`
	CreatedAt     string      `json:"created_at"`
	WhodunnitID   int64       `json:"whodunnit_id"`
	WhodunnitName string      `json:"whodunnit_name"`
	Activities    interface{} `json:"activities,omitempty"`
}

func makeGetListAuditLogEndpoint(s service.AuditService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetListAuditLogRequest)
		logs, total, err := s.GetListAuditLog(ctx, req.Page, req.Show, req.ItemID, req.ItemType)
		if err != nil {
			return nil, err
		}

		totalPage := int(math.Ceil(float64(total) / float64(req.Show)))
		
		var data []AuditData
		for _, l := range logs {
			data = append(data, AuditData{
				ItemID:        l.ItemID,
				Event:         l.Event,
				CreatedAt:     l.CreatedAt.Format("02 Jan 2006, 15:04:05"),
				WhodunnitID:   l.WhodunnitID,
				WhodunnitName: l.WhodunnitName,
				Activities:    l.Activities,
			})
		}

		return GetListAuditLogResponse{
			Meta: Meta{
				Page:      req.Page,
				Show:      req.Show,
				TotalPage: totalPage,
			},
			Data: data,
		}, nil
	}
}
