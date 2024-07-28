package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Expression struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`
	Expr      string  `gorm:"not null"`
	Status    string  `gorm:"not null"`
	Result    float64 `gorm:"not null"`
	Tasks     []*Task `gorm:"-"`
	TaskCount int     `gorm:"-"`
}

type Task struct {
	gorm.Model
	ExpressionID uint   `gorm:"not null"`
	Arg1         string `gorm:"not null"`
	Arg2         string `gorm:"not null"`
	Operation    string `gorm:"not null"`
	OpTime       int    `gorm:"not null"`
}

var tasks = make(chan *Task, 100)
var mutex = &sync.Mutex{}

func NewExpression(userID uint, expr string) *Expression {
	return &Expression{
		UserID:    userID,
		Expr:      expr,
		Status:    "in-progress",
		Result:    0,
		TaskCount: 0,
	}
}

func (e *Expression) ParseAndQueueTasks() {
	parts := strings.Fields(e.Expr)
	for i := 0; i < len(parts); i += 2 {
		task := &Task{
			ExpressionID: e.ID,
			Arg1:         parts[i],
			Arg2:         parts[i+2],
			Operation:    parts[i+1],
			OpTime:       1000,
		}
		e.Tasks = append(e.Tasks, task)
		tasks <- task
		e.TaskCount++
	}
	db.Save(e)
}

func AddExpression(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	expr := NewExpression(user.ID, req.Expression)
	if err := db.Create(&expr).Error; err != nil {
		http.Error(w, "Could not create expression", http.StatusInternalServerError)
		return
	}

	go expr.ParseAndQueueTasks()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": strconv.Itoa(int(expr.ID))})
}

func GetExpressions(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)
	var expressions []Expression
	if err := db.Where("user_id = ?", user.ID).Find(&expressions).Error; err != nil {
		http.Error(w, "Could not retrieve expressions", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string][]Expression{"expressions": expressions})
}

func GetExpression(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)
	vars := mux.Vars(r)
	id := vars["id"]
	var expr Expression
	if err := db.Where("user_id = ? AND id = ?", user.ID, id).First(&expr).Error; err != nil {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]Expression{"expression": expr})
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-tasks:
		json.NewEncoder(w).Encode(map[string]*Task{"task": task})
	default:
		http.Error(w, "No tasks available", http.StatusNotFound)
	}
}

func PostResult(w http.ResponseWriter, r *http.Request) {
	var result struct {
		ID     uint    `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	var task Task
	if err := db.First(&task, result.ID).Error; err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	var expr Expression
	if err := db.First(&expr, task.ExpressionID).Error; err != nil {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	expr.Result += result.Result
	expr.TaskCount--
	if expr.TaskCount == 0 {
		expr.Status = "completed"
	}

	if err := db.Save(&expr).Error; err != nil {
		http.Error(w, "Could not update expression", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
