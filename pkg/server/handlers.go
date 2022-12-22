package server

import (
	"errors"
	"github.com/PatricioNapoli/ethix/pkg/indexer"
	"github.com/PatricioNapoli/ethix/pkg/utils"
	"log"
	"net/http"
	"strings"
)

type NetHandler func(http.ResponseWriter, *http.Request)

func RegisterRoutes(sv *HTTPServer) {
	middlewares := []http.HandlerFunc{responseTypeMiddleware}

	sv.Mux.HandleFunc("/", HandlerWrap([]http.HandlerFunc{}, RootHandler))
	sv.Mux.HandleFunc("/block", HandlerWrap(middlewares, BlockHandler))
	sv.Mux.HandleFunc("/subscribe", HandlerWrap(middlewares, SubscribeHandler))
	sv.Mux.HandleFunc("/transactions", HandlerWrap(middlewares, TransactionsHandler))
}

func responseTypeMiddleware(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func HandlerWrap(middlewares []http.HandlerFunc, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				var err error
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unknown error")
				}
				log.Printf("panic: %s", err.Error())
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		for _, h := range middlewares {
			h(w, r)
		}
		handler(w, r)
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK\n"))
}

func BlockHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idx := ctx.Value("indexer").(*indexer.Indexer)

	blockRes := &BlockResponse{Block: idx.GetCurrentBlock()}
	json, err := utils.ToJSON(blockRes)
	if err != nil {
		log.Printf("error when serializing block response: %s", err)
		ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(json)
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idx := ctx.Value("indexer").(*indexer.Indexer)

	addr := r.URL.Query().Get("addr")

	if valid := ValidateAddr(w, addr); !valid {
		return
	}

	addr = strings.ToLower(addr)

	subRes := &SubscribeResponse{Subscribed: idx.Subscribe(addr)}
	json, err := utils.ToJSON(subRes)
	if err != nil {
		log.Printf("error when serializing sub response: %s", err)
		ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(json)
}

func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idx := ctx.Value("indexer").(*indexer.Indexer)

	addr := r.URL.Query().Get("addr")

	if valid := ValidateAddr(w, addr); !valid {
		return
	}

	addr = strings.ToLower(addr)

	txnsRes := &TransactionsResponse{Transactions: idx.GetTransactions(addr)}
	json, err := utils.ToJSON(txnsRes)
	if err != nil {
		log.Printf("error when serializing sub response: %s", err)
		ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(json)
}

func ValidateAddr(w http.ResponseWriter, addr string) bool {
	if len(addr) != 42 {
		err := errors.New("invalid address string")
		ErrorHandler(w, http.StatusBadRequest, err)
		return false
	}

	return true
}

func ErrorHandler(w http.ResponseWriter, status int, err error) {
	errorRes := &ErrorResponse{Error: err.Error()}
	json, err := utils.ToJSON(errorRes)

	w.Write(json)
	w.WriteHeader(status)
}
