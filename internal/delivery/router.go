package delivery

import (
	"database/sql"
	"log/slog"
	"strconv"
	"time"

	"github.com/anton-ag/todolist/internal/models"
	"github.com/anton-ag/todolist/internal/repository"
	util "github.com/anton-ag/todolist/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"errors"
	"net/http"
)

type Handler struct {
	db *repository.Storage
}

func NewHandler(db *repository.Storage, port string) {
	h := &Handler{db: db}
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	staticDir := "./web"
	e.Static("/*", staticDir)

	e.POST("/api/task", h.createTask)
	e.GET("/api/tasks", h.getTasks)
	e.GET("/api/task", h.getTask)
	e.PUT("api/task", h.updateTask)
	e.POST("api/task/done", h.completeTask)
	e.DELETE("/api/task", h.deleteTask)
	e.GET("api/nextdate", h.nextDate)

	if err := e.Start("0.0.0.0:" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (h *Handler) createTask(c echo.Context) error {
	var task models.Task

	now := time.Now()

	err := c.Bind(&task)
	if err != nil {
		log.Errorf("Ошибка декодирования: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверный формат данных"})
	}

	if task.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
	}

	if task.Date == "" {
		task.Date = now.Format(models.DateFormat)
	}

	if _, err = time.Parse(models.DateFormat, task.Date); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат времени"})
	}

	if task.Date < now.Format(models.DateFormat) {
		task.Date = now.Format(models.DateFormat)
	}

	if task.Repeat != "" {
		_, err := util.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат правила повтора"})
		}
	}

	log.Print(task)

	id, err := h.db.CreateTask(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при создании задачи"})
	}

	return c.JSON(http.StatusOK, map[string]int{"id": id})
}

func (h *Handler) getTasks(c echo.Context) error {
	tasks, err := h.db.GetTasks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении задач"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"tasks": tasks})
}

func (h *Handler) getTask(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID задачи не предоставлен"})
	}

	task, err := h.db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		}
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении задачи"})
	}

	return c.JSON(http.StatusOK, task)
}

func (h *Handler) completeTask(c echo.Context) error {

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID задачи не предоставлен"})
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный ID задачи"})
	}
	task, err := h.db.GetTask(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
	}

	if task.Repeat == "" {
		err := h.db.DeleteTask(task.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при удалении задачи"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{})
	}

	task.Date, err = util.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Ошибка при обновлении даты"})
	}

	err = h.db.UpdateTask(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении задачи"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func (h *Handler) deleteTask(c echo.Context) error {
	
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID задачи не предоставлен"})
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный ID задачи"})
	}
	_, err = h.db.GetTask(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
	}

	err = h.db.DeleteTask(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при удаление задачи"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func (h *Handler) updateTask(c echo.Context) error {
	var task models.Task
	now := time.Now()

	err := c.Bind(&task)
	if err != nil {
		log.Printf("Ошибка декодирования: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат данных"})
	}

	if task.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный ID задачи"})
	}

	_, err = strconv.Atoi(task.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Некорректный ID задачи"})
	}

	if task.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
	}

	if task.Date == "" {
		task.Date = now.Format(models.DateFormat)
	}
	parsedDate, err := time.Parse(models.DateFormat, task.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат даты"})
	}
	if parsedDate.Before(now) {
		task.Date = now.Format(models.DateFormat)
	}

	if task.Repeat != "" {
		_, err := util.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат правила повтора"})
		}
	}

	_, err = h.db.GetTask(task.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
	}

	log.Printf("Обновление задачи: %+v", task)

	err = h.db.UpdateTask(task)
	if err != nil {
		log.Printf("Ошибка обновления задачи: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении задачи"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Задача успешно обновлена"})
}


func (h *Handler) nextDate(c echo.Context) error {
	now := c.QueryParam("now")
	if now == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "начальное время не предоставлено"})
	}

	date := c.QueryParam("date")
	if now == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "дата не предоставлена"})
	}

	repeat := c.QueryParam("repeat")
	if now == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "правило повторения не предоставлено"})
	}
	
	formatNow, err := time.Parse(models.DateFormat, now)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "некорректный формат даты"})
	}

	nextDate, err := util.NextDate(formatNow, date, repeat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ошибка вычисления следующей даты"})
	}

	return c.String(http.StatusOK, nextDate)
}