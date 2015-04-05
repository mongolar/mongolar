package dynamicload

type DynamicLoad{
	Target string
	Controller string
	Id	string
	Template string
}

func New(t string, c string, i string, te string) *DynamicLoad{
	dl = new(DynamicLoad)
	dl.Target = t
	dl.Controller = c
	dl.Id = i
	dl.Template = te
}
