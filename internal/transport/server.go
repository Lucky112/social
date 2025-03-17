package transport

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/Lucky112/social/config"
	"github.com/Lucky112/social/internal/transport/auth"
	"github.com/Lucky112/social/internal/transport/friends"
	"github.com/Lucky112/social/internal/transport/jwt"
	"github.com/Lucky112/social/internal/transport/posts"
	"github.com/Lucky112/social/internal/transport/profiles"
)

type Server struct {
	server *fiber.App
	port   uint16
}

func NewServer(
	cfg *config.ServerConfig,
	authService auth.AuthService,
	profilesService profiles.ProfilesService,
	friendsService friends.FriendsService,
	postsService posts.PostsService,
) Server {
	jwtKey := []byte(cfg.JWTKey)

	authHandler := auth.NewAuthHandler(authService, jwtKey)
	profilesHandler := profiles.NewProfilesHandler(profilesService)
	friendsHandler := friends.NewFriendsHandler(friendsService, profilesService)
	postsHandler := posts.NewPostsHandler(postsService)

	server := fiber.New()

	server.Use(recover.New())

	publicGroup := server.Group("")
	publicGroup.Post("/register", authHandler.Register)
	publicGroup.Post("/login", authHandler.Login)

	authorizedGroup := server.Group("")
	authorizedGroup.Use(jwt.Middleware(jwtKey))

	authorizedGroup.Post("/profiles", profilesHandler.CreateProfile)
	authorizedGroup.Get("/profiles", profilesHandler.GetProfiles)
	authorizedGroup.Get("/profiles/search", profilesHandler.SearchProfile)
	authorizedGroup.Get("/profiles/:id", profilesHandler.GetProfileById)

	authorizedGroup.Post("/friends/:friend_id", friendsHandler.AddFriend)
	authorizedGroup.Delete("/friends/:friend_id", friendsHandler.DeleteFriend)
	authorizedGroup.Get("/friends", friendsHandler.GetFriends)

	authorizedGroup.Post("/posts", postsHandler.CreatePost)
	authorizedGroup.Get("/posts/:user_id", postsHandler.GetPosts)
	authorizedGroup.Get("/posts/:user_id/:post_id", postsHandler.GetPostById)

	return Server{
		server: server,
		port:   cfg.Port,
	}
}

func (s Server) Start() error {
	address := fmt.Sprintf(":%d", s.port)

	err := s.server.Listen(address)
	return err
}
