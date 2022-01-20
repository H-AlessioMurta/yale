package logger

// A custom modification for log, just adding same flavour color and specific intercation for our graph/model
import (
    "net/http"
    "log"
    "yale/borrowing/graph/model" 
    "strconv"
	"encoding/json"
	"strings"
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
    message = "\033[31m[FATAL]\033[0m:"+message+"\nError:"+e.Error()
	resp, err := http.Post("https://notificationsvc:3000", "application/json",strings.NewReader(e.Error()))
	CheckErr(err)
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
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
        "\033[35mBook id: "+
        b.ID+"\n"+
        "Book titile: "+
        b.Title+"\n"+
        "Book authors: "+
        b.Authors+"\n"+
        "\033[0m"
    log.Println(message)
}


func LogResponseCustomer(c *model.Customer) {   
    message := "\033[32m[Response]\033[0m:\n"+
        "\033[35mCustomer id: "+
        c.ID+"\n"+
        "Customer name: "+
        c.Name+"\n"+
        "Customer surname: "+
        c.Surname+"\n"+
        "Customer National identifier number: "+
        c.Nin+"\n"+
        "\033[0m"
    log.Println(message)
}

func LogResponseBorrowing(c *model.Borrowed) {   
    message := "\033[32m[Response]\033[0m:\n"+
        "\033[35mID Borrow :"+
        c.IDBorrowing+"\n"+
        "Id book: "+
        c.IDBook+"\n"+
        "Id customer: "+
        c.IDCustomer+"\n"+
        "When the borrow start : "+
        c.Starting.String()+"\n"+
        "Deadline for returning the book: "+
        c.Expiring.String()+"\n"+
        "Is the book returned: "+
        strconv.FormatBool(c.Returned)+"\n"+
        "\033[0m"
    log.Println(message)
}

//Cause a lot of errors need to be checked in Go, i just do it here for avoid to write the same "if check and panifc".
func CheckErr(err error) {
    if err != nil {
        LogFatal(err.Error(),err)
    }
}