package handler

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"io/ioutil"
)

type Hito struct {
	URI  string
	Title string
	fecha time.Time
}

type Response struct {
	Msg string `json:"text"`
	ChatID int64 `json:"chat_id"`
	Method string `json:"method"`
}

func format_date(diff time.Duration) string {
	// Función para convertir una diferencia de fechas
	// a una cadena del tipo dias horas minutos segundos

    // Definimos constantes para partir la diferencia
    const Decisecond = 100 * time.Millisecond
    const Day = 24*time.Hour

    // Extraemos la cantidad de dias y las quitamos de la diferencia
    d := diff / Day
    diff = diff % Day
    // Extraemos la cantidad de horas
    h := diff / time.Hour
    diff = diff % time.Hour
    // Extraemos la cantidad de minutos
    m := diff / time.Minute
    diff = diff % time.Minute
    // Extraemos la cantidad de secundos
    s := diff / time.Second
    diff = diff % time.Second
    // Nos quedamos con las partes de segundo
    f := diff / Decisecond
    return  fmt.Sprintf("%dd %dh %dm %d.%ds", d, h, m, s, f)

}


var hitos = []Hito {
	Hito {
		URI: "0.Repositorio",
		Title: "Datos básicos y repo",
		fecha: time.Date(2020, time.September, 29, 11, 30, 0, 0, time.UTC),
	},
	Hito {
		URI: "1.Infraestructura",
		Title: "HUs y entidad principal",
		fecha: time.Date(2020, time.October, 6, 11, 30, 0, 0, time.UTC),
	},
	Hito {
		URI: "2.Tests",
		Title: "Tests iniciales",
		fecha: time.Date(2020, time.October, 16, 11, 30, 0, 0, time.UTC),
	},
	Hito {
		URI: "3.Contenedores",
		Title: "Contenedores",
		fecha: time.Date(2020, time.October, 26, 11, 30, 0, 0, time.UTC),
	},
	Hito {
		URI: "4.CI",
		Title: "Integración continua",
		fecha: time.Date(2020, time.November, 6, 23, 59, 0, 0, time.UTC),
	},
	Hito {
		URI: "5.Serverless",
		Title: "Trabajando con funciones serverless",
		fecha: time.Date(2020, time.November, 24, 23, 59, 0, 0, time.UTC),
	},
	Hito {
		URI: "6.Microservicio",
		Title: "Diseñando un microservicio",
		fecha: time.Date(2020, time.December, 11, 23, 59, 0, 0, time.UTC),
	},
	Hito {
		URI: "7.PaaS",
		Title: "Desplegando en un PaaS",
		fecha: time.Date(2020, time.December, 23, 23, 59, 0, 0, time.UTC),
	},
	Hito {
		URI: "8.Despliegue en la nube",
		Title: "Desplegando en la nube",
		fecha: time.Date(2021, time.January, 14, 23, 59, 0, 0, time.UTC),
	},

}


func Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var update tgbotapi.Update
	if err := json.Unmarshal(body,&update); err != nil {
		log.Fatal("Error en el update →", err)
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	currentTime := time.Now()
	var next int
	var queda time.Duration
	for indice, hito := range hitos {
		if ( hito.fecha.After( currentTime ) ) {
			next = indice
			queda = hito.fecha.Sub( currentTime )
		}
	}
	if update.Message.IsCommand() {
		text := ""
		if ( next == 0 ) {
			text = "Ninguna entrega próxima"
		} else {

			switch update.Message.Command() {
			case "kk":
				text = format_date(queda)
			case "kekeda":
				text = fmt.Sprintf( "→ Próximo hito %s\n🔗 https://jj.github.io/IV/documentos/proyecto/%s\n📅 %s",
					hitos[next].Title,
					hitos[next].URI,
					hitos[next].fecha.String(),
				)
			default:
				text = "Usa /kk para lo que queda para el próximo hito, /kekeda para + detalle"
			}
		}
		data := Response{ Msg: text,
			Method: "sendMessage",
			ChatID: update.Message.Chat.ID }

		msg, _ := json.Marshal( data )
		log.Printf("Response %s", string(msg))
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w,string(msg))
	}
}
