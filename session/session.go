package session

import(
	"net/http"
)

type Session map[string]interface{}

func New() *Session{
	return new(Session)
}

func (s *Session) Set(w http.ResponseWriter, r *http.Request, si *site.SiteConfig){
	c, err := r.Cookie("m_session_id")
	if c == nil {
		expire := time.Now().AddDate(0, 0, 1)
		c := {
			"m_session_id",
			v,
			"/",
			r.Host,
			expire,
			expire,
			0,
			true,
			true,
			"m_session_id=" + v,
			[]string{"m_session_id=" + v}
		}
		http.SetCookie(w, c)
	}
}
