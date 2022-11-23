package sensor

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"timeseries/pkg/api"
	"timeseries/pkg/models"
	"timeseries/pkg/service/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const projectQuery = "PROJECT_ID as project_id, PROJECT_NAME as project_name"

func GetProjects(ctx *gin.Context) {
	var respBody []api.ProjectResp
	if err := mysql.GetClient().Model(&models.ProjectIdName{}).Select(projectQuery).Find(&respBody).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, api.ReplyError{Code: api.InternelError, Msg: err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, api.ReplyJson{Data: respBody})
}

const locationQuery = "select * from site_location_name where project_id=?"

func GetLocations(ctx *gin.Context) {
	projectId, err := strconv.Atoi(ctx.Query("project_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.QueryParamError, Msg: "project_id in query not found"})
		ctx.Abort()
		return
	}
	var _locations []models.SiteLocationName
	var respBody []*api.LocationNode

	if err := mysql.GetClient().Raw(locationQuery, projectId).Find(&_locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, api.ReplyJson{Data: respBody})
		} else {
			ctx.JSON(http.StatusInternalServerError, api.ReplyError{Code: api.InternelError, Msg: "unknow server error"})
		}
		ctx.Abort()
		return
	}

	// data transform

	var nodes []*api.LocationNode
	for _, l := range _locations {
		_node := &api.LocationNode{
			Location1Id:  l.Location1Id,
			Location2Id:  l.Location2Id,
			Location3Id:  l.Location3Id,
			Location4Id:  l.Location4Id,
			LocationName: l.LocationName,
		}
		if l.Location2Id == 0 {
			_node.Level = 1
			nodes = append(nodes, _node)
		} else if l.Location3Id == 0 {
			_node.Level = 2
			nodes = append(nodes, _node)
		} else if l.Location4Id == 0 {
			_node.Level = 3
			nodes = append(nodes, _node)
		} else {
			_node.Level = 4
			nodes = append(nodes, _node)
		}
	}

	respBody = buildIndexTree(nodes)
	ctx.JSON(http.StatusOK, api.ReplyJson{Data: respBody})
}

const sensorQuery = "select s.SENSOR_MAC as sensor_mac, s.TYPE_ID as type_id, st.TYPE_NAME as type_name, " +
	"sl_1.location_name as location1, sl_2.location_name as location2, sl_3.location_name as location3, sl_4.location_name as location4 " +
	"from sensor_location s  " +
	"left join sensor_type st on s.TYPE_ID = st.ID  " +
	"left join site_location_name sl_1 on s.project_id=sl_1.project_id and s.location_1_id=sl_1.LOCATION_1_ID and sl_1.location_2_id=0  " +
	"left join site_location_name sl_2 on s.project_id=sl_2.project_id and s.location_1_id=sl_2.LOCATION_1_ID and s.location_2_id=sl_2.LOCATION_2_ID and sl_2.location_3_id=0  " +
	"left join site_location_name sl_3 on s.project_id=sl_3.project_id and s.location_1_id=sl_3.LOCATION_1_ID and s.location_2_id=sl_3.LOCATION_2_ID and s.location_3_id=sl_3.LOCATION_3_ID and sl_3.location_4_id=0 " +
	"left join site_location_name sl_4 on s.project_id=sl_4.project_id and s.location_1_id=sl_4.LOCATION_1_ID and s.location_2_id=sl_4.LOCATION_2_ID and s.location_3_id=sl_4.LOCATION_3_ID and s.location_4_id=sl_4.LOCATION_4_ID  " +
	"where "

const sensorFilter = "s.%s=?"

func GetSensors(ctx *gin.Context) {
	var reqBody api.SensorReq

	if projectId := ctx.Query("project_id"); projectId != "" {
		reqBody.ProjectId = &projectId
	} else {
		ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.RequestBodyError, Msg: "the project_id is necessnary"})
		ctx.Abort()
		return
	}

	if location1 := ctx.Query("location1"); location1 != "" {
		reqBody.Location1 = &location1
	}

	if location2 := ctx.Query("location2"); location2 != "" {
		reqBody.Location2 = &location2
	}

	if location3 := ctx.Query("location3"); location3 != "" {
		reqBody.Location3 = &location3
	}

	if location4 := ctx.Query("location4"); location4 != "" {
		reqBody.Location4 = &location4
	}

	var filters []string     // 查询过滤器
	var values []interface{} // 过滤器对应值
	t := reflect.TypeOf(reqBody)
	v := reflect.ValueOf(reqBody)
	// 反射遍历查询字段
	for i := 0; i < t.NumField(); i++ {
		value, ok := v.Field(i).Interface().(*string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.RequestBodyError, Msg: "parse request body failed"})
			ctx.Abort()
			return
		}
		if value == nil {
			continue
		}
		filters = append(filters, fmt.Sprintf(sensorFilter, t.Field(i).Tag.Get("filter")))
		values = append(values, *value)
	}
	queryString := sensorQuery + strings.Join(filters, " and ")
	var respBody []api.SensorResp
	if err := mysql.GetClient().Raw(queryString, values...).Find(&respBody).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, api.ReplyJson{Data: respBody})
		} else {
			ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.InternelError, Msg: "unknow server error"})
		}
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, api.ReplyJson{Data: respBody})
}

const measurementQuery = "select receive_no, gather_type, gather_type_name, unit from sensor_gather_type where sensor_type_id=?"

func GetMeasurements(ctx *gin.Context) {
	typeId := ctx.Query("type_id")
	if typeId == "" {
		ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.QueryParamError, Msg: "type_id in query not found"})
		ctx.Abort()
		return
	}
	var respBody []api.MeasurementResp
	if err := mysql.GetClient().Raw(measurementQuery, typeId).Find(&respBody).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, api.ReplyJson{Data: respBody})
		} else {
			ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.InternelError, Msg: "unknow server error"})
		}
	}
	ctx.JSON(http.StatusOK, api.ReplyJson{Data: respBody})
}
