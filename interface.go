package envutils

type Defualter interface {
	SetDefaults()
}

type Initialler interface {
	Init()
}
