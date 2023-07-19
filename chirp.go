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

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)

    if err != nil {
        respondeWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
        return
    }


    if len(params.Body) > 140 {
        respondeWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }

    msg := removeBadWords(params.Body)

    respondeSuccess(w, http.StatusOK, msg)
    return
}

func respondeWithError(w http.ResponseWriter, code int, msg string) {
    type returnErrorVals struct {
        Error string `json:"error"`
    }
    
    respBody := returnErrorVals{Error: msg}
    data, _ := json.Marshal(respBody)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(data)
    return
}

func removeBadWords(body string) string {
    badWords := map[string]string{"kerfuffle": "kerfuffle", "sharbert": "sharbert", "fornax": "fornax"}
    words := strings.Split(body, " ")
    finalString := make([]string, 0)

    for _, word := range words {
        if _, ok := badWords[word]; ok {
            finalString = append(finalString, "****")
            continue
        }

        finalString = append(finalString, word)
    }

    return strings.Join(finalString, " ")
}

func respondeSuccess(w http.ResponseWriter, code int, msg string) {
    type returnSuccessVals struct {
        CleanedBody string `json:"cleaned_body"`
    }

    respBody := returnSuccessVals{CleanedBody: msg}
    data, _ := json.Marshal(respBody)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
    return
}
