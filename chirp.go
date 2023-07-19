package main

import (
    "encoding/json"
    "strings"
    "net/http"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body string `json:"body"`
    }

    type returnErrorVals struct {
        Error string `json:"error"`
    }

    type returnSuccessVals struct {
        CleanedBody string `json:"cleaned_body"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)


    if err != nil {
        respBody := returnErrorVals{Error: "Couldn't decode parameters"}
        code := http.StatusInternalServerError
        data, _ := json.Marshal(respBody)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(code)
        w.Write(data)
        return
    }


    if len(params.Body) > 140 {
        respBody := returnErrorVals{Error: "Chirp is too long"}
        data, _ := json.Marshal(respBody)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        w.Write(data)
        return
    }

    badWords := map[string]string{"kerfuffle": "kerfuffle", "sharbert": "sharbert", "fornax": "fornax"}
    words := strings.Split(params.Body, " ")
    finalString := make([]string, 0)

    for _, word := range words {
        if _, ok := badWords[word]; ok {
            finalString = append(finalString, "****")
            continue
        }

        finalString = append(finalString, word)
    }

    msg := strings.Join(finalString, " ")


    respBody := returnSuccessVals{CleanedBody: msg}
    data, _ := json.Marshal(respBody)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
    return
}

