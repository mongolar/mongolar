package form

import(
        "github.com/jasonrichardsmith/mongolar/wrapper"	
)

type Post map[string]string

type Form map[string]interface
type FormRecord struct{

}

type FormHandlerMap map[string]FormHandler

type FormHandler interface{
	Submit(*wrapper.Wrapper, Post)
}

func NewFormHandlerMap() FormHandlerMap { 
	return make(FormHandlerMap)
}

func BuildForm(h string, f map[string]interface{}) Form{
	
}



