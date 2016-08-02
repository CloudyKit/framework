package scheme

func New(entityName, primaryKey string) *Scheme {
	scheme := &Scheme{name: entityName, primaryKey: primaryKey, fields: make(map[string]*Field)}
	return scheme
}

func Init(scheme *Scheme, def func(*Def)) *Scheme {
	def((*Def)(scheme))
	return scheme
}

func NewInit(entityName, primaryKey string, def func(*Def)) *Scheme {
	return Init(New(entityName, primaryKey), def)
}
