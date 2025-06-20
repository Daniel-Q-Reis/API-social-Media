insert into usuarios (nome, nick, email, senha)
values
("Usuário 1", "usuario_1", "usuario1@gmail.com", "$2a$10$3H59Vv55y1LRA34tBJbRzOeFbAd8U14NaqplStlGMEAOWIWqn5pbK"), -- usuario1
("Usuário 2", "usuario_2", "usuario2@gmail.com", "$2a$10$3H59Vv55y1LRA34tBJbRzOeFbAd8U14NaqplStlGMEAOWIWqn5pbK"), -- usuario2
("Usuário 3", "usuario_3", "usuario3@gmail.com", "$2a$10$3H59Vv55y1LRA34tBJbRzOeFbAd8U14NaqplStlGMEAOWIWqn5pbK"); -- usuario3

insert into seguidores (usuario_id, seguidor_id)
values
(1, 2),
(3, 1),
(1, 3);