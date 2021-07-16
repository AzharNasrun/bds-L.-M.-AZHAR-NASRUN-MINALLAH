package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	ID       int    `json:"id"`
	NamaRsu  string `json:"nama_rsu"`
	JenisRsu string `json:"jenis_rsu"`
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Alamat    string   `json:"alamat"`
	KodePos   int      `json:"kode_pos"`
	Telepon   []string `json:"telepon"`
	Faximile  []string `json:"faximile"`
	Website   string   `json:"website"`
	Email     string   `json:"email"`
	Kelurahan struct {
		Kode int64  `json:"kode"`
		Nama string `json:"nama"`
	} `json:"kelurahan"`
	Kecamatan struct {
		Kode int    `json:"kode"`
		Nama string `json:"nama"`
	} `json:"kecamatan"`
	Kota struct {
		Kode int    `json:"kode"`
		Nama string `json:"nama"`
	} `json:"kota"`
}
type mainResponse struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
	Data   []Data `json:"data"`
}
type Kelurahan struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
	Data   []struct {
		KodeProvinsi  int    `json:"kode_provinsi"`
		NamaProvinsi  string `json:"nama_provinsi"`
		KodeKota      int    `json:"kode_kota"`
		NamaKota      string `json:"nama_kota"`
		KodeKecamatan int    `json:"kode_kecamatan"`
		NamaKecamatan string `json:"nama_kecamatan"`
		KodeKelurahan int64  `json:"kode_kelurahan"`
		NamaKelurahan string `json:"nama_kelurahan"`
	} `json:"data"`
}

type RumahSakit struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
	Data   []struct {
		ID       int    `json:"id"`
		NamaRsu  string `json:"nama_rsu"`
		JenisRsu string `json:"jenis_rsu"`
		Location struct {
			Alamat    string  `json:"alamat"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
		KodePos       int      `json:"kode_pos"`
		Telepon       []string `json:"telepon"`
		Faximile      []string `json:"faximile"`
		Website       string   `json:"website"`
		Email         string   `json:"email"`
		KodeKota      int      `json:"kode_kota"`
		KodeKecamatan int      `json:"kode_kecamatan"`
		KodeKelurahan int64    `json:"kode_kelurahan"`
		Latitude      float64  `json:"latitude"`
		Longitude     float64  `json:"longitude"`
	} `json:"data"`
}

func main() {

	http.HandleFunc("/data", getData)
	http.ListenAndServe(":8080", nil)
}

func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if r.Method == "GET" {
		kelurahan := getKelurahan()
		rmhSkt := getRumahSakit()
		mainresp := aggregation(kelurahan, rmhSkt)
		res, err := json.Marshal(mainresp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(res)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func getKelurahan() Kelurahan {
	req, err := http.NewRequest("GET", "http://api.jakarta.go.id/v1/kelurahan", nil)

	req.Header.Add("Authorization", "LdT23Q9rv8g9bVf8v/fQYsyIcuD14svaYL6Bi8f9uGhLBVlHA3ybTFjjqe+cQO8k")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	kel := Kelurahan{}

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &kel)

	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	return kel

}

func getRumahSakit() RumahSakit {
	req, err := http.NewRequest("GET", "http://api.jakarta.go.id/v1/rumahsakitumum", nil)

	req.Header.Add("Authorization", "LdT23Q9rv8g9bVf8v/fQYsyIcuD14svaYL6Bi8f9uGhLBVlHA3ybTFjjqe+cQO8k")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	rmh := RumahSakit{}

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &rmh)

	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	return rmh

}

func aggregation(kel Kelurahan, rmh RumahSakit) mainResponse {
	mainrsp := mainResponse{}
	mainrsp.Count = rmh.Count
	var data = make([]Data, rmh.Count)

	for i := 0; i < rmh.Count; i++ {
		kode := rmh.Data[i].KodeKelurahan
		data[i].Alamat = rmh.Data[i].Location.Alamat
		data[i].Email = rmh.Data[i].Email
		data[i].Faximile = rmh.Data[i].Faximile
		data[i].ID = rmh.Data[i].ID
		data[i].NamaRsu = rmh.Data[i].NamaRsu
		data[i].JenisRsu = rmh.Data[i].JenisRsu
		data[i].Location.Latitude = rmh.Data[i].Location.Latitude
		data[i].Location.Longitude = rmh.Data[i].Location.Longitude
		data[i].Telepon = rmh.Data[i].Telepon
		data[i].Website = rmh.Data[i].Website
		data[i].KodePos = rmh.Data[i].KodePos
		for j := 0; j < kel.Count; j++ {
			if kel.Data[j].KodeKelurahan == kode {
				data[i].Kecamatan.Nama = kel.Data[j].NamaKecamatan
				data[i].Kecamatan.Kode = kel.Data[j].KodeKecamatan
				data[i].Kelurahan.Kode = kel.Data[j].KodeKelurahan
				data[i].Kelurahan.Nama = kel.Data[j].NamaKelurahan
				data[i].Kota.Nama = kel.Data[j].NamaKota
				data[i].Kota.Kode = kel.Data[j].KodeKota
				break
			}
		}

	}
	mainrsp.Data = data
	return mainrsp
}
