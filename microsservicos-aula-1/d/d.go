package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Result struct {
	Status string
}

//type Coupon struct {
//	Code string
//}

//var string coupon

func main() {

	http.HandleFunc("/", process)
	http.ListenAndServe(":9093", nil)

}

//Check if the coupon is already used.
func process(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")

	result := Result{Status: ""}

	if coupon == "123456" {
		result.Status = "used"
	} else {
		result.Status = "notUsed"
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}
