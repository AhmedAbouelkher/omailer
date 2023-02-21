package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	ml "github.com/AhmedAbouelkher/omailer"
)

func main() {
	mux := http.NewServeMux()

	tmp, err := template.ParseFiles(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	n := ml.NewDialer("", 465, "", "")

	go func() {
		statsC := n.StatsC()
		for s := range statsC {
			if err, ok := s.(*ml.EmailError); ok {
				log.Println("EMAIL ERROR", err)
				continue
			}
		}
	}()

	mux.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		b := buildBody()
		data := struct {
			Body   string
			Footer string
		}{
			Body:   b.String(),
			Footer: "This is a footer",
		}
		buf := bytes.NewBuffer(nil)
		if err := tmp.Execute(buf, data); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		msg := &ml.EmailMessage{
			From:    "no-replay@example.net",
			To:      "user@example.net",
			Subject: "This is a test email, please ignore",
			Body:    buf.String(),
		}
		n.SendAsync(context.Background(), msg)
		w.Write([]byte("Sent async" + msg.Subject + "\tat\t" + time.Now().String()))
	})

	http.ListenAndServe(":8080", mux)
}

func buildBody() *ml.HTML {
	b := &ml.HTML{}
	img := ml.Img(
		"https://ci6.googleusercontent.com/proxy/Vny5Kmjia75JB3hbdJFXm46ezXRqc15ds-333eBovp-tTOyOc0nLiIX5YVlynFgiC-sRwpbR4d1nxLTvPs9Z8bw-Vmck0sKaxKR6m8PxSAAqeOwi-GQEYZNRGGoL_Ug2JSlutg=s0-d-e1-ft#https://info.cloudflare.com/rs/713-XSC-918/images/developer-week-announcement.png",
		&ml.ImgElem{
			Alt:    "this is image alt",
			Link:   "https://content.cloudflare.com/dc/X2R7oDDwEhxR9sYIShA9hMz3VBpgnVrL43k7tUZcVeLRW0iFol8jH-Rf0W1LzwXW1wty04tcoZiursehOQ5kuFQK44mjAygbhgj4eWZBEIEfaLllvYEgYKJBu_WkabnjsxPTA2NbqlEgaKh6rjNt5-JSD10Iipo9IZRYuseXnjug5Ctero24Pp1RdjnVC6YhejhlRtM0kfxcgdsqonteUQ==/NzEzLVhTQy05MTgAAAGIpa1TMnMf8H27SbkC5gL70BEd3f7ImfeY7TM7FNtAUIfriPyGH9S4XIZ7aEkEvn478hTonh4=",
			Height: 400,
			Width:  225,
		},
	)
	b.AddElem(
		img,
		ml.Padding(ml.P("Hello,"), 0, 5),
		ml.P("We were unable to process your renewal payment (invoice # CFUSAXXXXXXX). This typically happens when your bank issues a new card, your existing card has expired, or because of a billing error caused by your bank. We'll automatically try again on January 23."),
		ml.Padding(ml.Text("Will my service be impacted?", &ml.TextStyle{FontWeight: "bold", FontSize: 13}), 10, 10),
		ml.P("Yes, Cloudflare has resumed collecting payment on all outstanding invoices starting on Jan 18, 2022. If we do not receive payment 5 days after your invoice is issued, your paid services will be downgraded."),
		ml.Space(8, 0),
		ml.P("Continued payment failures may result in a service disruption. To avoid a service disruption, take a moment to review or update your billing information. If you have not yet set up a backup payment method, consider adding one to avoid future payment failures."),
		ml.Space(8, 0),
		ml.Center(ml.Btn("Update your billing information", "https://dash.cloudflare.com/?to=/:account/billing")),
		ml.Space(10, 0),
		ml.Text("Helpful resources", &ml.TextStyle{FontWeight: "bold", FontSize: 13}),
		ml.List(
			ml.A("Learn more about why a payment failed", "https://support.cloudflare.com/hc/en-us/articles/218344877"),
			ml.A("Update an existing payment method or add backup payment method", "https://support.cloudflare.com/hc/en-us/articles/200170236"),
			ml.A("Pay an outstanding balance", "https://dash.cloudflare.com/?to=/:account/billing"),
			ml.A("Be notified when usage exceeds your defined threshold value", "https://support.cloudflare.com/hc/en-us/articles/115004555148"),
		),
		ml.P("If you have additional questions, contact us at our", ml.A("Support portal", "https://dash.cloudflare.com/?to=/:account/support"), "."),
	)

	return b

}
