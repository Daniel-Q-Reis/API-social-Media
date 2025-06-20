package modelos

import "time"

// Publicacao representa uma publicação feita por um usuário
type Publicacao struct {
	ID        uint64    `json:"id, omitempty"`
	Titutlo   string    `json:"titulo, omitempty"`
	Conteudo  string    `json:"conteudo, omitempty"`
	AutorID   uint64    `json:"autorId, omitempty"`
	AutorNick uint64    `json:"autorNick, omitempty"`
	Curtidas  uint64    `json:"curtidas"`
	CriadaEm  time.Time `json:"criadaEm, omitempty`
}
