package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"errors"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/volssan/balances/db"
	"github.com/volssan/balances/models"
)

var userIdKey = "userId"

func balance(router chi.Router) {
	router.Route("/{userId}", func(router chi.Router) {
		router.Use(BalanceContext)
		router.Get("/", getBalance)
		router.Put("/", increaseBalance)
		router.Delete("/", decreaseBalance)
		router.Post("/transfer/", transfer)
	})
}

func BalanceContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "userId")
		if userId == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("user ID is required")))
			return
		}
		id, err := strconv.Atoi(userId)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid item ID")))
		}
		ctx := context.WithValue(r.Context(), userIdKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func getBalance(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdKey).(int)

	balance, err := dbInstance.GetBalanceByUserId(userId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}

	if err := render.Render(w, r, &balance); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func decreaseBalance(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdKey).(int)
	sum, err := strconv.ParseFloat(r.FormValue("sum"), 32)
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}

	balance, err := dbInstance.GetBalanceByUserId(userId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}

	if balance.Balance < sum {
		render.Render(w, r, ErrorRenderer(errors.New("not enough balance")))
		return
	}
	balance.Balance -= sum
	balance.UpdatedAt = time.Now()

	if err := dbInstance.UpdateBalance(&balance); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
}

func increaseBalance(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdKey).(int)
	sum, err := strconv.ParseFloat(r.FormValue("sum"), 32)
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}

	balance, err := dbInstance.GetBalanceByUserId(userId)
	if err != nil {
		newBalance := models.Balance{
			UserId: userId,
			Balance: sum,
		}
		 if err = dbInstance.AddBalance(&newBalance); err != nil {
		 	render.Render(w, r, ServerErrorRenderer(err))
		 }
		return
	}

	balance.Balance += sum
	balance.UpdatedAt = time.Now()

	if err := dbInstance.UpdateBalance(&balance); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
}

func transfer(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdKey).(int)
	sum, err := strconv.ParseFloat(r.FormValue("sum"), 32)
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	payeeUserId, err := strconv.Atoi(r.FormValue("payeeUserId"))
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}

	if payeeUserId == userId {
		render.Render(w, r, ErrorRenderer(errors.New("transfer to same user not allowed")))
		return
	}

	balance, err := dbInstance.GetBalanceByUserId(userId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if balance.Balance < sum {
		render.Render(w, r, ErrorRenderer(errors.New("not enough balance")))
		return
	}

	balance.Balance -= sum
	balance.UpdatedAt = time.Now()

	if err := dbInstance.UpdateBalance(&balance); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}

	payeeBalance, err := dbInstance.GetBalanceByUserId(payeeUserId)
	if err != nil {
		newBalance := models.Balance{
			UserId: payeeUserId,
			Balance: sum,
		}
		if err = dbInstance.AddBalance(&newBalance); err != nil {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}

	payeeBalance.Balance += sum
	payeeBalance.UpdatedAt = time.Now()

	if err := dbInstance.UpdateBalance(&payeeBalance); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
}