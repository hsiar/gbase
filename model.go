package gbase

type PageRequest struct {
	Page  int `json:"page"`  //页数
	Limit int `json:"limit"` //每页条数
}

// 只传一个ID参数的请求
type IdRequest struct {
	Id int64 `json:"id" vd:"$>0;msg:'id参数缺失'"`
}

// 只传一个ID参数的请求,用于snowflake的ID
type SnowFlakeIdRequest struct {
	Id int64 `json:"id,string" vd:"$>0;msg:'id参数缺失'"`
}
