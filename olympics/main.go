//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
)

var athletes []A

func Athletes(w http.ResponseWriter, r *http.Request) {
	var name string
	queries := r.URL.Query()

	if nameSlc, ok := queries["name"]; !ok || len(nameSlc[0]) <= 0 {
		http.Error(w, "no name param", http.StatusBadRequest)
		return
	} else {
		name = nameSlc[0]
	}

	filtered := Filter(athletes, func(athlete A) bool {
		return athlete.Athlete == name
	})

	if len(filtered) == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := athletToInfo(&filtered[0])

	for _, athlete := range filtered {
		_, ok := resp.MedalsByYear[athlete.Year]
		if !ok {
			resp.MedalsByYear[athlete.Year] = &Medals{0, 0, 0, 0}
		}
		resp.MedalsByYear[athlete.Year].Gold += athlete.Gold
		resp.MedalsByYear[athlete.Year].Silver += athlete.Silver
		resp.MedalsByYear[athlete.Year].Bronze += athlete.Bronze
		resp.MedalsByYear[athlete.Year].Total += athlete.Gold + athlete.Silver + athlete.Bronze

		resp.Medals.Gold += athlete.Gold
		resp.Medals.Silver += athlete.Silver
		resp.Medals.Bronze += athlete.Bronze
		resp.Medals.Total += athlete.Gold + athlete.Silver + athlete.Bronze
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func TopAthletes(w http.ResponseWriter, r *http.Request) {
	var sportParam string
	var limitParam int
	var err error
	queries := r.URL.Query()

	if sportSlc, ok := queries["sport"]; ok {
		if sportSlc[0] != "" {
			sportParam = sportSlc[0]
		} else {
			http.Error(w, "no sport param", http.StatusBadRequest)
		}
	}

	if limitSlc, ok := queries["limit"]; ok {
		if limitSlc[0] != "" {
			limitParam, err = strconv.Atoi(limitSlc[0])
			if err != nil {
				http.Error(w, "invalid limit param", http.StatusBadRequest)
				return
			}
		}
	} else {
		limitParam = 3
	}

	test := func(athlete A) bool {
		return athlete.Sport == sportParam
	}
	filtered := Filter(athletes, test)

	if len(filtered) == 0 {
		http.Error(w, "sport not found", http.StatusNotFound)
		return
	}

	filteredMap := GetAthlets(filtered)

	values := make([]*AInfo, 0, len(filteredMap))

	for _, v := range filteredMap {
		values = append(values, v)
	}

	sort.Slice(values, func(i, j int) bool {
		if values[i].Medals.Gold != values[j].Medals.Gold {
			return values[i].Medals.Gold > values[j].Medals.Gold
		}
		if values[i].Medals.Silver != values[j].Medals.Silver {
			return values[i].Medals.Silver > values[j].Medals.Silver
		}
		if values[i].Medals.Bronze != values[j].Medals.Bronze {
			return values[i].Medals.Bronze > values[j].Medals.Bronze
		}

		return values[i].Athlete < values[j].Athlete
	})

	limit := int(math.Min(float64(limitParam), float64(len(values))))
	result := values[:limit]
	b, err := json.Marshal(&result)

	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func TopCountries(w http.ResponseWriter, r *http.Request) {
	var yearParam int
	var limitParam int
	var err error
	queries := r.URL.Query()

	if yearSlc, ok := queries["year"]; ok {
		if yearSlc[0] != "" {
			yearParam, err = strconv.Atoi(yearSlc[0])
			if err != nil {
				http.Error(w, "invalid year param", http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "no sport param", http.StatusBadRequest)
		}
	}

	if limitSlc, ok := queries["limit"]; ok {
		if limitSlc[0] != "" {
			limitParam, err = strconv.Atoi(limitSlc[0])
			if err != nil {
				http.Error(w, "invalid limit param", http.StatusBadRequest)
				return
			}
		}
	} else {
		limitParam = 3
	}

	filtered := Filter(athletes, func(athlete A) bool {
		return athlete.Year == yearParam
	})

	if len(filtered) == 0 {
		http.Error(w, "year not found", http.StatusNotFound)
		return
	}

	fMap := GetCountries(filtered)

	vals := make([]*CInfo, 0, len(fMap))

	for _, v := range fMap {
		vals = append(vals, v)
	}

	sort.Slice(vals, func(i, j int) bool {
		if vals[i].Gold != vals[j].Gold {
			return vals[i].Gold > vals[j].Gold
		}
		if vals[i].Silver != vals[j].Silver {
			return vals[i].Silver > vals[j].Silver
		}
		if vals[i].Bronze != vals[j].Bronze {
			return vals[i].Bronze > vals[j].Bronze
		}

		return vals[i].Country < vals[j].Country
	})

	l := int(math.Min(float64(limitParam), float64(len(vals))))
	result := vals[:l]
	b, err := json.Marshal(&result)

	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func main() {
	port := flag.String("port", "80", "http server port")
	jsonP := flag.String("data", "./olympics/testdata/olympicWinners.json", "json path")
	flag.Parse()

	jF, err := os.Open(*jsonP)
	if err != nil {
		log.Fatal("file read err")
	}

	byteVal, _ := io.ReadAll(jF)
	jF.Close()
	err = json.Unmarshal(byteVal, &athletes)

	if err != nil {
		log.Fatal("json parse err")
	}

	http.HandleFunc("/athlete-info", Athletes)
	http.HandleFunc("/top-athletes-in-sport", TopAthletes)
	http.HandleFunc("/top-countries-in-year", TopCountries)

	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}
