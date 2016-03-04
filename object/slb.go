package object

import (
	"fmt"
	"reflect"
	"time"
	"third/gorm"
)

const (
	LOADBALANCER_TABLE_NAME = "slbs"
	DB_NAME = "kube_apiproxy"
)

type SLB struct {
	GroupName   string `gorm:"primary_key;column:groupname" sql:"type:varchar(128);not null" json:"groupname"`
	LoadBalancerId   string `gorm:"column:loadBalancerId" sql:"type:varchar(128);not null" json:"loadBalancerId"`
	LoadBalancerName string `gorm:"column:loadBalancerName" sql:"type:varchar(128)" json:"loadBalancerName"`
	Ip 	string `gorm:"column:Ip" sql:"type:varchar(128);not null" json:"Ip"`

	Created time.Time `gorm:"column:created" sql:"type:datetime;index:idx_created;not null;default:" json:"created"`
	Updated time.Time `grom:"column:updated" sql:"type:datetime;index:idx_updated;not null;default:" json:"updated"`
}

type SLBList []SLB

func (lb SLB) TableName() string {
	return fmt.Sprintf("%s.%s", DB_NAME, LOADBALANCER_TABLE_NAME)
}

func (lb SLB) String() string {
	lb_value := reflect.ValueOf(lb)
	lb_field_num := lb_value.NumField()
	lb_str := "loadBalancer:"

	for i := 0; i < lb_field_num; i++ {
		field_value := lb_value.Field(i)
		if !reflect.DeepEqual(field_value.Interface(), reflect.Zero(field_value.Type()).Interface()) {
			lb_str = fmt.Sprintf("%s[%s:%v]", lb_str, lb_value.Type().Field(i).Name,
				field_value.Interface())
		}
	}
	return lb_str
}

func (lb *SLB) Insert(db *gorm.DB) error {
	Info("Start insert new slb..")

	_, offset := time.Now().Zone()
	lb.Created = time.Unix(time.Now().Unix()+int64(offset), 0)
	lb.Updated = time.Unix(time.Now().Unix()+int64(offset), 0)

	if err := db.Create(lb).Error; err != nil {
		return DbError
	}
	return nil
}


func (lb *SLB) Fetch(db *gorm.DB) error {
	err := db.Where(lb).Find(&lb).Error
	if err == gorm.RecordNotFound {
		return RecordNotFoundError
	} else if err != nil {
		Warning("fetch loadBalancer failed![%s]", err.Error())
		return DbError
	}
	return nil
}

