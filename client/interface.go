package client

type Getter interface {
	Get(k string) (string, error)
}

type Setter interface {
	Set(k, v string) error
}

type Deller interface {
	Del(k string) error
}

type Client interface {
	Getter
	Setter
	Deller
}
