package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"microservice/audit-service/endpoint"
	"microservice/util/logger"
	"microservice/util/model"
	"microservice/util/requestid"

	httptransport "github.com/go-kit/kit/transport/http"
)

func RegisterHTTPServer(endpoints endpoint.AuditEndpoints, port string) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			return requestid.MiddlewareRequestIdHTTP(request)
		}),
	}

	getListHandler := httptransport.NewServer(
		endpoints.GetListAuditLogEndpoint,
		decodeGetListAuditLogRequest,
		encodeResponse,
		options...,
	)

	http.Handle("/v1/audits", getListHandler)

	logger.Info("Audit HTTP Service started on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.LogError(context.Background(), "http.ListenAndServe", err)
		os.Exit(1)
	}
}

func decodeGetListAuditLogRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.GetListAuditLogRequest
	// default values
	req.Page = 1
	req.Show = 10

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &req.Page)
	}
	if s := r.URL.Query().Get("show"); s != "" {
		fmt.Sscanf(s, "%d", &req.Show)
	}

	req.ItemID = r.URL.Query().Get("item_id")
	req.ItemType = r.URL.Query().Get("item_type")

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	reqID := requestid.GetRequestId(ctx)
	
	finalResp := model.Response{
		Code:      http.StatusOK,
		Message:   "Success",
		RequestID: reqID,
		Result:    response,
	}

	return json.NewEncoder(w).Encode(finalResp)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	
	reqID := requestid.GetRequestId(ctx)
	
	errorResp := model.Response{
		Code:      http.StatusInternalServerError,
		Message:   err.Error(),
		RequestID: reqID,
	}

	json.NewEncoder(w).Encode(errorResp)
}
