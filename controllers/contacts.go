package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/himanshukumar42/enterprise/forms"
	"github.com/himanshukumar42/enterprise/models"
)

type ContactsController struct{}

var contactsModel = new(models.ContactsModel)
var contactsForm = new(forms.ContactsForm)

func (ctrl ContactsController) Create(c *gin.Context) {
	userID := getUserID(c)

	var form forms.CreateContactsForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := contactsForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	id, err := contactsModel.Create(userID, form)
	fmt.Println("Error: ", err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Contacts could not be created"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Contacts created", "id": id})
}

func (ctrl ContactsController) All(c *gin.Context) {
	userID := getUserID(c)

	result, err := contactsModel.All(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "could not get contacts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": result})
}

func (ctrl ContactsController) One(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	getID, err := strconv.ParseInt(id, 10, 64)
	if getID == 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "invalid parameter"})
		return
	}
	result, err := contactsModel.One(userID, getID)
	fmt.Println("Error: ", err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "contact not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (ctrl ContactsController) Update(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	getID, err := strconv.ParseInt(id, 10, 64)
	if getID == 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "invalid parameter"})
		return
	}

	var form forms.CreateContactsForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := contactsForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	err = contactsModel.Update(userID, getID, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Contacts could not be updated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact Updated"})
}

func (ctrl ContactsController) PartialUpdate(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	getID, err := strconv.ParseInt(id, 10, 64)
	if getID == 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "invalid parameter"})
		return
	}

	var form forms.UpdateContactsForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := contactsForm.Update(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	err = contactsModel.PartialUpdate(userID, getID, form)
	fmt.Println("Error: ", err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Contacts could not be updated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact Updated"})
}

func (ctrl ContactsController) Delete(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	getID, err := strconv.ParseInt(id, 10, 64)
	if getID == 0 || err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "invalid parameter"})
		return
	}

	err = contactsModel.Delete(userID, getID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Contact could not be deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact Deleted"})
}
