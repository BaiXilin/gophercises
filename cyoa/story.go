package cyoa // Choose your own adventure

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	err := d.Decode(&story)
	if err != nil {
		return nil, err
	}
	return story, nil
}

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
    <head>
        <title>Choose Your Own Adventure</title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
            <p>{{.}}</p>
        {{end}}
        <ul>
            {{range .Options}}
                <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
            {{end}}
        </ul>
    </body>
</html>
`

// returned http.Handler will be feed to http.ListenAndServe
// allows ServeHTTP of handler to do its job
func NewHandler(s Story) http.Handler {
	return handler{s}
}

// struct handler implements http.Handler interface
type handler struct {
	s Story
}

// the only method under http.Handler interface is ServeHTTP
// specify whatever needs to be done by the web app in this func
// in this case, ServeHTTP needs to parse the html template, and execute it
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("").Parse(defaultHandlerTmpl))
	err := tpl.Execute(w, h.s["intro"])
	if err != nil {
		panic(err)
	}
}
