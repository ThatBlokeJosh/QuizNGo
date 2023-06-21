package scrape

import (
	"github.com/foolin/pagser"
	"log"
	"net/http"
)

func GetP(url string) []string {
  // "https://www.dejepis.com/ucebnice/prvni-svetova-valka-treti-a-ctvrta-etapa-valky-1917-1918/"
  type PageData struct {
	  Paragraphs []string  `pagser:"p->eachText()"`
  }
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	p := pagser.New()
	var data PageData

	err = p.ParseReader(&data, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
  
  return data.Paragraphs
}

// Test
/*func main() {
  log.Println(getP("https://www.dejepis.com/ucebnice/prvni-svetova-valka-treti-a-ctvrta-etapa-valky-1917-1918/"))
}*/
