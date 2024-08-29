package routers

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/quyld17/chat-app/handlers"
)

func RegisterAPIHandlers(router *echo.Echo, db *sql.DB) {
	//Authentication
	router.POST("/sign-up", func(c echo.Context) error {
		return handlers.SignUp(c, db)
	})
	router.POST("/sign-in", func(c echo.Context) error {
		return handlers.SignIn(c, db)
	})

	// Cart
	// router.GET("/cart-products", middlewares.JWTAuthorize(func(c echo.Context) error {
	// 	selected := c.QueryParam("selected")
	// 	return handlers.GetCartProducts(c, db, selected)
	// }))
	// router.POST("/cart-products", middlewares.JWTAuthorize(func(c echo.Context) error {
	// 	return handlers.AddProductToCart(c, db)
	// }))
	// router.PUT("/cart-products", middlewares.JWTAuthorize(func(c echo.Context) error {
	// 	return handlers.UpdateCartProducts(c, db)
	// }))
	// router.DELETE("/cart-products/:productID", middlewares.JWTAuthorize(func(c echo.Context) error {
	// 	productID := c.Param("productID")
	// 	return handlers.DeleteCartProduct(productID, c, db)
	// }))

}
