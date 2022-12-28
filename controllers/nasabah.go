package controllers

import (
	"net/http"
	"example.com/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type NasabahInput struct {
	Nim  string `json:"nim"`
	Nama string `json:"nama"`
}

// Tampil data Nasabah
func NasabahTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var nasabah []models.Nasabah
	db.Find(&nasabah)
	c.JSON(http.StatusOK, gin.H{"data": nasabah})
}

// Tambah data Nasabah
func NasabahTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	//validasi input/masukkan
	var dataInput NasabahInput
	if err := c.ShouldBindJSON(&dataInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//proses input
	nasabah := models.Nasabah{
		Nim:  dataInput.Nim,
		Nama: dataInput.Nama,
	}

	db.Create(&nasabah)

	c.JSON(http.StatusOK, gin.H{"data": nasabah})
}

// Ubah data Nasabah
func NasabahUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	//cek dulu datanya
	var nasabah models.Nasabah
	if err := db.Where("nim = ?", c.Param("nim")).First(&nasabah).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data Nasabah tidak ditemukan"})
		return
	}

	//validasi input/masukkan
	var dataInput NasabahInput
	if err := c.ShouldBindJSON(&dataInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//proses ubah data
	db.Model(&nasabah).Update(dataInput)

	c.JSON(http.StatusOK, gin.H{"data": nasabah})
}

// Hapus data Nasabah
func NasabahHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	//cek dulu datanya
	var nasabah models.Nasabah
	if err := db.Where("nim = ?", c.Param("nim")).First(&nasabah).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data Nasabah tidak ditemukan"})
		return
	}

	//proses hapus data
	db.Delete(&nasabah)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
