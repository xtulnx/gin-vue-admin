package apimodel

// 分页相关
//
// jason.liao 2020.07.01

// PageReq 分页
type PageReq struct {
	Page     int `form:"page" json:"page" query:"page"  example:"1"`             // 页号，从1开始
	PageSize int `form:"pageSize" json:"pageSize" query:"pageSize" example:"10"` // 分页大小
}

// FixPageSize 修正页码数，设置默认值
func (P *PageReq) FixPageSize(num, size int) {
	if P.Page <= 0 {
		P.Page = num
	}
	if P.PageSize <= 0 {
		P.PageSize = size
	}
}

type PageResp struct {
	Page      int `json:"page" example:"1"`      // 页号，从1开始
	PageSize  int `json:"pageSize" example:"10"` // 分页大小
	Total     int `json:"total" example:"51"`    // 总记录数
	TotalPage int `json:"totalPage" example:"6"` // 总页数
}

// SetPageSize 设置页码数，并计算返回 offset ，如果 offset == -1 表示不可用
func (P *PageResp) SetPageSize(num, size int, cnt int64) (offset int) {
	P.Page, P.PageSize = num, size
	if size > 0 {
		P.Total = int(cnt)
		P.TotalPage = (int(cnt) + size - 1) / size
		offset = (num - 1) * size
		if offset < 0 {
			offset = 0
		}
		if offset < P.Total {
			return offset
		}
	}
	return -1
}
