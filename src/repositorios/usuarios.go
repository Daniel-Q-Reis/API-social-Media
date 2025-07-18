package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

// Usuarios representa um repositório de usuarios
type Usuarios struct {
	db *sql.DB //ponteiro para uma conexão com o banco de dados
}

// NovoRepositorioDeUsuarios cria um repositório de usuarios: Essa função é um "construtor": ela cria uma nova instância de Usuarios com a conexão db recebida por parâmetro.
func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

// Metodo Criar insere um usuário no banco de dados
func (repositorio Usuarios) Criar(usuario modelos.Usuario) (uint64, error) {
	//(repositorio Usuarios) metodos vai estar dentro do repositorio de usuarios...Criar (vai receber um parametro e criar um usuario la dos modelos), e vai retornar um UINT64 e um ERROR
	//porque um Uint64? uint vai ser o id que ele vai inserir no banco
	statement, erro := repositorio.db.Prepare(
		"insert into usuarios (nome, nick, email, senha) values(?, ?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha) // usuario esta vindo do parametro criar
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId() //esse cara aqui é um int64 entao no final teremos que converter
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// Buscar traz todos os usuários que atendem um filtro de nome ou nick
func (repositorio Usuarios) Buscar(nomeOuNick string) ([]modelos.Usuario, error) {
	//Aqui estamos formatando a string para preparar uma busca com o operador SQL
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) //-> Equivale a %nomeOuNick%
	//Se o usuário buscar "jo", isso vira "%jo%", que casa com "João", "Major", "Joaquina", etc.

	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where nome LIKE ? or nick LIKE ?",
		nomeOuNick, nomeOuNick)
	//fazendo uma consulta SQL ao banco de dados usando LIKE, tanto no campo nome quanto no nick.

	if erro != nil {
		return nil, erro // Esse nil se refere ao slice []modelos.Usuario, e pode retornar um valor 0 (Nenhum usuario)
	}

	defer linhas.Close()
	//Isso garante que a conexão com os resultados da consulta será fechada quando a função terminar (seja com sucesso ou erro)
	//Evita vazamento de recursos.

	var usuarios []modelos.Usuario
	//Aqui é criado um slice vazio de usuários, que vai armazenar todos os usuários encontrados no banco.

	for linhas.Next() { //Essa linha começa um loop que percorre cada linha do resultado da consulta
		var usuario modelos.Usuario //linhas.Next() avança para a próxima linha — e retorna false quando não houver mais.
		//Dentro do loop, criamos uma estrutura de usuário que será preenchida com os dados do banco.
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome, //Aqui é feita a leitura da linha atual da consulta usando Scan, que mapeia os campos do banco de dados para os campos da struct.
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario) //Aqui adicionamos o usuario preenchido no slice usuarios
	}
	return usuarios, nil
	//Por fim, se tudo deu certo, retornamos:
	//O slice com todos os usuários encontrados.
	//nil para indicar que não houve erro.
}

// BuscarPorID traz um usuário do banco de dados
func (repositorio Usuarios) BuscarPorID(ID uint64) (modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEM from usuarios where id = ?",
		ID,
	)
	if erro != nil {
		return modelos.Usuario{}, erro //aqui ele entende que é um usuario tendo todos seus valores como zero, equivale a um nil
	}
	defer linhas.Close()

	var usuario modelos.Usuario

	if linhas.Next() { //observe que antes era um for, agora é if pois se trata apenas de 1 usuario
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome, //Aqui é feita a leitura da linha atual da consulta usando Scan, que mapeia os campos do banco de dados para os campos da struct.
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return modelos.Usuario{}, erro
		}

	}
	return usuario, nil

}

// Atualizar altera as informações de um usuário no banco de dados
func (repositorio Usuarios) Atualizar(ID uint64, usuario modelos.Usuario) error {
	statement, erro := repositorio.db.Prepare("update usuarios set nome = ?, nick = ?, email = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, ID); erro != nil {
		return erro
	}

	return nil
}

// Deletar exclui as informações de um usuário no banco de dados
func (repositorio Usuarios) Deletar(ID uint64) error {
	statement, erro := repositorio.db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

// BuscarPorEmail busca um usuário por email e retorna o seu id e senha com hash
func (repositorio Usuarios) BuscarPorEmail(email string) (modelos.Usuario, error) {
	linha, erro := repositorio.db.Query("select id, senha from usuarios where email = ?", email)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Senha); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

// Seguir
func (respositorio Usuarios) Seguir(usuarioID, seguidorID uint64) error {
	statement, erro := respositorio.db.Prepare(
		"insert ignore into seguidores (usuario_id, seguidor_id) values (?, ?)",
	) //ignore evita que eu coloque adicione um usuario que já está sendo seguido, se ja existir uma linha na tabela com aqueles valores, ele vai ignorar
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuarioID, seguidorID); erro != nil {
		return erro
	}

	return nil

}

// PararDeSeguir permite que um usuário pare de seguir o outro
func (repositorio Usuarios) PararDeSeguir(usuarioID, seguidorID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"delete from seguidores where usuario_id = ? and seguidor_id = ?",
	) // A linha que tiver esses dois IDs conforme estiver vindo no parametro, iremos deleta-la
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuarioID, seguidorID); erro != nil { //executar o nosso statment
		return erro
	}

	return nil
}

// BuscarSeguidores traz todos os seguidores de um usuário
func (repositorio Usuarios) BuscarSeguidores(usuarioID uint64) ([]modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(`
		select u.id, u.nome, u.nick, u.email, u.criadoEm
		from usuarios u inner join seguidores s on u.id = s.seguidor_id where s.usuario_id = ?
	`, usuarioID) // observe que declara o u depois de usar ele, é uma peculiariadade do sql
	if erro != nil {
		return nil, erro
	}

	var usuarios []modelos.Usuario
	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// BuscarSeguindo traz todos os usuários que um determinado usuário está seguindo
func (repositorio Usuarios) BuscarSeguindo(usuarioID uint64) ([]modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(`
	 	select u.id, u.nome, u.nick, u.email, u.criadoEm
		from usuarios u inner join seguidores s on u.id = s.usuario_id where s.seguidor_id = ?`,
		usuarioID,
	) // observe que declara o u depois de usar ele, é uma peculiariadade do sql
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario

	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil

}

// BuscarSenha traz a senha de um usuário pelo ID
func (repositorio Usuarios) BuscarSenha(usuarioID uint64) (string, error) {
	linha, erro := repositorio.db.Query("select senha from usuarios where id = ?", usuarioID)
	if erro != nil {
		return "", erro
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.Senha); erro != nil {
			return "", erro
		}
	}

	return usuario.Senha, nil
}

// AtualizarSenha traz a senha de um usuário pelo ID
func (repositorio Usuarios) AtualizarSenha(usuarioID uint64, senha string) error {
	statement, erro := repositorio.db.Prepare("update usuarios set senha = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(senha, usuarioID); erro != nil {
		return erro
	}

	return nil
}
