package common

import (
	"fmt"
	"net/http"
)

type CommonResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CommonListResponse struct {
	Status     int              `json:"status"`
	Message    string           `json:"message"`
	Data       interface{}      `json:"data"`
	Pagination CommonPagination `json:"pagination"`
}

type CommonPagination struct {
	Page      int `json:"page"`
	NextPage  int `json:"next_page"`
	PrevPage  int `json:"prev_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

func GetResponseSuccess(data interface{}) CommonResponse {
	return CommonResponse{
		Data:    data,
		Message: "Success",
		Status:  200,
	}
}

func GetListResponseSuccess(data interface{}, isListEmpty bool, pagination CommonPagination) CommonListResponse {
	if isListEmpty {
		return CommonListResponse{
			Data:       data,
			Message:    "No Data Found",
			Status:     401,
			Pagination: pagination,
		}
	}
	return CommonListResponse{
		Data:       data,
		Message:    "Success",
		Status:     200,
		Pagination: pagination,
	}
}

func GetResponseError(status int) CommonResponse {
	return CommonResponse{
		Data:    nil,
		Message: http.StatusText(status),
		Status:  status,
	}
}

func GetPagination(limit int, page int, totalData int) CommonPagination {
	totalPage := 1
	nextPage := page + 1
	prevPage := 0

	if limit > 0 {
		totalPage = totalData / limit
	}
	if page > 0 {
		prevPage = page - 1
	}

	fmt.Println("TotalData GetPagination = ", totalData)

	pagination := CommonPagination{
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
		Page:      page,
		TotalData: totalData,
	}
	return pagination
}
