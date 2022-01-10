package customersvc

import (
    "net/http"
    "encoding/json"
    "log"    
)

func logInfo(message string) {
    message = "\033[34m[INFO]\033[0m:"+message
    log.Println(message)
}

func logWarning(message string) {
    message = "\033[35m[WARNING]\033[0m:"+message
    log.Println(message)
}

func logError(message string) {
    message = "\033[33m[ERROR]\033[0m:"+message
    log.Println(message)
}

func logFatal(message string) {
    message = "\033[31m[FATAL]\033[0m:"+message
    log.Fatal(message)
}

func logRequest(r *http.Request) {
    
    message := "\033[32m[Request]\033[0m:"+
        "\033[33m"+    
        r.Method+
        "\033[0m URI:\033[35m "+
        r.RequestURI+
        "\033[0m IP:\033[35m "+
        r.RemoteAddr+
        "\033[0m"
    log.Println(message)
}
func logResponse(message string) {
    
    message = "\033[32m[Response]\033[0m:\n"+
        "\033[35m"+
        message+
        "\033[0m"
    log.Println(message)
}

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
        panic(err)
    }
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    logError(message)
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.MarshalIndent(payload,"","  ")
    logResponse(string (response))
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}