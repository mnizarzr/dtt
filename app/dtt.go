package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mnizarzr/dot-test/config"
	"github.com/mnizarzr/dot-test/db"
	_ "github.com/mnizarzr/dot-test/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func splash() {
	fmt.Print(`
  _____ _______ _______
 |  __ \__   __|__   __|
 | |  | | | |     | |
 | |  | | | |     | |
 | |__| | | |     | |
 |_____/  |_|     |_|
`)
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func Setup() {

	splash()

	config, err := config.LoadConfig(".")
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
	}

	_, err = db.NewPostgresGormDb(config.PgUri)

	r := gin.Default()

	r.GET("/", Home)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}

// Home godoc
//
//	@Summary		Show home
//	@Description	show app info
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/ [get]
func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app":         config.Configs.AppName,
		"env":         config.Configs.Env,
		"version":     "0.1.0",
		"status":      "running",
		"server_time": time.Now(),
	})
}
