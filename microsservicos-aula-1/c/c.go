package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}

func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

type Result struct {
	Status string
}

var coupons Coupons

func main() {
	coupon := Coupon{
		Code: "abc",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	valid := coupons.Check(coupon)

	result := Result{Status: valid}

	resultCouponUsed := makeHttpCall("http://localhost:9093", coupon)

	if resultCouponUsed.Status == "used" {
		result = Result{Status: "used"}
	} else {
		result = Result{Status: "valid"}
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}

func makeHttpCall(urlMicroservice string, coupon string) Result {

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	values := url.Values{}
	values.Add("coupon", coupon)

	res, err := retryClient.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Status: "Servidor fora do ar!"}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result

}
