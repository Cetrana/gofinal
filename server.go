package main

import (
	"fmt"
	"github.com/Cetrana/gofinal/customer"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Start(port string) {

	err := customer.InitServer()
	if err != nil {
		fmt.Println(err)
	}

	router := gin.Default()
	router.Use(checkAuthorization)
	customerList := router.Group("/customers")
	{
		customerList.POST("/", PostCustomers)
		customerList.GET("/", GetCustomers)
		customerList.GET("/:id", GetCustomer)
		customerList.PUT("/:id", PutCustomer)
		customerList.DELETE("/:id", DelCustomer)
	}
	router.Run(port)

}

func checkAuthorization(c *gin.Context) {
	token := c.Request.Header["Authorization"][0]
	fmt.Println(token)
	requiredToken := "November 10, 2009"
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message":"Authorization token required"})
		return
	}

	if token != requiredToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message":"Invalid Authorization token"})
		return
	}
	c.Next()

}

func GetCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	if customer, err := customer.Show(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, customer)
	}

}

func GetCustomers(c *gin.Context) {
	if customers, err := customer.Index(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, customers)
	}
}

func PostCustomers(c *gin.Context) {
	var cust customer.Customer
	if err := c.ShouldBind(&cust); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	if cust, err := customer.Insert(cust); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		fmt.Println(cust)
		c.JSON(http.StatusCreated, cust)
	}

}

func PutCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	var cust customer.Customer
	if err := c.ShouldBind(&cust); err != nil {
		c.String(http.StatusBadRequest, "bad request")
	} else {
		if customer, err := customer.Update(id, cust); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, customer)
		}
	}
}

func DelCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	if err := customer.Delete(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
	}
}
