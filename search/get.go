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

//PaginationGet is 分页调用

func initProducer(r mathcers.RestaurantResponse) <-chan mathcers.RestaurantInfo {
	out := make(chan mathcers.RestaurantInfo, 20)
	go func() {
		defer close(out)
		for _, v := range r.Pois {
			logger.Info.Printf("insert name=%s \n", v.Name)
			out <- v
		}
	}()
	return out
}

func getPageProducer(page int, limit int) <-chan mathcers.RestaurantInfo {
	out := make(chan mathcers.RestaurantInfo, 20)
	go func() {
		defer close(out)
		nextresult, err := requestGao(page, limit)
		if err != nil {
			logger.Error.Println(err.Error())
		}
		for _, v := range nextresult.Pois {
			logger.Info.Printf("insert%d name=%s \n", page, v.Name)
			out <- v
		}
	}()
	return out
}

func mergeResult(count int, cs ...<-chan mathcers.RestaurantInfo) <-chan mathcers.RestaurantInfo {
	out := make(chan mathcers.RestaurantInfo, count)
	var wg sync.WaitGroup
	collect := func(in <-chan mathcers.RestaurantInfo) {
		defer wg.Done()
		for n := range in {
			out <- n
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go collect(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func PaginationGet(maxpage int, count int) []<-chan mathcers.RestaurantInfo {
	logger.Info.Printf("maxpage:%d,count:%d", maxpage, count)
	var chanlist []<-chan mathcers.RestaurantInfo
	if count > 20 {
		for i := 2; i <= maxpage; i++ {
			func(index int) {
				o := getPageProducer(index, 20)
				chanlist = append(chanlist, o)
			}(i)
		}
	}
	return chanlist
}

func GetGaoResult() <-chan mathcers.RestaurantInfo {
	firstres, err := requestGao(1, 20)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	c, err := strconv.ParseInt(firstres.Count, 10, 16)
	count := int(c)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	maxpage := count / 20
	if (count - maxpage*20) > 0 {
		maxpage++
	}
	var outlist []<-chan mathcers.RestaurantInfo
	firstout := initProducer(firstres)
	outlist = append(outlist, firstout)
	remainoutlist := PaginationGet(maxpage, count)
	outlist = append(outlist, remainoutlist...)
	mergeout := mergeResult(count, outlist...)
	return mergeout
}
