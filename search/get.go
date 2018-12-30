package search

import (
	"net/http"
	"net/url"
	"sync"

	mathcers "eatinlife.com/gaodespider/matchers"

	"eatinlife.com/serverutil/logger"

	"strconv"
)

var requestkey = "8c56880afd72ba62ad85d8364e185d6a"

var params = map[string]string{
	"keywords":  "凌空soho",
	"types":     "050000",
	"city":      "上海/shanghai/021/310105",
	"citylimit": "true",
}

//RequestGao is 请求数据
func requestGao(page int, limit int) (mathcers.RestaurantResponse, error) {
	api, _ := url.Parse("https://restapi.amap.com/v3/place/text")
	query := api.Query()
	query.Set("key", requestkey)
	for i, value := range params {
		query.Set(i, value)
	}
	query.Set("page", strconv.Itoa(page))
	query.Set("offset", strconv.Itoa(limit))
	api.RawQuery = query.Encode()
	response, err := http.Get(api.String())
	if err != nil {
		logger.Error.Panicln(err.Error())
		return mathcers.RestaurantResponse{}, err
	}
	defer response.Body.Close()
	resultjson, err := mathcers.FormatResult(response)
	return resultjson, err
}

var wg sync.WaitGroup

//PaginationGet is 分页调用
func PaginationGet() []mathcers.RestaurantInfo {
	firstresult, err := requestGao(1, 20)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	count, err := strconv.ParseInt(firstresult.Count, 10, 16)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	resultchan := make(chan mathcers.RestaurantInfo, count)
	results := make([]mathcers.RestaurantInfo, count)
	maxpage := int(count) / 20
	if (int(count) - maxpage*20) > 0 {
		maxpage++
	}
	logger.Info.Printf("maxpage:%d,count:%d", maxpage, count)
	if count > 20 {
		wg.Add(maxpage)
		logger.Info.Println("wgadd", maxpage)
		go func() {
			defer wg.Done()
			for _, v := range firstresult.Pois {
				logger.Info.Println("insert", v.Name)
				resultchan <- v
			}
		}()
		for i := 2; i <= maxpage; i++ {
			go func(index int) {
				defer wg.Done()
				nextresult, err := requestGao(index, 20)
				if err != nil {
					logger.Error.Println(err.Error())
				}
				for _, v := range nextresult.Pois {
					logger.Info.Printf("insert%d name=%s \n", index, v.Name)
					resultchan <- v
				}
			}(i)
		}
		wg.Wait()
		close(resultchan)
		for i := 1; i <= int(count); i++ {
			t, ok := <-resultchan
			if ok {
				results = append(results, t)
			}
		}
	}
	return results
}
