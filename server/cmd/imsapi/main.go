package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ims-intro/pkg/common/app"
	"ims-intro/pkg/common/postgresql"
	"ims-intro/pkg/controller"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service"
	"log"
	"os"
	"path/filepath"
)