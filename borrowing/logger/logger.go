package logger

import (
    "net/http"
    "log"
    "yale/borrowing/graph/model" 
    "strconv"
)

func LogInfo(message string) {
    message = "\033[34m[INFO]\033[0m:"+message
    log.Println(message)
}

func LogWarning(message string) {
    message = "\033[35m[WARNING]\033[0m:"+message
    log.Println(message)
}

func LogError(message string) {
    message = "\033[33m[ERROR]\033[0m:"+message
    log.Println(message)
}

func LogFatal(message string, e error) {
    message = "\033[31m[FATAL]\033[0m:"+message
    log.Fatalf(message)
}

func LogGraph(message string) {
    message = "\033[32m[Graph Request]\033[0m:"+message
    log.Println(message)
}


func LogRequest(r *http.Request) {
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

func LogResponse(message string) {   
    message = "\033[32m[Response]\033[0m:\n"+
        "\033[35m"+
        message+
        "\033[0m"
    log.Println(message)
}

func LogResponseBook(b *model.Book) {   
    message := "\033[32m[Response]\033[0m:\n"+
        "\033[35m"+
        b.ID+"\n"+
        b.Title+"\n"+
        b.Authors+"\n"+
        "\033[0m"
    log.Println(message)
}


func LogResponseCustomer(c *model.Customer) {   
    message := "\033[32m[Response]\033[0m:\n"+
        "\033[35m"+
        c.ID+"\n"+
        c.Name+"\n"+
        c.Surname+"\n"+
        c.Nin+"\n"+
        "\033[0m"
    log.Println(message)
}

func LogResponseBorrowing(c *model.Borrowed) {   
    message := "\033[32m[Response]\033[0m:\n"+
        "\033[35m"+
        c.IDBorrowing+"\n"+
        c.IDBook+"\n"+
        c.IDCustomer+"\n"+
        c.Starting.String()+"\n"+
        c.Expiring.String()+"\n"+
        strconv.FormatBool(c.Returned)+"\n"+
        "\033[0m"
    log.Println(message)
}


func CheckErr(err error) {
    if err != nil {
        LogFatal(err.Error(),err)
    }
}
