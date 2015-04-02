package response

import()

type Response struct{
	Content map[string]interface{}
	Messages map[string]
	Redirects string
}

func New() *Response{
	r := new(Response)
	r.Content = make(map[string]interface{})
	r.Messages = make(map[string])
}

func (*Response) SetContent(c map[string]interface{}){
	r.Content = c
}

func (*Response) getMessages(){
//TODO get flash messages
}

func (*Response) getRedirects(){
//TODO get flash messages
}
