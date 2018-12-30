package mathcers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"eatinlife.com/serverutil/logger"
)

//RestaurantImage is 餐厅图片
type RestaurantImage struct {
	URL string `json:"url"` //餐厅url
}

//RestaurantInfo is 餐厅信息
type RestaurantInfo struct {
	Name     string            `json:"name"`
	Types    string            `json:"type"`
	Address  string            `json:"address"`
	Location string            `json:"location"`
	Tel      interface{}       `json:"tel"`
	Pname    string            `json:"pname"`
	Cityname string            `json:"cityname"`
	Adname   string            `json:"adname"`
	Photos   []RestaurantImage `json:"photos"`
}

//RestaurantResponse is 接口POIs
type RestaurantResponse struct {
	Pois  []RestaurantInfo `json:"pois"`
	Count string           `json:"count"`
}

//FormatResult is 输入接口响应，返回 RestaurantInfo list
func FormatResult(res *http.Response) (RestaurantResponse, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	result := buf.String()
	var jsonformat RestaurantResponse
	var err = json.Unmarshal([]byte(result), &jsonformat)
	if err != nil {
		logger.Info.Println(err.Error())
	}
	return jsonformat, nil
}

// "id": "B0FFGSRCFM",
//             "parent": "B0FFGTU5M6",
//             "name": "沙拉疯saladfun(凌空soho店)",
//             "type": "餐饮服务;快餐厅;快餐厅",
//             "typecode": "050300",
//             "biz_type": "diner",
//             "address": "金钟路968号凌空SOHO内17号楼一层商铺17-102号",
//             "location": "121.350954,31.222486",
//             "tel": "021-52728882",
//             "distance": [],
//             "biz_ext": {
//                 "rating": "4.5",
//                 "cost": "32.00",
//                 "meal_ordering": "0"
//             },
//             "pname": "上海市",
//             "cityname": "上海市",
//             "adname": "长宁区",
//             "importance": [],
//             "shopid": [],
//             "shopinfo": "0",
//             "poiweight": [],
//             "photos": [
//                 {
//                     "url": "http://store.is.autonavi.com/showpic/8db0c5c1f1e6e18a48045718bf998d5e",
//                     "title": [],
//                     "provider": []
//                 },
//                 {
//                     "url": "http://store.is.autonavi.com/showpic/00fd64eae917f7e7938ff99b5e50fc80",
//                     "title": [],
//                     "provider": []
//                 },
//                 {
//                     "url": "http://store.is.autonavi.com/showpic/8e29e4062e601f1d5968bf0199d31e34",
//                     "title": [],
//                     "provider": []
//                 }
//             ]
