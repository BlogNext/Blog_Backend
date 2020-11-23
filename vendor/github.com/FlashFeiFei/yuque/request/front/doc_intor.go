package front

import (
	"encoding/json"
	"fmt"
	"github.com/FlashFeiFei/yuque/response/front"
	"io/ioutil"
	"log"
	"net/http"
)

func GetDocIntorSerializer(slug string, book_id int64) *front.ResponseDocIntorSerializer {
	client := http.Client{}
	creq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.yuque.com/api/docs/%s?book_id=%d", slug, book_id), nil)
	creq.Header.Add("content-type", "application/json")
	resp, _ := client.Do(creq)
	body, _ := ioutil.ReadAll(resp.Body)
	response := new(front.ResponseDocIntorSerializer)
	err := json.Unmarshal(body, response)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	log.Println(response.Data.ID)
	log.Println(response.Data.Cover)
	log.Println(response.Data.CustomDescription)
	return response
}
