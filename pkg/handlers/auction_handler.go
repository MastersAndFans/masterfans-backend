package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MastersAndFans/masterfans-backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type AuctionHandler struct {
	AuctionRepo repository.IAuctionRepository
}

func NewAuctionHandler(auctionRepo repository.IAuctionRepository) *AuctionHandler {
	return &AuctionHandler{AuctionRepo: auctionRepo}
}

func (handler *AuctionHandler) ListAuctions(w http.ResponseWriter, r *http.Request) {
	auctions, err := handler.AuctionRepo.List(context.Background())

	auctions_json, err := json.Marshal(auctions)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to create JSON"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(auctions_json)
}

func (handler *AuctionHandler) GetAuctionById(w http.ResponseWriter, r *http.Request) {
	id_string := chi.URLParam(r, "id")

	// convert id to int64
	id, err := strconv.ParseInt(id_string, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "ID must contain only digits"})
		return
	}

	auction, err := handler.AuctionRepo.FindById(context.Background(), id)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Auction with ID %d not found", id)})
		return
	}

	auction_json, err := json.Marshal(auction)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to create JSON"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(auction_json)
}
